# G.5 Concorrenza e servizi HTTP

!!! obiettivi "Obiettivi"
    - Capire la differenza tra concorrenza e parallelismo.
    - Avviare lavori concorrenti con le goroutine e aspettarli con `sync.WaitGroup`.
    - Far comunicare le goroutine con i channel e usare `select`.
    - Proteggere i dati condivisi con `sync.Mutex` e trovare gli errori con il *race detector*.
    - Usare `context` per timeout e annullamenti.
    - Scrivere un server HTTP con `net/http`, log strutturati con `slog` e arresto pulito.

!!! info "Quando studiarla"
    Nella [Fase 2](../fase-2/index.md), dopo la [lezione 2.8](../fase-2/08-flask.md). Hai appena scritto un servizio
    HTTP in Python con Flask: qui scrivi un servizio in Go usando solo la libreria standard.

## Perché la concorrenza

Il `portcheck` della [lezione G.3](03-pacchetti-rete.md) controlla le porte una alla volta: con dieci porte filtrate
da un firewall aspetti dieci volte il timeout. Ma il programma, in quel tempo, **non fa nulla**: aspetta la rete.
Sarebbe molto meglio avviare tutti i controlli insieme e raccogliere i risultati man mano che arrivano.

- **Concorrenza** significa *gestire* più lavori nello stesso periodo, alternandoli mentre aspettano.
- **Parallelismo** significa *eseguire* più lavori nello stesso istante, su più core.

Go è nato per la concorrenza: un server deve servire migliaia di client che, quasi sempre, stanno aspettando
qualcosa. Il *runtime* di Go distribuisce i lavori concorrenti sui core disponibili, quindi spesso ottieni anche il
parallelismo senza fare nulla.

## Le goroutine

Una **goroutine** è una funzione che gira in modo concorrente con il resto del programma. Si avvia con la parola
`go`:

```go
go controlla("lab:22")     // parte e il programma prosegue subito, senza aspettare
```

Una goroutine costa pochissimo (pochi kilobyte di memoria): un programma può averne centinaia di migliaia. Non sono
thread del sistema operativo: è il runtime di Go a distribuirle su un numero ridotto di thread.

Il problema è sapere **quando hanno finito**. Se `main` termina, il programma termina e tutte le goroutine con lui.
Lo strumento più semplice per aspettare è `sync.WaitGroup`:

```go
var wg sync.WaitGroup
for _, a := range indirizzi {
	wg.Go(func() {             // (1)!
		fmt.Println(controlla(a, 2*time.Second))
	})
}
wg.Wait()                      // blocca finché tutte le funzioni avviate con wg.Go non sono finite
```

1.  `wg.Go` (da Go 1.25) avvia la funzione in una goroutine e la conta. Nel codice più vecchio troverai la forma
    equivalente `wg.Add(1)`, `go func() { defer wg.Done(); ... }()`.

Da Go 1.22 ogni giro del `for` ha la **sua** copia di `a`: la goroutine vede l'indirizzo giusto. Nelle versioni
precedenti era una trappola classica, e la troverai citata in molti articoli.

## I channel

Un **channel** è un tubo tipizzato attraverso cui le goroutine si passano valori. Il motto di Go è: *non comunicare
condividendo la memoria; condividi la memoria comunicando*.

```go
risultati := make(chan Esito)           // channel di Esito

for _, a := range indirizzi {
	go func() {
		risultati <- controlla(a, 2*time.Second)   // invia
	}()
}
for range indirizzi {
	e := <-risultati                    // riceve: blocca finché non arriva un valore
	fmt.Println(e)
}
```

Così stampi ogni risultato **appena arriva**, non nell'ordine in cui hai avviato i controlli.

| Operazione | Sintassi | Note |
|---|---|---|
| Creare | `make(chan T)` | senza buffer: chi invia aspetta che qualcuno riceva |
| Creare con buffer | `make(chan T, 10)` | chi invia aspetta solo se il buffer è pieno |
| Inviare | `ch <- v` | |
| Ricevere | `v := <-ch` oppure `v, ok := <-ch` | `ok` è `false` se il channel è chiuso e vuoto |
| Chiudere | `close(ch)` | solo chi invia; segnala "non arriverà altro" |
| Scorrere | `for v := range ch` | finisce quando il channel viene chiuso |

### Limitare il parallelismo

Mille controlli tutti insieme possono esaurire le connessioni o sembrare un attacco. Un channel con buffer funziona
da **semaforo**: ha N posti, e ogni goroutine ne occupa uno mentre lavora.

```go
sem := make(chan struct{}, 8)           // al massimo 8 controlli insieme
for i, t := range targets {
	wg.Go(func() {
		sem <- struct{}{}               // occupa un posto (aspetta se sono tutti occupati)
		defer func() { <-sem }()        // lo libera alla fine
		results[i] = c.One(ctx, t)      // ogni goroutine scrive in una posizione diversa: nessun conflitto
	})
}
wg.Wait()
```

`struct{}` è un tipo che non occupa memoria: si usa quando conta il segnale, non il valore.

### select

`select` aspetta su più channel insieme e procede con il primo pronto:

```go
select {
case e := <-risultati:
	fmt.Println(e)
case <-time.After(5 * time.Second):
	fmt.Println("troppo lento, smetto di aspettare")
}
```

## Dati condivisi e race condition

Se due goroutine leggono e scrivono la stessa variabile senza coordinarsi, hai una **race condition**: il risultato
dipende da chi arriva prima, e cambia da un'esecuzione all'altra. Sono gli errori più difficili da trovare, perché
spesso non si presentano mai durante le prove.

```go
type Stato struct {
	mu      sync.RWMutex          // RWMutex: tanti lettori insieme, oppure un solo scrittore
	results []Result
}

func (s *Stato) Aggiorna(r []Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results = r
}

func (s *Stato) Leggi() []Result {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.results
}
```

Go include un **race detector** che controlla, mentre il programma gira, gli accessi non protetti:

```console
$ go test -race ./...
==================
WARNING: DATA RACE
Write at 0x00c000124018 by goroutine 8:
  ...
```

Esegui **sempre** i test con `-race` nella pipeline: costa un po' di tempo e trova errori che altrimenti
arriverebbero in produzione.

## context: timeout e annullamenti

Un `context.Context` porta con sé una **scadenza** e un **segnale di annullamento**, e si passa come primo argomento
a tutte le funzioni che fanno I/O. Quando scade o viene annullato, tutte le operazioni che lo usano si fermano.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()                                  // libera le risorse anche se finisce prima

req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
resp, err := client.Do(req)                     // se passano 5 secondi: "context deadline exceeded"
```

Il caso più importante per un servizio: annullare tutto quando arriva il segnale di arresto.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
// ctx viene annullato con Ctrl+C (SIGINT) o con SIGTERM
```

`SIGTERM` è il segnale che `systemctl stop`, `docker stop` e Kubernetes inviano a un processo per chiedergli di
fermarsi ([lezione 1.4](../fase-1/04-processi-systemd.md)). Un buon servizio lo raccoglie, finisce le richieste in
corso e poi esce.

## Un server HTTP con la libreria standard

Il pacchetto `net/http` contiene un server HTTP pronto per la produzione: è lo stesso che usano molti dei progetti
cloud native che incontrerai.

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {   // (1)!
	w.Write([]byte("ok\n"))
})
mux.HandleFunc("GET /targets/{nome}", func(w http.ResponseWriter, r *http.Request) {
	nome := r.PathValue("nome")                                                // (2)!
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"target": nome})
})

srv := &http.Server{Addr: ":9101", Handler: mux, ReadHeaderTimeout: 5 * time.Second}   // (3)!
log.Fatal(srv.ListenAndServe())
```

1.  Da Go 1.22 il *pattern* può contenere il **metodo**: una `POST /healthz` riceve automaticamente
    `405 Method Not Allowed`.
2.  Le parti tra graffe del percorso sono **variabili**, lette con `r.PathValue`.
3.  `ReadHeaderTimeout` evita che un client lentissimo tenga occupata una connessione per sempre (l'attacco
    *Slowloris*). Il `http.ListenAndServe` "nudo" degli esempi rapidi non ha timeout.

Ogni richiesta viene gestita in una **goroutine sua**: il server è concorrente senza che tu scriva nulla. Per lo
stesso motivo, se gli handler leggono dati che un'altra goroutine aggiorna, servono il mutex o i channel.

### Log strutturati con slog

Nella [Fase 7](../fase-7/08-loki.md) i log finiranno in Loki, dove i log in **JSON** si filtrano molto meglio del
testo libero. Il pacchetto `log/slog` li produce direttamente:

```go
log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
log.Info("probe avviato", "addr", ":9101", "target", 3)
// {"time":"2026-11-01T10:00:00Z","level":"INFO","msg":"probe avviato","addr":":9101","target":3}
```

### Arresto pulito

```go
go func() { errCh <- srv.ListenAndServe() }()

select {
case err := <-errCh:                  // il server non è nemmeno partito (porta occupata?)
	return err
case <-ctx.Done():                    // arrivato SIGTERM o Ctrl+C
}
shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
return srv.Shutdown(shutdownCtx)      // smette di accettare connessioni, aspetta quelle in corso
```

## Laboratorio: probe serve

!!! laboratorio "Laboratorio · probe diventa un servizio"
    Porta avanti il progetto della [lezione G.4](04-cli-test.md#laboratorio-racoon-probe-primo-passo).

    1. Aggiungi a `internal/check` il metodo `All(ctx, targets, parallel) []Result`, che esegue i controlli in
       parallelo (al massimo `parallel` insieme) e restituisce i risultati **nello stesso ordine** dei target.
       Usalo anche in `probe check`.
    2. Crea il pacchetto `internal/server` con un tipo `Server` che contiene la configurazione e l'ultimo giro di
       risultati, protetto da un `sync.RWMutex`.
    3. Scrivi `Loop(ctx)`, che ripete i controlli ogni `interval` finché il contesto non viene annullato
       (usa `time.NewTicker`).
    4. Esponi `GET /healthz` (risponde `ok`) e `GET /results` (l'ultimo giro in JSON, `[]` prima del primo).
    5. Aggiungi il sottocomando `probe serve -addr :9101` con log JSON e arresto pulito su SIGTERM.
    6. Scrivi i test con `net/http/httptest`, che crea server HTTP finti per i test, ed eseguili con `-race`.
    7. Installa `probe` come servizio systemd sulla VM `lab` (o su `srv1` nel livello 2) e verifica che
       `systemctl stop probe` lo fermi in modo pulito, guardando i log con `journalctl -u probe`.

??? success "Soluzione guidata"
    Il codice completo è in
    [`labs/go/probe`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/probe):
    `internal/check/check.go` (metodo `All`), `internal/server/server.go` e `cmd/probe/main.go` (funzione
    `runServe`). Per i test, `httptest.NewServer` avvia un server vero su una porta libera:

    ```go title="internal/check/check_test.go (estratto)"
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    	if r.URL.Path == "/rotto" {
    		w.WriteHeader(http.StatusInternalServerError)
    	}
    }))
    defer srv.Close()

    r := New(time.Second).One(context.Background(),
    	config.Target{Name: "b", Type: "http", URL: srv.URL + "/rotto", ExpectStatus: 200})
    if r.OK {
    	t.Error("un 500 non deve risultare OK")
    }
    ```

    L'unità systemd, con le protezioni viste nella [lezione 1.4](../fase-1/04-processi-systemd.md):

    ```ini title="/etc/systemd/system/probe.service"
    [Unit]
    Description=racoon-probe
    After=network-online.target
    Wants=network-online.target

    [Service]
    ExecStart=/usr/local/bin/probe serve -config /etc/probe/probe.yaml -addr :9101
    DynamicUser=yes
    NoNewPrivileges=yes
    ProtectSystem=strict
    ProtectHome=yes
    Restart=on-failure

    [Install]
    WantedBy=multi-user.target
    ```

    ```bash
    go build -o probe ./cmd/probe && scp probe probe.yaml lab:
    ssh lab 'sudo install -m 0755 probe /usr/local/bin/ && sudo install -D -m 0644 probe.yaml /etc/probe/probe.yaml'
    # copia l'unità, poi:
    ssh lab 'sudo systemctl daemon-reload && sudo systemctl enable --now probe && curl -s localhost:9101/results'
    ssh lab 'sudo systemctl stop probe && journalctl -u probe -n 3 -o cat'
    # {"level":"INFO","msg":"arresto in corso"}
    ```

    **Spiegazione.** `DynamicUser=yes` fa girare il servizio con un utente creato al volo e senza privilegi: un
    binario statico come `probe` non ha bisogno d'altro. Il messaggio "arresto in corso" nei log conferma che
    SIGTERM è stato raccolto: senza `signal.NotifyContext` il processo verrebbe terminato a metà.

## Esercizi

!!! esercizio "Esercizio 1 · portcheck in parallelo"
    Modifica il `portcheck` della [lezione G.3](03-pacchetti-rete.md) perché controlli tutte le porte in parallelo e
    stampi i risultati appena arrivano, usando un channel. Misura la differenza con `time ./portcheck` su cinque
    indirizzi filtrati (ad esempio `10.255.255.1:22`, che non risponde).

??? success "Soluzione"
    ```go
    risultati := make(chan Esito)
    for _, a := range indirizzi {
    	go func() { risultati <- controlla(a, 2*time.Second) }()
    }
    chiuse := 0
    for range indirizzi {
    	e := <-risultati
    	fmt.Println(e)
    	if !e.Aperta {
    		chiuse++
    	}
    }
    ```
    **Spiegazione.** In sequenza cinque timeout costano circa 10 secondi; in parallelo circa 2, il tempo del più
    lento. Si ricevono esattamente tanti valori quante goroutine sono partite, quindi non serve un `WaitGroup` e non
    serve chiudere il channel.

!!! esercizio "Esercizio 2 · Trova la race"
    Questo contatore di richieste ha un errore. Scrivi un test che chiama l'handler da 100 goroutine, eseguilo con
    `go test -race` e correggi il codice.

    ```go
    type Contatore struct{ n int }

    func (c *Contatore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    	c.n++
    	fmt.Fprintln(w, c.n)
    }
    ```

??? success "Soluzione"
    ```go title="contatore_test.go"
    func TestContatore(t *testing.T) {
    	c := &Contatore{}
    	var wg sync.WaitGroup
    	for range 100 {
    		wg.Go(func() {
    			c.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
    		})
    	}
    	wg.Wait()
    	if c.n != 100 {
    		t.Errorf("n = %d, voglio 100", c.n)
    	}
    }
    ```
    ```go title="correzione"
    type Contatore struct{ n atomic.Int64 }

    func (c *Contatore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    	fmt.Fprintln(w, c.n.Add(1))
    }
    ```
    **Spiegazione.** `c.n++` sembra un'operazione sola, ma sono tre: leggi, somma, scrivi. Due goroutine possono
    leggere lo stesso valore e perdere un incremento. Per un singolo numero il pacchetto `sync/atomic` è più
    semplice di un mutex; nel test il controllo va adattato a `c.n.Load()`. `httptest.NewRecorder` registra la
    risposta senza bisogno di una rete.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa succede alle goroutine ancora attive quando main termina?","options":["Continuano in background","Vengono terminate con il programma","main aspetta automaticamente","Diventano processi separati"],"answer":[1],"explanation":"Per aspettarle si usa un WaitGroup o un channel."},
 {"q":"A cosa serve `make(chan struct{}, 8)` nel metodo All di probe?","options":["A raccogliere i risultati","A limitare a 8 i controlli contemporanei","A ordinare i target","A gestire i timeout"],"answer":[1],"explanation":"È un semaforo con 8 posti."},
 {"q":"Quale comando trova gli accessi concorrenti non protetti?","options":["go vet","go test -race","go build -v","gofmt -d"],"answer":[1],"explanation":"Il race detector controlla gli accessi mentre il codice gira."},
 {"q":"Cosa succede a una richiesta HTTP creata con un context scaduto dopo 5 secondi?","options":["Continua finché il server risponde","Si interrompe con l'errore context deadline exceeded","Viene ripetuta","Il programma va in panic"],"answer":[1],"explanation":"Il context propaga la scadenza a tutte le operazioni che lo usano."},
 {"q":"Quale segnale inviano docker stop e Kubernetes per chiedere a un processo di fermarsi?","options":["SIGKILL","SIGTERM","SIGHUP","SIGSTOP"],"answer":[1],"explanation":"Dopo un periodo di attesa, se il processo non esce, arriva SIGKILL."},
 {"q":"Con `mux.HandleFunc(\"GET /healthz\", ...)`, cosa riceve una richiesta POST /healthz?","options":["200","404","405 Method Not Allowed","500"],"answer":[2],"explanation":"Da Go 1.22 i pattern di ServeMux includono il metodo."}
]}
</script>

## Fonti

- [A Tour of Go: Concurrency](https://go.dev/tour/concurrency/1)
- [Effective Go: Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go Concurrency Patterns: Context](https://go.dev/blog/context) e [Pipelines and cancellation](https://go.dev/blog/pipelines)
- [Data Race Detector](https://go.dev/doc/articles/race_detector)
- [Routing Enhancements for Go 1.22](https://go.dev/blog/routing-enhancements) e [Fixing For Loops in Go 1.22](https://go.dev/blog/loopvar-preview)
- [Structured Logging with slog](https://go.dev/blog/slog)
- Pacchetti [sync](https://pkg.go.dev/sync), [context](https://pkg.go.dev/context), [net/http](https://pkg.go.dev/net/http), [net/http/httptest](https://pkg.go.dev/net/http/httptest)

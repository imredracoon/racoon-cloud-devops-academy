# G.4 Una CLI vera, con i test

!!! obiettivi "Obiettivi"
    - Leggere opzioni e sottocomandi da riga di comando con il pacchetto `flag`.
    - Leggere e scrivere JSON e YAML con le *struct tag*.
    - Eseguire comandi esterni in modo sicuro con `os/exec`.
    - Scrivere test automatici con `go test`: test a tabella, sottotest, file temporanei, copertura.
    - Dare a un progetto la struttura `cmd/` + `internal/` e iniziare **racoon-probe**.

!!! info "Quando studiarla"
    Nella [Fase 2](../fase-2/index.md), dopo la [lezione 2.7](../fase-2/07-python-automazione.md). Vedrai le stesse idee
    della lezione Python (argomenti, JSON e YAML, `subprocess`) scritte in Go: confrontarle aiuta a fissarle.

## Opzioni da riga di comando

Il pacchetto `flag` della libreria standard legge opzioni come `-config file.yaml` o `-v`:

```go
config := flag.String("config", "probe.yaml", "file di configurazione")   // nome, predefinito, descrizione
timeout := flag.Duration("timeout", 5*time.Second, "tempo massimo per controllo")
verboso := flag.Bool("v", false, "stampa più dettagli")
flag.Parse()

fmt.Println(*config, *timeout, *verboso)   // flag.String restituisce un puntatore: si legge con *
fmt.Println(flag.Args())                    // ciò che resta dopo le opzioni
```

```console
$ ./probe -config /etc/probe.yaml -timeout 2s -v extra
/etc/probe.yaml 2s true
[extra]
$ ./probe -h
Usage of ./probe:
  -config string
    	file di configurazione (default "probe.yaml")
  ...
```

L'aiuto con `-h` è generato in automatico. `flag.Duration` capisce direttamente `500ms`, `2s`, `1m30s`.

### Sottocomandi

Molti tool DevOps hanno **sottocomandi**: `git commit`, `kubectl get`, `terraform plan`. Con la libreria standard
si ottengono creando un `flag.FlagSet` per ciascuno:

```go
switch os.Args[1] {
case "check":
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	formato := fs.String("o", "table", "table o json")
	fs.Parse(os.Args[2:])
	// ...
case "serve":
	// ...
}
```

Per CLI molto ricche esistono librerie come [Cobra](https://cobra.dev/) (la usano `kubectl`, `helm` e `gh`), ma
per i nostri strumenti `flag` basta.

## JSON e YAML

Go legge JSON e YAML **dentro struct**, guidato dalle **struct tag**: piccole etichette tra apici inversi che dicono
quale chiave del file corrisponde a quale campo.

```go
type Target struct {
	Name         string `json:"name" yaml:"name"`
	Type         string `json:"type" yaml:"type"`
	URL          string `json:"url,omitempty" yaml:"url,omitempty"`     // omitempty: non scrivere se vuoto
	ExpectStatus int    `json:"expect_status,omitempty" yaml:"expect_status,omitempty"`
}
```

```go
// JSON → struct
var t Target
err := json.Unmarshal([]byte(`{"name":"api","type":"http","url":"http://localhost:8000"}`), &t)

// struct → JSON
out, err := json.MarshalIndent(t, "", "  ")
fmt.Println(string(out))
```

Il JSON fa parte della libreria standard (`encoding/json`); per lo YAML si usa il modulo
`go.yaml.in/yaml/v3`, mantenuto dal progetto YAML ufficiale, che funziona allo stesso modo:

```bash
go get go.yaml.in/yaml/v3
```

```go
var cfg Config
err := yaml.Unmarshal(dati, &cfg)   // i campi time.Duration accettano valori come "30s"
```

!!! tip "Le struct sono lo schema"
    In Python un file YAML diventa un dizionario in cui puoi trovare qualsiasi cosa. In Go diventa una struct con
    tipi precisi: se nel file c'è `cpu: "due"` dove ti aspetti un numero, l'errore arriva **subito**, alla
    lettura. Una chiave scritta male (`expect_staus`) invece viene ignorata in silenzio: per questo la
    configurazione si **valida** sempre dopo averla letta.

## Eseguire comandi esterni

`os/exec` lancia un programma passando gli argomenti come **elenco**, senza passare dalla shell. È l'equivalente di
`subprocess.run([...])` della [lezione 2.7](../fase-2/07-python-automazione.md) e, per lo stesso motivo, è sicuro
anche con input che arriva dall'utente.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

out, err := exec.CommandContext(ctx, "systemctl", "is-active", servizio).Output()
stato := strings.TrimSpace(string(out))

var exitErr *exec.ExitError
if errors.As(err, &exitErr) {
	fmt.Println("codice di uscita:", exitErr.ExitCode())   // systemctl is-active esce con 3 se il servizio è fermo
}
```

`errors.As` è il fratello di `errors.Is`: invece di confrontare con un errore preciso, controlla se l'errore è di
un certo **tipo** e, se sì, te lo dà per leggerne i dettagli.

## I test

In Go i test sono parte del linguaggio. Un file che finisce in `_test.go` contiene funzioni `TestQualcosa`, e
`go test` le trova e le esegue.

```go title="livello_test.go"
package main

import "testing"

func TestLivello(t *testing.T) {
	casi := []struct {            // (1)!
		uso  int
		vuoi string
	}{
		{0, "ok"},
		{79, "ok"},
		{80, "attenzione"},
		{90, "critico"},
		{100, "critico"},
	}
	for _, c := range casi {
		got := livello(c.uso)
		if got != c.vuoi {
			t.Errorf("livello(%d) = %q, voglio %q", c.uso, got, c.vuoi)   // (2)!
		}
	}
}
```

1.  Un **test a tabella** (*table-driven test*): un elenco di casi, ognuno con l'input e il risultato atteso, e un
    unico ciclo che li prova tutti. È lo stile più diffuso nei progetti Go: aggiungere un caso è una riga.
2.  `t.Errorf` segnala il fallimento e **continua** con gli altri casi; `t.Fatalf` si ferma subito. Il messaggio
    segue la convenzione "ho ottenuto X, voglio Y".

```console
$ go test ./...
ok  	disco	0.002s
$ go test -v -run TestLivello ./...       # dettagliato, solo i test che corrispondono al nome
$ go test -cover ./...                    # percentuale di codice coperta dai test
ok  	disco	0.002s	coverage: 45.5% of statements
```

Altri strumenti di `testing` che userai:

| Strumento | A cosa serve |
|---|---|
| `t.Run("nome", func(t *testing.T) {...})` | un **sottotest** con un nome, eseguibile da solo con `-run 'TestX/nome'` |
| `t.TempDir()` | una cartella temporanea, cancellata alla fine del test |
| `t.Helper()` | segna una funzione di supporto: gli errori indicano la riga del test, non quella dell'helper |
| `t.Parallel()` | il test può girare in parallelo agli altri |
| `go test -race` | cerca accessi concorrenti non protetti ([lezione G.5](05-concorrenza-http.md)) |

## La struttura di un progetto

Per un programma con più pacchetti, la comunità Go usa di solito questa struttura (descritta nella guida
ufficiale [Organizing a Go module](https://go.dev/doc/modules/layout)):

```text
probe/
├── go.mod
├── go.sum
├── probe.example.yaml
├── cmd/
│   └── probe/
│       └── main.go          # solo argomenti, input/output e codici di uscita
└── internal/
    ├── config/              # leggere e validare probe.yaml
    │   ├── config.go
    │   └── config_test.go
    └── check/               # eseguire i controlli
        ├── check.go
        └── check_test.go
```

`main.go` resta sottile: tutta la logica sta nei pacchetti sotto `internal`, dove è facile da testare.

## Laboratorio: racoon-probe, primo passo

!!! laboratorio "Laboratorio · probe check"
    Inizia il progetto del percorso: un comando `probe check` che legge un file YAML con un elenco di target e li
    controlla.

    ```yaml title="probe.yaml"
    interval: 30s
    timeout: 5s
    targets:
      - name: tasks-api
        type: http
        url: http://localhost:8000/health
        expect_status: 200
      - name: postgres
        type: tcp
        address: localhost:5432
    ```

    1. Crea il repository `racoon-probe` su GitHub e il modulo con `go mod init github.com/TUONOME/racoon-probe`.
    2. Nel pacchetto `internal/config` scrivi le struct `Config` e `Target`, una funzione `Parse(dati []byte)` che
       applica i predefiniti (`interval` 30s, `timeout` 5s, `expect_status` 200) e un metodo `Validate()` che
       restituisce **tutti** gli errori insieme con `errors.Join`: nessun target, nome mancante o duplicato, tipo
       sconosciuto, `url` o `address` mancante, `timeout` non minore di `interval`.
    3. Scrivi `config_test.go` con un test a tabella che copra ogni errore.
    4. Nel pacchetto `internal/check` scrivi `Checker.One(ctx, target) Result`: per `http` una richiesta GET che
       confronta il codice di stato, per `tcp` una connessione.
    5. In `cmd/probe/main.go` scrivi il sottocomando `check` con le opzioni `-config` e `-o table|json`. Esce con
       codice 1 se un controllo fallisce.
    6. Avvia la Tasks API della [Fase 2](../fase-2/08-flask.md) e prova `probe check`. Poi fermala e riprova.

??? success "Soluzione guidata"
    La versione di riferimento completa è in
    [`labs/go/probe`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/go/probe). I punti
    chiave:

    ```go title="internal/config/config.go (estratto)"
    func (c Config) Validate() error {
    	var errs []error
    	if len(c.Targets) == 0 {
    		errs = append(errs, errors.New("nessun target definito"))
    	}
    	visti := map[string]bool{}
    	for i, t := range c.Targets {
    		switch {
    		case t.Name == "":
    			errs = append(errs, fmt.Errorf("target %d: name mancante", i+1))
    		case visti[t.Name]:
    			errs = append(errs, fmt.Errorf("target %q: nome duplicato", t.Name))
    		}
    		visti[t.Name] = true
    		// ... controllo di type, url e address
    	}
    	return errors.Join(errs...)   // nil se errs è vuota
    }
    ```

    ```go title="internal/check/check.go (estratto)"
    func (c *Checker) http(ctx context.Context, t config.Target) error {
    	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.URL, nil)
    	if err != nil {
    		return err
    	}
    	resp, err := c.Client.Do(req)
    	if err != nil {
    		return err
    	}
    	defer resp.Body.Close()   // il corpo va sempre chiuso, anche se non lo leggi
    	if resp.StatusCode != t.ExpectStatus {
    		return fmt.Errorf("stato HTTP %d, atteso %d", resp.StatusCode, t.ExpectStatus)
    	}
    	return nil
    }
    ```

    ```go title="cmd/probe/main.go (estratto)"
    for _, t := range cfg.Targets {          // in sequenza: nella lezione G.5 diventerà parallelo
    	res = append(res, checker.One(context.Background(), t))
    }
    ```

    ```console
    $ go run ./cmd/probe check -config probe.yaml
    TARGET     TIPO  ESITO  DURATA  ERRORE
    tasks-api  http  OK     3ms
    postgres   tcp   KO     0s      dial tcp 127.0.0.1:5432: connect: connection refused
    errore: almeno un controllo è fallito
    $ echo $?
    1
    ```

    **Spiegazione.** `Validate` raccoglie tutti gli errori così chi scrive il file li corregge in un colpo solo.
    La tabella si allinea con `text/tabwriter`. Nella versione di riferimento `main` chiama una funzione
    `run(args, out io.Writer) error`: in questo modo anche il comando intero si testa senza lanciare processi
    (vedi `cmd/probe/main_test.go`).

## Esercizi

!!! esercizio "Esercizio 1 · Testare il disco"
    Riprendi il programma `disco` della [lezione G.1](01-primi-passi.md#esercizi) e scrivi un test a tabella per la
    funzione `livello`, con un sottotest per ogni caso (`t.Run`). Misura la copertura con `go test -cover`. Perché
    non arriva al 100%?

??? success "Soluzione"
    ```go title="main_test.go"
    func TestLivello(t *testing.T) {
    	casi := map[string]struct {
    		uso  int
    		vuoi string
    	}{
    		"zero":            {0, "ok"},
    		"sotto soglia":    {79, "ok"},
    		"soglia gialla":   {80, "attenzione"},
    		"soglia rossa":    {90, "critico"},
    		"pieno":           {100, "critico"},
    	}
    	for nome, c := range casi {
    		t.Run(nome, func(t *testing.T) {
    			if got := livello(c.uso); got != c.vuoi {
    				t.Errorf("livello(%d) = %q, voglio %q", c.uso, got, c.vuoi)
    			}
    		})
    	}
    }
    ```
    **Spiegazione.** La copertura non tiene conto di `main`, che legge `os.Args` e chiama `os.Exit`: per testarla
    andrebbe trasformata in una funzione che riceve gli argomenti e restituisce un errore, come la `run` di probe.
    È il motivo per cui la logica si tiene fuori da `main`. I casi sulle soglie (79, 80, 90) sono i più preziosi:
    gli errori di "minore o minore-uguale" si nascondono lì.

!!! esercizio "Esercizio 2 · Stato dei servizi"
    Scrivi `servizi`, che riceve dei nomi di unità systemd e stampa per ognuno lo stato restituito da
    `systemctl is-active`, in formato JSON se si passa `-json`. Provalo sulla VM `lab` con `ssh nginx cron
    inesistente`. Attenzione: `systemctl is-active` esce con un codice diverso da zero quando il servizio non è
    attivo, ma l'output resta utile.

??? success "Soluzione"
    ```go title="main.go"
    package main

    import (
    	"encoding/json"
    	"errors"
    	"flag"
    	"fmt"
    	"os"
    	"os/exec"
    	"strings"
    )

    type Stato struct {
    	Unita string `json:"unit"`
    	Stato string `json:"state"`
    }

    func stato(unita string) (string, error) {
    	out, err := exec.Command("systemctl", "is-active", unita).Output()
    	var exitErr *exec.ExitError
    	if err != nil && !errors.As(err, &exitErr) {
    		return "", err          // systemctl non trovato o non eseguibile: errore vero
    	}
    	return strings.TrimSpace(string(out)), nil
    }

    func main() {
    	comeJSON := flag.Bool("json", false, "uscita in JSON")
    	flag.Parse()

    	var stati []Stato
    	for _, u := range flag.Args() {
    		s, err := stato(u)
    		if err != nil {
    			fmt.Fprintln(os.Stderr, "errore:", err)
    			os.Exit(2)
    		}
    		stati = append(stati, Stato{u, s})
    	}
    	if *comeJSON {
    		json.NewEncoder(os.Stdout).Encode(stati)
    		return
    	}
    	for _, s := range stati {
    		fmt.Printf("%-15s %s\n", s.Unita, s.Stato)
    	}
    }
    ```
    ```console
    $ ./servizi -json ssh nginx inesistente
    [{"unit":"ssh","state":"active"},{"unit":"nginx","state":"active"},{"unit":"inesistente","state":"inactive"}]
    ```
    **Spiegazione.** `exec.ExitError` significa "il comando è partito ed è uscito con un codice diverso da zero": per
    `is-active` è un'informazione, non un guasto. Se invece `systemctl` non esiste (nel Dev Container, per
    esempio) l'errore è di un altro tipo e va segnalato. Su AlmaLinux l'unità si chiama `sshd`, non `ssh`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa restituisce `flag.String(\"config\", \"probe.yaml\", \"...\")`?","options":["Una stringa","Un puntatore a stringa, valorizzato dopo flag.Parse()","Un errore","Una slice"],"answer":[1],"explanation":"Si legge con *config dopo aver chiamato flag.Parse()."},
 {"q":"A cosa servono le struct tag come `yaml:\"expect_status\"`?","options":["A documentare il campo","A dire quale chiave del file corrisponde al campo","A validare il valore","A rendere il campo privato"],"answer":[1],"explanation":"Le usano encoding/json, yaml e molte altre librerie."},
 {"q":"Perché `exec.Command(\"systemctl\", \"is-active\", nome)` è sicuro anche se nome arriva dall'utente?","options":["Perché systemctl valida l'input","Perché gli argomenti sono passati come elenco, senza interpretazione della shell","Perché Go filtra i caratteri speciali","Non è sicuro"],"answer":[1],"explanation":"Nessuna shell, quindi nessuna command injection."},
 {"q":"Qual è la differenza tra t.Errorf e t.Fatalf?","options":["Nessuna","Errorf segnala e continua, Fatalf segnala e interrompe il test","Fatalf ferma tutti i test del pacchetto","Errorf non fa fallire il test"],"answer":[1],"explanation":"Fatalf si usa quando proseguire non ha senso (ad esempio un file non creato)."},
 {"q":"Come esegui solo il sottotest \"soglia rossa\" di TestLivello?","options":["go test -only soglia","go test -run 'TestLivello/soglia_rossa'","go test TestLivello soglia","go test -v"],"answer":[1],"explanation":"Gli spazi nei nomi dei sottotest diventano trattini bassi."},
 {"q":"Nella struttura cmd/ + internal/, dove va la logica principale?","options":["Tutta in main.go","Nei pacchetti sotto internal, dove si testa facilmente","In cmd/","Nei file _test.go"],"answer":[1],"explanation":"main resta sottile: argomenti, output e codici di uscita."}
]}
</script>

## Fonti

- Pacchetti [flag](https://pkg.go.dev/flag), [encoding/json](https://pkg.go.dev/encoding/json), [os/exec](https://pkg.go.dev/os/exec) e [testing](https://pkg.go.dev/testing)
- [go.yaml.in/yaml/v3](https://pkg.go.dev/go.yaml.in/yaml/v3), il modulo YAML del progetto [yaml/go-yaml](https://github.com/yaml/go-yaml)
- [Tutorial: Add a test](https://go.dev/doc/tutorial/add-a-test)
- [Go Wiki: TableDrivenTests](https://go.dev/wiki/TableDrivenTests)
- [Organizing a Go module](https://go.dev/doc/modules/layout)
- [The cover story](https://go.dev/blog/cover), sulla copertura dei test

# G.2 Collezioni, struct ed errori

!!! obiettivi "Obiettivi"
    - Usare array, slice e map, e scorrerli con `range`.
    - Definire tipi propri con le `struct` e aggiungere loro dei metodi.
    - Capire i puntatori quel tanto che serve per usarli bene.
    - Gestire gli errori in modo idiomatico: creare, arricchire e riconoscere un errore.
    - Leggere file di testo riga per riga e usare `defer`.

!!! info "Quando studiarla"
    Durante la [Fase 1](../fase-1/index.md), dopo la [lezione 1.4](../fase-1/04-processi-systemd.md). Nell'esercizio
    finale leggerai da Go gli stessi file di `/proc` e `/etc` che hai esplorato da terminale.

## Slice: elenchi che crescono

Un **array** in Go ha una lunghezza fissa che fa parte del tipo (`[3]string`). Si usa poco. Quasi sempre userai
una **slice**, un elenco ordinato che può crescere:

```go
server := []string{"web01", "web02", "db01"}
fmt.Println(server[0], len(server))   // web01 3
server = append(server, "cache01")    // append restituisce la slice aggiornata: va riassegnata
primi := server[0:2]                  // [web01 web02]: dall'indice 0 al 2 escluso

for i, s := range server {            // range dà indice e valore
	fmt.Println(i, s)
}
for _, s := range server {            // _ ignora un valore che non ti serve
	fmt.Println(s)
}
```

Una slice è una "finestra" su un array nascosto. Per questo `primi` e `server` **condividono gli stessi dati**:
se modifichi `primi[0]`, cambia anche `server[0]`. Se ti serve una copia indipendente usa `slices.Clone(server)`.

Il pacchetto `slices` della libreria standard offre le operazioni più comuni:

```go
slices.Sort(server)                       // ordina sul posto
slices.Contains(server, "db01")           // true
i := slices.Index(server, "web02")        // posizione, oppure -1
s := make([]int, 0, 100)                  // slice vuota con spazio già riservato per 100 elementi
```

## Map: chiavi e valori

Una **map** associa chiavi a valori, come un dizionario Python:

```go
ram := map[string]int{"web01": 4, "db01": 16}
ram["cache01"] = 2                // aggiunge o modifica
delete(ram, "web01")              // rimuove

gb := ram["sconosciuto"]          // 0: una chiave assente restituisce il valore zero, senza errore
gb, ok := ram["db01"]             // ok dice se la chiave c'era davvero
if !ok {
	fmt.Println("non trovato")
}

for nome, gb := range ram {       // l'ordine di una map NON è garantito
	fmt.Println(nome, gb)
}
```

!!! warning "Una map va creata prima di scriverci"
    `var m map[string]int` dichiara una map **nil**: leggerla funziona, ma scriverci manda il programma in
    *panic* (un arresto immediato). Creala sempre con `make(map[string]int)` o con `map[string]int{}`.

Se ti serve un ordine stabile, ordina le chiavi: `for _, k := range slices.Sorted(maps.Keys(ram))`.

## Struct: i tuoi tipi

Una **struct** raggruppa più valori sotto un solo nome. È il modo di Go per descrivere "una cosa" del tuo problema:

```go
type VM struct {
	Nome   string
	CPU    int
	RAMGB  int
	Tag    []string
}

web := VM{Nome: "web01", CPU: 2, RAMGB: 4, Tag: []string{"prod"}}
web.CPU = 4
fmt.Printf("%+v\n", web)   // {Nome:web01 CPU:4 RAMGB:4 Tag:[prod]}

flotta := []VM{web, {Nome: "db01", CPU: 8, RAMGB: 32}}
```

I nomi che iniziano con la **maiuscola** sono **esportati**, cioè visibili dagli altri pacchetti; quelli con la
minuscola restano privati. Vale per struct, campi, funzioni e costanti: in Go non esistono le parole `public` e
`private`, conta solo l'iniziale.

### Metodi

Un **metodo** è una funzione legata a un tipo:

```go
func (v VM) Descrizione() string {
	return fmt.Sprintf("%s (%d vCPU, %d GB)", v.Nome, v.CPU, v.RAMGB)
}

fmt.Println(web.Descrizione())
```

`(v VM)` è il **ricevitore**: dentro il metodo, `v` è la VM su cui l'hai chiamato.

## Puntatori, quanto basta

In Go, quando passi un valore a una funzione, la funzione ne riceve una **copia**. Se vuoi che la funzione modifichi
l'originale, le passi un **puntatore**, cioè l'indirizzo in memoria del valore:

```go
func raddoppia(v *VM) {   // *VM: "puntatore a una VM"
	v.RAMGB *= 2          // Go segue il puntatore da solo: non serve scrivere (*v).RAMGB
}

raddoppia(&web)           // &web: "l'indirizzo di web"
fmt.Println(web.RAMGB)    // 8
```

La stessa regola vale per i metodi: un metodo che deve **modificare** il valore usa un ricevitore puntatore.

```go
func (v *VM) AggiungiTag(t string) { v.Tag = append(v.Tag, t) }
```

Regola pratica: usa ricevitori puntatore quando il metodo modifica la struct o quando la struct è grande; altrimenti
va bene il valore. Un puntatore che non punta a niente vale `nil`, e usarlo provoca un *panic*: per questo le
funzioni controllano `if v == nil` quando un `nil` è possibile.

## Gli errori, sul serio

Hai visto nella [lezione G.1](01-primi-passi.md#funzioni) che un errore è un valore restituito. Ecco gli strumenti
per usarlo bene.

```go
// Creare un errore semplice
err := errors.New("configurazione vuota")

// Creare un errore con dettagli
err = fmt.Errorf("porta %d fuori intervallo", 70000)

// Arricchire un errore ricevuto, conservando l'originale: %w ("wrap")
dati, err := os.ReadFile("/etc/probe.yaml")
if err != nil {
	return fmt.Errorf("lettura della configurazione: %w", err)
}
```

Arricchire gli errori a ogni livello produce messaggi che raccontano la strada percorsa, come
`avvio: lettura della configurazione: open /etc/probe.yaml: no such file or directory`.

Per **riconoscere** un errore, anche quando è stato avvolto con `%w`, si usa `errors.Is`:

```go
if errors.Is(err, os.ErrNotExist) {
	fmt.Println("il file non c'è: uso i valori predefiniti")
}
```

E per raccogliere più errori insieme (ad esempio tutti i problemi di un file di configurazione) c'è
`errors.Join(err1, err2, ...)`: la [lezione G.4](04-cli-test.md) lo usa.

!!! tip "panic non è il modo di gestire gli errori"
    `panic` ferma il programma e stampa la pila delle chiamate. È riservato agli errori di programmazione
    (un indice fuori dalla slice, una map nil), non ai problemi prevedibili come un file mancante o una rete
    irraggiungibile: quelli si restituiscono come `error`.

## Leggere file e `defer`

```go
dati, err := os.ReadFile("/proc/loadavg")   // tutto il file in memoria: va bene per file piccoli
if err != nil {
	return err
}
fmt.Println(strings.Fields(string(dati))[0]) // il carico dell'ultimo minuto
```

Per file grandi, o per lavorare riga per riga, si apre il file e lo si legge con uno **scanner**:

```go
f, err := os.Open("/etc/passwd")
if err != nil {
	return err
}
defer f.Close()   // (1)!

sc := bufio.NewScanner(f)
for sc.Scan() {
	utente, _, _ := strings.Cut(sc.Text(), ":")   // divide in due al primo ":"
	fmt.Println(utente)
}
if err := sc.Err(); err != nil {
	return err
}
```

1.  `defer` rimanda l'esecuzione di `f.Close()` al momento in cui la funzione termina, **comunque termini**: con un
    `return` normale, con un errore o con un panic. Si scrive subito dopo aver aperto la risorsa, così non te ne
    dimentichi.

Le funzioni del pacchetto `strings` che userai di più: `Fields` (divide sugli spazi), `Split`, `Cut`,
`TrimSpace`, `HasPrefix`, `Contains`, `ToLower`, `Join`.

## Esercizi

!!! esercizio "Esercizio 1 · Contare i servizi per stato"
    Data la slice `stati := []string{"active", "failed", "active", "inactive", "active", "failed"}`, costruisci una
    map con il numero di servizi per ciascuno stato e stampala **in ordine alfabetico**.

??? success "Soluzione"
    ```go
    conteggio := make(map[string]int)
    for _, s := range stati {
    	conteggio[s]++          // una chiave assente vale 0: non serve controllarla
    }
    for _, s := range slices.Sorted(maps.Keys(conteggio)) {
    	fmt.Printf("%-10s %d\n", s, conteggio[s])
    }
    ```
    **Spiegazione.** Il valore zero di `int` rende il conteggio una riga. `maps.Keys` restituisce le chiavi come
    sequenza e `slices.Sorted` le raccoglie in una slice ordinata (pacchetti `maps` e `slices`).

!!! esercizio "Esercizio 2 · sysinfo"
    Scrivi `sysinfo`, un programma che stampa:

    - il nome della distribuzione, letto dalla riga `PRETTY_NAME` di `/etc/os-release`;
    - il carico dell'ultimo minuto, il primo campo di `/proc/loadavg`;
    - la percentuale di memoria in uso, calcolata da `MemTotal` e `MemAvailable` in `/proc/meminfo`.

    Raccogli i dati in una struct `Sistema` con un metodo `MemUsata() float64`, e restituisci errori arricchiti con
    `%w`. Prova il binario nel Dev Container e poi copialo con `scp` sulla VM `rhel-lab`: funziona senza installare
    nulla?

??? success "Soluzione"
    ```go title="main.go"
    package main

    import (
    	"bufio"
    	"fmt"
    	"os"
    	"strconv"
    	"strings"
    )

    type Sistema struct {
    	Distro    string
    	Carico1m  float64
    	MemTotKB  int
    	MemDispKB int
    }

    func (s Sistema) MemUsata() float64 {
    	if s.MemTotKB == 0 {
    		return 0
    	}
    	return float64(s.MemTotKB-s.MemDispKB) / float64(s.MemTotKB) * 100
    }

    // leggiChiaveValore legge righe del tipo CHIAVE=valore (sep "=") o "Chiave:  valore" (sep ":").
    func leggiChiaveValore(percorso, sep string) (map[string]string, error) {
    	f, err := os.Open(percorso)
    	if err != nil {
    		return nil, err
    	}
    	defer f.Close()

    	valori := make(map[string]string)
    	sc := bufio.NewScanner(f)
    	for sc.Scan() {
    		chiave, valore, ok := strings.Cut(sc.Text(), sep)
    		if !ok {
    			continue
    		}
    		valori[strings.TrimSpace(chiave)] = strings.Trim(strings.TrimSpace(valore), `"`)
    	}
    	return valori, sc.Err()
    }

    // primoNumero trasforma "16318520 kB" in 16318520; 0 se manca.
    func primoNumero(s string) int {
    	campi := strings.Fields(s)
    	if len(campi) == 0 {
    		return 0
    	}
    	n, _ := strconv.Atoi(campi[0])
    	return n
    }

    func leggiSistema() (Sistema, error) {
    	var s Sistema

    	osr, err := leggiChiaveValore("/etc/os-release", "=")
    	if err != nil {
    		return s, fmt.Errorf("os-release: %w", err)
    	}
    	s.Distro = osr["PRETTY_NAME"]

    	load, err := os.ReadFile("/proc/loadavg")
    	if err != nil {
    		return s, fmt.Errorf("loadavg: %w", err)
    	}
    	if s.Carico1m, err = strconv.ParseFloat(strings.Fields(string(load))[0], 64); err != nil {
    		return s, fmt.Errorf("loadavg: %w", err)
    	}

    	mem, err := leggiChiaveValore("/proc/meminfo", ":")
    	if err != nil {
    		return s, fmt.Errorf("meminfo: %w", err)
    	}
    	s.MemTotKB = primoNumero(mem["MemTotal"])
    	s.MemDispKB = primoNumero(mem["MemAvailable"])
    	return s, nil
    }

    func main() {
    	s, err := leggiSistema()
    	if err != nil {
    		fmt.Fprintln(os.Stderr, "errore:", err)
    		os.Exit(1)
    	}
    	fmt.Printf("Sistema:  %s\n", s.Distro)
    	fmt.Printf("Carico:   %.2f (1 minuto)\n", s.Carico1m)
    	fmt.Printf("Memoria:  %.1f%% in uso di %d MB\n", s.MemUsata(), s.MemTotKB/1024)
    }
    ```
    ```console
    $ go build -o sysinfo . && ./sysinfo
    Sistema:  Ubuntu 24.04.3 LTS
    Carico:   0.60 (1 minuto)
    Memoria:  23.7% in uso di 16091 MB
    $ scp sysinfo rhel-lab: && ssh rhel-lab ./sysinfo
    Sistema:  AlmaLinux 10.0 (Purple Lion)
    ...
    ```
    **Spiegazione.** Lo stesso binario gira su Ubuntu e su AlmaLinux perché non dipende da librerie del sistema.
    `MemAvailable` è la stima del kernel della memoria utilizzabile senza ricorrere allo swap: è più onesta di
    `MemFree`, che non conta la cache liberabile (è la colonna *available* di `free -h`, che hai usato nella
    [lezione 1.13](../fase-1/13-troubleshooting.md#i-primi-60-secondi-su-un-server)). Un'unica funzione `leggiChiaveValore` gestisce due formati
    diversi grazie al parametro `sep`.

!!! esercizio "Esercizio 3 · Riconoscere l'errore"
    Modifica `sysinfo` perché, se `/etc/os-release` non esiste, stampi `Sistema: sconosciuto` invece di fermarsi.
    Usa `errors.Is` sull'errore restituito da `leggiChiaveValore`.

??? success "Soluzione"
    ```go
    osr, err := leggiChiaveValore("/etc/os-release", "=")
    switch {
    case errors.Is(err, os.ErrNotExist):
    	s.Distro = "sconosciuto"
    case err != nil:
    	return s, fmt.Errorf("os-release: %w", err)
    default:
    	s.Distro = osr["PRETTY_NAME"]
    }
    ```
    **Spiegazione.** `os.Open` restituisce un `*fs.PathError` che "avvolge" `os.ErrNotExist`: `errors.Is` lo
    riconosce anche se lo avessi avvolto di nuovo con `%w`. Un permesso negato, invece, resta un errore vero.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Perché si scrive `server = append(server, \"x\")` e non solo `append(server, \"x\")`?","options":["È solo una convenzione","append restituisce la slice aggiornata, che può puntare a un nuovo array","append non funziona senza assegnazione per motivi di sintassi","Per evitare un panic"],"answer":[1],"explanation":"Se l'array sottostante è pieno, append ne crea uno nuovo."},
 {"q":"Cosa restituisce `m[\"assente\"]` se m è una map[string]int?","options":["Un errore","nil","0","Un panic"],"answer":[2],"explanation":"Il valore zero del tipo; con `v, ok := m[k]` sai se la chiave c'era."},
 {"q":"Cosa succede scrivendo in una map dichiarata con `var m map[string]int`?","options":["Funziona","Panic: la map è nil","Errore di compilazione","Viene creata automaticamente"],"answer":[1],"explanation":"Va creata con make o con un letterale."},
 {"q":"Un campo di una struct si chiama `porta`. Un altro pacchetto può leggerlo?","options":["Sì","No: i nomi con l'iniziale minuscola non sono esportati","Solo con un puntatore","Solo con reflect"],"answer":[1],"explanation":"In Go la visibilità dipende dall'iniziale."},
 {"q":"Quale verbo di fmt.Errorf conserva l'errore originale perché errors.Is lo possa riconoscere?","options":["%v","%s","%w","%e"],"answer":[2],"explanation":"%w avvolge (wrap) l'errore."},
 {"q":"Quando viene eseguita un'istruzione `defer f.Close()`?","options":["Subito","Alla fine del programma","Quando la funzione che la contiene termina, in qualunque modo","Solo se c'è un errore"],"answer":[2],"explanation":"Per questo si scrive subito dopo aver aperto la risorsa."}
]}
</script>

## Fonti

- [A Tour of Go: More types](https://go.dev/tour/moretypes/1) e [Methods](https://go.dev/tour/methods/1)
- [Go Slices: usage and internals](https://go.dev/blog/slices-intro)
- [Error handling and Go](https://go.dev/blog/error-handling-and-go) e [Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- [Defer, Panic, and Recover](https://go.dev/blog/defer-panic-and-recover)
- Pacchetti [slices](https://pkg.go.dev/slices), [maps](https://pkg.go.dev/maps), [strings](https://pkg.go.dev/strings), [bufio](https://pkg.go.dev/bufio), [errors](https://pkg.go.dev/errors)
- [proc_meminfo(5)](https://man7.org/linux/man-pages/man5/proc_meminfo.5.html) e [os-release(5)](https://www.freedesktop.org/software/systemd/man/latest/os-release.html)

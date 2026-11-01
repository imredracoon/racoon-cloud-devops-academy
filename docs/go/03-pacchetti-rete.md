# G.3 Pacchetti, interfacce e rete

!!! obiettivi "Obiettivi"
    - Organizzare un programma in più pacchetti e capire `go.mod`, `go.sum` e la cartella `internal`.
    - Aggiungere e aggiornare una dipendenza esterna con `go get`.
    - Capire le interfacce e usare quelle della libreria standard (`io.Reader`, `io.Writer`, `fmt.Stringer`).
    - Usare il pacchetto `net` per risolvere nomi, aprire connessioni TCP e lavorare con indirizzi e reti.

!!! info "Quando studiarla"
    Alla fine della [Fase 1](../fase-1/index.md), dopo la [lezione 1.9](../fase-1/09-diagnostica.md). Scriverai in Go
    una versione semplificata di quello che fanno `nc -zv` e `ipcalc`.

## Pacchetti e moduli

Un **pacchetto** è una cartella di file `.go` che iniziano con la stessa riga `package nome`. Un **modulo** è un
insieme di pacchetti con un unico `go.mod`. Un progetto tipico ha questa forma:

```text
portcheck/
├── go.mod                 # nome del modulo e versione di Go
├── go.sum                 # impronte crittografiche delle dipendenze
├── main.go                # package main: il programma
└── internal/
    └── scan/
        ├── scan.go        # package scan
        └── scan_test.go
```

```go title="go.mod"
module github.com/tuonome/portcheck

go 1.27
```

Il **nome del modulo** di solito coincide con l'indirizzo del repository: così chiunque può scaricarlo con
`go install github.com/tuonome/portcheck@latest`. Un pacchetto si importa con il nome del modulo seguito dal
percorso della cartella:

```go
import "github.com/tuonome/portcheck/internal/scan"

func main() {
	scan.Porta("localhost:22")   // solo i nomi con la maiuscola sono accessibili da fuori
}
```

La cartella **`internal`** ha un significato speciale: i suoi pacchetti possono essere importati solo dal codice
dello stesso modulo. È il modo di dire "questo è un dettaglio, non contarci".

### Le dipendenze esterne

```bash
go get go.yaml.in/yaml/v3@latest    # aggiunge (o aggiorna) la dipendenza in go.mod e go.sum
go mod tidy                         # aggiunge ciò che manca e toglie ciò che non usi più
go list -m all                      # elenca tutte le dipendenze, anche indirette
go list -m -u all                   # mostra quali hanno aggiornamenti disponibili
```

`go.sum` contiene l'impronta di ogni versione scaricata: se qualcuno modificasse un pacchetto già pubblicato, il
download fallirebbe. Le impronte sono verificate anche contro un database pubblico di Google, `sum.golang.org`.
**Vanno nel repository** sia `go.mod` sia `go.sum`.

!!! tip "Meno dipendenze, meglio è"
    La libreria standard di Go è molto ricca: HTTP, JSON, crittografia, TLS, test, log strutturati. Prima di
    aggiungere una dipendenza chiediti se la standard library basta. Ogni dipendenza è codice di qualcun altro da
    tenere aggiornato e da controllare per le vulnerabilità (`govulncheck`, che vedrai nella
    [lezione G.9](09-ci-metriche.md)).

## Le interfacce

Un'**interfaccia** descrive un comportamento: un elenco di metodi. Qualsiasi tipo che abbia quei metodi
**soddisfa** l'interfaccia automaticamente, senza dichiararlo.

```go
type Controllo interface {
	Esegui() error
	Nome() string
}

type ControlloTCP struct{ Indirizzo string }

func (c ControlloTCP) Nome() string  { return "tcp " + c.Indirizzo }
func (c ControlloTCP) Esegui() error { /* ... */ return nil }

type ControlloDisco struct{ Soglia int }

func (c ControlloDisco) Nome() string  { return "disco" }
func (c ControlloDisco) Esegui() error { /* ... */ return nil }

func main() {
	controlli := []Controllo{ControlloTCP{"localhost:22"}, ControlloDisco{90}}
	for _, c := range controlli {
		if err := c.Esegui(); err != nil {
			fmt.Println(c.Nome(), "KO:", err)
		}
	}
}
```

Il codice che usa `Controllo` non sa e non deve sapere quale tipo concreto ha davanti. Così puoi aggiungere nuovi
controlli senza toccarlo, e nei test puoi sostituire un controllo vero con uno finto (lo farai nella
[lezione G.7](07-sdk-cloud.md)).

### Le interfacce della libreria standard

Tre interfacce piccolissime reggono buona parte della libreria standard:

```go
type Reader interface { Read(p []byte) (n int, err error) }   // io.Reader: qualcosa da cui leggere
type Writer interface { Write(p []byte) (n int, err error) }  // io.Writer: qualcosa su cui scrivere
type Stringer interface { String() string }                   // fmt.Stringer: come stamparmi
```

Un file, una connessione di rete, il corpo di una risposta HTTP, `os.Stdin` e una stringa in memoria
(`strings.NewReader`) sono tutti `io.Reader`. Una funzione che accetta un `io.Reader` funziona con **tutti**:

```go
func contaRighe(r io.Reader) (int, error) {
	n := 0
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		n++
	}
	return n, sc.Err()
}

contaRighe(os.Stdin)                          // da tastiera o da una pipe
contaRighe(strings.NewReader("a\nb\nc\n"))    // da una stringa: comodo nei test
```

Se un tuo tipo ha il metodo `String() string`, `fmt.Println` lo usa per stamparlo.

## Il pacchetto net

### Risolvere nomi

```go
ip, err := net.LookupHost("go.dev")         // come dig +short go.dev
mx, err := net.LookupMX("gmail.com")
nomi, err := net.LookupAddr("8.8.8.8")      // risoluzione inversa (PTR)
```

Su Linux Go usa normalmente un proprio risolutore che legge `/etc/resolv.conf` e `/etc/hosts`, gli stessi file della
[lezione 1.8](../fase-1/08-dns-routing.md).

### Aprire una connessione TCP

```go
conn, err := net.DialTimeout("tcp", "localhost:22", 2*time.Second)
if err != nil {
	return err          // "connection refused": porta chiusa; "i/o timeout": filtrata da un firewall
}
defer conn.Close()

conn.SetReadDeadline(time.Now().Add(2 * time.Second))
banner := make([]byte, 256)
n, _ := conn.Read(banner)
fmt.Printf("%s", banner[:n])   // SSH-2.0-OpenSSH_9.6p1 Ubuntu-...
```

I due errori più comuni corrispondono a quello che hai visto con `nc` e `ss`: **connection refused** significa che
la macchina ha risposto ma nessun processo ascolta su quella porta; **i/o timeout** significa che nessuno ha
risposto, di solito perché un firewall scarta i pacchetti ([lezione 1.11](../fase-1/11-firewall.md)).

!!! warning "Sempre un timeout"
    `net.Dial` senza timeout può aspettare a lungo, decide il sistema operativo. In qualsiasi programma di rete
    imposta un tempo massimo: `DialTimeout`, o un `context` con scadenza (lo vedrai nella
    [lezione G.5](05-concorrenza-http.md)).

### Indirizzi e reti

Il pacchetto `net/netip` rappresenta indirizzi IP e reti in modo efficiente e sicuro:

```go
rete := netip.MustParsePrefix("10.0.1.0/24")
ip := netip.MustParseAddr("10.0.1.77")
rete.Contains(ip)          // true
rete.Bits()                // 24
ip.IsPrivate()             // true: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16 (RFC 1918)

host, porta, err := net.SplitHostPort("db01.lab:5432")   // "db01.lab", "5432"
net.JoinHostPort("::1", "8080")                           // "[::1]:8080": le parentesi per IPv6 le mette lui
```

## Esercizi

!!! esercizio "Esercizio 1 · portcheck"
    Scrivi `portcheck`, che riceve uno o più indirizzi `host:porta` come argomenti e per ognuno stampa se la porta è
    aperta e in quanto tempo ha risposto. Se non riceve argomenti, legge gli indirizzi da standard input, uno per
    riga, ignorando righe vuote e commenti (`#`). Esce con codice 1 se almeno una porta è chiusa.

    Requisiti: un tipo `Esito` con un metodo `String()`, una funzione `leggiIndirizzi(r io.Reader)` e un timeout di
    2 secondi. Provalo sulle VM: `portcheck lab:22 lab:80 rhel-lab:22`.

??? success "Soluzione"
    ```go title="main.go"
    package main

    import (
    	"bufio"
    	"fmt"
    	"io"
    	"net"
    	"os"
    	"strings"
    	"time"
    )

    type Esito struct {
    	Indirizzo string
    	Aperta    bool
    	Durata    time.Duration
    	Err       error
    }

    func (e Esito) String() string {
    	if e.Aperta {
    		return fmt.Sprintf("%-25s APERTA  %s", e.Indirizzo, e.Durata.Round(time.Millisecond))
    	}
    	return fmt.Sprintf("%-25s CHIUSA  %v", e.Indirizzo, e.Err)
    }

    func controlla(indirizzo string, timeout time.Duration) Esito {
    	if _, _, err := net.SplitHostPort(indirizzo); err != nil {
    		return Esito{Indirizzo: indirizzo, Err: err}
    	}
    	inizio := time.Now()
    	conn, err := net.DialTimeout("tcp", indirizzo, timeout)
    	e := Esito{Indirizzo: indirizzo, Durata: time.Since(inizio), Err: err}
    	if err == nil {
    		e.Aperta = true
    		conn.Close()
    	}
    	return e
    }

    func leggiIndirizzi(r io.Reader) ([]string, error) {
    	var out []string
    	sc := bufio.NewScanner(r)
    	for sc.Scan() {
    		riga := strings.TrimSpace(sc.Text())
    		if riga == "" || strings.HasPrefix(riga, "#") {
    			continue
    		}
    		out = append(out, riga)
    	}
    	return out, sc.Err()
    }

    func main() {
    	indirizzi := os.Args[1:]
    	if len(indirizzi) == 0 {
    		var err error
    		if indirizzi, err = leggiIndirizzi(os.Stdin); err != nil {
    			fmt.Fprintln(os.Stderr, "errore:", err)
    			os.Exit(2)
    		}
    	}
    	chiuse := 0
    	for _, a := range indirizzi {
    		e := controlla(a, 2*time.Second)
    		fmt.Println(e)
    		if !e.Aperta {
    			chiuse++
    		}
    	}
    	if chiuse > 0 {
    		os.Exit(1)
    	}
    }
    ```
    ```console
    $ ./portcheck lab:22 lab:5432 nonsense
    lab:22                    APERTA  1ms
    lab:5432                  CHIUSA  dial tcp 192.168.56.10:5432: connect: connection refused
    nonsense                  CHIUSA  address nonsense: missing port in address
    ```
    **Spiegazione.** Controllare il formato con `SplitHostPort` prima di connettersi dà un errore chiaro invece di
    un tentativo inutile. I controlli sono **uno dopo l'altro**: con cento indirizzi filtrati aspetteresti cento
    volte due secondi. Nella [lezione G.5](05-concorrenza-http.md) li farai in parallelo.

!!! esercizio "Esercizio 2 · Una piccola calcolatrice di reti"
    Scrivi `cidr RETE IP` che stampa la rete normalizzata (ad esempio `10.0.1.5/24` diventa `10.0.1.0/24`), quanti
    indirizzi contiene, quanti sono utilizzabili per gli host e se l'IP indicato ne fa parte. Limita il programma a
    IPv4. Verifica i risultati con gli esempi della [lezione 1.7](../fase-1/07-tcpip.md).

??? success "Soluzione"
    ```go title="main.go"
    package main

    import (
    	"fmt"
    	"net/netip"
    	"os"
    )

    func main() {
    	if len(os.Args) != 3 {
    		fmt.Fprintln(os.Stderr, "uso: cidr RETE IP   (es. cidr 10.0.1.0/24 10.0.1.77)")
    		os.Exit(2)
    	}
    	rete, err := netip.ParsePrefix(os.Args[1])
    	if err != nil || !rete.Addr().Is4() {
    		fmt.Fprintln(os.Stderr, "rete IPv4 non valida:", os.Args[1])
    		os.Exit(2)
    	}
    	ip, err := netip.ParseAddr(os.Args[2])
    	if err != nil {
    		fmt.Fprintln(os.Stderr, "IP non valido:", err)
    		os.Exit(2)
    	}
    	rete = rete.Masked()
    	totale := 1 << (32 - rete.Bits())
    	fmt.Printf("rete:       %s\n", rete)
    	fmt.Printf("indirizzi:  %d (utilizzabili %d)\n", totale, max(totale-2, 0))
    	fmt.Printf("contiene %s: %t\n", ip, rete.Contains(ip))
    }
    ```
    ```console
    $ ./cidr 10.0.1.5/24 10.0.2.1
    rete:       10.0.1.0/24
    indirizzi:  256 (utilizzabili 254)
    contiene 10.0.2.1: false
    ```
    **Spiegazione.** `1 << n` vale 2 elevato a `n`: con 24 bit di rete restano 8 bit per gli host, cioè 256
    indirizzi, meno quello di rete e quello di broadcast. `max` è una funzione predefinita da Go 1.21. Nel cloud i
    conti cambiano: AWS riserva **5** indirizzi per subnet, come vedrai nella [lezione 4.7](../fase-4/07-rete.md).

!!! esercizio "Esercizio 3 · Dividere in pacchetti"
    Riorganizza `portcheck` come modulo `github.com/TUONOME/portcheck`: sposta `Esito`, `controlla` (rinominata
    `Controlla`) e `leggiIndirizzi` (rinominata `LeggiIndirizzi`) nel pacchetto `internal/scan` e lascia in
    `main.go` solo la gestione di argomenti e codici di uscita. Pubblica il modulo su GitHub.

??? success "Soluzione"
    ```bash
    go mod init github.com/TUONOME/portcheck
    mkdir -p internal/scan
    # in internal/scan/scan.go: "package scan" con Esito, Controlla e LeggiIndirizzi
    ```
    ```go title="main.go (estratto)"
    import "github.com/TUONOME/portcheck/internal/scan"

    e := scan.Controlla(a, 2*time.Second)
    ```
    ```bash
    go build ./... && go vet ./...
    git init && git add . && git commit -m "portcheck: primo rilascio"
    gh repo create portcheck --public --source . --push
    go install github.com/TUONOME/portcheck@latest   # da qualsiasi macchina con Go: finisce in ~/go/bin
    ```
    **Spiegazione.** Le funzioni devono iniziare con la maiuscola per essere usate da `main`, che è un altro
    pacchetto. `go install` scarica, compila e installa il programma in un solo comando: è così che si distribuiscono
    molti tool Go open source.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quali file vanno salvati nel repository insieme al codice?","options":["Solo go.mod","go.mod e go.sum","Solo go.sum","Nessuno dei due: si rigenerano"],"answer":[1],"explanation":"go.sum contiene le impronte che proteggono da dipendenze manomesse."},
 {"q":"Chi può importare il pacchetto `github.com/tu/app/internal/db`?","options":["Chiunque","Solo il codice del modulo github.com/tu/app","Solo il pacchetto main","Nessuno"],"answer":[1],"explanation":"La cartella internal limita l'import allo stesso modulo."},
 {"q":"Come dichiara un tipo di soddisfare un'interfaccia?","options":["Con la parola implements","Non lo dichiara: basta che abbia i metodi richiesti","Con un'annotazione","Ereditando da una classe base"],"answer":[1],"explanation":"Le interfacce in Go sono soddisfatte implicitamente."},
 {"q":"Una funzione accetta un io.Reader. Quale di questi NON puoi passarle?","options":["os.Stdin","Un *os.File aperto","strings.NewReader(\"testo\")","Un int"],"answer":[3],"explanation":"Un int non ha il metodo Read."},
 {"q":"`net.DialTimeout` restituisce `connection refused`. Cosa significa più probabilmente?","options":["Un firewall scarta i pacchetti","La macchina ha risposto ma nessun processo ascolta su quella porta","Il DNS non funziona","Il timeout è troppo breve"],"answer":[1],"explanation":"Un firewall che scarta produce invece un timeout."},
 {"q":"Quale comando toglie da go.mod le dipendenze che il codice non usa più?","options":["go get -u","go mod tidy","go clean","go vet"],"answer":[1],"explanation":"go mod tidy allinea go.mod e go.sum al codice."}
]}
</script>

## Fonti

- [Tutorial: Create a Go module](https://go.dev/doc/tutorial/create-module)
- [Organizing a Go module](https://go.dev/doc/modules/layout)
- [Go Modules Reference](https://go.dev/ref/mod) e [Managing dependencies](https://go.dev/doc/modules/managing-dependencies)
- [A Tour of Go: Interfaces](https://go.dev/tour/methods/9)
- Pacchetti [io](https://pkg.go.dev/io), [net](https://pkg.go.dev/net) e [net/netip](https://pkg.go.dev/net/netip)
- [RFC 1918: Address Allocation for Private Internets](https://www.rfc-editor.org/rfc/rfc1918)

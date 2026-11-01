# G.1 Primi passi con Go

!!! obiettivi "Obiettivi"
    - Capire cos'è un linguaggio compilato e cosa lo distingue da uno interpretato.
    - Usare gli strumenti di Go: `go run`, `go build`, `go mod init`, `gofmt` e `go vet`.
    - Dichiarare variabili e costanti, usare i tipi di base e convertirli.
    - Controllare il flusso con `if`, `for` e `switch`.
    - Scrivere funzioni, anche con più valori di ritorno.

!!! info "Quando studiarla"
    Dopo la [lezione 0.7](../fase-0/07-laboratorio.md), quando il Dev Container funziona. Non serve sapere
    programmare: questa lezione parte da zero. Se hai già fatto la [lezione 2.6](../fase-2/06-python-base.md) su
    Python, troverai molte idee familiari.

## Che cos'è un programma

Un **programma** è un elenco di istruzioni scritte in un **linguaggio di programmazione**, cioè in una forma
precisa che una persona può leggere e un computer può eseguire. Il processore però capisce solo il **codice
macchina** (la [lezione 0.1](../fase-0/01-computer.md) spiega cos'è): serve qualcuno che traduca.

| | Linguaggio interpretato (Python, Bash) | Linguaggio compilato (Go, C, Rust) |
|---|---|---|
| **Traduzione** | un interprete legge ed esegue il codice riga per riga, ogni volta | un compilatore traduce tutto il programma **una volta**, prima dell'esecuzione |
| **Cosa distribuisci** | il codice sorgente, più l'interprete giusto sulla macchina | un **file eseguibile** già pronto |
| **Errori di tipo** | li scopri quando quella riga viene eseguita | li scopri **prima**, durante la compilazione |
| **Velocità** | più lento | più veloce |

Go è stato creato in Google nel 2009 e pubblicato come open source. È stato progettato per essere **semplice**
(poche parole chiave, un solo modo di formattare il codice), **veloce da compilare** e adatto a programmi che
fanno molte cose insieme, come i server.

## Gli strumenti

Nel Dev Container Go è già installato:

```console
$ go version
go version go1.27.3 linux/amd64
```

Il comando `go` fa quasi tutto da solo: compila, esegue, scarica le dipendenze, lancia i test e formatta il codice.

??? info "Installare Go sulla tua workstation (facoltativo)"
    Nel corso usi sempre il Dev Container. Se vuoi Go anche fuori dal container, la strada consigliata dal sito
    ufficiale è l'archivio di [go.dev/dl](https://go.dev/dl/), uguale su ogni distribuzione:

    === "Ubuntu"

        ```bash
        curl -LO https://go.dev/dl/go1.27.3.linux-amd64.tar.gz   # controlla l'ultima versione su go.dev/dl
        sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.27.3.linux-amd64.tar.gz
        echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc && source ~/.bashrc
        go version
        ```

        Il pacchetto `golang-go` di Ubuntu LTS è spesso indietro di qualche versione: meglio l'archivio ufficiale.

    === "Fedora"

        ```bash
        sudo dnf install golang     # Fedora segue da vicino le versioni di Go
        go version
        ```

    Non preoccuparti troppo della versione: da Go 1.21, se un progetto chiede una versione più recente di quella
    installata, il comando `go` la **scarica da solo** (meccanismo *toolchain*).

## Il primo programma

Crea una cartella e un **modulo**. Un modulo è un progetto Go: una cartella con un file `go.mod` che ne indica il
nome e la versione di Go richiesta.

```bash
mkdir -p ~/go-lab/ciao && cd ~/go-lab/ciao
go mod init ciao        # crea go.mod
```

```go title="main.go"
package main // (1)!

import "fmt" // (2)!

func main() { // (3)!
	fmt.Println("Ciao, mondo!")
}
```

1.  Ogni file Go appartiene a un **pacchetto**. Il pacchetto `main` è speciale: diventa un programma eseguibile.
2.  `import` carica un pacchetto della **libreria standard**: `fmt` serve a stampare e formattare testo.
3.  La funzione `main` è il punto di partenza: il programma inizia da qui e finisce quando `main` termina.

```console
$ go run .               # compila in una cartella temporanea ed esegue
Ciao, mondo!
$ go build -o ciao .     # compila e lascia il file eseguibile
$ ls -lh ciao
-rwxr-xr-x 1 vscode vscode 2,2M ciao
$ ./ciao
Ciao, mondo!
$ file ciao
ciao: ELF 64-bit LSB executable, x86-64, ... statically linked ...
```

Quei 2 MB contengono **tutto**: il tuo codice, la parte di libreria standard che usi e il *runtime* di Go. Puoi
copiare `ciao` su qualsiasi Linux a 64 bit e funziona, senza installare nulla. È la caratteristica che ha reso Go
così popolare negli strumenti DevOps.

### Il formato è uno solo

In Go non si discute su spazi e parentesi: il formato lo decide `gofmt`. Le indentazioni si fanno con il tasto
Tab e la graffa aperta va **sulla stessa riga** di `func`, `if` e `for` (altrimenti il codice non compila).

```bash
gofmt -l .       # elenca i file non formattati
gofmt -w .       # li formatta
go vet ./...     # cerca errori probabili che il compilatore non segnala
```

VS Code, con l'estensione Go del Dev Container, formatta il file a ogni salvataggio.

## Variabili e tipi

Go è a **tipizzazione statica**: ogni variabile ha un tipo deciso prima dell'esecuzione, e non può cambiare.

```go
var nome string = "web01"   // dichiarazione completa
var porta = 8080            // il tipo (int) si deduce dal valore
carico := 0.75              // forma breve, solo dentro le funzioni: la più usata
var attivo bool             // senza valore: vale il "valore zero", qui false

const MaxTentativi = 3      // costante: non cambia mai
```

Ogni tipo ha un **valore zero**: `0` per i numeri, `""` per le stringhe, `false` per i booleani. In Go non esistono
variabili "non inizializzate".

| Tipo | Esempio | Note |
|---|---|---|
| `int` | `8080` | intero a 64 bit sui sistemi a 64 bit; esistono anche `int8`...`int64`, `uint` |
| `float64` | `0.75` | numero decimale |
| `string` | `"web01"` | testo UTF-8, immutabile |
| `bool` | `true` | vero o falso |
| `byte` | `'A'` | un byte (alias di `uint8`) |
| `rune` | `'è'` | un carattere Unicode (alias di `int32`) |

Le **conversioni sono sempre esplicite**: Go non mescola tipi diversi da solo.

```go
cpu := 4
ram := 7.5
totale := float64(cpu) * ram      // senza float64() è un errore di compilazione
testo := strconv.Itoa(porta)      // int → string (pacchetto strconv)
n, err := strconv.Atoi("42")      // string → int: può fallire, quindi restituisce anche un errore
```

!!! warning "Variabili e import inutilizzati"
    Se dichiari una variabile e non la usi, o importi un pacchetto che non usi, **il programma non compila**.
    All'inizio sembra fastidioso; in realtà elimina un'intera famiglia di errori.

## Stampare con `fmt`

```go
fmt.Println("porta", porta)                    // valori separati da spazi, poi a capo
fmt.Printf("%s ascolta su %d\n", nome, porta)  // con un modello
fmt.Printf("carico %.1f%%\n", carico*100)      // 75.0%
s := fmt.Sprintf("%s:%d", nome, porta)         // come Printf, ma restituisce una stringa
fmt.Printf("%v | %T\n", porta, porta)          // %v: il valore in forma predefinita, %T: il tipo
```

## Il controllo del flusso

### if

```go
uso := 87
if uso >= 90 {
	fmt.Println("critico")
} else if uso >= 80 {
	fmt.Println("attenzione")
} else {
	fmt.Println("ok")
}

// Un if può iniziare con una breve istruzione: la variabile esiste solo dentro l'if.
if n, err := strconv.Atoi("12"); err == nil {
	fmt.Println("numero", n)
}
```

Niente parentesi tonde attorno alla condizione; le graffe sono **obbligatorie**.

### for: l'unico ciclo

Go ha un solo ciclo, `for`, che copre tutti i casi:

```go
for i := 0; i < 3; i++ {      // classico: inizio; condizione; passo
	fmt.Println(i)
}

for i := range 3 {            // da Go 1.22: 0, 1, 2
	fmt.Println(i)
}

tentativi := 0
for tentativi < 3 {           // come un "while"
	tentativi++
}

for {                         // ciclo infinito: si esce con break o return
	break
}
```

### switch

```go
switch giorno := time.Now().Weekday(); giorno {
case time.Saturday, time.Sunday:
	fmt.Println("weekend: niente rilasci")
default:
	fmt.Println("giorno lavorativo")
}

switch {                      // senza valore: è un elenco di if più leggibile
case uso >= 90:
	fmt.Println("critico")
case uso >= 80:
	fmt.Println("attenzione")
}
```

Al contrario di C e Java, un `case` **non prosegue** nel successivo: non serve `break`.

## Funzioni

```go
func percentuale(usato, totale float64) float64 {
	if totale == 0 {
		return 0
	}
	return usato / totale * 100
}
```

I tipi si scrivono **dopo** il nome. Una funzione può restituire **più valori**, ed è così che Go gestisce gli
errori: l'ultimo valore restituito è di tipo `error`, e vale `nil` (cioè "niente") se tutto è andato bene.

```go
func dividi(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("divisione per zero")
	}
	return a / b, nil
}

func main() {
	q, err := dividi(10, 0)
	if err != nil {               // lo schema che vedrai in ogni programma Go
		fmt.Println("errore:", err)
		return
	}
	fmt.Println(q)
}
```

In Go non ci sono eccezioni come in Python: un errore è un **valore** che la funzione restituisce e che tu
controlli subito. È più verboso, ma guardando il codice vedi sempre cosa succede quando qualcosa va storto.
La [lezione G.2](02-collezioni-errori.md) approfondisce.

## Leggere gli argomenti e l'ambiente

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {                         // os.Args[0] è il nome del programma
		fmt.Fprintln(os.Stderr, "uso: saluta NOME")
		os.Exit(2)                                // codice di uscita, come in Bash
	}
	utente := os.Getenv("USER")                   // variabile d'ambiente
	fmt.Printf("Ciao %s, sono %s\n", os.Args[1], utente)
}
```

I messaggi d'errore vanno su **standard error** e il codice di uscita diverso da zero segnala il fallimento:
sono le stesse regole che hai visto per i comandi Linux nella [lezione 0.4](../fase-0/04-terminale.md).

## Esercizi

!!! esercizio "Esercizio 1 · Classificare l'uso del disco"
    Scrivi un programma `disco` che riceve come argomento una percentuale d'uso (ad esempio `./disco 87`) e stampa
    `ok`, `attenzione` (da 80) o `critico` (da 90). Se l'argomento manca o non è un numero tra 0 e 100, stampa un
    messaggio su standard error ed esce con codice 2. Con `critico` esce con codice 1, così puoi usarlo in uno
    script.

??? success "Soluzione"
    ```go title="main.go"
    package main

    import (
    	"fmt"
    	"os"
    	"strconv"
    )

    func livello(uso int) string {
    	switch {
    	case uso >= 90:
    		return "critico"
    	case uso >= 80:
    		return "attenzione"
    	default:
    		return "ok"
    	}
    }

    func main() {
    	if len(os.Args) != 2 {
    		fmt.Fprintln(os.Stderr, "uso: disco PERCENTUALE")
    		os.Exit(2)
    	}
    	uso, err := strconv.Atoi(os.Args[1])
    	if err != nil || uso < 0 || uso > 100 {
    		fmt.Fprintf(os.Stderr, "%q non è una percentuale valida\n", os.Args[1])
    		os.Exit(2)
    	}
    	l := livello(uso)
    	fmt.Println(l)
    	if l == "critico" {
    		os.Exit(1)
    	}
    }
    ```
    ```console
    $ go build -o disco . && ./disco 95; echo "codice $?"
    critico
    codice 1
    ```
    **Spiegazione.** Separare la logica (`livello`) dall'input/output (`main`) rende la funzione facile da
    provare: nella [lezione G.4](04-cli-test.md) le scriverai un test. `%q` stampa la stringa tra virgolette, utile
    per vedere spazi nascosti.

!!! esercizio "Esercizio 2 · FizzBuzz dei server"
    Stampa i numeri da 1 a 20 con un `for`. Al posto dei multipli di 3 scrivi `web`, al posto dei multipli di 5
    `db` e al posto dei multipli di entrambi `web+db`.

??? success "Soluzione"
    ```go
    for i := 1; i <= 20; i++ {
    	switch {
    	case i%15 == 0:
    		fmt.Println("web+db")
    	case i%3 == 0:
    		fmt.Println("web")
    	case i%5 == 0:
    		fmt.Println("db")
    	default:
    		fmt.Println(i)
    	}
    }
    ```
    **Spiegazione.** L'ordine dei `case` conta: il controllo su 15 deve venire prima, altrimenti il 15 verrebbe
    catturato dal `case` del 3.

!!! esercizio "Esercizio 3 · Quanto pesa un binario"
    Compila il programma dell'esercizio 1 in tre modi e confronta le dimensioni con `ls -lh`: normale,
    con `-ldflags "-s -w"` e per un Raspberry Pi (`GOOS=linux GOARCH=arm64`). Prova a eseguire il binario arm64
    sul Dev Container: cosa succede?

??? success "Soluzione"
    ```bash
    go build -o disco .
    go build -ldflags "-s -w" -o disco-piccolo .
    GOOS=linux GOARCH=arm64 go build -o disco-arm64 .
    ls -lh disco*
    file disco-arm64          # ELF 64-bit LSB executable, ARM aarch64
    ./disco-arm64 50          # cannot execute binary file: Exec format error
    ```
    **Spiegazione.** `-s -w` toglie le informazioni di debug e riduce il file di circa un terzo. `GOOS` e `GOARCH`
    scelgono sistema operativo e processore di destinazione: Go **compila per un'altra architettura senza
    installare nulla** (*cross-compilazione*). Il binario arm64 non parte su un processore x86-64, ma funziona su
    un Raspberry Pi o su un'istanza cloud Graviton/Axion.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa produce `go build`?","options":["Un file .py da eseguire con un interprete","Un file eseguibile che contiene il programma e il runtime di Go","Un'immagine Docker","Un pacchetto .deb"],"answer":[1],"explanation":"Il binario è autonomo: non serve installare Go sulla macchina di destinazione."},
 {"q":"Quale riga dichiara una variabile con la forma breve?","options":["var x := 5","x := 5","let x = 5","x = 5 (senza dichiarazione)"],"answer":[1],"explanation":":= dichiara e assegna, solo dentro le funzioni."},
 {"q":"Quanto vale una variabile `var n int` appena dichiarata?","options":["nil","0","Un valore casuale","Il programma non compila"],"answer":[1],"explanation":"Ogni tipo ha il suo valore zero."},
 {"q":"Come segnala un errore una funzione Go idiomatica?","options":["Lancia un'eccezione","Restituisce un valore di tipo error come ultimo risultato","Stampa su standard output","Chiama os.Exit(1)"],"answer":[1],"explanation":"Chi chiama controlla subito `if err != nil`."},
 {"q":"Hai importato \"strings\" ma non lo usi. Cosa succede?","options":["Niente","Un avviso a runtime","Il programma non compila","Il pacchetto viene ignorato"],"answer":[2],"explanation":"Import e variabili inutilizzati sono errori di compilazione."},
 {"q":"Come compili per un server ARM a 64 bit dal Dev Container x86-64?","options":["Non si può","GOOS=linux GOARCH=arm64 go build","go build --arm","Serve una VM ARM"],"answer":[1],"explanation":"La cross-compilazione è integrata nel comando go."}
]}
</script>

## Fonti

- [A Tour of Go](https://go.dev/tour/)
- [Tutorial: Get started with Go](https://go.dev/doc/tutorial/getting-started)
- [Download and install](https://go.dev/doc/install) e [Go Toolchains](https://go.dev/doc/toolchain)
- [Effective Go](https://go.dev/doc/effective_go)
- [The Go Programming Language Specification](https://go.dev/ref/spec)
- [Command gofmt](https://pkg.go.dev/cmd/gofmt) e [Command vet](https://pkg.go.dev/cmd/vet)

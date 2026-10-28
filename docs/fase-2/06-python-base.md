# 2.6 Python: le basi

!!! obiettivi "Obiettivi"
    - Eseguire codice Python in modo interattivo e da file.
    - Usare i tipi di dati fondamentali: numeri, stringhe, booleani, liste, dizionari, tuple e insiemi.
    - Controllare il flusso con `if`, `for` e `while`.
    - Scrivere funzioni, gestire le eccezioni e organizzare il codice in moduli.
    - Creare ambienti virtuali e installare pacchetti con `pip`.

## Perché Python

Bash è perfetto per mettere in fila comandi; quando la logica cresce (dati strutturati, chiamate a servizi web,
calcoli, errori da gestire) diventa difficile da leggere. **Python** è il linguaggio più usato nell'automazione
di sistemi e cloud: Ansible è scritto in Python, gli SDK di AWS (`boto3`) e Google Cloud hanno librerie Python, e
moltissimi strumenti interni delle aziende sono script Python.

Python è **interpretato** (non serve compilare), ha una sintassi leggibile e una libreria standard ricchissima.
Nel corso usiamo **Python 3.12**, già presente nel Dev Container.

## Eseguire Python

```bash
python3 --version
python3                  # interprete interattivo (REPL): scrivi e vedi subito il risultato; esci con exit()
python3 script.py        # esegue un file
python3 -c 'print(2**10)'
```

Il REPL è ottimo per provare. Per qualsiasi cosa che vuoi conservare, scrivi un file `.py`.

```python title="ciao.py"
nome = input("Come ti chiami? ")
print(f"Ciao, {nome}!")
```

## L'indentazione è sintassi

In Python i blocchi di codice non si delimitano con parentesi graffe o con `fi`/`done`, ma con l'**indentazione**:
quattro spazi per livello, per convenzione.

```python
if 5 > 3:
    print("vero")         # dentro l'if
    print("ancora dentro")
print("fuori")            # fuori dall'if
```

Un'indentazione sbagliata produce un `IndentationError` o, peggio, cambia il significato del programma.
Configura l'editor perché inserisca 4 spazi quando premi Tab (VS Code lo fa già per i file `.py`).

## Variabili e tipi

```python
nome = "web01"          # str, stringa
porta = 8080            # int, intero
carico = 0.75           # float, decimale
attivo = True           # bool: True o False (con la maiuscola)
nulla = None            # assenza di valore

type(porta)             # <class 'int'>
int("42") + 1           # 43: conversione da stringa a intero
str(porta)              # "8080"
```

Python è a **tipizzazione dinamica**: la variabile non ha un tipo fisso, il valore sì. È a **tipizzazione forte**:
`"porta " + 8080` è un errore (`TypeError`), devi convertire esplicitamente.

### Operatori

```python
7 + 3, 7 - 3, 7 * 3     # 10, 4, 21
7 / 2                   # 3.5  (divisione sempre decimale)
7 // 2                  # 3    (divisione intera)
7 % 2                   # 1    (resto)
2 ** 10                 # 1024 (potenza)

porta == 8080, porta != 80, porta >= 1024    # confronti: True / False
attivo and porta > 1024, not attivo          # logica
"web" in nome                                # True: appartenenza
```

### Stringhe

```python
host = "web01.example.com"
len(host)                     # 17
host.upper()                  # "WEB01.EXAMPLE.COM"
host.split(".")               # ['web01', 'example', 'com']
host.startswith("web")        # True
host.replace("web", "db")     # "db01.example.com"
"  testo  ".strip()           # "testo"
host[0:5]                     # "web01" (slicing: dall'indice 0 al 5 escluso)
host[-3:]                     # "com"

f"Server {host} sulla porta {porta}"     # f-string: il modo moderno di comporre stringhe
f"{carico:.0%}"                          # "75%": formattazione
```

Le stringhe sono **immutabili**: i metodi restituiscono una nuova stringa, non modificano l'originale.

## Le collezioni

### Liste

Una **lista** è una sequenza ordinata e modificabile:

```python
server = ["web01", "web02", "db01"]
server[0]                # "web01" (gli indici partono da 0)
server[-1]               # "db01" (l'ultimo)
server.append("cache01") # aggiunge in fondo
server.remove("web02")   # rimuove per valore
len(server)              # 3
"db01" in server         # True
sorted(server)           # nuova lista ordinata
```

### Dizionari

Un **dizionario** associa **chiavi** a **valori**. È la struttura più importante per l'automazione, perché JSON e
YAML (prossima lezione) si traducono in dizionari.

```python
vm = {"nome": "web01", "cpu": 2, "ram_gb": 4, "tag": ["prod", "web"]}
vm["nome"]                 # "web01"
vm["cpu"] = 4              # modifica
vm["zona"] = "eu-south-1"  # aggiunge
vm.get("disco", 20)        # 20: valore predefinito se la chiave manca (vm["disco"] darebbe KeyError)
"ram_gb" in vm             # True: controlla le chiavi
for chiave, valore in vm.items():
    print(chiave, "→", valore)
```

### Tuple e insiemi

```python
coordinate = ("eu-south-1", "a")     # tupla: come una lista, ma immutabile
regione, zona = coordinate           # "spacchettamento"

porte = {22, 80, 443, 80}            # set (insieme): niente duplicati, niente ordine
porte                                # {80, 443, 22}
{22, 80} & {80, 443}                 # {80}: intersezione
```

| Tipo | Ordinato | Modificabile | Duplicati | Esempio d'uso |
|---|---|---|---|---|
| `list` | sì | sì | sì | un elenco di server |
| `tuple` | sì | no | sì | una coppia di valori fissa |
| `dict` | sì (ordine di inserimento) | sì | chiavi uniche | la configurazione di una VM |
| `set` | no | sì | no | porte aperte, confronti tra elenchi |

## Controllo del flusso

```python
uso_disco = 87

if uso_disco >= 90:
    livello = "critico"
elif uso_disco >= 80:
    livello = "attenzione"
else:
    livello = "ok"
```

```python
for s in server:
    print(s)

for i, s in enumerate(server, start=1):   # con il numero d'ordine
    print(i, s)

for n in range(3):                         # 0, 1, 2
    print(n)

tentativi = 0
while tentativi < 3:
    tentativi += 1
```

### Comprensioni

Un modo compatto e molto usato per costruire liste e dizionari:

```python
nomi = ["web01", "web02", "db01"]
web = [n for n in nomi if n.startswith("web")]          # ['web01', 'web02']
lunghezze = {n: len(n) for n in nomi}                   # {'web01': 5, ...}
```

## Funzioni

```python
def percentuale(usato, totale, decimali=1):
    """Restituisce la percentuale di usato su totale."""
    if totale == 0:
        return 0.0
    return round(usato / totale * 100, decimali)

percentuale(30, 120)                 # 25.0
percentuale(1, 3, decimali=3)        # 33.333
```

- `def` definisce una funzione; `return` restituisce un valore (senza `return` restituisce `None`).
- I parametri possono avere un **valore predefinito** e si possono passare **per nome**.
- La stringa tra triple virgolette subito dopo `def` è la **docstring**: documenta la funzione.

### Annotazioni di tipo

Puoi indicare i tipi attesi. Python non li controlla durante l'esecuzione, ma rendono il codice più chiaro e
gli editor e strumenti come `mypy` li usano per trovare errori:

```python
def percentuale(usato: float, totale: float, decimali: int = 1) -> float:
    ...
```

## Gestire gli errori: le eccezioni

Quando qualcosa va storto Python solleva un'**eccezione**. Se non la gestisci, il programma si ferma e stampa un
**traceback**: leggilo dal basso, l'ultima riga dice il tipo di errore e il messaggio.

```python
try:
    with open("/etc/inesistente") as f:
        dati = f.read()
except FileNotFoundError:
    print("File non trovato")
except PermissionError as e:
    print(f"Permesso negato: {e}")
else:
    print("Letto correttamente")       # solo se non ci sono stati errori
finally:
    print("Eseguito in ogni caso")
```

Cattura sempre eccezioni **specifiche**. Un `except:` generico nasconde anche gli errori di programmazione e
rende impossibile capire cosa succede. Per segnalare tu un errore: `raise ValueError("porta non valida")`.

## Moduli e libreria standard

Un **modulo** è un file `.py` che si può importare. La libreria standard ne contiene centinaia:

```python
import os
import sys
from pathlib import Path
from datetime import datetime, timedelta

os.environ.get("HOME")          # variabili d'ambiente
sys.argv                        # argomenti della riga di comando
Path("/etc/hostname").read_text().strip()
datetime.now() - timedelta(days=7)
```

Ogni tuo file è a sua volta un modulo. Lo schema che vedrai ovunque:

```python title="strumenti.py"
def saluta(nome: str) -> str:
    return f"Ciao, {nome}"

if __name__ == "__main__":
    # eseguito solo se lanci il file direttamente, NON quando lo importi
    print(saluta("mondo"))
```

## Ambienti virtuali e pacchetti

I pacchetti esterni si installano con **pip** dal [Python Package Index (PyPI)](https://pypi.org/). Installarli
nel Python di sistema è una cattiva idea: progetti diversi richiedono versioni diverse, e su Ubuntu e Debian recenti
`pip install` fuori da un ambiente virtuale è **bloccato** (errore `externally-managed-environment`, PEP 668) per
non rompere i pacchetti del sistema.

La soluzione è un **ambiente virtuale**: una cartella con una copia isolata di Python e dei suoi pacchetti, una per
progetto.

```bash
cd ~/mio-progetto
python3 -m venv .venv            # crea l'ambiente nella cartella .venv
source .venv/bin/activate        # lo attiva: il prompt mostra (.venv)
pip install requests             # installa nell'ambiente
pip freeze > requirements.txt    # salva l'elenco esatto delle versioni
deactivate                       # disattiva

# su un'altra macchina, o per un collega:
python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
```

=== "Debian/Ubuntu"

    Il modulo `venv` è in un pacchetto separato: `sudo apt install -y python3-venv python3-pip`.

=== "RHEL (Rocky/AlmaLinux)"

    `venv` è incluso in `python3`; pip si installa con `sudo dnf install -y python3-pip`.

Ricorda di aggiungere `.venv/` al `.gitignore`: l'ambiente non si versiona, si ricrea da `requirements.txt`.

!!! tip "Lo stile: PEP 8"
    La [PEP 8](https://peps.python.org/pep-0008/) è la guida di stile ufficiale: nomi in `minuscolo_con_trattini`
    per variabili e funzioni, `MAIUSCOLO` per le costanti, `CamelCase` per le classi, righe non troppo lunghe.
    Strumenti come **Ruff** controllano e formattano il codice in automatico (`pip install ruff`, poi
    `ruff check .` e `ruff format .`).

## Esercizi

!!! esercizio "Esercizio 1 · Inventario"
    Data questa lista di dizionari:

    ```python
    vm = [
        {"nome": "web01", "ambiente": "prod", "cpu": 2, "ram_gb": 4},
        {"nome": "web02", "ambiente": "prod", "cpu": 2, "ram_gb": 4},
        {"nome": "db01",  "ambiente": "prod", "cpu": 4, "ram_gb": 16},
        {"nome": "test01", "ambiente": "test", "cpu": 1, "ram_gb": 2},
    ]
    ```

    scrivi un programma che stampi: il numero di VM per ambiente, il totale di CPU e RAM di produzione, e i nomi
    delle VM con almeno 4 GB di RAM, in ordine alfabetico.

??? success "Soluzione"
    ```python
    per_ambiente = {}
    for v in vm:
        per_ambiente[v["ambiente"]] = per_ambiente.get(v["ambiente"], 0) + 1
    print("VM per ambiente:", per_ambiente)

    prod = [v for v in vm if v["ambiente"] == "prod"]
    print("CPU prod:", sum(v["cpu"] for v in prod))
    print("RAM prod:", sum(v["ram_gb"] for v in prod), "GB")

    grandi = sorted(v["nome"] for v in vm if v["ram_gb"] >= 4)
    print("Almeno 4 GB:", ", ".join(grandi))
    ```
    **Spiegazione.** `dict.get(chiave, 0)` evita il `KeyError` la prima volta che incontri un ambiente. La libreria
    standard offre anche `collections.Counter`, che fa lo stesso conteggio in una riga:
    `Counter(v["ambiente"] for v in vm)`. `sum()` e `sorted()` accettano direttamente un'espressione generatrice.

!!! esercizio "Esercizio 2 · Validare una porta"
    Scrivi la funzione `valida_porta(valore)` che riceve una stringa e restituisce l'intero corrispondente se è una
    porta valida (1-65535); altrimenti solleva `ValueError` con un messaggio chiaro. Poi scrivi un ciclo che chiede
    una porta all'utente finché non ne inserisce una valida.

??? success "Soluzione"
    ```python
    def valida_porta(valore: str) -> int:
        try:
            porta = int(valore)
        except ValueError:
            raise ValueError(f"'{valore}' non è un numero") from None
        if not 1 <= porta <= 65535:
            raise ValueError(f"{porta} è fuori dall'intervallo 1-65535")
        return porta


    if __name__ == "__main__":
        while True:
            try:
                porta = valida_porta(input("Porta: "))
                break
            except ValueError as e:
                print(f"Errore: {e}")
        print(f"Userò la porta {porta}")
    ```
    **Spiegazione.** Python permette i confronti concatenati: `1 <= porta <= 65535`. `from None` nasconde
    l'eccezione originale di `int()` e lascia solo il nostro messaggio, più chiaro per chi legge.

!!! esercizio "Esercizio 3 · Ambiente virtuale"
    Nel Dev Container crea la cartella `~/py-prova`, un ambiente virtuale, installa il pacchetto `requests` e salva
    le dipendenze. Verifica con `which python` che stai usando il Python dell'ambiente. Poi scrivi uno script che
    stampi la versione di `requests` installata.

??? success "Soluzione"
    ```bash
    mkdir ~/py-prova && cd ~/py-prova
    python3 -m venv .venv
    source .venv/bin/activate
    which python                 # ~/py-prova/.venv/bin/python
    pip install requests
    pip freeze > requirements.txt
    cat requirements.txt         # requests e le sue dipendenze, con le versioni esatte
    python -c 'import requests; print(requests.__version__)'
    echo ".venv/" > .gitignore
    ```
    **Spiegazione.** `pip freeze` registra anche le dipendenze indirette (come `urllib3` e `certifi`), con le
    versioni esatte: è ciò che serve per riprodurre lo stesso ambiente altrove.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa delimita un blocco di codice in Python?","options":["Le parentesi graffe","L'indentazione","Le parole begin/end","Il punto e virgola"],"answer":[1],"explanation":"Per convenzione 4 spazi per livello."},
 {"q":"Quanto vale `7 // 2`?","options":["3.5","3","4","1"],"answer":[1],"explanation":"// è la divisione intera; / restituisce sempre un decimale."},
 {"q":"Come leggi una chiave che potrebbe non esistere senza provocare un errore?","options":["d[\"chiave\"]","d.get(\"chiave\", predefinito)","d.chiave","d.read(\"chiave\")"],"answer":[1],"explanation":"d[\"chiave\"] solleva KeyError se la chiave manca."},
 {"q":"Quale tipo usi per un elenco senza duplicati su cui fare intersezioni?","options":["list","tuple","set","str"],"answer":[2],"explanation":"Il set ignora i duplicati e supporta &, |, -."},
 {"q":"Perché si usa un ambiente virtuale?","options":["Per eseguire Python più velocemente","Per isolare i pacchetti di ogni progetto dal sistema e dagli altri progetti","Per usare Python 2","Per compilare il codice"],"answer":[1],"explanation":"Ogni progetto ha le sue versioni, ricreabili da requirements.txt."},
 {"q":"Quando viene eseguito il codice dentro `if __name__ == \"__main__\":`?","options":["Sempre","Solo quando il file viene eseguito direttamente, non quando viene importato","Solo quando viene importato","Mai"],"answer":[1],"explanation":"Permette di usare lo stesso file sia come script sia come modulo."}
]}
</script>

## Fonti

- [The Python Tutorial (documentazione ufficiale)](https://docs.python.org/3.12/tutorial/index.html)
- [Python: Built-in Types](https://docs.python.org/3.12/library/stdtypes.html)
- [Python: Errors and Exceptions](https://docs.python.org/3.12/tutorial/errors.html)
- [Python: venv](https://docs.python.org/3.12/library/venv.html)
- [Python Packaging User Guide: Installing packages using pip and virtual environments](https://packaging.python.org/en/latest/guides/installing-using-pip-and-virtual-environments/)
- [PEP 8: Style Guide for Python Code](https://peps.python.org/pep-0008/)
- [PEP 668: externally managed environments](https://peps.python.org/pep-0668/)

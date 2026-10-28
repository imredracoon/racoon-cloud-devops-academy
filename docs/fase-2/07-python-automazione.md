# 2.7 Python per l'automazione

!!! obiettivi "Obiettivi"
    - Lavorare con file e percorsi usando `pathlib`.
    - Leggere e scrivere JSON, YAML e CSV.
    - Eseguire comandi di sistema con `subprocess` in modo sicuro.
    - Chiamare API HTTP con `requests`.
    - Scrivere strumenti da riga di comando con `argparse` e `logging`.
    - Scrivere i primi test con `pytest`.

## File e percorsi con pathlib

`pathlib` rappresenta i percorsi come oggetti, in modo più leggibile e portabile delle stringhe:

```python
from pathlib import Path

home = Path.home()
config = home / ".config" / "app" / "config.yaml"     # l'operatore / unisce i percorsi
config.parent.mkdir(parents=True, exist_ok=True)     # come mkdir -p
config.exists(), config.is_file(), config.suffix     # False, False, '.yaml'

testo = Path("/etc/hostname").read_text().strip()
Path("note.txt").write_text("prima riga\n")

for f in Path("/var/log").glob("*.log"):             # rglob() cerca anche nelle sottocartelle
    print(f.name, f.stat().st_size)
```

Per i file grandi, o per scrivere poco alla volta, si usa `open` con il costrutto `with`, che chiude il file
automaticamente anche in caso di errore:

```python
with open("/var/log/syslog", encoding="utf-8", errors="replace") as f:
    for riga in f:                      # legge una riga alla volta, senza caricare tutto in memoria
        if "error" in riga.lower():
            print(riga.rstrip())
```

Per copiare, spostare e archiviare: il modulo `shutil` (`shutil.copy2`, `shutil.move`, `shutil.make_archive`).

## I formati dei dati

### JSON

**JSON** è il formato delle API web e di molti strumenti (`aws`, `gcloud`, `kubectl` possono tutti produrre JSON).
Si traduce direttamente in dizionari e liste.

```python
import json

dati = json.loads('{"nome": "web01", "porte": [80, 443]}')   # da stringa
dati["porte"][0]                                               # 80

with open("inventario.json") as f:
    inventario = json.load(f)                                  # da file

with open("output.json", "w") as f:
    json.dump(inventario, f, indent=2, ensure_ascii=False)     # su file, leggibile

print(json.dumps(dati, indent=2))                              # in stringa
```

### YAML

**YAML** è il formato delle configurazioni: Ansible, Kubernetes, GitHub Actions, Docker Compose. È più leggibile
di JSON ma ha regole precise: l'**indentazione con spazi** (mai tab) definisce la struttura.

```yaml title="server.yaml"
ambiente: produzione
server:
  - nome: web01
    porte: [80, 443]
  - nome: db01
    porte:
      - 5432
monitoraggio: true
```

Python non legge YAML nativamente: serve la libreria **PyYAML** (`pip install pyyaml`).

```python
import yaml

with open("server.yaml") as f:
    conf = yaml.safe_load(f)          # SEMPRE safe_load: yaml.load può eseguire codice arbitrario
for s in conf["server"]:
    print(s["nome"], s["porte"])

print(yaml.safe_dump(conf, sort_keys=False, allow_unicode=True))
```

!!! warning "Le trappole di YAML"
    - Un tab al posto degli spazi è un errore.
    - Alcuni valori vengono interpretati in modo inatteso: in YAML 1.1 (usato da PyYAML) `yes`, `no`, `on`, `off`
      diventano booleani, e `08` può diventare una stringa o un numero a seconda del parser. Se intendi una
      stringa, mettila tra virgolette: `paese: "no"`, `versione: "3.10"` (senza virgolette `3.10` diventa il
      numero `3.1`).

### CSV

```python
import csv

with open("utenti.csv", newline="") as f:
    for riga in csv.DictReader(f):          # ogni riga è un dizionario: colonna → valore
        print(riga["nome"], riga["email"])

with open("report.csv", "w", newline="") as f:
    w = csv.DictWriter(f, fieldnames=["host", "stato"])
    w.writeheader()
    w.writerow({"host": "web01", "stato": "ok"})
```

## Eseguire comandi con subprocess

```python
import subprocess

r = subprocess.run(
    ["df", "--output=pcent", "/"],    # il comando come LISTA di argomenti
    capture_output=True,              # cattura stdout e stderr
    text=True,                        # come stringhe, non byte
    check=True,                       # solleva CalledProcessError se il codice di uscita non è 0
    timeout=10,                       # secondi
)
print(r.stdout)
print(r.returncode)
```

!!! danger "Mai shell=True con dati esterni"
    Con `shell=True` il comando passa attraverso la shell, e se contiene dati forniti da un utente si apre la
    porta alla **command injection**: un nome file come `x; rm -rf ~` verrebbe eseguito. Passa sempre una
    **lista** di argomenti: ogni elemento arriva al programma così com'è, senza interpretazioni.

Prima di lanciare un comando esterno, chiediti se Python lo fa già da solo: per copiare file c'è `shutil`, per
la dimensione dei dischi `shutil.disk_usage("/")`, per i processi c'è la libreria `psutil`. Il codice è più
portabile e gli errori sono più facili da gestire.

## Chiamare API HTTP con requests

La libreria **requests** (`pip install requests`) è lo standard per HTTP in Python.

```python
import requests

r = requests.get("https://api.github.com/repos/python/cpython", timeout=10)
r.raise_for_status()                 # solleva un'eccezione se lo stato è 4xx o 5xx
dati = r.json()
print(dati["stargazers_count"])

r = requests.post(
    "http://localhost:8000/tasks",
    json={"title": "Studiare Python"},          # serializza in JSON e imposta Content-Type
    headers={"Authorization": f"Bearer {token}"},
    timeout=5,
)
print(r.status_code, r.json())
```

- Imposta **sempre** un `timeout`: senza, una richiesta può restare bloccata per sempre.
- Controlla **sempre** lo stato della risposta (`raise_for_status()` o `r.ok`).
- I token non si scrivono nel codice: si leggono da una variabile d'ambiente
  (`os.environ["GITHUB_TOKEN"]`).

## Strumenti da riga di comando

### argparse

```python
import argparse

parser = argparse.ArgumentParser(description="Controlla lo spazio su disco.")
parser.add_argument("percorsi", nargs="+", help="uno o più punti di montaggio")
parser.add_argument("-s", "--soglia", type=int, default=80, help="percentuale di allarme (predefinita: 80)")
parser.add_argument("-v", "--verbose", action="store_true", help="output dettagliato")
args = parser.parse_args()

print(args.percorsi, args.soglia, args.verbose)
```

`argparse` genera automaticamente l'aiuto (`-h`), controlla i tipi e segnala gli errori. È il `getopts` di Python,
molto più potente.

### logging

Negli script che girano senza nessuno davanti (cron, timer, container), `print` non basta: servono livelli,
orari e la possibilità di scegliere dove scrivere. Il modulo `logging`:

```python
import logging

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s: %(message)s",
)
log = logging.getLogger("controllo-disco")

log.debug("non mostrato con livello INFO")
log.info("controllo di %s", "/")
log.warning("spazio quasi esaurito")
log.error("impossibile leggere %s", "/mnt/dati")
```

Per impostazione predefinita i log vanno su **standard error**: in un servizio systemd o in un container finiscono
direttamente nel journal o nei log del container, che è esattamente ciò che vuoi.

## Uno strumento completo

```python title="controllo_disco.py"
#!/usr/bin/env python3
"""Controlla l'uso dei filesystem e segnala quelli oltre soglia."""
import argparse
import json
import logging
import shutil
import sys

log = logging.getLogger("controllo-disco")


def uso_percentuale(percorso: str) -> float:
    totale, usato, _libero = shutil.disk_usage(percorso)
    return round(usato / totale * 100, 1)


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("percorsi", nargs="+")
    parser.add_argument("-s", "--soglia", type=float, default=80)
    parser.add_argument("--json", action="store_true", help="output in formato JSON")
    args = parser.parse_args(argv)

    logging.basicConfig(level=logging.INFO, format="%(levelname)s %(message)s")

    risultati = []
    for p in args.percorsi:
        try:
            uso = uso_percentuale(p)
        except FileNotFoundError:
            log.error("%s non esiste", p)
            return 2
        risultati.append({"percorso": p, "uso": uso, "allarme": uso >= args.soglia})

    if args.json:
        print(json.dumps(risultati, indent=2))
    else:
        for r in risultati:
            stato = "ALLARME" if r["allarme"] else "ok"
            print(f"{r['percorso']:<20} {r['uso']:>5}%  {stato}")

    return 1 if any(r["allarme"] for r in risultati) else 0


if __name__ == "__main__":
    sys.exit(main())
```

Lo schema `main()` che **restituisce** il codice di uscita, chiamato con `sys.exit(main())`, ha un vantaggio
importante: puoi chiamare `main([...])` dai test, passando gli argomenti che vuoi.

## I primi test con pytest

Un **test automatico** è codice che verifica che altro codice funzioni. Si esegue in un secondo, sempre allo stesso
modo, e ti avvisa quando una modifica rompe qualcosa. Nella [Fase 7](../fase-7/index.md) i test gireranno da soli a ogni push.

**pytest** (`pip install pytest`) trova ed esegue le funzioni che iniziano con `test_` nei file `test_*.py`:

```python title="test_controllo_disco.py"
import json

import controllo_disco


def test_uso_percentuale_nella_gamma():
    uso = controllo_disco.uso_percentuale("/")
    assert 0 <= uso <= 100


def test_soglia_zero_fa_scattare_allarme(capsys):
    codice = controllo_disco.main(["/", "--soglia", "0", "--json"])
    assert codice == 1
    dati = json.loads(capsys.readouterr().out)
    assert dati[0]["allarme"] is True


def test_percorso_inesistente():
    assert controllo_disco.main(["/non/esiste"]) == 2
```

```bash
pytest -v
```

- `assert` verifica una condizione: se è falsa, il test fallisce e pytest mostra i valori coinvolti.
- `capsys` è una **fixture** di pytest: un oggetto che pytest prepara per il test, in questo caso per catturare
  l'output.
- Un buon test verifica **un comportamento** ed ha un nome che lo descrive.

## Esercizi

!!! esercizio "Esercizio 1 · Da YAML a report"
    Crea `server.yaml` con l'esempio della lezione (aggiungi un terzo server). Scrivi `report.py` che lo legga e
    produca `report.csv` con le colonne `nome`, `numero_porte`, `porte` (le porte separate da `;`). Il nome del
    file YAML si passa come argomento con `argparse`.

??? success "Soluzione"
    ```python
    #!/usr/bin/env python3
    import argparse
    import csv

    import yaml

    parser = argparse.ArgumentParser()
    parser.add_argument("file_yaml")
    parser.add_argument("-o", "--output", default="report.csv")
    args = parser.parse_args()

    with open(args.file_yaml) as f:
        conf = yaml.safe_load(f)

    with open(args.output, "w", newline="") as f:
        w = csv.DictWriter(f, fieldnames=["nome", "numero_porte", "porte"])
        w.writeheader()
        for s in conf["server"]:
            porte = s.get("porte", [])
            w.writerow({
                "nome": s["nome"],
                "numero_porte": len(porte),
                "porte": ";".join(str(p) for p in porte),
            })
    print(f"Scritto {args.output}")
    ```
    **Spiegazione.** `";".join()` vuole stringhe, per questo convertiamo ogni porta con `str()`. `s.get("porte", [])`
    rende lo script tollerante verso un server senza porte.

!!! esercizio "Esercizio 2 · Interrogare un'API"
    Scrivi `repo_info.py PROPRIETARIO/REPOSITORY` che usi l'API pubblica di GitHub
    (`https://api.github.com/repos/PROPRIETARIO/REPOSITORY`) per stampare: descrizione, linguaggio principale,
    numero di stelle e data dell'ultimo push. Se il repository non esiste deve stampare un messaggio chiaro e
    uscire con codice 1. Se è impostata la variabile d'ambiente `GITHUB_TOKEN`, usala per autenticarti.

??? success "Soluzione"
    ```python
    #!/usr/bin/env python3
    import argparse
    import os
    import sys

    import requests

    parser = argparse.ArgumentParser()
    parser.add_argument("repo", help="nel formato proprietario/nome")
    args = parser.parse_args()

    headers = {"Accept": "application/vnd.github+json"}
    token = os.environ.get("GITHUB_TOKEN")
    if token:
        headers["Authorization"] = f"Bearer {token}"

    try:
        r = requests.get(f"https://api.github.com/repos/{args.repo}", headers=headers, timeout=10)
    except requests.RequestException as e:
        sys.exit(f"Errore di rete: {e}")

    if r.status_code == 404:
        print(f"Repository {args.repo} non trovato", file=sys.stderr)
        sys.exit(1)
    r.raise_for_status()

    d = r.json()
    print(f"Descrizione: {d['description']}")
    print(f"Linguaggio:  {d['language']}")
    print(f"Stelle:      {d['stargazers_count']}")
    print(f"Ultimo push: {d['pushed_at']}")
    ```
    **Spiegazione.** `sys.exit("messaggio")` stampa il messaggio su standard error ed esce con codice 1. Senza token
    l'API di GitHub consente un numero limitato di richieste all'ora per indirizzo IP: con il token il limite è
    molto più alto. `requests.RequestException` è la classe base di tutti gli errori di rete di `requests`.

!!! esercizio "Esercizio 3 · Testare il proprio codice"
    Riprendi la funzione `valida_porta` della [lezione 2.6](06-python-base.md), mettila in `porte.py` e scrivi `test_porte.py` con
    almeno quattro test: una porta valida, una stringa non numerica, una porta 0 e una porta 70000. Per verificare
    che venga sollevata un'eccezione usa `pytest.raises`.

??? success "Soluzione"
    ```python
    import pytest

    from porte import valida_porta


    def test_porta_valida():
        assert valida_porta("443") == 443


    @pytest.mark.parametrize("valore", ["abc", "", "0", "70000", "-1"])
    def test_valori_non_validi(valore):
        with pytest.raises(ValueError):
            valida_porta(valore)
    ```
    **Spiegazione.** `pytest.raises` fa passare il test **solo se** il blocco solleva l'eccezione indicata.
    `@pytest.mark.parametrize` esegue lo stesso test con valori diversi: qui genera cinque test da una sola funzione,
    e se uno fallisce pytest dice quale valore lo ha causato.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale funzione usi per leggere YAML da una fonte non fidata?","options":["yaml.load","yaml.safe_load","json.load","yaml.read"],"answer":[1],"explanation":"yaml.load con il loader completo può costruire oggetti Python arbitrari."},
 {"q":"Perché è pericoloso `subprocess.run(f\"ls {nome}\", shell=True)` con un nome fornito dall'utente?","options":["È lento","Permette la command injection","Non cattura l'output","Non funziona su Linux"],"answer":[1],"explanation":"Passa una lista di argomenti senza shell=True."},
 {"q":"Cosa succede a una richiesta `requests.get(url)` senza timeout se il server non risponde?","options":["Fallisce dopo 30 secondi","Può restare bloccata indefinitamente","Riprova tre volte","Restituisce None"],"answer":[1],"explanation":"requests non ha un timeout predefinito."},
 {"q":"In YAML, cosa diventa `versione: 3.10` senza virgolette?","options":["La stringa \"3.10\"","Il numero 3.1","Un errore","La lista [3, 10]"],"answer":[1],"explanation":"Viene interpretato come numero decimale: usa le virgolette per le stringhe."},
 {"q":"Quali affermazioni su pytest sono corrette?","options":["Esegue le funzioni il cui nome inizia con test_","Richiede di scrivere una classe per ogni test","Con pytest.raises verifica che venga sollevata un'eccezione","Funziona solo con Flask","Esegue i test solo in ordine alfabetico obbligatorio"],"answer":[0,2],"explanation":"Le classi sono facoltative; pytest funziona con qualsiasi codice Python."},
 {"q":"Dove scrive `logging` per impostazione predefinita?","options":["Su standard output","Su standard error","In /var/log/python.log","Nel journal direttamente"],"answer":[1],"explanation":"Sotto systemd o in un container, standard error finisce comunque nei log."}
]}
</script>

## Fonti

- [Python: pathlib](https://docs.python.org/3.12/library/pathlib.html)
- [Python: json](https://docs.python.org/3.12/library/json.html), [csv](https://docs.python.org/3.12/library/csv.html)
- [Python: subprocess](https://docs.python.org/3.12/library/subprocess.html) e [Security Considerations](https://docs.python.org/3.12/library/subprocess.html#security-considerations)
- [Python: argparse tutorial](https://docs.python.org/3.12/howto/argparse.html), [Logging HOWTO](https://docs.python.org/3.12/howto/logging.html)
- [PyYAML documentation](https://pyyaml.org/wiki/PyYAMLDocumentation)
- [Requests: Quickstart](https://requests.readthedocs.io/en/latest/user/quickstart/)
- [pytest: Get started](https://docs.pytest.org/en/stable/getting-started.html)
- [GitHub REST API: Repositories](https://docs.github.com/en/rest/repos/repos#get-a-repository)

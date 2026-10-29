# 2.8 API HTTP con Flask

!!! obiettivi "Obiettivi"
    - Capire cos'è un'API REST: risorse, metodi, codici di stato, JSON.
    - Scrivere un'applicazione Flask con più endpoint.
    - Salvare i dati in un database con SQLAlchemy, configurato da una variabile d'ambiente.
    - Testare l'API con il client di test di Flask e pytest.
    - Eseguire l'applicazione con un server WSGI di produzione (Gunicorn).

In questa lezione costruisci la **Tasks API**, l'applicazione che accompagnerà il resto del corso.

## Le API REST

Un'**API** (*Application Programming Interface*) è un'interfaccia che un programma offre ad altri programmi.
Un'**API web** si usa via HTTP ([lezione 1.12](../fase-1/12-http-nginx.md)). Lo stile più diffuso è **REST**, basato su alcune idee semplici:

- i dati sono **risorse** identificate da un URL: `/tasks` è la collezione, `/tasks/3` la singola attività;
- il **metodo HTTP** dice cosa fare sulla risorsa;
- i dati viaggiano di solito in **JSON**;
- il **codice di stato** dice com'è andata;
- ogni richiesta è **indipendente** (*stateless*): il server non ricorda le richieste precedenti.

| Metodo | Su `/tasks` | Su `/tasks/3` |
|---|---|---|
| `GET` | elenca le attività | legge l'attività 3 |
| `POST` | crea una nuova attività | — |
| `PUT` | — | sostituisce l'attività 3 |
| `PATCH` | — | modifica alcuni campi dell'attività 3 |
| `DELETE` | — | cancella l'attività 3 |

| Codice | Significato | Quando |
|---|---|---|
| `200 OK` | successo | lettura o modifica riuscita |
| `201 Created` | risorsa creata | dopo un `POST` riuscito |
| `204 No Content` | successo, nessun contenuto | dopo un `DELETE` |
| `400 Bad Request` | richiesta non valida | dati mancanti o sbagliati |
| `404 Not Found` | risorsa inesistente | `/tasks/999` |
| `405 Method Not Allowed` | metodo non supportato | `DELETE /tasks` |
| `500 Internal Server Error` | errore del server | un bug |
| `503 Service Unavailable` | servizio non disponibile | il database non risponde |

`GET`, `PUT` e `DELETE` sono **idempotenti**: ripeterli più volte ha lo stesso effetto di farlo una volta. `POST`
non lo è: due `POST` creano due attività.

## Flask in dieci righe

**Flask** è un *micro-framework* web per Python: fornisce il minimo indispensabile (instradamento delle richieste,
gestione di richieste e risposte) e lascia a te il resto. È semplice da imparare e molto usato per API e servizi
interni.

```bash
mkdir -p ~/tasks-api && cd ~/tasks-api
python3 -m venv .venv && source .venv/bin/activate
pip install Flask
```

```python title="app.py"
from flask import Flask, jsonify

app = Flask(__name__)


@app.get("/health")
def health():
    return jsonify(status="ok")
```

```bash
flask --app app run --port 8000 --debug
# in un altro terminale:
curl -i http://127.0.0.1:8000/health
```

- `@app.get("/health")` è un **decoratore**: registra la funzione come gestore delle richieste `GET` su `/health`.
  Esistono anche `@app.post`, `@app.patch`, `@app.delete` e la forma generica `@app.route(..., methods=[...])`.
- `jsonify` trasforma i dati in una risposta JSON con l'intestazione `Content-Type: application/json`.
- `--debug` ricarica il codice a ogni modifica e mostra gli errori nel browser. **Mai in produzione**: la pagina di
  debug permette di eseguire codice sul server.

## Primo passo: i dati in memoria

Iniziamo tenendo le attività in una lista Python, per concentrarci sull'API:

```python title="app.py"
from flask import Flask, jsonify, request

app = Flask(__name__)
tasks = []          # ogni attività è un dizionario: {"id", "title", "done"}
next_id = 1


@app.get("/tasks")
def list_tasks():
    return jsonify(tasks)


@app.post("/tasks")
def create_task():
    global next_id
    data = request.get_json(silent=True) or {}
    title = str(data.get("title", "")).strip()
    if not title:
        return jsonify(error="il campo 'title' è obbligatorio"), 400
    task = {"id": next_id, "title": title, "done": False}
    next_id += 1
    tasks.append(task)
    return jsonify(task), 201


@app.get("/tasks/<int:task_id>")
def get_task(task_id):
    for t in tasks:
        if t["id"] == task_id:
            return jsonify(t)
    return jsonify(error="attività non trovata"), 404
```

- `request.get_json(silent=True)` legge il corpo JSON della richiesta; con `silent=True` restituisce `None` invece
  di un errore se il corpo non è JSON valido. `or {}` gestisce quel caso.
- Una funzione può restituire una **tupla** `(corpo, codice)`: è il modo di impostare il codice di stato.
- `<int:task_id>` è una **variabile nel percorso**: Flask la converte in intero e la passa alla funzione; con un
  valore non numerico risponde da solo `404`.
- **Valida sempre l'input**: mai fidarsi di ciò che arriva dal client.

```bash
curl -s -X POST http://127.0.0.1:8000/tasks -H 'Content-Type: application/json' -d '{"title": "Studiare Flask"}'
curl -s http://127.0.0.1:8000/tasks
curl -s -i http://127.0.0.1:8000/tasks/42
```

Il limite è evidente: al riavvio dell'applicazione i dati spariscono. Serve un database.

## Secondo passo: un database con SQLAlchemy

**SQLAlchemy** è la libreria di riferimento per i database in Python; **Flask-SQLAlchemy** la integra con Flask.
Permette di descrivere le tabelle come classi Python (un **ORM**, *Object-Relational Mapper*) e di cambiare
database modificando solo l'**URL di connessione**:

| Database | URL |
|---|---|
| SQLite (un file locale, nessun server) | `sqlite:///tasks.db` |
| SQLite in memoria (per i test) | `sqlite:///:memory:` |
| PostgreSQL (dalla [Fase 3](../fase-3/index.md)) | `postgresql+psycopg://utente:password@host:5432/nomedb` |

```bash
pip install Flask-SQLAlchemy
```

```python
from flask_sqlalchemy import SQLAlchemy

db = SQLAlchemy()


class Task(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.String(200), nullable=False)
    done = db.Column(db.Boolean, nullable=False, default=False)

    def to_dict(self):
        return {"id": self.id, "title": self.title, "done": self.done}
```

Le operazioni principali:

```python
task = Task(title="Studiare SQL")
db.session.add(task)              # prepara l'inserimento
db.session.commit()               # lo scrive nel database (ora task.id è assegnato)

db.session.get(Task, 1)           # per chiave primaria, None se non esiste
Task.query.order_by(Task.id).all()
Task.query.filter_by(done=True).all()

task.done = True
db.session.commit()               # salva la modifica

db.session.delete(task)
db.session.commit()
```

## La configurazione: la factory e le variabili d'ambiente

Un'applicazione ben scritta non contiene nel codice l'indirizzo del database né le password: li legge
dall'**ambiente**. È uno dei principi della metodologia [Twelve-Factor App](https://12factor.net/it/config), che
ritroverai con i container, il cloud e Kubernetes: lo stesso codice gira in sviluppo, test e produzione, cambia solo
la configurazione.

Inoltre, invece di creare l'applicazione a livello di modulo, si usa una **application factory**: una funzione
`create_app()` che la costruisce. I test possono così crearne una nuova, con un database in memoria, per ogni test.

Ecco la versione completa. Scrivila tu, un endpoint alla volta, provandolo con `curl`:

```python title="app.py"
"""Tasks API: l'applicazione d'esempio del corso."""
import os
import socket

from flask import Flask, jsonify, request
from flask_sqlalchemy import SQLAlchemy

db = SQLAlchemy()


class Task(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.String(200), nullable=False)
    done = db.Column(db.Boolean, nullable=False, default=False)

    def to_dict(self):
        return {"id": self.id, "title": self.title, "done": self.done}


def create_app(database_url=None):
    app = Flask(__name__)
    app.config["SQLALCHEMY_DATABASE_URI"] = (
        database_url or os.environ.get("DATABASE_URL") or "sqlite:///tasks.db"
    )
    db.init_app(app)
    with app.app_context():
        db.create_all()             # crea le tabelle se non esistono

    @app.get("/health")
    def health():
        try:
            db.session.execute(db.text("SELECT 1"))
            return jsonify(status="ok", host=socket.gethostname())
        except Exception:           # il database non risponde
            return jsonify(status="error"), 503

    @app.get("/tasks")
    def list_tasks():
        return jsonify([t.to_dict() for t in Task.query.order_by(Task.id)])

    @app.post("/tasks")
    def create_task():
        data = request.get_json(silent=True) or {}
        title = str(data.get("title", "")).strip()
        if not title:
            return jsonify(error="il campo 'title' è obbligatorio"), 400
        task = Task(title=title[:200])
        db.session.add(task)
        db.session.commit()
        return jsonify(task.to_dict()), 201

    @app.get("/tasks/<int:task_id>")
    def get_task(task_id):
        task = db.session.get(Task, task_id)
        if task is None:
            return jsonify(error="attività non trovata"), 404
        return jsonify(task.to_dict())

    @app.patch("/tasks/<int:task_id>")
    def update_task(task_id):
        task = db.session.get(Task, task_id)
        if task is None:
            return jsonify(error="attività non trovata"), 404
        data = request.get_json(silent=True) or {}
        if "title" in data:
            title = str(data["title"]).strip()
            if not title:
                return jsonify(error="il titolo non può essere vuoto"), 400
            task.title = title[:200]
        if "done" in data:
            task.done = bool(data["done"])
        db.session.commit()
        return jsonify(task.to_dict())

    @app.delete("/tasks/<int:task_id>")
    def delete_task(task_id):
        task = db.session.get(Task, task_id)
        if task is None:
            return jsonify(error="attività non trovata"), 404
        db.session.delete(task)
        db.session.commit()
        return "", 204

    return app


if __name__ == "__main__":
    create_app().run(host="127.0.0.1", port=8000, debug=True)
```

L'endpoint `/health` merita attenzione: non risponde semplicemente "ok", ma **verifica il database**. Se il
database non risponde, restituisce `503`. Load balancer ([Fase 5](../fase-5/index.md)) e Kubernetes ([Fase 8](../fase-8/index.md)) useranno questo endpoint
per decidere se mandare traffico all'applicazione. Il nome host nella risposta ti servirà per vedere quale copia
dell'applicazione ha risposto.

```bash
python app.py
curl -s -X POST localhost:8000/tasks -H 'Content-Type: application/json' -d '{"title":"Imparare SQLAlchemy"}'
curl -s -X PATCH localhost:8000/tasks/1 -H 'Content-Type: application/json' -d '{"done":true}'
curl -s localhost:8000/tasks
# Ctrl+C, poi riavvia: i dati ci sono ancora, nel file instance/tasks.db
```

Con Flask-SQLAlchemy, un percorso SQLite relativo viene creato nella cartella `instance/` del progetto. Aggiungila al
`.gitignore`, insieme a `.venv/`.

## I test

Flask fornisce un **client di test** che simula le richieste HTTP senza avviare un server. Con una **fixture** di
pytest prepariamo un'applicazione nuova, con database in memoria, per ogni test:

```python title="tests/test_app.py"
import pytest

from app import create_app


@pytest.fixture
def client():
    app = create_app("sqlite:///:memory:")
    app.config["TESTING"] = True
    with app.test_client() as client:
        yield client


def test_health(client):
    r = client.get("/health")
    assert r.status_code == 200
    assert r.get_json()["status"] == "ok"


def test_create_and_list(client):
    r = client.post("/tasks", json={"title": "Studiare Git"})
    assert r.status_code == 201
    assert r.get_json() == {"id": 1, "title": "Studiare Git", "done": False}
    assert client.get("/tasks").get_json() == [{"id": 1, "title": "Studiare Git", "done": False}]


def test_title_required(client):
    assert client.post("/tasks", json={}).status_code == 400
    assert client.post("/tasks", json={"title": "   "}).status_code == 400


def test_update_and_delete(client):
    client.post("/tasks", json={"title": "Scrivere test"})
    r = client.patch("/tasks/1", json={"done": True})
    assert r.get_json()["done"] is True
    assert client.delete("/tasks/1").status_code == 204
    assert client.get("/tasks/1").status_code == 404
```

Perché pytest trovi `app.py` partendo dalla cartella `tests/`, crea nella radice del progetto un file `conftest.py`
(anche vuoto): pytest aggiunge quella cartella ai percorsi di importazione.

```bash
pip install pytest
pytest -v
```

!!! tip "Cosa testare"
    Il **percorso felice** (tutto va bene) e soprattutto i **casi di errore**: dati mancanti, risorse inesistenti,
    valori limite. I bug si nascondono lì.

## Le dipendenze

Separa le dipendenze dell'applicazione da quelle che servono solo per sviluppare:

```text title="requirements.txt"
Flask==3.1.2
Flask-SQLAlchemy==3.1.1
gunicorn==23.0.0
```

```text title="requirements-dev.txt"
-r requirements.txt
pytest==8.4.2
```

Le versioni **bloccate** (`==`) garantiscono che l'applicazione si installi sempre allo stesso modo: in sviluppo,
nella pipeline e in produzione.

## In produzione: Gunicorn

Il server integrato di Flask (`flask run`, `app.run()`) è pensato **solo per lo sviluppo**: gestisce male il carico
e non è progettato per essere sicuro. In produzione un'applicazione Python web si esegue con un **server WSGI**.
**Gunicorn** è il più diffuso:

```bash
pip install gunicorn
gunicorn --bind 127.0.0.1:8000 --workers 2 "app:create_app()"
```

- `"app:create_app()"` dice a Gunicorn: nel modulo `app`, chiama `create_app()` per ottenere l'applicazione.
- `--workers` è il numero di processi che servono le richieste in parallelo. La documentazione di Gunicorn suggerisce
  come punto di partenza `(2 × numero di CPU) + 1`.
- Davanti a Gunicorn si mette di solito un **reverse proxy** come Nginx ([lezione 1.12](../fase-1/12-http-nginx.md)), che gestisce HTTPS e i
  client lenti.

```mermaid
flowchart LR
  C[Client] -->|HTTPS 443| N[Nginx<br>reverse proxy]
  N -->|HTTP 127.0.0.1:8000| G[Gunicorn<br>2 worker]
  G --> A[Tasks API<br>Flask]
  A --> D[(Database)]
```

È esattamente l'architettura del progetto della [Fase 1](../fase-1/index.md), con la tua API al posto di `python3 -m http.server`.

## Esercizi

!!! esercizio "Esercizio 1 · Costruisci la Tasks API"
    Nel Dev Container crea il repository `tasks-api` con: `app.py`, `conftest.py`, `tests/test_app.py`,
    `requirements.txt`, `requirements-dev.txt`, `.gitignore`, `README.md` (con la tabella degli endpoint). Tutti i
    test devono passare. Pubblicalo su GitHub. Confrontalo solo alla fine con
    [la versione di riferimento](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/tasks-api).

??? success "Soluzione"
    ```bash
    mkdir -p ~/tasks-api/tests && cd ~/tasks-api
    git init
    python3 -m venv .venv && source .venv/bin/activate
    # scrivi i file come nella lezione
    printf '.venv/\n__pycache__/\ninstance/\n*.db\n' > .gitignore
    touch conftest.py
    pip install -r requirements-dev.txt
    pytest -v                       # 4 passed
    git add . && git commit -m "feat: prima versione della Tasks API"
    gh repo create tasks-api --public --source=. --push
    ```
    **Spiegazione.** `instance/` e `*.db` sono dati, non codice: non vanno versionati. Il `README.md` deve spiegare a
    un estraneo come installare, avviare e testare l'applicazione in meno di un minuto.

!!! esercizio "Esercizio 2 · Un filtro"
    Aggiungi a `GET /tasks` il parametro facoltativo `done`: `/tasks?done=true` restituisce solo le attività
    completate, `/tasks?done=false` solo quelle da fare. Un valore diverso da `true` o `false` deve restituire
    `400`. Scrivi i test **prima** del codice, verifica che falliscano, poi implementa.

??? success "Soluzione"
    ```python
    @app.get("/tasks")
    def list_tasks():
        query = Task.query.order_by(Task.id)
        done = request.args.get("done")
        if done is not None:
            if done not in ("true", "false"):
                return jsonify(error="il parametro 'done' deve essere true o false"), 400
            query = query.filter_by(done=(done == "true"))
        return jsonify([t.to_dict() for t in query])
    ```
    ```python
    def test_filter_done(client):
        client.post("/tasks", json={"title": "A"})
        client.post("/tasks", json={"title": "B"})
        client.patch("/tasks/2", json={"done": True})
        assert [t["title"] for t in client.get("/tasks?done=true").get_json()] == ["B"]
        assert [t["title"] for t in client.get("/tasks?done=false").get_json()] == ["A"]
        assert client.get("/tasks?done=forse").status_code == 400
    ```
    **Spiegazione.** `request.args` contiene i parametri della query string, sempre come stringhe. Scrivere il test
    prima del codice (*Test-Driven Development*) ti obbliga a definire con precisione il comportamento atteso.

!!! esercizio "Esercizio 3 · Dietro Nginx, come servizio"
    Sulla VM `lab`, installa la Tasks API in `/srv/tasks-api` con un ambiente virtuale, avviala con Gunicorn tramite
    un servizio systemd che giri come utente `webapp`, con `DATABASE_URL=sqlite:////srv/tasks-api/data/tasks.db`, e
    configura Nginx come reverse proxy. Verifica dal `client` che `POST /tasks` funzioni e che i dati sopravvivano a
    un riavvio.

??? success "Soluzione"
    ```bash
    sudo useradd --system --home /srv/tasks-api --shell /usr/sbin/nologin webapp 2>/dev/null || true
    sudo git clone https://github.com/TUO-UTENTE/tasks-api.git /srv/tasks-api
    sudo python3 -m venv /srv/tasks-api/.venv
    sudo /srv/tasks-api/.venv/bin/pip install -r /srv/tasks-api/requirements.txt
    sudo install -d -o webapp -g webapp /srv/tasks-api/data
    ```
    ```ini title="/etc/systemd/system/tasks-api.service"
    [Unit]
    Description=Tasks API
    After=network.target

    [Service]
    User=webapp
    Group=webapp
    WorkingDirectory=/srv/tasks-api
    Environment=DATABASE_URL=sqlite:////srv/tasks-api/data/tasks.db
    ExecStart=/srv/tasks-api/.venv/bin/gunicorn --bind 127.0.0.1:8000 --workers 2 "app:create_app()"
    Restart=on-failure
    NoNewPrivileges=true
    ProtectSystem=strict
    ReadWritePaths=/srv/tasks-api/data
    ProtectHome=true
    PrivateTmp=true

    [Install]
    WantedBy=multi-user.target
    ```
    ```bash
    sudo systemctl daemon-reload && sudo systemctl enable --now tasks-api
    curl -s localhost:8000/health
    ```
    Nella configurazione Nginx della [lezione 1.12](../fase-1/12-http-nginx.md) basta che `proxy_pass` punti a `http://127.0.0.1:8000`. Dal client:
    ```bash
    curl -sk -X POST https://web.lab/tasks -H 'Content-Type: application/json' -d '{"title":"Dal client"}'
    ```
    **Spiegazione.** Gli URL SQLite assoluti hanno **quattro** barre: `sqlite:///` più il percorso che inizia con `/`.
    `ProtectSystem=strict` rende tutto il filesystem in sola lettura per il servizio, quindi la cartella del database
    va concessa con `ReadWritePaths`. Sul percorso RHEL ricorda il booleano SELinux `httpd_can_network_connect`
    ([lezione 1.12](../fase-1/12-http-nginx.md)), altrimenti Nginx non può collegarsi a Gunicorn e risponde `502 Bad Gateway`.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale codice di stato restituisci dopo aver creato una risorsa con POST?","options":["200","201","204","302"],"answer":[1],"explanation":"201 Created."},
 {"q":"Quali metodi HTTP sono idempotenti?","options":["GET","POST","DELETE","Nessuno","Solo POST e PATCH"],"answer":[0,2],"explanation":"Ripetere GET o DELETE ha lo stesso effetto; due POST creano due risorse."},
 {"q":"Perché l'URL del database si legge da una variabile d'ambiente?","options":["Perché Flask lo impone","Perché lo stesso codice deve girare in ambienti diversi senza modifiche, e i segreti non devono stare nel codice","Per velocità","Perché SQLite lo richiede"],"answer":[1],"explanation":"È il principio Config della Twelve-Factor App."},
 {"q":"Perché non si usa `flask run` in produzione?","options":["Perché non supporta JSON","Perché è un server di sviluppo, non progettato per carico e sicurezza","Perché è a pagamento","Perché funziona solo con Python 2"],"answer":[1],"explanation":"In produzione si usa un server WSGI come Gunicorn, dietro un reverse proxy."},
 {"q":"Cosa restituisce l'endpoint /health della Tasks API se il database non risponde?","options":["200","404","503","500"],"answer":[2],"explanation":"503 Service Unavailable: load balancer e Kubernetes smetteranno di mandarle traffico."},
 {"q":"A cosa serve l'application factory create_app()?","options":["A creare il database","A costruire l'applicazione su richiesta, per esempio con configurazioni diverse nei test","A generare la documentazione","A creare i container"],"answer":[1],"explanation":"Ogni test può creare un'applicazione con un database in memoria pulito."}
]}
</script>

## Fonti

- [Flask: Quickstart](https://flask.palletsprojects.com/en/stable/quickstart/)
- [Flask: Application Factories](https://flask.palletsprojects.com/en/stable/patterns/appfactories/)
- [Flask: Testing Flask Applications](https://flask.palletsprojects.com/en/stable/testing/)
- [Flask: Deploying to Production](https://flask.palletsprojects.com/en/stable/deploying/)
- [Flask-SQLAlchemy: Quick Start](https://flask-sqlalchemy.readthedocs.io/en/stable/quickstart/)
- [Gunicorn: Design e numero di worker](https://docs.gunicorn.org/en/stable/design.html)
- [MDN: HTTP response status codes](https://developer.mozilla.org/it/docs/Web/HTTP/Status)
- [The Twelve-Factor App](https://12factor.net/it/)

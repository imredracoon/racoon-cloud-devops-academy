# Tasks API

L'applicazione d'esempio della Racoon Cloud & DevOps Academy: una piccola API REST scritta in Python con Flask,
che cresce fase dopo fase (Fase 2: codice e test; Fase 3: container; Fase 5: cloud; Fase 7: pipeline; Fase 8: Kubernetes).

```bash
python3 -m venv .venv && . .venv/bin/activate
pip install -r requirements-dev.txt
pytest
python app.py            # http://127.0.0.1:8000
```

| Metodo | Percorso | Descrizione |
|---|---|---|
| GET | `/health` | Stato dell'applicazione e del database |
| GET | `/tasks` | Elenco delle attività |
| POST | `/tasks` | Crea un'attività: `{"title": "..."}` |
| GET | `/tasks/<id>` | Una attività |
| PATCH | `/tasks/<id>` | Modifica `title` e/o `done` |
| DELETE | `/tasks/<id>` | Cancella |

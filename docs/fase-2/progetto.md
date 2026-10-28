# Progetto 2 · Toolkit di automazione e Tasks API

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai trasformare operazioni manuali in codice affidabile: script Bash robusti, uno strumento Python con test,
    un'API web pronta per la produzione e un'installazione ripetibile su un server. Insieme al progetto della
    Fase 1, mostra a un selezionatore che sai sia gestire un sistema sia automatizzarlo.

| | |
|---|---|
| **Durata** | 12-18 ore |
| **Consegna** | Due repository pubblici: `tasks-api` (iniziato nella lezione 2.8) e `ops-toolkit` |
| **Ambiente** | Dev Container per lo sviluppo, VM `web` del Progetto 1 per l'installazione |

## Il compito

### Repository `tasks-api`

1. La Tasks API della lezione 2.8, con il filtro `?done=` dell'esercizio 2.
2. Almeno **8 test** pytest, inclusi i casi di errore (400, 404) e il filtro.
3. Un endpoint in più a tua scelta, con i suoi test. Per esempio `GET /stats`, che restituisce
   `{"total": N, "done": N, "todo": N}`.
4. `README.md` con: descrizione, tabella degli endpoint con esempi `curl`, come eseguire i test, come avviarla in
   sviluppo e con Gunicorn, le variabili d'ambiente supportate.
5. Storia Git pulita: lavora con un branch e una pull request per ogni funzionalità, con messaggi chiari.

### Repository `ops-toolkit`

1. `bin/backup.sh`: backup con rotazione (lezione 2.4), con opzioni `-k NUMERO` (archivi da conservare) e `-h`.
2. `bin/check-sites.sh`: verifica di più URL da un file (lezione 2.4).
3. `toolkit/healthcheck.py`: uno strumento Python che interroga l'endpoint `/health` di una o più istanze della
   Tasks API, con:
    - argomenti con `argparse`: gli URL, `--timeout`, `--json`;
    - log con `logging`;
    - codice di uscita 0 se tutte le istanze sono sane, 1 altrimenti, 2 per un uso sbagliato;
    - test pytest che **non** richiedono una rete (vedi i suggerimenti).
4. `systemd/`: un servizio e un timer che eseguono `backup.sh` ogni notte con `Persistent=true`.
5. `deploy/install-tasks-api.sh`: uno script idempotente che installa la Tasks API sulla VM `web` (utente di
   servizio, ambiente virtuale, servizio systemd con Gunicorn, configurazione Nginx), così che si possa rilanciare
   senza errori.
6. `README.md` con lo scopo di ogni strumento ed esempi d'uso; `Makefile` o sezione "Sviluppo" con i comandi per
   lint e test.

Tutti gli script Bash devono passare `shellcheck` senza avvisi.

### Sul server

La VM `web` del Progetto 1 deve servire la Tasks API su `https://web.lab/` al posto della pagina statica, con i
dati salvati in modo persistente e il backup notturno del database attivo.

## Suggerimenti

- Per testare `healthcheck.py` senza rete, separa la logica in una funzione che riceve una **funzione per fare la
  richiesta** come parametro, oppure usa la fixture `monkeypatch` di pytest per sostituire `requests.get` con una
  funzione finta:

    ```python
    class RispostaFinta:
        def __init__(self, codice, dati):
            self.status_code = codice
            self._dati = dati
        def json(self):
            return self._dati

    def test_istanza_sana(monkeypatch):
        monkeypatch.setattr(healthcheck.requests, "get",
                            lambda url, timeout: RispostaFinta(200, {"status": "ok"}))
        assert healthcheck.main(["http://finto"]) == 0
    ```

- Un `Makefile` minimo rende i comandi facili da ricordare (e la pipeline della Fase 7 li riuserà):

    ```makefile
    .PHONY: lint test
    lint:
    	shellcheck bin/*.sh deploy/*.sh
    	ruff check .
    test:
    	pytest -v
    ```

    Attenzione: in un Makefile le righe dei comandi devono iniziare con un **tab**, non con spazi.

- Per l'idempotenza di `install-tasks-api.sh`: `id webapp || useradd ...`, `[ -d .venv ] || python3 -m venv ...`,
  `git -C /srv/tasks-api pull` se la cartella esiste già, `install` per copiare i file di configurazione,
  `nginx -t` prima di `systemctl reload nginx`.

## Come verificare

| Test | Risultato atteso |
|---|---|
| `pytest` in entrambi i repository | tutti i test passano |
| `shellcheck bin/*.sh deploy/*.sh` | nessun avviso |
| `install-tasks-api.sh` lanciato due volte di fila sulla VM | nessun errore, nessuna modifica la seconda volta |
| dal client: `curl -sk https://web.lab/health` | `{"status": "ok", ...}` |
| dal client: `POST /tasks`, poi riavvio di `web`, poi `GET /tasks` | l'attività c'è ancora |
| `systemctl list-timers` su `web` | il timer del backup è presente con la prossima esecuzione |
| `python toolkit/healthcheck.py https://web.lab --json` (con `web.lab` raggiungibile) | JSON con lo stato e codice di uscita 0 |
| fermata la Tasks API, lo stesso comando | codice di uscita 1 |

## Criteri di valutazione

- [ ] Tutte le verifiche della tabella passano.
- [ ] I README permettono a un estraneo di usare gli strumenti senza chiederti nulla.
- [ ] Nessun segreto, ambiente virtuale o database nei repository.
- [ ] Ogni funzionalità è arrivata in `main` con una pull request.
- [ ] Gli script gestiscono gli errori: argomenti mancanti, file inesistenti, servizi irraggiungibili.

??? success "Una possibile soluzione: healthcheck.py"
    Prova a scriverlo da solo prima di aprire questa sezione.

    ```python title="toolkit/healthcheck.py"
    #!/usr/bin/env python3
    """Verifica lo stato di una o più istanze della Tasks API."""
    import argparse
    import json
    import logging
    import sys

    import requests

    log = logging.getLogger("healthcheck")


    def controlla(url: str, timeout: float) -> dict:
        try:
            r = requests.get(url.rstrip("/") + "/health", timeout=timeout)
            sano = r.status_code == 200 and r.json().get("status") == "ok"
            return {"url": url, "codice": r.status_code, "sano": sano}
        except (requests.RequestException, ValueError) as e:
            log.warning("%s: %s", url, e)
            return {"url": url, "codice": None, "sano": False}


    def main(argv: list[str] | None = None) -> int:
        parser = argparse.ArgumentParser(description=__doc__)
        parser.add_argument("urls", nargs="+", help="URL di base delle istanze")
        parser.add_argument("--timeout", type=float, default=5.0)
        parser.add_argument("--json", action="store_true")
        args = parser.parse_args(argv)
        logging.basicConfig(level=logging.INFO, format="%(levelname)s %(message)s")

        risultati = [controlla(u, args.timeout) for u in args.urls]
        if args.json:
            print(json.dumps(risultati, indent=2))
        else:
            for r in risultati:
                print(f"{'OK ' if r['sano'] else 'KO '} {r['url']} ({r['codice']})")
        return 0 if all(r["sano"] for r in risultati) else 1


    if __name__ == "__main__":
        sys.exit(main())
    ```

    **Spiegazione.** `argparse` esce da solo con codice 2 se gli argomenti sono sbagliati, quindi il requisito "2 per
    un uso sbagliato" è soddisfatto senza codice in più. `ValueError` copre il caso di una risposta che non è JSON
    valido. Con certificati autofirmati come quello di `web.lab`, `requests` rifiuta la connessione: puoi aggiungere
    un'opzione `--ca-file` e passarla come `verify=` a `requests.get`, che è più sicuro che disattivare la verifica.

## Fonti

- [pytest: monkeypatch](https://docs.pytest.org/en/stable/how-to/monkeypatch.html)
- [GNU make: manuale](https://www.gnu.org/software/make/manual/make.html)
- [Gunicorn: Deploying (con Nginx)](https://docs.gunicorn.org/en/stable/deploy.html)
- [ShellCheck](https://www.shellcheck.net/)
- [Requests: SSL Cert Verification](https://requests.readthedocs.io/en/latest/user/advanced/#ssl-cert-verification)

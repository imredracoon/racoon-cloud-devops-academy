# Progetto 7 · Pipeline e osservabilità

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai portare un'applicazione dal commit alla produzione in modo automatico e sicuro, e sapere in ogni momento
    se funziona. È il cuore del lavoro DevOps, e un progetto che si mostra bene in un colloquio: un link al
    repository con la pipeline verde, una dashboard, un incidente simulato documentato.

| | |
|---|---|
| **Durata** | 20-25 ore |
| **Consegna** | i repository `tasks-api` e `tasks-api-iac` aggiornati, con un documento `docs/osservabilita.md` |
| **Costo** | GitHub Actions gratuito per i repository pubblici; cloud solo durante le prove di deploy |

## Il compito

### Parte A · Pipeline completa

Nel repository `tasks-api`:

1. **CI** come nella lezione 7.3: lint, formattazione, test, build, scansione Trivy, pubblicazione su GHCR con tag
   `sha-...` e versioni semantiche.
2. Una **matrice** di test con Python 3.12 e 3.13 e il **report di copertura** come artefatto; la pipeline fallisce
   se la copertura scende sotto l'80% (`pytest --cov=app --cov-fail-under=80`).
3. **Protezione del ramo** `main` con controlli obbligatori.
4. **Deploy** con OIDC (lezione 7.4) sul servizio di container gestito della Fase 5 (ECS Fargate o Cloud Run),
   tramite l'ambiente `produzione` con approvazione manuale. Nessuna chiave di accesso nei segreti del repository.
5. **Dependabot** attivo per le azioni GitHub e le dipendenze Python (`.github/dependabot.yml`).

Nel repository `tasks-api-iac` (Fase 6):

6. Un workflow che, su ogni PR, esegue `fmt`, `validate`, `terraform test` dei moduli e un **piano** commentato nella
   PR, con un ruolo di sola lettura via OIDC; dopo il merge, l'**apply** con approvazione.

### Parte B · Osservabilità

7. Lo stack della lezione 7.6 (Prometheus, Alertmanager, node-exporter, Loki, Alloy, Grafana) avviato con la Tasks
   API in locale.
8. Una dashboard **"Tasks API · SLO"** versionata nel repository, con:
    - SLI di disponibilità: percentuale di richieste non 5xx, su 5 minuti e su 30 giorni;
    - budget di errore residuo per uno SLO del 99,5%;
    - p95 della latenza per endpoint;
    - pannello dei log degli errori da Loki.
9. Tre regole di alert in Prometheus, ognuna con `summary`, `description` e un link a un **runbook** (un file
   Markdown nel repository con i passi di diagnosi).
10. Alertmanager configurato con un ricevitore reale a tua scelta (per esempio un canale Slack o Discord gratuito, via
    webhook, o un'email) e un raggruppamento per `alertname`.

### Parte C · Il post-mortem

11. Simula **un incidente** (database fermo, oppure un rilascio che introduce un errore 500 su un endpoint) e scrivi
    un **post-mortem senza colpe** (*blameless*) in `docs/postmortem-01.md`: cronologia, impatto (quante richieste
    fallite, quanto budget di errore consumato), come è stato rilevato, causa, azioni correttive. Allega gli
    screenshot della dashboard e delle query.

## Come verificare

| Test | Risultato atteso |
|---|---|
| PR con un test che fallisce | pipeline rossa, merge bloccato |
| PR con codice non formattato | pipeline rossa al passo di formattazione |
| Merge su `main` | immagine `sha-...` su GHCR, deploy in attesa di approvazione |
| Approvazione | il servizio nel cloud usa la nuova immagine (`/metrics` risponde) |
| Segreti del repository | nessuna chiave AWS o Google |
| Database fermo per 2 minuti | alert `firing` e notifica ricevuta; dashboard SLO che mostra il consumo di budget |
| `docker compose ... down` e `up` | dashboard e sorgenti dati ricompaiono (provisioning) |

## Criteri di valutazione

- [ ] Pipeline CI veloce (meno di 5 minuti con la cache), con controlli di qualità, sicurezza e copertura.
- [ ] Deploy senza credenziali di lunga durata, con ambienti protetti e privilegi minimi.
- [ ] Infrastruttura con piano in PR e apply approvato.
- [ ] Dashboard e alert versionati, con SLO e runbook.
- [ ] Post-mortem chiaro, basato sui dati.

??? success "Una possibile soluzione: le query dello SLO"
    ```text
    # SLI di disponibilità sugli ultimi 5 minuti (0-1)
    1 - (
      sum(rate(http_requests_total{job="tasks-api", status=~"5.."}[5m]))
        /
      sum(rate(http_requests_total{job="tasks-api"}[5m]))
    )

    # SLI su 30 giorni
    1 - (
      sum(increase(http_requests_total{job="tasks-api", status=~"5.."}[30d]))
        /
      sum(increase(http_requests_total{job="tasks-api"}[30d]))
    )

    # Budget di errore residuo per uno SLO del 99,5% (1 = intatto, 0 = esaurito, negativo = superato)
    1 - (
      (
        sum(increase(http_requests_total{job="tasks-api", status=~"5.."}[30d]))
          /
        sum(increase(http_requests_total{job="tasks-api"}[30d]))
      ) / 0.005
    )
    ```

    **Spiegazione.** Il budget è la frazione di errori ammessa (0,5%); dividendo il tasso di errore osservato per il
    budget ottieni quanto ne hai consumato. Su 30 giorni la query è costosa: in produzione si usano **recording rule**
    che precalcolano i tassi a intervalli brevi. Il laboratorio conserva le metriche per 7 giorni
    (`--storage.tsdb.retention.time=7d`): la finestra di 30 giorni mostrerà quindi i dati disponibili.

??? success "Una possibile soluzione: dependabot.yml"
    ```yaml
    version: 2
    updates:
      - package-ecosystem: github-actions
        directory: /
        schedule:
          interval: weekly
      - package-ecosystem: pip
        directory: /
        schedule:
          interval: weekly
      - package-ecosystem: docker
        directory: /
        schedule:
          interval: weekly
    ```

    Dependabot apre una pull request per ogni aggiornamento: la pipeline la verifica, tu decidi se unirla.

## Fonti

- [GitHub Docs: Configuring Dependabot version updates](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/configuring-dependabot-version-updates)
- [pytest-cov](https://pytest-cov.readthedocs.io/)
- [Google SRE Workbook: Implementing SLOs](https://sre.google/workbook/implementing-slos/)
- [Google SRE Book: Postmortem culture](https://sre.google/sre-book/postmortem-culture/)
- [Prometheus: Recording rules](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)
- [Alertmanager: Configuration](https://prometheus.io/docs/alerting/latest/configuration/)

# Fase 2 · Automazione: Git, Bash e Python

!!! obiettivi "Al termine di questa fase saprai"
    - versionare il tuo lavoro con Git e collaborare su GitHub con branch e pull request;
    - scrivere script Bash robusti, con gestione degli errori e degli argomenti;
    - pianificare attività con cron e con i timer di systemd;
    - programmare in Python: tipi di dati, controllo del flusso, funzioni, moduli e ambienti virtuali;
    - automatizzare con Python: file, JSON e YAML, chiamate HTTP, riga di comando;
    - scrivere e testare una piccola API web con Flask.

| | |
|---|---|
| **Durata indicativa** | 1,5 mesi, circa 55 ore |
| **Prerequisiti** | [Fase 1](../fase-1/index.md) |
| **Ambiente** | Dev Container (lezioni [2.1](01-git-base.md)-[2.4](04-bash-avanzato.md) e [2.6](06-python-base.md)-[2.8](08-flask.md)); VM `lab` ([lezione 2.5](05-pianificazione.md)) |
| **Progetto** | [Toolkit di automazione](progetto.md) |

## Perché automatizzare

Nella [Fase 1](../fase-1/index.md) hai configurato un server a mano, comando dopo comando. Funziona per un server; non funziona per
cento, né per rifare lo stesso lavoro ogni settimana senza errori. Il cuore del lavoro DevOps è **trasformare le
operazioni manuali in codice**: codice che si può rileggere, versionare, testare e rieseguire.

Questa fase ti dà i tre strumenti di base:

- **Git**, per tenere traccia di ogni modifica e collaborare;
- **Bash**, per automatizzare il lavoro sul sistema;
- **Python**, per l'automazione più complessa e per le applicazioni.

## L'applicazione del corso

Nella [lezione 2.8](08-flask.md) scriverai la **Tasks API**, una piccola API REST in Flask. È l'applicazione che porterai nei
container ([Fase 3](../fase-3/index.md)), nel cloud ([Fase 5](../fase-5/index.md)), in una pipeline automatica ([Fase 7](../fase-7/index.md)) e su Kubernetes ([Fase 8](../fase-8/index.md)). La versione
di riferimento si trova nel repository del corso, nella cartella
[`labs/tasks-api`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/tasks-api):
consultala **dopo** aver scritto la tua.

## Le lezioni

1. [Git: i fondamenti](01-git-base.md)
2. [Git: branch e collaborazione](02-git-collaborazione.md)
3. [Bash: variabili e condizioni](03-bash-base.md)
4. [Bash: cicli, funzioni e robustezza](04-bash-avanzato.md)
5. [Pianificare i job](05-pianificazione.md)
6. [Python: le basi](06-python-base.md)
7. [Python per l'automazione](07-python-automazione.md)
8. [API HTTP con Flask](08-flask.md)

Poi il [progetto](progetto.md) e la [verifica di fine fase](verifica.md).

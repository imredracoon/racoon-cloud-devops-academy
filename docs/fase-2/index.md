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

## Nei tre livelli

Il codice si scrive sempre nel Dev Container della workstation e si salva su GitHub. Cambia dove lo fai girare.

=== "Livello 1 · Un PC"

    Dev Container per Git, Bash e Python; VM `lab` per i job pianificati della [lezione 2.5](05-pianificazione.md) e VM `web` per
    l'installazione del progetto.

=== "Livello 2 · Workstation + server"

    Gli script li provi nel Dev Container e li **installi su `srv1`**, che diventa il tuo primo server "di
    produzione": il backup notturno gira lì e salva sul disco USB della [lezione 1.6](../fase-1/06-storage.md#nei-tre-livelli),
    la Tasks API del progetto è servita da `srv1`. In più scrivi `fleet.sh`, uno script che lancia un comando su
    `srv1` e sulle sue VM ([lezione 2.4](04-bash-avanzato.md#nei-tre-livelli)).

=== "Livello 3 · Mini data center"

    Gli script agiscono su **tutti i nodi**: `fleet.sh` controlla il data center in un colpo solo, il backup
    notturno di `nodo2` e `nodo3` finisce sullo storage NFS di `nodo1`, la Tasks API gira su due nodi e
    `healthcheck.py` li controlla entrambi. È l'automazione "fatta a mano" che nella [Fase 6](../fase-6/index.md)
    sostituirai con Ansible.

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

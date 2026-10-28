# Fase 3 · Container

!!! obiettivi "Al termine di questa fase saprai"
    - spiegare cosa sono i container, in cosa differiscono dalle VM e quali meccanismi di Linux li rendono possibili;
    - usare Docker e Podman per scaricare immagini, avviare, ispezionare e fermare container;
    - scrivere Dockerfile efficienti e sicuri, con build multi-stage e utenti senza privilegi;
    - gestire dati persistenti con i volumi e far comunicare i container in reti dedicate;
    - descrivere un'applicazione a più servizi con Docker Compose;
    - pubblicare immagini in un registry e verificarne le vulnerabilità.

| | |
|---|---|
| **Durata indicativa** | 1 mese, circa 40 ore |
| **Prerequisiti** | Fasi 1 e 2 |
| **Ambiente** | Dev Container (Docker è già disponibile) oppure VM `lab` / `rhel-lab` |
| **Progetto** | [Stack containerizzato](progetto.md) |

## Perché i container

Nella Fase 2 hai installato la Tasks API su un server: ambiente virtuale, dipendenze, servizio systemd, utente
dedicato. Funziona, ma è fragile: su un altro server con un'altra versione di Python o di una libreria potrebbe non
funzionare più. Il classico "sul mio computer funziona".

Un **container** impacchetta l'applicazione **con tutto ciò che le serve** (interprete, librerie, file di
configurazione) in un'**immagine** che gira allo stesso modo sul tuo portatile, nella pipeline di test, nel cloud e
in Kubernetes. Oggi è il formato standard con cui si distribuisce il software lato server: tutto il resto del corso
lo dà per acquisito.

## Docker e Podman

=== "Debian/Ubuntu"

    Nel percorso Debian/Ubuntu usiamo **Docker Engine**, lo strumento più diffuso. Nel Dev Container è già
    disponibile; per installarlo sulla VM `lab` segui la lezione 3.2.

=== "RHEL (Rocky/AlmaLinux)"

    Red Hat include nella distribuzione **Podman**, compatibile con i comandi di Docker ma senza un demone centrale e
    pensato per girare senza privilegi di root. Nella lezione 3.2 lo installi su `rhel-lab`. Quasi tutti i comandi
    di questa fase funzionano sostituendo `docker` con `podman`; le differenze sono segnalate.

## Le lezioni

1. [Cosa sono i container](01-concetti.md)
2. [Immagini e container](02-immagini-container.md)
3. [Dockerfile](03-dockerfile.md)
4. [Volumi e reti](04-volumi-reti.md)
5. [Docker Compose](05-compose.md)
6. [Registry e sicurezza](06-registry-sicurezza.md)

Poi il [progetto](progetto.md) e la [verifica di fine fase](verifica.md).

# Fase 3 · Container

!!! obiettivi "Al termine di questa fase saprai"
    - spiegare cosa sono i container, in cosa differiscono dalle VM e quali meccanismi di Linux li rendono possibili;
    - usare Docker e Podman per scaricare immagini, avviare, ispezionare e fermare container;
    - scrivere Dockerfile efficienti e sicuri, con build multi-stage e utenti senza privilegi;
    - gestire dati persistenti con i volumi e far comunicare i container in reti dedicate;
    - descrivere un'applicazione a più servizi con Docker Compose;
    - pubblicare immagini in un registry e verificarne le vulnerabilità;
    - usare Podman rootless, i pod e i Quadlet per far gestire i container a systemd.

| | |
|---|---|
| **Durata indicativa** | 1 mese, circa 45 ore |
| **Prerequisiti** | Fasi 1 e 2 |
| **Ambiente** | Dev Container (Docker è già disponibile) oppure VM `lab` / `rhel-lab` |
| **Progetto** | [Stack containerizzato](progetto.md) |

## Perché i container

Nella [Fase 2](../fase-2/index.md) hai installato la Tasks API su un server: ambiente virtuale, dipendenze, servizio systemd, utente
dedicato. Funziona, ma è fragile: su un altro server con un'altra versione di Python o di una libreria potrebbe non
funzionare più. Il classico "sul mio computer funziona".

Un **container** impacchetta l'applicazione **con tutto ciò che le serve** (interprete, librerie, file di
configurazione) in un'**immagine** che gira allo stesso modo sul tuo portatile, nella pipeline di test, nel cloud e
in Kubernetes. Oggi è il formato standard con cui si distribuisce il software lato server: tutto il resto del corso
lo dà per acquisito.

## Docker e Podman

=== "Debian/Ubuntu"

    Nel percorso Debian/Ubuntu usiamo **Docker Engine**, lo strumento più diffuso. Nel Dev Container è già
    disponibile; per installarlo sulla VM `lab` segui la [lezione 3.2](02-immagini-container.md).

=== "RHEL (Rocky/AlmaLinux)"

    Red Hat include nella distribuzione **Podman**, compatibile con i comandi di Docker ma senza un demone centrale e
    pensato per girare senza privilegi di root. Nella [lezione 3.2](02-immagini-container.md) lo installi su `rhel-lab`. Quasi tutti i comandi
    di questa fase funzionano sostituendo `docker` con `podman`; le differenze sono segnalate.

## Nei tre livelli

Le lezioni 3.1-3.6 si fanno nel Dev Container in tutti i livelli: le immagini si costruiscono lì e si pubblicano
sul registry. Cambia dove fai girare lo stack.

=== "Livello 1 · Un PC"

    Dev Container per le lezioni, VM `web` (Docker) o `rhel-lab` (Podman) per il progetto.

=== "Livello 2 · Workstation + server"

    Installa Docker su `srv1` con Ubuntu (come nella VM, [lezione 3.2](02-immagini-container.md)) oppure usa Podman su `srv1` con
    AlmaLinux. Lo stack del progetto gira su `srv1`, dietro il Nginx dell'host che hai configurato nei progetti 1 e 2.
    La [lezione 3.7](07-podman.md#nei-tre-livelli) la provi anche sul server, con un Quadlet che parte all'avvio.

=== "Livello 3 · Mini data center"

    I due percorsi diventano due motori sullo stesso data center: **`nodo2` con Docker e Compose**, **`nodo3` con
    Podman e i Quadlet**. Entrambi scaricano la **stessa immagine** dal registry: è la prova più concreta che un
    container è portabile.

!!! warning "Docker, firewall e VM sulla stessa macchina"
    Su `srv1`, `nodo2` e `nodo3` tieni presenti due cose:

    - **le porte pubblicate da Docker scavalcano il firewall dell'host**: il traffico verso un container passa per
      la catena `forward`, non per `input`, quindi la tua politica `drop` non lo vede. Pubblica le porte solo su
      `127.0.0.1` (`"127.0.0.1:8080:80"`) e metti davanti il Nginx dell'host, che è l'unico a esporsi;
    - Docker crea regole proprie, come libvirt: un `/etc/nftables.conf` con `flush ruleset` le cancellerebbe al
      primo `systemctl restart nftables`. Usa la tabella dedicata `homelab` vista nella
      [lezione 1.11](../fase-1/11-firewall.md#nei-tre-livelli) anche su `nodo2`.

## Le lezioni

1. [Cosa sono i container](01-concetti.md)
2. [Immagini e container](02-immagini-container.md)
3. [Dockerfile](03-dockerfile.md)
4. [Volumi e reti](04-volumi-reti.md)
5. [Docker Compose](05-compose.md)
6. [Registry e sicurezza](06-registry-sicurezza.md)
7. [Podman in profondità](07-podman.md)

Poi il [progetto](progetto.md) e la [verifica di fine fase](verifica.md).

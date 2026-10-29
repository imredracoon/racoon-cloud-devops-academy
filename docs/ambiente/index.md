# L'ambiente di laboratorio

Per imparare questo mestiere bisogna **fare**: scrivere comandi, rompere configurazioni, ripararle. Per farlo
senza rischi e nello stesso modo per tutti, il corso usa tre ambienti di laboratorio.

## Il tuo computer: Ubuntu o Fedora

Il corso dà per scontato che sul tuo computer ci sia **Linux**, in una di queste due distribuzioni desktop:

| Distribuzione | Versione | Famiglia | Perché sceglierla |
|---|---|---|---|
| **Ubuntu Desktop** | 24.04 LTS o 26.04 LTS | Debian | la più diffusa, tantissima documentazione; stessa famiglia della VM Ubuntu |
| **Fedora Workstation** | l'ultima versione stabile | Red Hat | è la distribuzione "a monte" di RHEL: ciò che arriva in RHEL nasce in Fedora |

Le istruzioni di installazione sono in **schede affiancate** con il nome della distribuzione, come le schede
Debian/Ubuntu e RHEL della [Fase 1](../fase-1/index.md). Lavorare ogni giorno su Linux è già allenamento: i comandi che usi sul
tuo computer sono gli stessi che userai sui server.

!!! tip "Non hai ancora Linux?"
    Nessun problema: la pagina [Non hai ancora Linux?](linux-sul-pc.md) spiega come installare Ubuntu o Fedora,
    oppure come seguire il corso dal browser con GitHub Codespaces e Killercoda.

!!! info "Uno, due o più computer"
    Il corso si fa con un solo computer, ma puoi anche costruire un piccolo laboratorio fisico: una workstation più
    un server, oppure un mini data center con tre mini PC usati. La pagina [I tre livelli del laboratorio](livelli.md)
    spiega le tre strade, l'hardware e i costi.

## I tre ambienti di laboratorio

| Ambiente | Cos'è | Dove si usa | Costo |
|---|---|---|---|
| [**Dev Container**](devcontainer.md) | Un "contenitore" con tutti gli strumenti già installati: Git, Python, Docker, Terraform, kubectl, AWS CLI, gcloud e altri. | Fasi 2-8 | Gratis |
| [**Macchina virtuale Linux**](vm-linux.md) | Un computer Linux simulato dentro il tuo, completo di tutto: servizi, dischi, firewall. | Fasi 0-1 e appendice LFCS | Gratis |
| [**Account cloud**](account-cloud.md) | Il tuo account personale su AWS e su Google Cloud. | Fasi 4-8 | Gratis con le offerte per i nuovi clienti, se segui le regole di sicurezza |

!!! info "Non preoccuparti dei termini"
    Container, macchina virtuale e cloud sono spiegati con calma nella [Fase 0](../fase-0/index.md) e nelle fasi successive. Per ora ti
    basta sapere che sono "spazi di lavoro" separati dal tuo computer: qualunque errore tu faccia lì dentro, il tuo
    computer resta al sicuro.

## Quando preparare cosa

Non serve preparare tutto adesso. Il corso ti dice quando è il momento:

1. **Prima di iniziare**: hai Ubuntu o Fedora sul tuo computer ([Non hai ancora Linux?](linux-sul-pc.md)).
2. **[Fase 0](../fase-0/index.md), [lezione 0.7](../fase-0/07-laboratorio.md)**: installi la macchina virtuale Linux e il Dev Container.
3. **[Fase 4](../fase-4/index.md), [lezione 4.3](../fase-4/03-account.md)**: crei gli account cloud, con budget e allarmi di spesa.

## Requisiti del computer

| Componente | Minimo | Consigliato |
|---|---|---|
| Memoria (RAM) | 8 GB | 16 GB o più |
| Spazio libero su disco | 40 GB | 80 GB |
| Processore | 4 core, con virtualizzazione attiva | 8 core |
| Sistema operativo | Ubuntu Desktop 24.04 LTS o Fedora Workstation | Ubuntu 26.04 LTS o l'ultima Fedora |

!!! tip "Il computer non basta?"
    Puoi svolgere le fasi 2-8 interamente dal browser con **GitHub Codespaces**, usando la quota gratuita mensile
    di GitHub, e la [Fase 1](../fase-1/index.md) con i **playground Linux gratuiti** di Killercoda. Le istruzioni sono nelle pagine
    [Dev Container](devcontainer.md) e [Macchina virtuale Linux](vm-linux.md).

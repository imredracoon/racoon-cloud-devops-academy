# L'ambiente di laboratorio

Per imparare questo mestiere bisogna **fare**: scrivere comandi, rompere configurazioni, ripararle. Per farlo
senza rischi e nello stesso modo su qualsiasi computer, il corso usa tre ambienti di laboratorio.

| Ambiente | Cos'è | Dove si usa | Costo |
|---|---|---|---|
| [**Dev Container**](devcontainer.md) | Un "contenitore" con tutti gli strumenti già installati: Git, Python, Docker, Terraform, kubectl, AWS CLI, gcloud e altri. | Fasi 2-8 | Gratis |
| [**Macchina virtuale Linux**](vm-linux.md) | Un computer Linux simulato dentro il tuo, completo di tutto: servizi, dischi, firewall. | Fasi 0-1 e appendice LFCS | Gratis |
| [**Account cloud**](account-cloud.md) | Il tuo account personale su AWS e su Google Cloud. | Fasi 4-8 | Gratis con le offerte per i nuovi clienti, se segui le regole di sicurezza |

!!! info "Non preoccuparti dei termini"
    Container, macchina virtuale e cloud sono spiegati con calma nella Fase 0 e nelle fasi successive. Per ora ti
    basta sapere che sono "spazi di lavoro" separati dal tuo computer: qualunque errore tu faccia lì dentro, il tuo
    computer resta al sicuro.

## Quando preparare cosa

Non serve preparare tutto adesso. Il corso ti dice quando è il momento:

1. **Fase 0, lezione 0.7**: installi la macchina virtuale Linux e il Dev Container.
2. **Fase 4, lezione 4.3**: crei gli account cloud, con budget e allarmi di spesa.

## Requisiti del computer

| Componente | Minimo | Consigliato |
|---|---|---|
| Memoria (RAM) | 8 GB | 16 GB o più |
| Spazio libero su disco | 40 GB | 80 GB |
| Processore | 4 core, con virtualizzazione attiva | 8 core |
| Sistema operativo | Windows 10/11, macOS 14 o successivo, una distribuzione Linux recente | — |

!!! tip "Il computer non basta?"
    Puoi svolgere le fasi 2-8 interamente dal browser con **GitHub Codespaces**, usando la quota gratuita mensile
    di GitHub, e la Fase 1 con i **playground Linux gratuiti** di Killercoda. Le istruzioni sono nelle pagine
    [Dev Container](devcontainer.md) e [Macchina virtuale Linux](vm-linux.md).

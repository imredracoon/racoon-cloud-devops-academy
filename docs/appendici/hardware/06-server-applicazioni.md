# 6. Server per calcolo e applicazioni

Questi server **eseguono il lavoro**: le applicazioni degli utenti, le macchine virtuali, i container, i modelli di
intelligenza artificiale, le pipeline. Sono la parte più grande di ogni data center.

!!! abstract "Già visto nel corso"
    - **Nginx**, reverse proxy e HTTPS: [lezione 1.12](../../fase-1/12-http-nginx.md); un'applicazione Flask in produzione con Gunicorn: [lezione 2.8](../../fase-2/08-flask.md#in-produzione-gunicorn).
    - **Container e macchine virtuali** a confronto: [lezione 3.1](../../fase-3/01-concetti.md#container-e-macchine-virtuali); installare Docker o Podman: [lezione 3.2](../../fase-3/02-immagini-container.md#installazione).
    - **Registry** di immagini: [lezione 3.6](../../fase-3/06-registry-sicurezza.md#i-registry).
    - **Kubernetes** e cluster con kubeadm su macchine proprie: [Fase 8](../../fase-8/index.md) e [lezione 8.12](../../fase-8/12-kubeadm.md);
      riservare i nodi GPU: [lezione 8.8](../../fase-8/08-scheduling.md#scegliere-i-nodi).
    - **Runner** delle pipeline: [lezione 7.2](../../fase-7/02-github-actions.md#i-concetti).
    - **Metriche e log** centralizzati: [Prometheus (7.6)](../../fase-7/06-prometheus.md), [Grafana (7.7)](../../fase-7/07-grafana.md), [Loki (7.8)](../../fase-7/08-loki.md).

## Il quadro d'insieme

| Server | Cosa fa | Software su Ubuntu | Software su RHEL |
|---|---|---|---|
| [Web e applicazioni](#web-e-applicazioni) | pubblica siti e API | Nginx, Apache, Gunicorn, Tomcat | Nginx, Apache (httpd), JBoss EAP |
| [Virtualizzazione](#virtualizzazione) | ospita macchine virtuali | KVM, libvirt, LXD | KVM, libvirt, Cockpit, OpenShift Virtualization |
| [Container e Kubernetes](#container-e-kubernetes) | esegue container e cluster | Docker, Canonical Kubernetes | Podman, OpenShift |
| [GPU e AI](#gpu-e-intelligenza-artificiale) | addestra ed esegue modelli | driver NVIDIA, CUDA | driver NVIDIA, CUDA |
| [HPC](#hpc) | calcolo scientifico in parallelo | Slurm, Open MPI | Slurm, Open MPI |
| [CI/CD](#cicd-e-git) | compila, testa e rilascia | runner GitHub o GitLab, Jenkins | runner GitHub o GitLab, Jenkins |
| [Git e registry](#cicd-e-git) | ospita codice e immagini | GitLab, Forgejo, Harbor | GitLab, Forgejo, Quay |
| [Mirror dei pacchetti](#mirror-dei-pacchetti) | copia locale dei repository | Landscape, debmirror | Satellite, reposync |
| [Monitoraggio e log](#monitoraggio-e-log) | raccoglie metriche e log | Prometheus, Grafana, Loki, rsyslog | PCP, Prometheus, Grafana, rsyslog |
| [Media e giochi](#media-e-giochi) | streaming e server di gioco | Jellyfin, server di gioco | in container |

## Web e applicazioni

Il **server web** risponde alle richieste HTTP e restituisce pagine o inoltra le richieste a un **application
server** che esegue il codice (Gunicorn per Python, Tomcat o WildFly per Java, Node.js). Hai costruito questa
catena nella [lezione 1.12](../../fase-1/12-http-nginx.md#reverse-proxy) e nella [lezione 2.8](../../fase-2/08-flask.md#in-produzione-gunicorn).

**Hardware tipico**: server 1U con molti core e RAM media, SSD, 10-25 Gbit/s. Si scala **orizzontalmente**:
più server uguali dietro un load balancer ([pagina 1](01-dispositivi-rete.md#load-balancer)).

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y nginx          # oppure apache2
    sudo apt install -y tomcat10       # application server Java, se serve
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y nginx          # oppure httpd (Apache)
    sudo systemctl enable --now nginx
    sudo firewall-cmd --permanent --add-service={http,https} && sudo firewall-cmd --reload
    # per Java Red Hat propone JBoss EAP (prodotto), la cui versione della comunità è WildFly
    ```

## Virtualizzazione

Un **hypervisor** divide un server fisico in molte macchine virtuali. È esattamente ciò che fanno i provider cloud
con le VM che hai creato nella [Fase 4](../../fase-4/05-compute.md#le-macchine-virtuali). Su Linux l'hypervisor è
**KVM**, integrato nel kernel, con **QEMU** che emula l'hardware e **libvirt** che offre comandi e API comuni
(`virsh`, `virt-install`). Più host di virtualizzazione con uno storage condiviso permettono di spostare le VM da
un server all'altro senza spegnerle (*live migration*).

**Hardware tipico**: server 2U con **molti core e molta RAM** (si sommano le risorse di tutte le VM), rete a
25-100 Gbit/s, dischi NVMe locali o storage condiviso ([pagina 5](05-server-dati.md#storage-a-blocchi-e-san)).
Il processore deve avere le estensioni di virtualizzazione (Intel VT-x o AMD-V) attive nel firmware.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y qemu-kvm libvirt-daemon-system virtinst
    sudo usermod -aG libvirt "$USER"
    virsh list --all
    # LXD: VM e container di sistema con un solo strumento
    sudo snap install lxd && sudo lxd init --minimal
    lxc launch ubuntu:24.04 prova --vm
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y qemu-kvm libvirt virt-install
    sudo systemctl enable --now virtqemud.socket
    sudo virt-host-validate        # controlla che l'host sia pronto
    # interfaccia web per le VM:
    sudo dnf install -y cockpit cockpit-machines && sudo systemctl enable --now cockpit.socket
    ```

    Per centinaia di VM Red Hat propone **OpenShift Virtualization**, che gestisce le VM con Kubernetes.

Altri hypervisor che incontrerai: **VMware vSphere/ESXi** (commerciale, molto diffuso nelle aziende),
**Microsoft Hyper-V**, **Proxmox VE** (open source, basato su Debian, con KVM e container LXC).

## Container e Kubernetes

Un server può eseguire container direttamente (Docker o Podman, [Fase 3](../../fase-3/index.md)) o essere un
**nodo** di un cluster Kubernetes ([Fase 8](../../fase-8/index.md)). In un data center i nodi Kubernetes sono
spesso **bare metal**, cioè server fisici senza hypervisor, per non perdere prestazioni.

**Hardware tipico**: come per la virtualizzazione; tre server piccoli per il control plane e molti worker.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y docker.io      # Docker dai repository di Ubuntu (lezione 3.2)
    # Canonical Kubernetes, distribuito come snap:
    sudo snap install k8s --classic
    sudo k8s bootstrap
    sudo k8s status
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y podman         # il motore di container di RHEL (lezione 3.2)
    ```

    Il Kubernetes di Red Hat è **OpenShift**, un prodotto con il suo installatore. Per imparare si può costruire
    un cluster con kubeadm, come nella [lezione 8.12](../../fase-8/12-kubeadm.md).

## GPU e intelligenza artificiale

I **server GPU** addestrano modelli di intelligenza artificiale (*training*) e li usano per rispondere
(*inferenza*). Nel cloud corrispondono alle famiglie di istanze accelerate della [lezione 4.5](../../fase-4/05-compute.md#le-macchine-virtuali).

**Hardware tipico**:

- **inferenza** e lavori piccoli: server 2U con una o più schede PCIe (per esempio NVIDIA L4 o L40S);
- **training**: server con otto GPU SXM collegate da NVLink, come il Dell XE9780
  ([pagina 3](03-hardware-server.md#modelli-reali-a-confronto)), centinaia di GB di RAM, SSD NVMe veloci e
  reti InfiniBand o Ethernet a 400 Gbit/s per collegare molti server in un unico cluster;
- **energia e raffreddamento** dimensionati apposta: un server con otto GPU può assorbire oltre 10 kW, quanto
  un intero armadio di server normali ([pagina 9](09-data-center.md#raffreddamento)).

Il software si costruisce a strati: **driver** della GPU, librerie **CUDA**, il **container toolkit** per usare la
GPU dentro i container, e in Kubernetes il **GPU Operator** che installa tutto sui nodi.

=== "Debian/Ubuntu"

    ```bash
    sudo ubuntu-drivers list --gpgpu          # driver adatti ai server
    sudo ubuntu-drivers install --gpgpu
    sudo reboot
    nvidia-smi                                # GPU, memoria, temperatura, processi
    ```

=== "RHEL (Rocky/AlmaLinux)"

    Su RHEL il driver si installa dal **repository CUDA di NVIDIA** per RHEL, seguendo la guida ufficiale di
    installazione di NVIDIA per la tua versione: aggiungi il repository con `dnf config-manager`, installa il
    pacchetto del driver e riavvia. Con il Secure Boot attivo servono moduli firmati.

    ```bash
    nvidia-smi                                # verifica, uguale su entrambe le famiglie
    ```

Per i container, su entrambe le famiglie, si installa **NVIDIA Container Toolkit** e si avvia il container con
accesso alla GPU (`docker run --gpus all ...` oppure, con Podman, `--device nvidia.com/gpu=all`).

## HPC

L'**HPC** (*High Performance Computing*) divide un grande calcolo scientifico (meteo, fluidodinamica, genomica)
tra centinaia o migliaia di server che lavorano insieme come un unico supercomputer. Gli ingredienti:

- un **gestore di code** che assegna i lavori ai nodi: **Slurm** è lo standard;
- **MPI** (Open MPI) per far comunicare i processi su server diversi;
- una rete a bassissima latenza (**InfiniBand**) e un **filesystem parallelo** (Lustre, IBM Storage Scale) che
  tutti i nodi leggono contemporaneamente.

Ubuntu ha Slurm nei repository (`slurm-wlm`); su RHEL si compila dai sorgenti del progetto o si usano i pacchetti
dei repository della comunità.

## CI/CD e Git

Le aziende che non usano servizi in cloud ospitano in casa il codice e le pipeline:

- **server Git**: **GitLab** (pacchetti ufficiali per Ubuntu e per la famiglia RHEL), **Forgejo** o **Gitea**,
  più leggeri;
- **runner**: le macchine che eseguono i job ([lezione 7.2](../../fase-7/02-github-actions.md#i-concetti)). Un
  *self-hosted runner* di GitHub Actions o un GitLab Runner si installa su qualsiasi server Linux; **Jenkins** è lo
  strumento storico;
- **registry di immagini** ([lezione 3.6](../../fase-3/06-registry-sicurezza.md#i-registry)): **Harbor** o il
  registry integrato in GitLab; nel mondo Red Hat **Quay**.

**Hardware tipico**: i runner vogliono molti core e SSD veloci per compilare; il server Git vuole RAM e backup
affidabili.

## Mirror dei pacchetti

Con decine di server conviene una **copia locale dei repository**: gli aggiornamenti si scaricano una volta sola,
i server senza accesso a Internet possono aggiornarsi e si decide **quando** una nuova versione arriva in
produzione (prima ai server di test, poi agli altri).

=== "Debian/Ubuntu"

    **Landscape** di Canonical gestisce mirror, aggiornamenti e inventario di molti server Ubuntu. Per un
    mirror semplice esiste `debmirror`.

=== "RHEL (Rocky/AlmaLinux)"

    **Satellite** gestisce i contenuti per migliaia di sistemi. Senza Satellite, `dnf reposync` scarica un
    repository e `createrepo_c` crea i metadati, poi un server web lo pubblica.

    ```bash
    sudo dnf install -y dnf-plugins-core createrepo_c
    sudo dnf reposync --repoid=rhel-10-for-x86_64-baseos-rpms -p /srv/mirror --download-metadata
    ```

## Monitoraggio e log

Raccolgono metriche e log da tutti i server e avvisano quando qualcosa va storto. Hai costruito questa pila nella
[Fase 7](../../fase-7/index.md): **Prometheus** e i suoi *exporter* per le metriche, **Grafana** per i grafici,
**Loki** per i log. Sui server fisici si aggiungono le metriche dell'hardware (temperature, ventole, dischi, BMC),
vedi la [pagina 8](08-manutenzione.md#la-salute-dellhardware).

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y prometheus-node-exporter    # metriche del sistema per Prometheus
    # log centralizzati con rsyslog: sui client, in /etc/rsyslog.d/90-remoto.conf
    #   *.* @@log.esempio.lan:514
    ```

=== "RHEL (Rocky/AlmaLinux)"

    RHEL include **PCP** (*Performance Co-Pilot*), che raccoglie e archivia le metriche e le esporta anche verso
    Grafana; la console web Cockpit le mostra in grafici.

    ```bash
    sudo dnf install -y pcp-zeroconf     # PCP con una configurazione pronta
    pmstat                               # riepilogo in tempo reale
    ```

    Il log centralizzato si fa con **rsyslog**, come su Ubuntu.

## Media e giochi

Anche lo svago gira su server:

- **streaming multimediale**: **Jellyfin** (open source) o Plex organizzano film e musica e li convertono al
  volo (*transcoding*) per ogni dispositivo. Una GPU o la grafica integrata del processore accelerano molto il
  transcoding. Jellyfin ha repository ufficiali per Ubuntu; su RHEL si usa in container;
- **server di gioco** (Minecraft, Valheim e simili): preferiscono **pochi core molto veloci** a molti core lenti,
  perché il ciclo di gioco gira su un thread solo, e una rete con bassa latenza verso i giocatori.

## Esercizi

!!! esercizio "Esercizio 1"
    Un'azienda vuole ospitare 60 VM con in media 4 vCPU e 16 GB di RAM ciascuna, con la possibilità di spegnere un
    host per manutenzione senza fermare le VM. Quanti host da 64 core e 1 TB di RAM servono, e cosa serve oltre
    agli host?

??? success "Soluzione"
    In totale 240 vCPU e 960 GB di RAM. Con un rapporto prudente di 2 vCPU per core bastano 120 core, e la RAM
    (che non si condivide) richiede almeno 1 TB: **due host** reggono il carico, ma per spegnerne uno ne servono
    **tre** (regola N+1). Servono anche uno **storage condiviso** (o replicato) per la live migration e una rete a
    25 Gbit/s o più.

!!! esercizio "Esercizio 2"
    Perché un cluster per addestrare modelli di AI ha bisogno di una rete InfiniBand o Ethernet a 400 Gbit/s,
    mentre un cluster di server web no?

??? success "Soluzione"
    Nel training distribuito le GPU di server diversi si scambiano continuamente grandi quantità di dati (i
    gradienti) e la più lenta rallenta tutte. I server web invece lavorano ognuno per conto suo sulle richieste
    degli utenti e comunicano poco tra loro.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è l'hypervisor integrato nel kernel Linux?","options":["VMware ESXi","KVM","Hyper-V","VirtualBox"],"answer":[1],"explanation":"KVM è un modulo del kernel; QEMU e libvirt lo completano."},
 {"q":"Quale comando installa su Ubuntu i driver NVIDIA adatti a un server di calcolo?","options":["apt install cuda-server","ubuntu-drivers install --gpgpu","nvidia-smi --install","snap install nvidia"],"answer":[1],"explanation":"L'opzione --gpgpu sceglie i driver per il calcolo invece di quelli per il desktop."},
 {"q":"Quali sono ingredienti tipici di un cluster HPC?","options":["Un gestore di code come Slurm","MPI","Rete a bassa latenza come InfiniBand","Un filesystem parallelo","Un proxy Squid"],"answer":[0,1,2,3],"explanation":"Squid è un proxy web e non c'entra con il calcolo parallelo."},
 {"q":"Cosa offre RHEL per raccogliere e archiviare le metriche di sistema?","options":["PCP (Performance Co-Pilot)","Landscape","MAAS","LXD"],"answer":[0],"explanation":"Landscape, MAAS e LXD sono strumenti di Canonical."},
 {"q":"Perché un server di gioco preferisce pochi core veloci?","options":["Perché i giochi usano solo la GPU","Perché il ciclo di gioco gira principalmente su un solo thread","Perché i core lenti consumano di più","Non è vero, servono molti core"],"answer":[1],"explanation":"La frequenza del singolo core conta più del numero di core."}
]}
</script>

## Fonti

- [Ubuntu Server: Libvirt](https://documentation.ubuntu.com/server/how-to/virtualisation/libvirt/)
- [Ubuntu Server: NVIDIA drivers installation](https://documentation.ubuntu.com/server/how-to/graphics/install-nvidia-drivers/)
- [Ubuntu Server: Docker for system admins](https://ubuntu.com/server/docs/how-to/containers/docker-for-system-admins/)
- [Canonical Kubernetes](https://documentation.ubuntu.com/canonical-kubernetes/)
- [RHEL 10: Configuring and managing Linux virtual machines](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html-single/configuring_and_managing_linux_virtual_machines/index)
- [RHEL 10: Building, running, and managing containers](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html-single/building_running_and_managing_containers/index)
- [RHEL 10: Deploying web servers and reverse proxies](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/deploying_web_servers_and_reverse_proxies/index)
- [RHEL 10: Monitoring and managing system status and performance (PCP)](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html-single/monitoring_and_managing_system_status_and_performance/index)
- [NVIDIA CUDA Installation Guide for Linux](https://docs.nvidia.com/cuda/cuda-installation-guide-linux/)
- [NVIDIA Container Toolkit, installazione](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/install-guide.html)
- [Slurm Workload Manager](https://slurm.schedmd.com/)
- [Jellyfin documentation](https://jellyfin.org/docs/)

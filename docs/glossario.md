# Glossario

I termini principali del corso, in ordine alfabetico. Ogni voce rimanda alla lezione in cui il concetto è spiegato
per esteso. I nomi dei prodotti restano in inglese, come nella documentazione ufficiale.

## A

ACL (Access Control List)
:   Elenco di permessi aggiuntivi su file e directory, oltre a proprietario, gruppo e altri. Nel cloud indica anche le regole di accesso di rete delle subnet (network ACL). [Lezione](fase-1/03-permessi.md)

Alta disponibilità
:   Capacità di un servizio di restare utilizzabile anche se un componente si guasta, di solito duplicandolo in più zone. [Lezione](fase-5/02-ha-scaling.md)

Ansible
:   Strumento di automazione che configura i server via SSH a partire da playbook YAML, senza agenti installati. [Lezione](fase-6/10-ansible.md)

API
:   Interfaccia con cui un programma offre funzioni ad altri programmi; nel web di solito tramite richieste HTTP. [Lezione](fase-2/08-flask.md)

Autoscaling
:   Aumento o riduzione automatica delle risorse (istanze, Pod) in base al carico. [Lezione](fase-5/02-ha-scaling.md)

Availability Zone (zona)
:   Uno o più data center isolati all'interno di una regione cloud. [Lezione](fase-4/02-infrastruttura-globale.md)

## B

Backend (Terraform)
:   Luogo in cui Terraform conserva lo state: file locale o servizio remoto con lock. [Lezione](fase-6/06-state.md)

Bilanciatore di carico
:   Componente che distribuisce le richieste tra più server e smette di inviarle a quelli non sani. [Lezione](fase-5/02-ha-scaling.md)

Bucket
:   Contenitore di oggetti nello storage a oggetti (Amazon S3, Cloud Storage). [Lezione](fase-4/06-storage.md)

## C

CI/CD
:   Integrazione continua (ogni modifica viene compilata e testata in automatico) e rilascio continuo (le modifiche approvate arrivano in produzione in automatico). [Lezione](fase-7/01-concetti.md)

CIDR
:   Notazione per indicare un intervallo di indirizzi IP, per esempio 10.0.0.0/16. [Lezione](fase-1/07-tcpip.md)

Cloud-init
:   Strumento che configura una macchina al primo avvio (utenti, pacchetti, file) a partire da dati forniti dal cloud. [Lezione](fase-4/05-compute.md)

Cluster
:   Insieme di macchine che lavorano come un unico sistema, per esempio i nodi di Kubernetes. [Lezione](fase-8/01-architettura.md)

CNI, CSI, CRI
:   Le interfacce standard di Kubernetes per la rete dei Pod, lo storage e il runtime dei container. [Lezione](fase-8/01-architettura.md)

Commit
:   Fotografia delle modifiche registrata nella storia di un repository Git. [Lezione](fase-2/01-git-base.md)

ConfigMap e Secret
:   Oggetti Kubernetes che forniscono configurazione e dati sensibili ai Pod. [Lezione](fase-8/06-config.md)

Container
:   Processo isolato che porta con sé l'applicazione e le sue dipendenze, condividendo il kernel dell'host. [Lezione](fase-3/01-concetti.md)

Control plane
:   Componenti che gestiscono un cluster Kubernetes: API server, etcd, scheduler, controller manager. [Lezione](fase-8/01-architettura.md)

Cron
:   Servizio che esegue comandi a orari pianificati. [Lezione](fase-2/05-pianificazione.md)

## D

DaemonSet
:   Controller Kubernetes che esegue una copia di un Pod su ogni nodo. [Lezione](fase-8/04-deployment.md)

Data source (Terraform)
:   Blocco che legge informazioni su risorse esistenti senza gestirle. [Lezione](fase-6/04-risorse.md)

Deployment
:   Controller Kubernetes che mantiene un numero di repliche di un Pod e ne gestisce gli aggiornamenti. [Lezione](fase-8/04-deployment.md)

Dev Container
:   Ambiente di sviluppo dentro un container, descritto in un file e riproducibile su qualunque computer. [Lezione](fase-0/07-laboratorio.md)

DevOps
:   Insieme di pratiche e cultura che unisce sviluppo e operazioni per rilasciare software in modo frequente e affidabile. [Lezione](fase-0/06-cloud-devops.md)

DHCP
:   Protocollo che assegna automaticamente indirizzo IP, gateway e DNS ai dispositivi di una rete. [Lezione](fase-1/08-dns-routing.md)

DNS
:   Sistema che traduce i nomi (esempio.it) in indirizzi IP. [Lezione](fase-1/08-dns-routing.md)

Dockerfile
:   File di istruzioni per costruire un'immagine container. [Lezione](fase-3/03-dockerfile.md)

Drift (deriva)
:   Differenza tra l'infrastruttura reale e quella descritta nel codice, di solito per modifiche manuali. [Lezione](fase-6/06-state.md)

## E

Elasticità
:   Capacità del cloud di aggiungere e togliere risorse rapidamente secondo la domanda. [Lezione](fase-4/01-cloud-computing.md)

Endpoint
:   Indirizzo a cui si raggiunge un servizio; in Kubernetes, gli indirizzi dei Pod dietro un Service. [Lezione](fase-8/05-service.md)

etcd
:   Database chiave-valore distribuito in cui Kubernetes salva lo stato del cluster. [Lezione](fase-8/12-kubeadm.md)

## F

Filesystem
:   Modo in cui i dati sono organizzati su un disco (ext4, XFS, VFAT). [Lezione](fase-1/06-storage.md)

Firewall
:   Filtro che consente o blocca il traffico di rete secondo regole. [Lezione](fase-1/11-firewall.md)

## G

Gateway API
:   API di Kubernetes per gestire il traffico in ingresso, successore di Ingress. [Lezione](fase-8/09-ingress.md)

GitOps
:   Pratica in cui lo stato desiderato di un sistema sta in Git e un agente lo applica e lo mantiene. [Lezione](fase-8/15-gitops-cloud.md)

## H

Hash
:   Impronta di dimensione fissa calcolata da un dato; cambia completamente se il dato cambia. [Lezione](fase-2/01-git-base.md)

Health check (controllo di salute)
:   Verifica periodica che un'istanza o un Pod risponda correttamente. [Lezione](fase-5/02-ha-scaling.md)

Helm
:   Gestore di pacchetti per Kubernetes; un pacchetto si chiama chart. [Lezione](fase-8/13-helm-kustomize.md)

HTTP e HTTPS
:   Protocollo del web; HTTPS è HTTP cifrato con TLS. [Lezione](fase-1/12-http-nginx.md)

## I

IaC (Infrastructure as Code)
:   Descrivere l'infrastruttura in file di codice versionati e applicati in automatico. [Lezione](fase-6/01-iac.md)

IAM
:   Sistema che decide chi può fare cosa sulle risorse cloud: identità, ruoli, policy. [Lezione](fase-4/04-iam.md)

Idempotenza
:   Proprietà di un'operazione che, ripetuta, produce lo stesso risultato senza nuovi cambiamenti. [Lezione](fase-6/10-ansible.md)

Immagine (container)
:   Pacchetto immutabile con file system e configurazione da cui si avviano i container. [Lezione](fase-3/02-immagini-container.md)

Ingress
:   Risorsa Kubernetes che espone servizi HTTP all'esterno tramite un controller. [Lezione](fase-8/09-ingress.md)

Inventario (Ansible)
:   Elenco degli host gestiti da Ansible, divisi in gruppi. [Lezione](fase-6/10-ansible.md)

Istanza
:   Macchina virtuale nel cloud (EC2, Compute Engine). [Lezione](fase-4/05-compute.md)

## J

JSON e YAML
:   Formati di testo per dati strutturati; YAML è più leggibile ed è usato per configurazioni. [Lezione](fase-2/07-python-automazione.md)

## K

Kernel
:   Il nucleo del sistema operativo, che gestisce hardware, memoria e processi. [Lezione](fase-0/02-sistemi-operativi.md)

kubectl
:   Il client da riga di comando di Kubernetes. [Lezione](fase-8/02-kubectl.md)

## L

Latenza
:   Tempo che impiega un dato ad andare e tornare; si misura in millisecondi. [Lezione](fase-4/02-infrastruttura-globale.md)

Least privilege (privilegio minimo)
:   Principio per cui ogni identità riceve solo i permessi indispensabili. [Lezione](fase-4/04-iam.md)

LVM
:   Gestore di volumi logici di Linux: unisce dischi in gruppi e crea volumi ridimensionabili. [Lezione](fase-1/06-storage.md)

## M

Metrica
:   Valore numerico misurato nel tempo, come l'uso della CPU o il numero di richieste. [Lezione](fase-7/06-prometheus.md)

Modulo (Terraform)
:   Insieme riutilizzabile di risorse con variabili in ingresso e output. [Lezione](fase-6/07-moduli.md)

## N

Namespace
:   In Kubernetes, spazio che separa gli oggetti di squadre o applicazioni diverse; in Linux, meccanismo di isolamento dei container. [Lezione](fase-8/03-pod.md)

NAT
:   Traduzione degli indirizzi che permette a reti private di uscire su Internet con un indirizzo pubblico. [Lezione](fase-1/08-dns-routing.md)

Network Policy
:   Regole Kubernetes che limitano il traffico tra Pod. [Lezione](fase-8/10-networkpolicy.md)

Nodo
:   Macchina di un cluster; in Kubernetes esegue i Pod. [Lezione](fase-8/01-architettura.md)

## O

Object storage
:   Storage che conserva file come oggetti in bucket, accessibili via API (S3, Cloud Storage). [Lezione](fase-4/06-storage.md)

OIDC
:   Protocollo di identità che permette, per esempio, a una pipeline di ottenere credenziali cloud temporanee senza chiavi salvate. [Lezione](fase-7/04-deploy-oidc.md)

Osservabilità
:   Capacità di capire lo stato interno di un sistema da metriche, log e tracce. [Lezione](fase-7/05-osservabilita.md)

## P

Pacchetto
:   Software pronto da installare con le sue informazioni e dipendenze (deb, rpm, snap, Flatpak). [Lezione](fase-1/05-pacchetti.md)

Pipeline
:   Sequenza automatica di passi (build, test, rilascio) eseguita a ogni modifica. [Lezione](fase-7/03-pipeline.md)

Playbook
:   File YAML di Ansible che descrive lo stato da applicare a un gruppo di host. [Lezione](fase-6/10-ansible.md)

Pod
:   La più piccola unità eseguibile in Kubernetes: uno o più container che condividono rete e volumi. [Lezione](fase-8/03-pod.md)

Porta
:   Numero che identifica un servizio su un host (22 per SSH, 443 per HTTPS). [Lezione](fase-1/07-tcpip.md)

Probe
:   Controllo di Kubernetes che verifica se un container è avviato, vivo o pronto. [Lezione](fase-8/03-pod.md)

Processo
:   Programma in esecuzione, con un identificativo (PID). [Lezione](fase-1/04-processi-systemd.md)

Provider (Terraform)
:   Plugin che collega Terraform alle API di una piattaforma. [Lezione](fase-6/02-primo-progetto.md)

Pull request
:   Proposta di modifica su GitHub, da rivedere e approvare prima dell'unione. [Lezione](fase-2/02-git-collaborazione.md)

## R

RBAC
:   Controllo degli accessi basato sui ruoli, per esempio in Kubernetes. [Lezione](fase-8/11-rbac.md)

Regione
:   Area geografica del cloud con più zone, scelta per latenza, costi e leggi sui dati. [Lezione](fase-4/02-infrastruttura-globale.md)

Registry
:   Archivio di immagini container (Docker Hub, GitHub Container Registry, ECR, Artifact Registry). [Lezione](fase-3/06-registry-sicurezza.md)

Repository
:   Archivio di un progetto con la sua storia (Git) o archivio di pacchetti (apt, dnf). [Lezione](fase-2/01-git-base.md)

Reverse proxy
:   Server che riceve le richieste dei client e le inoltra ai server applicativi. [Lezione](fase-1/12-http-nginx.md)

Rollback
:   Ritorno a una versione precedente dopo un rilascio problematico. [Lezione](fase-8/04-deployment.md)

Routing
:   Scelta del percorso dei pacchetti tra reti diverse. [Lezione](fase-1/08-dns-routing.md)

RPO e RTO
:   Quanti dati si possono perdere (RPO) e quanto tempo può durare il fermo (RTO) in caso di disastro. [Lezione](fase-5/10-well-architected.md)

## S

SELinux e AppArmor
:   Sistemi di controllo d'accesso obbligatorio che limitano cosa possono fare i processi, anche se eseguiti da root. [Lezione](fase-1/14-selinux-apparmor.md)

Serverless
:   Modello in cui il fornitore gestisce i server e si paga solo l'esecuzione effettiva (Lambda, Cloud Run). [Lezione](fase-5/09-container-serverless.md)

Service (Kubernetes)
:   Indirizzo stabile che bilancia il traffico verso un gruppo di Pod. [Lezione](fase-8/05-service.md)

Service account
:   Identità usata da programmi e macchine invece che da persone. [Lezione](fase-4/04-iam.md)

Shell
:   Programma che interpreta i comandi digitati nel terminale (per esempio Bash). [Lezione](fase-0/04-terminale.md)

SLA, SLO, SLI
:   Accordo con il cliente (SLA), obiettivo interno di affidabilità (SLO) e misura su cui si valuta (SLI). [Lezione](fase-7/05-osservabilita.md)

Snapshot
:   Copia di un disco o di un volume in un momento preciso. [Lezione](fase-5/03-storage.md)

SSH
:   Protocollo per collegarsi in modo sicuro a un computer remoto. [Lezione](fase-1/10-ssh.md)

State (Terraform)
:   File in cui Terraform registra le risorse che gestisce e i loro attributi. [Lezione](fase-6/06-state.md)

StatefulSet
:   Controller Kubernetes per applicazioni con identità e dischi stabili, come i database. [Lezione](fase-8/07-storage.md)

Subnet
:   Porzione di una rete con un proprio intervallo di indirizzi. [Lezione](fase-5/01-vpc.md)

sudo
:   Comando per eseguire un'operazione con i privilegi di un altro utente, di solito root. [Lezione](fase-1/03-permessi.md)

systemd
:   Il sistema di avvio e gestione dei servizi della maggior parte delle distribuzioni Linux. [Lezione](fase-1/04-processi-systemd.md)

## T

TCP e UDP
:   Protocolli di trasporto: TCP garantisce consegna e ordine, UDP è più semplice e veloce. [Lezione](fase-1/07-tcpip.md)

TLS
:   Protocollo che cifra le connessioni e verifica l'identità dei server con certificati. [Lezione](fase-1/12-http-nginx.md)

Traccia (trace)
:   Percorso di una richiesta attraverso più servizi, con i tempi di ogni passaggio. [Lezione](fase-7/05-osservabilita.md)

## V

Variabile d'ambiente
:   Valore passato a un processo dall'ambiente in cui viene avviato. [Lezione](fase-2/03-bash-base.md)

Virtualizzazione
:   Esecuzione di più macchine virtuali isolate su un unico computer fisico. [Lezione](fase-0/02-sistemi-operativi.md)

Volume
:   Spazio di archiviazione collegato a un container, un Pod o una macchina virtuale. [Lezione](fase-3/04-volumi-reti.md)

VPC
:   Rete privata virtuale nel cloud, con subnet, tabelle di routing e regole di accesso. [Lezione](fase-5/01-vpc.md)

## W

Webhook
:   Chiamata HTTP automatica inviata da un servizio quando accade un evento. [Lezione](fase-7/03-pipeline.md)

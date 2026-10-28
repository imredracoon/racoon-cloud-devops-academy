# 5. Server per i dati

Questi server **custodiscono e servono dati**: file, dischi di rete, copie di sicurezza, database. Sono quelli
in cui l'hardware conta di più, perché dischi, memoria e rete decidono le prestazioni e la sicurezza dei dati.

!!! abstract "Già visto nel corso"
    - **I tre tipi di storage**, a blocchi, a file e a oggetti: [lezione 4.6](../../fase-4/06-storage.md#tre-tipi-di-storage).
    - **NFS** e autofs, con server e client su entrambe le famiglie: [lezione 1.6](../../fase-1/06-storage.md#filesystem-di-rete-nfs-e-autofs).
    - **LVM**, filesystem e montaggio: [lezione 1.6](../../fase-1/06-storage.md#lvm-volumi-logici).
    - **RAID** e dischi del server: [pagina 3](03-hardware-server.md#raid-controller-hardware-o-software).
    - Database **relazionali e non relazionali**: [lezione 4.8](../../fase-4/08-database-ai.md#relazionali-e-non-relazionali);
      repliche e **cache**: [lezione 5.4](../../fase-5/04-database.md).
    - **Code** e messaggistica: [lezione 5.5](../../fase-5/05-messaggistica.md#code-pubsub-ed-eventi).
    - **Backup** centralizzati nel cloud: [lezione 5.3](../../fase-5/03-storage.md#backup-centralizzati).

## Il quadro d'insieme

| Server | Cosa offre | Software su Ubuntu | Software su RHEL |
|---|---|---|---|
| [File server](#file-server) | cartelle condivise | Samba, NFS | Samba, NFS |
| [NAS](#nas) | file server "pronto", in una scatola | (apparecchio dedicato) o Ubuntu con ZFS | (apparecchio dedicato) o RHEL con XFS e LVM |
| [Storage a blocchi (SAN)](#storage-a-blocchi-e-san) | dischi virtuali per altri server | targetcli (iSCSI) | targetcli (iSCSI) |
| [Storage distribuito e a oggetti](#storage-distribuito-e-a-oggetti) | spazio che cresce aggiungendo nodi | Ceph (MicroCeph) | Red Hat Ceph Storage |
| [Backup](#server-di-backup) | copie di sicurezza | Bacula, rsnapshot, restic | ReaR, Bacula |
| [Database](#server-di-database) | dati strutturati | PostgreSQL, MySQL | PostgreSQL, MySQL, MariaDB |
| [Cache](#server-di-cache) | dati in memoria per risposte rapide | Valkey o Redis, Memcached | Valkey, Memcached |
| [Code di messaggi](#code-di-messaggi) | comunicazione asincrona tra applicazioni | RabbitMQ, Kafka | Kafka (Streams for Apache Kafka), RabbitMQ |

## File server

Condivide cartelle in rete. Due protocolli:

- **SMB** (Windows e macOS, ma anche Linux), offerto su Linux da **Samba**;
- **NFS** (Linux e Unix), che hai configurato nella [lezione 1.6](../../fase-1/06-storage.md#filesystem-di-rete-nfs-e-autofs).

**Hardware tipico**: 4-8 core, 16-64 GB di RAM (fa da cache), molti dischi in RAID 6 con qualche SSD di cache,
rete a 10-25 Gbit/s.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y samba
    sudo mkdir -p /srv/condivisa && sudo chgrp sambashare /srv/condivisa && sudo chmod 2770 /srv/condivisa
    sudo tee -a /etc/samba/smb.conf <<'EOF'
    [condivisa]
       path = /srv/condivisa
       read only = no
       valid users = @sambashare
    EOF
    sudo smbpasswd -a mario          # password Samba per un utente Linux esistente
    sudo systemctl restart smbd
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y samba
    sudo mkdir -p /srv/condivisa
    sudo semanage fcontext -a -t samba_share_t "/srv/condivisa(/.*)?"   # lezione 1.14
    sudo restorecon -Rv /srv/condivisa
    # stessa sezione [condivisa] in /etc/samba/smb.conf, con valid users = @gruppo
    sudo smbpasswd -a mario
    sudo systemctl enable --now smb
    sudo firewall-cmd --permanent --add-service=samba && sudo firewall-cmd --reload
    ```

Su RHEL ricorda le etichette SELinux, viste nella [lezione 1.14](../../fase-1/14-selinux-apparmor.md#i-contesti):
senza `samba_share_t` Samba non può leggere la cartella anche se i permessi sono giusti.

## NAS

Un **NAS** (*Network Attached Storage*) è un file server venduto come **apparecchio**: hardware, dischi e un
sistema operativo con interfaccia web per creare condivisioni, utenti, snapshot e backup senza usare il terminale.
Esempi: Synology (RackStation, come l'RS2423RP+ della [pagina 3](03-hardware-server.md#modelli-reali-a-confronto)),
QNAP, NetApp e Dell PowerScale per le grandi aziende. Molti NAS sono basati su Linux.

Costruirne uno con un server standard dà più controllo:

=== "Debian/Ubuntu"

    Ubuntu supporta **ZFS**, un filesystem che unisce RAID, gestione dei volumi, controllo di integrità dei dati
    (*checksum*), compressione e **snapshot** istantanei.

    ```bash
    sudo apt install -y zfsutils-linux
    # pool con doppia parità (equivalente a RAID 6) su sei dischi
    sudo zpool create -o ashift=12 dati raidz2 /dev/sdb /dev/sdc /dev/sdd /dev/sde /dev/sdf /dev/sdg
    sudo zfs create -o compression=lz4 dati/condivisa
    sudo zfs snapshot dati/condivisa@prima-della-modifica
    sudo zpool status
    ```

=== "RHEL (Rocky/AlmaLinux)"

    RHEL non include ZFS. La combinazione supportata è **RAID software con mdadm**, **LVM** (con snapshot) e
    il filesystem **XFS**. Esiste anche **Stratis**, che gestisce pool e snapshot sopra XFS con comandi semplici:
    Red Hat lo ha classificato come anteprima tecnologica, quindi verifica lo stato di supporto nella tua versione.

    ```bash
    sudo mdadm --create /dev/md0 --level=6 --raid-devices=6 /dev/sd{b,c,d,e,f,g}
    sudo pvcreate /dev/md0 && sudo vgcreate dati /dev/md0
    sudo lvcreate -L 10T -n condivisa dati && sudo mkfs.xfs /dev/dati/condivisa
    # oppure Stratis:
    sudo dnf install -y stratisd stratis-cli && sudo systemctl enable --now stratisd
    ```

Distribuzioni specializzate come **TrueNAS** (con ZFS) offrono un'interfaccia simile a quella dei NAS commerciali
su hardware qualsiasi.

## Storage a blocchi e SAN

Un file server condivide **file**. Uno storage a blocchi presenta invece a un altro server un **disco intero**
(chiamato LUN), che quel server formatta e usa come se fosse interno. È la versione fisica dei dischi delle VM
nel cloud ([lezione 4.6](../../fase-4/06-storage.md#i-dischi-delle-vm)). La rete dedicata a questo traffico si
chiama **SAN** (*Storage Area Network*) e usa:

- **Fibre Channel**: rete separata con switch e schede (HBA) propri, molto affidabile, costosa;
- **iSCSI**: comandi SCSI su normale rete IP, più economico;
- **NVMe over Fabrics** (NVMe/TCP, NVMe/RoCE): il più recente, per SSD NVMe in rete con latenza bassissima.

Gli **array di storage** aziendali (Dell PowerStore, HPE Alletra, NetApp AFF, Pure Storage FlashArray) sono
apparati con due controller ridondati che offrono LUN, snapshot, replica e deduplica. Con Linux un server diventa
un **target iSCSI**:

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y targetcli-fb
    sudo targetcli              # shell interattiva: backstores, iscsi, luns, acls
    # sul client (initiator):
    sudo apt install -y open-iscsi
    sudo iscsiadm -m discovery -t sendtargets -p IP-TARGET
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y targetcli
    sudo systemctl enable --now target
    sudo targetcli
    sudo firewall-cmd --permanent --add-service=iscsi-target && sudo firewall-cmd --reload
    # sul client:
    sudo dnf install -y iscsi-initiator-utils
    sudo iscsiadm -m discovery -t sendtargets -p IP-TARGET
    ```

Quando un server raggiunge lo stesso LUN da due percorsi (due schede, due switch) si usa il **multipath**
(`multipath-tools` su Ubuntu, `device-mapper-multipath` su RHEL) per vederlo come un disco solo e sopravvivere
alla caduta di un percorso.

## Storage distribuito e a oggetti

Invece di un grande array, molti server normali pieni di dischi formano un unico sistema di storage che cresce
aggiungendo nodi e sopravvive alla perdita di dischi e di interi server. Il progetto open source di riferimento è
**Ceph**, che offre insieme storage a blocchi, a file e **a oggetti compatibile con S3**
([lezione 4.6](../../fase-4/06-storage.md#lo-storage-a-oggetti)). È lo stesso principio con cui i provider cloud
costruiscono i loro servizi di storage.

**Hardware tipico**: almeno tre nodi (meglio cinque), ognuno con 12-24 dischi, SSD NVMe per i metadati, due reti
a 25-100 Gbit/s (una per i client, una per la replica interna).

=== "Debian/Ubuntu"

    Canonical distribuisce **MicroCeph**, un pacchetto snap che installa e configura Ceph in pochi comandi:

    ```bash
    sudo snap install microceph
    sudo microceph cluster bootstrap
    sudo microceph disk add /dev/sdb --wipe
    ```

=== "RHEL (Rocky/AlmaLinux)"

    Red Hat offre **Red Hat Ceph Storage**, un prodotto separato con abbonamento proprio, installato con lo
    strumento `cephadm` in container Podman. Per i test si può usare la versione della comunità.

## Server di backup

Il backup protegge da ciò che il RAID non copre: cancellazioni, ransomware, guasti del sito. La regola di base è
il **3-2-1**: **3** copie dei dati, su **2** supporti diversi, di cui **1** fuori sede (oggi spesso anche una copia
**immutabile** o scollegata, che un attacco non può cancellare). Un backup non provato non è un backup: il
**ripristino** va verificato con regolarità.

**Hardware tipico**: pochi core, molta capacità su HDD, in un luogo diverso da quello dei dati; a volte una
**libreria a nastro** (LTO), ancora molto usata per l'archivio a lungo termine perché economica e scollegata.

=== "Debian/Ubuntu"

    - **Bacula**: sistema client-server completo per molti server, con pianificazioni e cataloghi.
    - **rsnapshot**: copie incrementali basate su rsync, semplici e leggibili.
    - **restic** e **Borg**: backup cifrati e deduplicati, anche verso storage a oggetti.

    ```bash
    sudo apt install -y restic
    restic -r /mnt/backup/repo init
    restic -r /mnt/backup/repo backup /etc /srv
    restic -r /mnt/backup/repo snapshots
    ```

=== "RHEL (Rocky/AlmaLinux)"

    - **ReaR** (*Relax-and-Recover*), incluso in RHEL: crea un'immagine di ripristino avviabile dell'intero
      sistema per il *bare-metal recovery*, cioè ricostruire un server da zero su hardware nuovo.
    - Per i dati si usano software commerciali (Veeam, Commvault e altri) o Bacula.

    ```bash
    sudo dnf install -y rear
    sudoedit /etc/rear/local.conf     # OUTPUT=ISO, BACKUP=NETFS, BACKUP_URL=nfs://...
    sudo rear -v mkbackup
    ```

Gli snapshot di **LVM** e **ZFS** sono utili per avere un'immagine coerente da copiare, ma da soli non sono un
backup: stanno sugli stessi dischi.

## Server di database

Ospita database relazionali come **PostgreSQL** e **MySQL/MariaDB** o non relazionali. Nel cloud li hai usati
gestiti ([lezione 5.4](../../fase-5/04-database.md#il-database-relazionale-gestito)); on-premises sei tu a
installarli, configurarli, replicarli e fare i backup.

**Hardware tipico**: processori con core veloci, **molta RAM** (il database tiene in memoria i dati più usati),
SSD **NVMe in RAID 10**, rete a 25 Gbit/s, due server (primario e replica) o tre per un cluster.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y postgresql
    sudo -u postgres psql -c "SELECT version();"
    # oppure
    sudo apt install -y mysql-server
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y postgresql-server
    sudo postgresql-setup --initdb
    sudo systemctl enable --now postgresql
    # oppure mysql-server o mariadb-server
    sudo firewall-cmd --permanent --add-service=postgresql && sudo firewall-cmd --reload
    ```

## Server di cache

Tiene in **memoria** dati letti spesso (sessioni, risultati di query, pagine) per rispondere in meno di un
millisecondo e alleggerire il database, come visto nella [lezione 5.4](../../fase-5/04-database.md#le-cache).
**Valkey** è il proseguimento open source di Redis, nato nel 2024 quando Redis ha cambiato licenza; **Memcached**
è un'alternativa più semplice.

**Hardware tipico**: tanta RAM, pochi dischi, rete veloce.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y valkey-server    # nelle versioni recenti; nelle precedenti: redis-server
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y valkey           # RHEL 10 include Valkey al posto di Redis
    sudo systemctl enable --now valkey
    ```

## Code di messaggi

Disaccoppiano le applicazioni: chi produce un messaggio non aspetta chi lo elabora
([lezione 5.5](../../fase-5/05-messaggistica.md)). **RabbitMQ** è il broker di code più diffuso; **Apache Kafka**
è una piattaforma di flussi di eventi per volumi molto grandi.

**Hardware tipico**: cluster di tre nodi, SSD veloci (Kafka scrive tutto su disco), molta rete.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y rabbitmq-server
    sudo rabbitmqctl status
    ```

=== "RHEL (Rocky/AlmaLinux)"

    RabbitMQ non è nei repository di RHEL: si installa dai repository del progetto o in container. Red Hat
    offre Kafka come prodotto supportato (**Streams for Apache Kafka**), di solito su OpenShift.

## Esercizi

!!! esercizio "Esercizio 1"
    Un'azienda tiene tutti i dati su un NAS in RAID 6 e fa ogni notte uno snapshot sullo stesso NAS. Rispetta la
    regola 3-2-1? Cosa manca?

??? success "Soluzione"
    No. Gli snapshot stanno sugli stessi dischi: un incendio, un furto o un ransomware che cifra il NAS cancella
    tutto. Servono almeno una **seconda copia su un altro supporto** (un server di backup) e **una copia fuori
    sede**, meglio se immutabile (nel cloud o su nastro).

!!! esercizio "Esercizio 2"
    Per ognuno indica la tecnologia giusta: (a) cartelle condivise per PC Windows; (b) un disco per un server di
    virtualizzazione; (c) archiviare milioni di foto con un'API compatibile S3; (d) ricostruire un server RHEL da
    zero dopo un guasto.

??? success "Soluzione"
    (a) **Samba** (SMB). (b) Un **LUN** su storage a blocchi via iSCSI, Fibre Channel o NVMe/TCP. (c) **Storage a
    oggetti**, per esempio Ceph con il gateway S3. (d) **ReaR**.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è la differenza principale tra un NAS e una SAN?","options":["Il NAS condivide file, la SAN offre dischi a blocchi","Il NAS è più veloce","La SAN funziona solo nel cloud","Non c'è differenza"],"answer":[0],"explanation":"Il NAS parla SMB/NFS (file); la SAN presenta LUN che il server formatta (blocchi)."},
 {"q":"Cosa prevede la regola 3-2-1 dei backup?","options":["3 copie dei dati","2 supporti diversi","1 copia fuori sede","3 backup al giorno","2 amministratori diversi"],"answer":[0,1,2],"explanation":"Tre copie, due supporti, una fuori sede."},
 {"q":"Su RHEL, quale combinazione si usa al posto di ZFS?","options":["mdadm, LVM e XFS (o Stratis)","Btrfs","NTFS","FAT32"],"answer":[0],"explanation":"RHEL non include ZFS; XFS con LVM e RAID software sono la scelta supportata."},
 {"q":"Quale protocollo trasporta comandi a blocchi su una normale rete IP?","options":["SMB","NFS","iSCSI","HTTP"],"answer":[2],"explanation":"iSCSI porta i comandi SCSI su TCP/IP."},
 {"q":"Quale strumento incluso in RHEL crea un'immagine per ricostruire un server da zero?","options":["rsnapshot","ReaR (Relax-and-Recover)","Valkey","targetcli"],"answer":[1],"explanation":"ReaR è pensato per il bare-metal recovery."},
 {"q":"Quale hardware conta di più per un server di database?","options":["Molta RAM","SSD NVMe veloci","Una GPU potente","Molti HDD da 3,5\" lenti","Processori con core veloci"],"answer":[0,1,4],"explanation":"RAM come cache, NVMe per la latenza, core veloci per le query; GPU e HDD lenti non aiutano un database tradizionale."}
]}
</script>

## Fonti

- [Ubuntu Server: Samba file server](https://ubuntu.com/server/docs/how-to/samba/file-server/)
- [Ubuntu Server: Backups (Bacula, rsnapshot)](https://ubuntu.com/server/docs/how-to/backups/)
- [Ubuntu Server: iSCSI initiator](https://ubuntu.com/server/docs/how-to/storage/iscsi-initiator-or-client/)
- [Ubuntu Server: Databases](https://ubuntu.com/server/docs/how-to/databases/)
- [RHEL 10: Configuring and using network file services (Samba)](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html-single/configuring_and_using_network_file_services/index)
- [RHEL 10: Managing storage devices](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/managing_storage_devices/index)
- [RHEL 10: Configuring and using database servers](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html-single/configuring_and_using_database_servers/index)
- [OpenZFS documentation](https://openzfs.github.io/openzfs-docs/)
- [Ceph documentation](https://docs.ceph.com/)
- [MicroCeph documentation](https://documentation.ubuntu.com/microceph/)
- [Relax-and-Recover (ReaR)](https://relax-and-recover.org/)
- [Valkey](https://valkey.io/)

# 3. L'hardware di un server

!!! obiettivi "Obiettivi"
    - Capire in cosa un server differisce da un PC.
    - Leggere la scheda tecnica di un server reale.
    - Scegliere CPU, memoria, dischi e rete in base al lavoro che il server dovrà fare.
    - Riconoscere l'hardware di una macchina Linux dal terminale.

!!! abstract "Già visto nel corso"
    - **CPU, RAM, disco e scheda di rete**, cosa fa ciascuno: [lezione 0.1](../../fase-0/01-computer.md#i-componenti-principali).
    - Cos'è un **server** e perché lavora senza schermo: [lezione 0.1](../../fase-0/01-computer.md#il-server).
    - **Dischi a blocchi**, partizioni e filesystem: [lezione 1.6](../../fase-1/06-storage.md#dischi-a-blocchi-e-partizioni).
    - Le **famiglie di macchine virtuali** nel cloud (uso generale, calcolo, memoria, GPU): [lezione 4.5](../../fase-4/05-compute.md#le-macchine-virtuali).
      Le stesse famiglie esistono nell'hardware fisico.

## Cosa rende un server un server

Un server contiene gli stessi componenti di un PC, ma progettati per lavorare **giorno e notte per anni** e per
essere riparati senza spegnerli:

| Caratteristica | Perché conta |
|---|---|
| **Memoria ECC** | corregge da sola gli errori di un bit nella RAM, che su anni di lavoro sono inevitabili |
| **Componenti ridondati** | due alimentatori, più ventole, dischi in RAID: un guasto non ferma il servizio |
| **Sostituzione a caldo** (*hot-swap*) | dischi, alimentatori e ventole si cambiano a server acceso |
| **BMC** | un piccolo computer indipendente per accendere, spegnere e controllare il server da remoto |
| **Formato rack** | si monta in un armadio standard ([pagina 2](02-rack.md)) |
| **Certificazioni dei sistemi operativi** | il produttore garantisce che Ubuntu, RHEL e altri funzionano e sono supportati |
| **Supporto e garanzia** | ricambi in giornata, anche per 5-7 anni |

## I componenti, uno per uno

### Processori

I server hanno spesso **due socket**, cioè due processori sulla stessa scheda madre. Le famiglie attuali sono
**Intel Xeon 6** e **AMD EPYC**: arrivano a decine o oltre cento core per processore. Si sceglie tra:

- **molti core** per virtualizzazione, container e servizi web con tanti utenti contemporanei;
- **meno core ma più veloci** (frequenza alta) per database e software con licenza per core;
- **un solo socket** quando basta, per risparmiare corrente e licenze.

Esistono anche server **ARM** (per esempio con processori Ampere), più efficienti nei consumi; Ubuntu e RHEL
hanno versioni per ARM (architettura `aarch64` o `arm64`).

### Memoria

I server usano moduli **RDIMM DDR5 con ECC**. Le regole pratiche:

- popola **tutti i canali** di memoria in modo uniforme: ogni processore ne ha 8 o 12, e lasciarne vuoti riduce
  la velocità;
- dimensiona in base al lavoro: pochi GB per un DNS, centinaia di GB per un database o un host di virtualizzazione.

### Archiviazione

| Tipo | Interfaccia | Uso |
|---|---|---|
| SSD **NVMe** (formati U.2, U.3, E3.S) | PCIe, collegato direttamente al processore | database, virtualizzazione, tutto ciò che deve essere veloce |
| SSD **SAS/SATA** | controller SAS | uso generale, a costo minore |
| HDD **SAS/SATA** (3,5") | controller SAS | grandi volumi a basso costo per TB: backup, archivio, NAS |
| **M.2** in coppia (es. Dell BOSS, HPE NS204i) | dedicato | solo per il sistema operativo, separato dai dati |

### RAID: controller hardware o software

Il **RAID** unisce più dischi per avere ridondanza, velocità o entrambe:

| Livello | Come funziona | Dischi persi senza perdere dati | Spazio utile con 4 dischi da 4 TB |
|---|---|---|---|
| RAID 0 | divide i dati sui dischi | **nessuno** | 16 TB |
| RAID 1 | copia identica su due dischi | 1 della coppia | 4 TB (con 2 dischi) |
| RAID 5 | dati più una parità distribuita | 1 | 12 TB |
| RAID 6 | dati più due parità | 2 | 8 TB |
| RAID 10 | coppie in RAID 1, poi in RAID 0 | 1 per coppia | 8 TB |

Il RAID può essere fatto da un **controller hardware** (Dell PERC, HPE Smart Array e MR, Broadcom MegaRAID), con
una sua cache protetta da batteria, oppure **dal sistema operativo** con un semplice **HBA** che passa i dischi
così come sono. Su Linux il RAID software si fa con `mdadm` (entrambe le famiglie), con **LVM**
([lezione 1.6](../../fase-1/06-storage.md#lvm-volumi-logici)) o, su Ubuntu, con **ZFS**. Con gli SSD NVMe oggi si
preferisce spesso il RAID software.

!!! warning "Il RAID non è un backup"
    Il RAID protegge dal guasto di un disco, non da un file cancellato per errore, da un ransomware o da un
    incendio. Per quelli serve un backup, vedi il [server di backup](05-server-dati.md#server-di-backup).

### Rete

I server hanno di solito due o quattro porte a **10, 25 o 100 Gbit/s**, spesso su una scheda in formato
**OCP 3.0** (uno slot standard, sostituibile dal retro). Due porte collegate a due switch diversi, unite con il
*bonding* (LACP, vedi la [pagina 1](01-dispositivi-rete.md#switch)), danno velocità e ridondanza. Per AI e HPC si
aggiungono schede **InfiniBand** o Ethernet a 400 Gbit/s con **RDMA**, che trasferiscono dati tra le memorie dei
server senza passare dal processore.

### GPU e acceleratori

Le **GPU** fanno in parallelo migliaia di calcoli semplici: servono per l'addestramento e l'uso di modelli di
intelligenza artificiale, il rendering e il calcolo scientifico. Nei server si trovano in due forme:

- **schede PCIe** (per esempio NVIDIA L4, L40S, H100 NVL): da una a otto in un server standard;
- **moduli SXM** su una scheda base (piattaforme **NVIDIA HGX**): otto GPU collegate tra loro con **NVLink**,
  in server dedicati alti da 4U a 10U.

AMD (Instinct) e Intel (Gaudi) producono acceleratori concorrenti.

### Alimentazione

Due **alimentatori ridondanti** sostituibili a caldo, con efficienza certificata (**80 PLUS Platinum** o
**Titanium**). La potenza va da 800 W per un server leggero a oltre 3 kW per alimentatore nei server GPU.

### Il BMC: la gestione fuori banda

Il **BMC** (*Baseboard Management Controller*) è un piccolo computer dentro il server, con la sua porta di rete,
che resta acceso anche a server spento. Permette di:

- accendere, spegnere e riavviare il server da remoto;
- vedere lo schermo e usare tastiera e mouse (**console remota**), anche durante l'avvio e nel BIOS;
- montare un'immagine ISO da remoto per installare il sistema operativo (*virtual media*);
- leggere sensori, registro degli eventi e stato di ogni componente;
- aggiornare il firmware.

Ogni marca ha il suo: **iDRAC** (Dell), **iLO** (HPE), **XClarity Controller** (Lenovo). Tutti parlano il
protocollo standard **Redfish** (API REST in JSON, definita dal consorzio DMTF) e quasi tutti ancora il vecchio
**IPMI**. La porta del BMC va su una **rete di gestione separata**, mai su Internet: chi controlla il BMC
controlla il server.

```bash
# Interrogare un BMC con Redfish: stato di alimentazione e salute del server
curl -sk -u utente:password https://IP-DEL-BMC/redfish/v1/Systems/ | jq
# Con IPMI (pacchetto ipmitool, stesso nome su Ubuntu e RHEL)
ipmitool -I lanplus -H IP-DEL-BMC -U utente -P password chassis power status
ipmitool -I lanplus -H IP-DEL-BMC -U utente -P password sensor list
```

### TPM e avvio sicuro

Il **TPM 2.0** è un chip che custodisce chiavi crittografiche e misura cosa si avvia. Con il **Secure Boot** UEFI
il server avvia solo bootloader e kernel firmati: Ubuntu e RHEL sono firmati e lo supportano. Il TPM permette
anche di sbloccare in automatico i dischi cifrati all'avvio.

## Formati

| Formato | Caratteristiche | Quando si usa |
|---|---|---|
| **Tower** | come un PC grande, silenzioso | piccoli uffici senza armadio |
| **Rack 1U** | sottile, pochi dischi, molti server per armadio | web, applicazioni, nodi di cluster |
| **Rack 2U** | più dischi, più schede, più GPU | il formato più versatile: database, virtualizzazione, storage |
| **Rack 4U e oltre** | molti dischi o molte GPU | storage ad alta densità, server GPU |
| **Blade e modulari** | molti server in un telaio con alimentazione e rete condivise | grandi installazioni |
| **Edge** | compatti, robusti, poco profondi | negozi, fabbriche, antenne |

## Modelli reali a confronto

Le tre marche principali producono server equivalenti con nomi diversi. I numeri cambiano a ogni generazione;
verificati il 28 ottobre 2026 sulle pagine dei produttori.

| Ruolo | Dell | HPE | Lenovo |
|---|---|---|---|
| Rack 1U, 2 socket, uso generale | PowerEdge **R670** | ProLiant Compute **DL360 Gen12** | ThinkSystem **SR630 V4** |
| Rack 2U, 2 socket, versatile | PowerEdge **R770** | ProLiant Compute **DL380 Gen12** | ThinkSystem **SR650 V4** |
| Controller di gestione | iDRAC10 | iLO 7 | XClarity Controller 3 |

**Esempio: un 2U di uso generale.** Il Dell PowerEdge R770 ha due processori Intel Xeon 6, fino a 8 TB di
memoria DDR5 ECC in 32 slot, configurazioni fino a 40 SSD NVMe E3.S o 24 dischi SAS/SATA da 2,5", fino a quattro
schede di rete OCP 3.0, due alimentatori sostituibili a caldo e iDRAC10. È certificato per Ubuntu Server LTS e
RHEL. L'HPE DL380 Gen12 e il Lenovo SR650 V4 sono costruiti per lo stesso ruolo con caratteristiche simili.

**Esempio: un server GPU.** Il Dell PowerEdge XE9780 ha due processori Intel Xeon 6 e **otto GPU NVIDIA** in
formato SXM (piattaforma HGX B200 o B300) collegate con NVLink, raffreddate ad aria in un telaio da circa 10U.
È fatto per addestrare modelli di intelligenza artificiale.

**Esempio: un NAS da rack.** Il Synology RackStation RS2423RP+ è un 2U con 12 alloggiamenti per dischi
(espandibili a 24 con un'unità aggiuntiva), una porta 10GbE e due 1GbE, alimentatori ridondanti (le lettere
**RP** nel nome) e un sistema operativo proprio basato su Linux. È un'alternativa "pronta" a un file server
costruito con Ubuntu o RHEL, vedi la [pagina 5](05-server-dati.md#nas).

## Dimensionare: dal servizio all'hardware

| Servizio | CPU | RAM | Dischi | Rete | Note |
|---|---|---|---|---|---|
| DNS, DHCP, NTP | pochi core | 4-8 GB | piccoli SSD in RAID 1 | 1-10 Gbit/s | meglio due server piccoli che uno grande |
| Web e applicazioni | molti core | medi | SSD | 10-25 Gbit/s | scala aggiungendo server |
| Database | core veloci | **molta** | **NVMe** in RAID 10 | 25 Gbit/s | la RAM fa da cache dei dati |
| Virtualizzazione | **molti core** | **molta** | NVMe o storage condiviso | 25-100 Gbit/s | la somma delle VM che ospita |
| File server e NAS | pochi core | medi | **molti HDD** in RAID 6 + SSD di cache | 10-25 Gbit/s | capacità prima di tutto |
| Backup | pochi core | medi | molti HDD | 10-25 Gbit/s | in un luogo diverso dai dati |
| GPU e AI | molti core | **molta** | NVMe veloci | 100-400 Gbit/s | GPU, energia e raffreddamento decidono tutto |

## Riconoscere l'hardware da Linux

Questi comandi funzionano allo stesso modo su Ubuntu e RHEL:

```bash
lscpu                       # processori, core, thread, architettura
free -h                     # memoria
lsblk -o NAME,SIZE,TYPE,MODEL,TRAN   # dischi e tipo di collegamento (nvme, sas, sata)
lspci | grep -Ei 'ethernet|raid|nvme|vga|3d'   # schede sul bus PCIe
ip -br link                 # interfacce di rete (lezione 1.7)
sudo ethtool eno1 | grep -E 'Speed|Link'       # velocità negoziata di una porta
sudo dmidecode -t system    # produttore, modello e numero di serie
sudo dmidecode -t memory | grep -E 'Size|Type:|Speed' | sort | uniq -c
```

Il numero di serie (*service tag* per Dell) serve per l'assistenza e per controllare la garanzia.

## Esercizi

!!! esercizio "Esercizio 1"
    Hai 6 dischi da 8 TB. Quanto spazio utile ottieni con RAID 5, RAID 6 e RAID 10? Quale scegli per un file
    server con molti dati e poche scritture, e quale per un database?

??? success "Soluzione"
    RAID 5: 40 TB (perdi un disco). RAID 6: 32 TB (perdi due dischi). RAID 10: 24 TB (metà). Per il **file
    server** RAID 6: con dischi grandi la ricostruzione dura ore o giorni e un secondo guasto durante la
    ricostruzione sarebbe fatale con RAID 5. Per il **database** RAID 10: scritture più veloci e ricostruzione
    rapida.

!!! esercizio "Esercizio 2"
    Sulla tua VM del laboratorio esegui i comandi della sezione precedente. Cosa noti nell'output di
    `dmidecode -t system`?

??? success "Soluzione"
    Il produttore e il modello non sono quelli di un server fisico ma dell'**hypervisor** (per esempio "QEMU",
    "VirtualBox" o il nome del provider cloud). È il segno che il sistema gira in una macchina virtuale.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"A cosa serve la memoria ECC?","options":["A rendere la RAM più veloce","A correggere automaticamente gli errori di un bit","A usare meno corrente","A collegare più GPU"],"answer":[1],"explanation":"Su server accesi per anni gli errori di memoria sono inevitabili: l'ECC li corregge."},
 {"q":"Cosa permette di fare il BMC?","options":["Accendere e spegnere il server da remoto","Vedere la console anche durante l'avvio","Montare un'ISO da remoto","Aumentare la RAM","Leggere i sensori anche a server spento"],"answer":[0,1,2,4],"explanation":"Il BMC è un computer indipendente per la gestione fuori banda; non cambia l'hardware."},
 {"q":"Quale standard aperto usano i BMC moderni per l'API REST di gestione?","options":["SNMP","Redfish","SMTP","LDAP"],"answer":[1],"explanation":"Redfish, definito dal DMTF, ha sostituito in gran parte IPMI."},
 {"q":"Con 4 dischi da 4 TB in RAID 6, quanto spazio utile ottieni?","options":["4 TB","8 TB","12 TB","16 TB"],"answer":[1],"explanation":"RAID 6 usa l'equivalente di due dischi per la parità."},
 {"q":"Quale affermazione è vera?","options":["Il RAID sostituisce il backup","Il RAID protegge dal guasto di un disco ma non dalla cancellazione di un file","RAID 0 tollera il guasto di un disco","RAID 1 raddoppia lo spazio"],"answer":[1],"explanation":"Il RAID è disponibilità, non protezione dei dati da errori, attacchi o disastri."},
 {"q":"Quale comando mostra modello e numero di serie del server?","options":["lscpu","dmidecode -t system","lsblk","ip link"],"answer":[1],"explanation":"dmidecode legge le tabelle SMBIOS del firmware."}
]}
</script>

## Fonti

- [Dell PowerEdge R770, spec sheet](https://delltechnologies.com/asset/en-ph/products/servers/technical-support/poweredge-r770-spec-sheet.pdf)
- [Dell PowerEdge R670, spec sheet](https://www.delltechnologies.com/asset/en-gb/products/servers/technical-support/poweredge-r670-spec-sheet.pdf)
- [Dell PowerEdge XE9780](https://www.dell.com/en-us/shop/ipovw/poweredge-xe9780)
- [HPE ProLiant Compute DL380 Gen12](https://www.hpe.com/us/en/compute/proliant/dl380-gen12.html)
- [HPE ProLiant Compute DL360 Gen12](https://www.hpe.com/us/en/compute/hpe-proliant-compute/dl360-gen12.html)
- [Lenovo ThinkSystem SR650 V4, product guide](https://lenovopress.lenovo.com/lp2127-thinksystem-sr650-v4-server)
- [Lenovo ThinkSystem SR630 V4, product guide](https://lenovopress.lenovo.com/lp1971-thinksystem-sr630-v4-server)
- [Synology RackStation RS2423+ / RS2423RP+](https://www.synology.com/en-global/company/news/article/rs2423plus)
- [DMTF Redfish](https://www.dmtf.org/standards/redfish)
- [RHEL 10: Managing storage devices](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/managing_storage_devices/index)

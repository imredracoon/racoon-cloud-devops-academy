# Verifica finale

Venti domande su tutta l'appendice. Se sbagli, la spiegazione indica la pagina da ripassare. Non è un esame di
certificazione: serve a capire se hai un quadro chiaro di come è fatta un'infrastruttura fisica.

<script type="application/json" class="quiz-data">
{"title":"Hardware, server e data center · 20 domande","questions":[
 {"q":"Un dispositivo collega reti IP diverse e decide il percorso dei pacchetti. Di quale si tratta?","options":["Switch","Router","Access point","Patch panel"],"answer":[1],"explanation":"Il router lavora al livello 3 (pagina 1)."},
 {"q":"Vuoi separare telefoni IP, ospiti e uffici sugli stessi switch. Cosa usi?","options":["Più UPS","VLAN","RAID","Un KVM"],"answer":[1],"explanation":"Le VLAN 802.1Q creano reti logiche separate (pagina 1)."},
 {"q":"Quali funzioni offre uno switch gestito da ufficio?","options":["VLAN","Aggregazione dei collegamenti (LACP)","PoE, nei modelli che lo prevedono","Backup dei database","Spanning Tree"],"answer":[0,1,2,4],"explanation":"Il backup non è un compito dello switch (pagina 1)."},
 {"q":"Un server da 2U quanto spazio occupa in altezza?","options":["44,45 mm","88,9 mm","19 pollici","2 cm"],"answer":[1],"explanation":"1U = 44,45 mm, quindi 2U = 88,9 mm (pagina 2)."},
 {"q":"A cosa serve un console server in un rack?","options":["A raggiungere da remoto le porte seriali di switch e firewall","A mostrare le pagine web","A raffreddare i server","A distribuire la corrente"],"answer":[0],"explanation":"È l'accesso di emergenza agli apparati di rete (pagina 2)."},
 {"q":"Perché i server si dispongono con corridoi caldi e freddi?","options":["Per estetica","Per non mescolare l'aria fresca in ingresso con quella calda in uscita","Per ridurre i cavi","Per la sicurezza fisica"],"answer":[1],"explanation":"Il contenimento rende il raffreddamento più efficiente (pagina 2)."},
 {"q":"Come si chiama il controller di gestione dei server HPE?","options":["iDRAC","iLO","XClarity Controller","Cockpit"],"answer":[1],"explanation":"iDRAC è di Dell, XClarity Controller di Lenovo; Cockpit è una console web del sistema operativo (pagina 3)."},
 {"q":"Con 8 dischi da 2 TB in RAID 10, quanto spazio utile ottieni?","options":["4 TB","8 TB","12 TB","16 TB"],"answer":[1],"explanation":"RAID 10 usa metà dei dischi per le copie (pagina 3)."},
 {"q":"Quale formato e collegamento usano le 8 GPU di un server per il training come il Dell XE9780?","options":["Schede PCIe senza collegamenti tra loro","Moduli SXM collegati con NVLink","GPU esterne via USB","GPU integrate nel processore"],"answer":[1],"explanation":"Le piattaforme HGX collegano otto GPU SXM con NVLink (pagina 3)."},
 {"q":"Su RHEL, quale prodotto gestisce identità, Kerberos, DNS e CA in modo integrato?","options":["Samba AD","Identity Management (FreeIPA)","OpenLDAP","MAAS"],"answer":[1],"explanation":"Su Ubuntu le alternative sono Samba AD o OpenLDAP (pagina 4)."},
 {"q":"Quale servizio deve essere sempre presente in almeno due istanze?","options":["DNS","DHCP","Controller di dominio","Server di stampa di un piccolo ufficio","NTP"],"answer":[0,1,2,4],"explanation":"I servizi di infrastruttura critici vanno ridondati; la stampa di un piccolo ufficio può aspettare (pagina 4)."},
 {"q":"Su Ubuntu, quale filesystem unisce RAID, volumi, checksum e snapshot?","options":["XFS","ext4","ZFS","FAT32"],"answer":[2],"explanation":"Su RHEL si usano invece mdadm, LVM e XFS (pagina 5)."},
 {"q":"Quale tecnologia presenta a un server un disco remoto a blocchi su rete IP?","options":["SMB","iSCSI","IMAP","SNMP"],"answer":[1],"explanation":"SMB e NFS condividono file; iSCSI offre LUN a blocchi (pagina 5)."},
 {"q":"Un backup su un disco USB nello stesso armadio del server. Cosa manca per la regola 3-2-1?","options":["Una copia fuori sede","Un controller RAID","Una GPU","Un secondo alimentatore"],"answer":[0],"explanation":"Un incendio o un furto distruggerebbero dati e backup insieme (pagina 5)."},
 {"q":"Quali componenti formano la virtualizzazione standard su Linux?","options":["KVM","QEMU","libvirt","Squid","Kea"],"answer":[0,1,2],"explanation":"Squid è un proxy, Kea un server DHCP (pagina 6)."},
 {"q":"Quale file usa RHEL per l'installazione automatica?","options":["autoinstall","Kickstart","Dockerfile","playbook"],"answer":[1],"explanation":"Ubuntu Server usa autoinstall (pagina 7)."},
 {"q":"Prima di installare il sistema su un server nuovo, cosa configuri per lavorare da remoto?","options":["Il BMC","Il proxy","La stampante","Il database"],"answer":[0],"explanation":"Dal BMC hai console, supporti virtuali e accensione remota (pagina 7)."},
 {"q":"Quale strumento controlla la salute di un disco SATA o SAS?","options":["smartctl","fwupdmgr","kpatch","nmcli"],"answer":[0],"explanation":"smartctl, del pacchetto smartmontools, legge gli attributi SMART (pagina 8)."},
 {"q":"Quale topologia di rete garantisce al massimo due salti tra due server qualsiasi?","options":["Tre livelli core-aggregazione-accesso","Spine-leaf","Anello","Stella con un solo switch"],"answer":[1],"explanation":"Ogni leaf è collegato a tutti gli spine (pagina 9)."},
 {"q":"Quale tipo di data center è pensato per portare il calcolo vicino agli utenti o ai dispositivi?","options":["Hyperscale","Edge","Colocation","HPC"],"answer":[1],"explanation":"I micro data center edge riducono la latenza in negozi, fabbriche e antenne (pagina 9)."}
]}
</script>

## E adesso?

Torna all'[indice dell'appendice](index.md) per ripassare, oppure prova a disegnare l'infrastruttura del tuo
[progetto finale](../../fase-8/progetto.md) come se dovessi ospitarla in un tuo data center: quanti server, quali
servizi di infrastruttura, quali reti e quale ridondanza.

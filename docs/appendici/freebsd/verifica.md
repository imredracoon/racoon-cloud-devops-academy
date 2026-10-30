# Verifica finale

Dodici domande su tutta l'appendice. Se sbagli, la spiegazione indica la pagina da ripassare.

<script type="application/json" class="quiz-data">
{"title":"FreeBSD · 12 domande","questions":[
 {"q":"Qual è la differenza principale tra FreeBSD e una distribuzione Linux?","options":["FreeBSD non ha un kernel","FreeBSD sviluppa kernel e sistema di base insieme, in un solo progetto","FreeBSD non è un sistema Unix","FreeBSD funziona solo su server ARM"],"answer":[1],"explanation":"Una distribuzione Linux unisce componenti di progetti diversi (pagina 1)."},
 {"q":"Hai installato nginx con pkg. Dove trovi nginx.conf?","options":["/etc/nginx/","/usr/local/etc/nginx/","/opt/nginx/","/var/nginx/"],"answer":[1],"explanation":"Il software aggiunto vive in /usr/local (pagina 1)."},
 {"q":"Quali comandi attivano nginx all'avvio e lo avviano subito?","options":["systemctl enable --now nginx","sysrc nginx_enable=YES e service nginx start","pkg enable nginx","rc-update add nginx"],"answer":[1],"explanation":"I servizi si accendono in /etc/rc.conf (pagina 2)."},
 {"q":"Come aggiorni le correzioni di sicurezza del sistema di base installato con i distribution sets?","options":["freebsd-update fetch e freebsd-update install","apt upgrade","make world","pkg audit"],"answer":[0],"explanation":"pkg upgrade aggiorna i pacchetti; con pkgbase anche il sistema di base (pagina 2)."},
 {"q":"Quali compiti, su Linux divisi tra mdadm, LVM e filesystem, ZFS svolge da solo?","options":["Ridondanza dei dischi","Volumi che condividono lo spazio","Filesystem con checksum","Firewall"],"answer":[0,1,2],"explanation":"ZFS unisce pool, dataset e filesystem (pagina 3)."},
 {"q":"Cosa fa zfs rollback dati/documenti@prima?","options":["Copia lo snapshot su un altro server","Riporta il dataset allo stato dello snapshot","Cancella il pool","Crea un boot environment"],"answer":[1],"explanation":"Per copiare si usano zfs send e zfs receive (pagina 3)."},
 {"q":"Cosa ti permette di fare un boot environment?","options":["Avviare un'altra distribuzione","Tornare al sistema com'era prima di un aggiornamento, scegliendolo all'avvio","Avviare le jail","Cifrare il disco"],"answer":[1],"explanation":"bectl crea e attiva i boot environment (pagina 3)."},
 {"q":"Cosa condividono una jail e il suo host?","options":["Il kernel","La radice del filesystem","Gli utenti","L'indirizzo IP"],"answer":[0],"explanation":"Come i container Linux, le jail usano il kernel dell'host (pagina 4)."},
 {"q":"Quale comando esegue una shell dentro la jail web?","options":["jail -c web","jexec web sh","jls web","service jail web"],"answer":[1],"explanation":"jexec è l'equivalente di docker exec (pagina 4)."},
 {"q":"In pf, senza la parola quick, quale regola decide per un pacchetto?","options":["La prima che corrisponde","L'ultima che corrisponde","La più specifica","Sempre la regola block"],"answer":[1],"explanation":"È il contrario di nftables (pagina 5)."},
 {"q":"Qual è l'equivalente di KVM su FreeBSD?","options":["bhyve","jail","pf","ZFS"],"answer":[0],"explanation":"bhyve si usa di solito con un gestore come vm-bhyve (pagina 5)."},
 {"q":"Nel laboratorio, perché il port forwarding funziona solo dopo aver disattivato Block private networks sulla WAN?","options":["Perché OPNsense non supporta il NAT","Perché nel laboratorio la WAN è una rete privata, 192.168.56.0/24","Perché il client non ha Nginx","Perché manca un alias"],"answer":[1],"explanation":"Su Internet gli indirizzi privati non devono arrivare: l'opzione li blocca (pagina 6)."}
]}
</script>

## E adesso?

Torna all'[indice dell'appendice](index.md) per ripassare. Se il tuo homelab ha un mini PC libero, prova OPNsense su
una macchina fisica, come nella [pagina 6](06-opnsense.md#per-andare-oltre-un-firewall-fisico-livello-3).

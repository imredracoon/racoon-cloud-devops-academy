# Verifica di fine fase 1

Passa alla Fase 2 quando ottieni **almeno l'80%** (16 risposte corrette su 20). Le spiegazioni indicano la lezione
da ripassare. Se vuoi prepararti alla LFCS, prosegui con l'[appendice dedicata](../appendici/lfcs/index.md).

<script type="application/json" class="quiz-data">
{"title":"Verifica di fine fase 1 · 20 domande","questions":[
 {"q":"In quale directory si trovano i log tradizionali di sistema?","options":["/etc/log","/var/log","/usr/log","/tmp/log"],"answer":[1],"explanation":"Lezione 1.1."},
 {"q":"Hai un hard link e un link simbolico verso lo stesso file. Cancelli il file originale. Cosa succede?","options":["Entrambi smettono di funzionare","L'hard link funziona ancora, il simbolico è rotto","Il simbolico funziona ancora, l'hard link no","Entrambi funzionano"],"answer":[1],"explanation":"L'hard link punta all'inode, il simbolico al nome (lezione 1.1)."},
 {"q":"Quale pipe conta quante volte compare ogni valore della prima colonna di un file, dal più frequente?","options":["cut -f1 file | uniq -c | sort","awk '{print $1}' file | sort | uniq -c | sort -rn","sort file | wc -l","grep -c file"],"answer":[1],"explanation":"uniq richiede input ordinato (lezione 1.2)."},
 {"q":"Come mandi sia stdout sia stderr di un comando nel file out.log?","options":["cmd 2> out.log","cmd > out.log 2>&1","cmd 2>&1 > out.log","cmd | out.log"],"answer":[1],"explanation":"L'ordine delle redirezioni conta (lezione 1.2)."},
 {"q":"Quali permessi ottali corrispondono a rw-r-----?","options":["640","644","604","460"],"answer":[0],"explanation":"Lezione 1.3."},
 {"q":"In una directory condivisa vuoi che i nuovi file ereditino il gruppo della directory. Quale permesso speciale usi?","options":["setuid","setgid","sticky bit","umask"],"answer":[1],"explanation":"Lezione 1.3."},
 {"q":"Quale combinazione avvia subito un servizio e lo abilita all'avvio?","options":["systemctl start","systemctl enable","systemctl enable --now","systemctl reload"],"answer":[2],"explanation":"Lezione 1.4."},
 {"q":"Un servizio è in stato failed con `status=203/EXEC`. Cosa controlli per primo?","options":["La rete","Il percorso e i permessi del programma in ExecStart","Lo spazio su disco","Il firewall"],"answer":[1],"explanation":"203/EXEC: systemd non è riuscito a eseguire il programma (lezione 1.4)."},
 {"q":"Qual è il gestore di pacchetti di Amazon Linux 2023?","options":["apt","dnf","zypper","pacman"],"answer":[1],"explanation":"Famiglia Red Hat (lezioni 0.2 e 1.5)."},
 {"q":"Hai aggiunto un disco in fstab senza `nofail` e il disco viene rimosso. Cosa può succedere al riavvio?","options":["Nulla","Il sistema può bloccarsi in modalità di emergenza","Il disco viene ricreato","fstab viene cancellato"],"answer":[1],"explanation":"Lezione 1.6."},
 {"q":"Qual è l'ordine corretto degli strati LVM, dal disco al filesystem?","options":["VG → PV → LV","PV → VG → LV","LV → VG → PV","PV → LV → VG"],"answer":[1],"explanation":"Lezione 1.6."},
 {"q":"Quanti host utilizzabili ha una /28?","options":["16","14","30","6"],"answer":[1],"explanation":"2^4 − 2 (lezione 1.7)."},
 {"q":"Gli host 10.0.1.10/23 e 10.0.0.200/23 sono nella stessa rete?","options":["Sì","No"],"answer":[0],"explanation":"Una /23 raggruppa due /24 consecutive con il terzo ottetto pari come inizio: 10.0.0.0 – 10.0.1.255 (lezione 1.7)."},
 {"q":"Quale record DNS crea un alias verso un altro nome?","options":["A","MX","CNAME","TXT"],"answer":[2],"explanation":"Lezione 1.8."},
 {"q":"Una tabella ha le rotte 10.0.0.0/8 via A e 10.1.0.0/16 via B. Verso quale gateway va un pacchetto per 10.1.2.3?","options":["A","B","Nessuno","Entrambi"],"answer":[1],"explanation":"Vince il prefisso più lungo (lezione 1.8)."},
 {"q":"`ss -tlpn` mostra un servizio su 127.0.0.1:3000. Cosa significa?","options":["È raggiungibile da tutti","È raggiungibile solo dalla macchina stessa","È spento","Usa UDP"],"answer":[1],"explanation":"Lezione 1.9."},
 {"q":"Quale impostazione di sshd impedisce l'accesso con password?","options":["PermitRootLogin no","PasswordAuthentication no","PubkeyAuthentication no","UsePAM no"],"answer":[1],"explanation":"Lezione 1.10."},
 {"q":"Quale regola nftables permette le risposte alle connessioni già accettate?","options":["ct state established,related accept","tcp dport 22 accept","policy accept","iif lo accept"],"answer":[0],"explanation":"Lezione 1.11."},
 {"q":"Nginx risponde 502. Il log dice `connect() failed (111: Connection refused) while connecting to upstream`. Causa più probabile?","options":["Certificato scaduto","Il backend non è in ascolto sulla porta configurata","Il DNS del client è sbagliato","Il firewall blocca la 443"],"answer":[1],"explanation":"Lezioni 1.12 e 1.13."},
 {"q":"Un servizio è sparito e `dmesg` mostra `Out of memory: Killed process`. Cosa è successo?","options":["Il disco è pieno","Il kernel ha terminato il processo per mancanza di RAM","Il servizio è stato aggiornato","Un utente ha usato kill"],"answer":[1],"explanation":"È l'OOM killer (lezione 1.13)."}
]}
</script>

## E adesso?

Hai completato anche il [progetto](progetto.md)? Passa alla [Fase 2 · Automazione](../fase-2/index.md).

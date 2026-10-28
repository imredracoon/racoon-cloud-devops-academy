# Verifica di fine fase 2

Passa alla [Fase 3](../fase-3/index.md) quando ottieni **almeno l'80%** (16 risposte corrette su 20). Le spiegazioni indicano la lezione
da ripassare.

<script type="application/json" class="quiz-data">
{"title":"Verifica di fine fase 2 · 20 domande","questions":[
 {"q":"Hai modificato tre file ma vuoi includerne solo due nel prossimo commit. Cosa fai?","options":["git commit -a","git add dei due file, poi git commit","git push dei due file","git stash"],"answer":[1],"explanation":"La staging area serve proprio a scegliere cosa registrare (lezione 2.1)."},
 {"q":"Un collega ha già scaricato il tuo branch. Vuoi annullare un commit che contiene un errore. Quale comando usi?","options":["git reset --hard HEAD~1 e git push --force","git revert <hash>","git commit --amend","git rebase -i"],"answer":[1],"explanation":"revert non riscrive la storia condivisa (lezioni 2.1 e 2.2)."},
 {"q":"Durante un merge compare un conflitto in config.txt. Qual è la sequenza corretta?","options":["git merge --continue, poi modifica il file","Modifica il file togliendo i marcatori, git add config.txt, git commit","git push --force","git restore config.txt e git commit"],"answer":[1],"explanation":"Lezione 2.2."},
 {"q":"Qual è il vantaggio principale di una pull request rispetto a un push diretto su main?","options":["È più veloce","Le modifiche vengono riviste e verificate dai controlli automatici prima di entrare in main","Non richiede Git","Evita i conflitti"],"answer":[1],"explanation":"Lezione 2.2."},
 {"q":"Hai cancellato per errore un commit con git reset --hard. Quale comando ti aiuta a ritrovarlo?","options":["git log","git reflog","git fsck --lost","git status"],"answer":[1],"explanation":"Il reflog registra tutti gli spostamenti di HEAD (lezione 2.2)."},
 {"q":"Il file si chiama `report finale.txt`. Quale comando lo cancella correttamente in uno script?","options":["rm $file","rm \"$file\"","rm '$file'","rm ${file}"],"answer":[1],"explanation":"Senza virgolette doppie la variabile viene divisa sugli spazi (lezione 2.3)."},
 {"q":"Cosa stampa `grep -q root /etc/passwd; echo $?` se root esiste?","options":["root","0","1","true"],"answer":[1],"explanation":"0 significa successo; -q non stampa nulla (lezione 2.3)."},
 {"q":"Su Ubuntu uno script con `#!/bin/sh` usa `[[ ]]` e fallisce. Perché?","options":["[[ ]] non esiste","Su Ubuntu /bin/sh è dash, che non supporta [[ ]]","Manca il permesso di esecuzione","Ubuntu non supporta gli script"],"answer":[1],"explanation":"Su RHEL /bin/sh è bash e lo stesso script funzionerebbe (lezione 2.3)."},
 {"q":"Quali opzioni rendono uno script Bash più robusto?","options":["set -e, per fermarsi al primo errore","set -o pipefail, perché una pipeline fallisca se fallisce un suo comando","set +x, per nascondere gli errori","shopt -s nocasematch","set -o vi"],"answer":[0,1],"explanation":"Insieme a -u formano il classico set -euo pipefail (lezione 2.4)."},
 {"q":"Perché in uno script `for riga in $(cat file)` è sbagliato per leggere un file riga per riga?","options":["È lento","Divide il contenuto sugli spazi, non sulle righe","Non legge l'ultima riga","Funziona solo con file piccoli"],"answer":[1],"explanation":"Si usa while IFS= read -r riga (lezione 2.4)."},
 {"q":"Quale riga di crontab esegue un comando ogni giorno feriale alle 7:30?","options":["7 30 * * 1-5","30 7 * * 1-5","30 7 1-5 * *","* 7:30 * * 1-5"],"answer":[1],"explanation":"minuto ora giorno mese giorno-della-settimana (lezione 2.5)."},
 {"q":"Un job pianificato deve recuperare l'esecuzione persa se il server era spento. Quale soluzione scegli?","options":["Un crontab con @reboot","Un timer systemd con Persistent=true","Un comando at","Un ciclo while con sleep"],"answer":[1],"explanation":"Lezione 2.5."},
 {"q":"Hai creato backup.service e backup.timer. Quale comando attiva la pianificazione?","options":["systemctl enable --now backup.service","systemctl enable --now backup.timer","systemctl start backup.service","crontab -e"],"answer":[1],"explanation":"Si abilita il timer; il servizio viene avviato da lui (lezione 2.5)."},
 {"q":"Quale struttura Python rappresenta meglio la configurazione di una VM con nome, CPU e RAM?","options":["list","dict","set","str"],"answer":[1],"explanation":"Un dizionario associa chiavi a valori ed è ciò in cui si traducono JSON e YAML (lezione 2.6)."},
 {"q":"Su Ubuntu 24.04 `pip install requests` fuori da un ambiente virtuale restituisce `externally-managed-environment`. Qual è la soluzione corretta?","options":["Usare sudo pip install","Creare un ambiente virtuale con python3 -m venv e installare lì","Cancellare il file EXTERNALLY-MANAGED","Passare a Python 2"],"answer":[1],"explanation":"La PEP 668 protegge i pacchetti del sistema (lezione 2.6)."},
 {"q":"Quale chiamata subprocess è sicura con un nome file fornito dall'utente?","options":["subprocess.run(f\"wc -l {nome}\", shell=True)","subprocess.run([\"wc\", \"-l\", nome])","os.system(\"wc -l \" + nome)","eval(nome)"],"answer":[1],"explanation":"Con una lista di argomenti non c'è interpretazione della shell (lezione 2.7)."},
 {"q":"Perché si imposta sempre `timeout=` nelle chiamate requests?","options":["Per obbligo di legge","Perché senza timeout una richiesta può restare bloccata per sempre","Per ridurre il traffico","Per abilitare HTTPS"],"answer":[1],"explanation":"Lezione 2.7."},
 {"q":"Un client invia POST /tasks senza il campo title. Quale codice deve restituire la Tasks API?","options":["200","201","400","500"],"answer":[2],"explanation":"400 Bad Request: l'errore è nella richiesta del client (lezione 2.8)."},
 {"q":"Come si esegue la Tasks API in produzione?","options":["python app.py con debug=True","flask run","gunicorn \"app:create_app()\" dietro un reverse proxy come Nginx","Con cron ogni minuto"],"answer":[2],"explanation":"Il server di Flask è solo per lo sviluppo (lezione 2.8)."},
 {"q":"Perché la Tasks API legge DATABASE_URL dall'ambiente invece di scriverla nel codice?","options":["Per rendere il codice più corto","Perché lo stesso codice deve girare con database diversi in ambienti diversi, senza segreti nel codice","Perché SQLAlchemy non accetta stringhe","Per velocità"],"answer":[1],"explanation":"Twelve-Factor App, principio Config (lezione 2.8)."}
]}
</script>

## E adesso?

Hai completato anche il [progetto](progetto.md)? Passa alla [Fase 3 · Container](../fase-3/index.md).

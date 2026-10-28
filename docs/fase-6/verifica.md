# Verifica di fine fase 6

Passa alla [Fase 7](../fase-7/index.md) quando ottieni **almeno l'80%** (20 risposte corrette su 25). Le spiegazioni indicano la lezione
da ripassare. Per l'esame Terraform Associate continua con l'[appendice dedicata](../appendici/terraform/index.md),
che contiene una simulazione completa; per la RHCE con l'[appendice RHCE](../appendici/rhce/index.md).

<script type="application/json" class="quiz-data">
{"title":"Verifica di fine fase 6 · 25 domande","questions":[
 {"q":"Quale caratteristica rende Terraform idempotente?","options":["È scritto in Go","È dichiarativo: confronta lo stato voluto con quello reale e applica solo le differenze","Usa SSH","Salva i log"],"answer":[1],"explanation":"Lezione 6.1."},
 {"q":"Quale strumento è più adatto a installare pacchetti e gestire file di configurazione dentro le VM?","options":["Terraform","Ansible","CloudFormation","HCP Terraform"],"answer":[1],"explanation":"Configuration management (lezioni 6.1 e 6.10)."},
 {"q":"Cosa fa terraform init?","options":["Crea le risorse","Inizializza la cartella: scarica provider e moduli, configura il backend","Formatta il codice","Distrugge lo state"],"answer":[1],"explanation":"Lezioni 6.2 e 6.3."},
 {"q":"Il vincolo version = \"~> 6.2\" per un provider accetta...","options":["Solo la 6.2.0","Le versioni da 6.2 in su, ma inferiori alla 7.0","Le versioni da 6.2.0 a 6.2.x","Qualunque versione"],"answer":[1],"explanation":"Con ~> 6.2.0 sarebbero ammesse solo le 6.2.x (lezione 6.2)."},
 {"q":"Come gestisci con Terraform risorse in due regioni AWS diverse nella stessa configurazione?","options":["Due blocchi provider aws, uno con alias, e l'argomento provider nelle risorse","Non è possibile","Due file di lock","Con i workspace"],"answer":[0],"explanation":"Lezione 6.2."},
 {"q":"Quale file registra le versioni esatte dei provider e va versionato con Git?","options":["terraform.tfstate",".terraform.lock.hcl","terraform.tfvars",".terraform/"],"answer":[1],"explanation":"Lezione 6.3."},
 {"q":"Come applichi esattamente il piano approvato in revisione?","options":["terraform apply -auto-approve","terraform plan -out=piano.tfplan e poi terraform apply piano.tfplan","terraform apply -target","terraform refresh"],"answer":[1],"explanation":"Lezione 6.3."},
 {"q":"Quale variabile d'ambiente attiva i log dettagliati?","options":["TF_LOG","TF_VERBOSE","TF_DEBUG_LEVEL","TF_VAR_log"],"answer":[0],"explanation":"TRACE è il livello più dettagliato (lezione 6.3)."},
 {"q":"Cosa fa terraform fmt?","options":["Valida la configurazione","Riformatta i file secondo lo stile canonico","Crea il piano","Aggiorna i provider"],"answer":[1],"explanation":"Con -check verifica senza modificare (lezione 6.3)."},
 {"q":"Ti serve l'ID dell'immagine Amazon Linux più recente. Cosa usi?","options":["Un blocco resource","Un data source","Un output","Una variabile senza default"],"answer":[1],"explanation":"Lezione 6.4."},
 {"q":"Una risorsa A usa aws_vpc.principale.id. Che tipo di dipendenza si crea?","options":["Esplicita","Implicita","Nessuna","Circolare"],"answer":[1],"explanation":"depends_on crea dipendenze esplicite (lezione 6.4)."},
 {"q":"Quale argomento di lifecycle fa creare la nuova risorsa prima di distruggere la vecchia?","options":["prevent_destroy","create_before_destroy","ignore_changes","replace_triggered_by"],"answer":[1],"explanation":"Lezione 6.4."},
 {"q":"Ordina per precedenza crescente: default, TF_VAR_, terraform.tfvars, -var. Quale vince?","options":["default","TF_VAR_","terraform.tfvars","-var"],"answer":[3],"explanation":"Lezione 6.5."},
 {"q":"Perché preferire for_each a count per un elenco di utenti?","options":["È obbligatorio","Togliere un elemento non sposta gli indirizzi degli altri","count non supporta le stringhe","for_each crea risorse più velocemente"],"answer":[1],"explanation":"Lezione 6.5."},
 {"q":"Quale costrutto genera blocchi annidati ripetuti, come le regole di un security group?","options":["Un blocco dynamic","Un modulo","Un local","Un blocco check"],"answer":[0],"explanation":"Lezione 6.5."},
 {"q":"Quali meccanismi fanno fallire l'operazione se la condizione non è soddisfatta?","options":["validation di una variabile","precondition","assert di un blocco check","Una description","Un commento"],"answer":[0,1],"explanation":"I blocchi check producono solo avvisi (lezione 6.5)."},
 {"q":"Cosa contiene lo state?","options":["Solo i nomi delle risorse","La corrispondenza tra risorse della configurazione e oggetti reali, con i loro attributi, a volte segreti in chiaro","Il codice sorgente","Le credenziali del provider"],"answer":[1],"explanation":"Lezione 6.6."},
 {"q":"Cosa impedisce il blocco dello state?","options":["La lettura","Che due esecuzioni scrivano contemporaneamente lo stesso state","L'uso dei moduli","L'importazione"],"answer":[1],"explanation":"Lezione 6.6."},
 {"q":"Qualcuno ha modificato a mano un'istanza gestita da Terraform. Cosa mostra terraform plan?","options":["Niente","Le azioni per riportare l'istanza a quanto descritto nella configurazione","Un errore di blocco","La distruzione dello state"],"answer":[1],"explanation":"Con -refresh-only vedi solo le differenze (lezione 6.6)."},
 {"q":"Come porti sotto Terraform un bucket creato a mano?","options":["Con un blocco import (eventualmente con -generate-config-out)","Con terraform state rm","Con un blocco moved","Ricreandolo"],"answer":[0],"explanation":"Lezione 6.6."},
 {"q":"Un modulo figlio può leggere direttamente var.ambiente del modulo root?","options":["Sì","No: il valore va passato come argomento nel blocco module","Solo se è sensitive","Solo con terraform_remote_state"],"answer":[1],"explanation":"Lezione 6.7."},
 {"q":"Come fissi la versione di un modulo preso da un repository Git?","options":["Con l'argomento version","Con ?ref=TAG nel source","Con il file di lock","Non si può"],"answer":[1],"explanation":"version vale solo per i registry (lezione 6.7)."},
 {"q":"Cosa garantisce sensitive = true su un output?","options":["Che il valore sia cifrato nello state","Che il valore non venga mostrato nell'output della CLI","Che il valore non venga salvato","Che il valore venga ruotato"],"answer":[1],"explanation":"Per non salvarlo servono valori effimeri e argomenti write-only (lezione 6.8)."},
 {"q":"In HCP Terraform, cosa avvia un piano speculativo?","options":["Una pull request sul repository collegato a un workspace VCS-driven","Il comando terraform fmt","Un variable set","Un blocco moved"],"answer":[0],"explanation":"Lezione 6.9."},
 {"q":"In un playbook Ansible, come riavvii un servizio solo quando il suo file di configurazione cambia?","options":["Con un task shell alla fine","Con un handler notificato dal task che modifica il file","Con --check","Con un ciclo"],"answer":[1],"explanation":"Lezione 6.10."}
]}
</script>

## E adesso?

Hai completato anche il [progetto](progetto.md)? Passa alla [Fase 7 · CI/CD e osservabilità](../fase-7/index.md): il
codice dell'infrastruttura e dell'applicazione verrà verificato e rilasciato automaticamente da una pipeline.

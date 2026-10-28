# 4.5 Calcolo

!!! obiettivi "Obiettivi"
    - Creare, configurare e distruggere una macchina virtuale nel cloud.
    - Scegliere tipo di macchina, immagine e disco.
    - Automatizzare la configurazione all'avvio con user data e script di avvio.
    - Conoscere i modelli di prezzo: on demand, impegni, spot.
    - Orientarsi tra le altre opzioni di calcolo: container gestiti, Kubernetes, funzioni serverless.

## Le macchine virtuali

=== "AWS"

    Il servizio è **Amazon EC2** (*Elastic Compute Cloud*). Una VM si chiama **istanza**. Gli elementi:

    - **AMI** (*Amazon Machine Image*): l'immagine del sistema operativo (Amazon Linux 2023, Ubuntu, RHEL, Windows...).
    - **Tipo di istanza**: CPU e memoria. Il nome ha una struttura precisa: in `t3.micro`, `t` è la **famiglia**
      (uso generale con crediti di CPU), `3` la **generazione**, `micro` la **dimensione**. Una `g` dopo la
      generazione (`t4g`, `m7g`) indica i processori **Graviton** (ARM) di AWS, più economici.
    - **Volume EBS**: il disco, persistente e indipendente dall'istanza (lezione 4.6).
    - **Coppia di chiavi** per SSH, oppure **EC2 Instance Connect** / **Session Manager** per collegarsi senza
      gestire chiavi.
    - **Security group**: il firewall dell'istanza (lezione 4.7).
    - **User data**: uno script eseguito al primo avvio.

    | Famiglia | Esempi | Uso |
    |---|---|---|
    | Uso generale | `t3`, `t4g`, `m7i`, `m8g` | web server, applicazioni |
    | Ottimizzate per il calcolo | `c7i`, `c8g` | elaborazioni intensive |
    | Ottimizzate per la memoria | `r7i`, `x2` | database, cache in memoria |
    | Calcolo accelerato | `p5`, `g6` | GPU per AI e grafica |
    | Ottimizzate per lo storage | `i4i` | database con molto I/O locale |

=== "Google Cloud"

    Il servizio è **Compute Engine**. Gli elementi:

    - **Immagine**: Debian, Ubuntu, RHEL, Rocky Linux, Container-Optimized OS, Windows...
    - **Tipo di macchina**: in `e2-medium`, `e2` è la **serie**, `medium` la dimensione. Esistono tipi predefiniti
      (`e2-standard-2`: 2 vCPU, 8 GB) e **personalizzati**, in cui scegli CPU e memoria.
    - **Disco di avvio**: Persistent Disk o Hyperdisk (lezione 4.6).
    - **Accesso SSH**: `gcloud compute ssh` gestisce le chiavi in automatico; con **OS Login** l'accesso è legato
      all'identità IAM.
    - **Regole firewall** della rete VPC, applicate tramite **tag di rete** o service account (lezione 4.7).
    - **Script di avvio** (*startup script*), eseguito a **ogni** avvio.

    | Famiglia | Serie | Uso |
    |---|---|---|
    | Uso generale | `e2`, `n4`, `c4a` (ARM Axion), `t2d` | web server, applicazioni |
    | Ottimizzate per il calcolo | `c3`, `h3` | elaborazioni intensive |
    | Ottimizzate per la memoria | `m3`, `x4` | database molto grandi in memoria |
    | Acceleratori | `a3`, `g2` | GPU e AI |

Le famiglie e le generazioni cambiano spesso: per scegliere usa la documentazione ufficiale e il calcolatore dei
prezzi. Per il corso bastano sempre le macchine più piccole: `t3.micro` / `t4g.micro` su AWS, `e2-micro` su Google
Cloud.

## La configurazione all'avvio

Una VM nel cloud non si configura a mano via SSH: è lento, non ripetibile e non scala. Si passa uno **script** che
viene eseguito all'avvio. Sulla maggior parte delle immagini Linux lo esegue **cloud-init**, lo stesso strumento che
usi con Multipass (Fase 0).

```bash title="avvio.sh"
#!/bin/bash
# Installa Nginx e pubblica una pagina con il nome della macchina
apt-get update
apt-get install -y nginx
echo "<h1>Servito da $(hostname)</h1>" > /var/www/html/index.html
systemctl enable --now nginx
```

## I modelli di prezzo

=== "AWS"

    | Modello | Sconto indicativo | Quando |
    |---|---|---|
    | **On-Demand** | — | carichi brevi, imprevedibili, prove. Si paga al secondo (Linux), senza impegno. |
    | **Savings Plans** | fino a circa 72% | impegno di spesa oraria per 1 o 3 anni, flessibile su famiglia, regione e anche Fargate e Lambda (Compute Savings Plans) |
    | **Reserved Instances** | fino a circa 72% | impegno per 1 o 3 anni su un tipo di istanza in una regione |
    | **Spot** | fino a circa 90% | capacità inutilizzata che AWS può **riprendersi con 2 minuti di preavviso**: per lavori interrompibili (elaborazioni batch, CI, test) |
    | **Dedicated Hosts** / **Dedicated Instances** | — | server fisici dedicati, per licenze o requisiti di conformità |

=== "Google Cloud"

    | Modello | Sconto indicativo | Quando |
    |---|---|---|
    | **On demand** | — | fatturazione al secondo (con un minimo di un minuto) |
    | **Sconti per utilizzo sostenuto** | automatici | applicati alle VM di alcune serie (come N1 e N2) che restano accese per gran parte del mese, senza impegno |
    | **Sconti per impegno di utilizzo** (*Committed Use Discounts*) | fino a circa 55-70% | impegno per 1 o 3 anni su risorse o su una spesa |
    | **Spot VM** | 60-91% | capacità inutilizzata che Google può riprendersi con 30 secondi di preavviso |
    | **Sole-tenant node** | — | server fisici dedicati |

Gli sconti esatti dipendono da servizio, regione e durata: controllali sempre nelle pagine dei prezzi. Una strategia
tipica: la base stabile coperta da impegni, i picchi con on demand, i lavori interrompibili con spot.

## Oltre le VM

| Esigenza | AWS | Google Cloud |
|---|---|---|
| Eseguire container senza gestire server | **ECS** con **Fargate**, App Runner | **Cloud Run** |
| Kubernetes gestito | **EKS** | **GKE** |
| Funzioni eseguite a evento | **Lambda** | **Cloud Run functions** (ex Cloud Functions) |
| Piattaforma per applicazioni web | Elastic Beanstalk | App Engine |
| Elaborazioni batch | AWS Batch | Batch |
| VPS semplice a prezzo fisso | Lightsail | — |

Le funzioni serverless (Lambda, Cloud Run functions) eseguono il tuo codice in risposta a un **evento** (una richiesta
HTTP, un file caricato, un messaggio in coda), scalano da zero a migliaia di esecuzioni e si pagano solo per il tempo
di esecuzione. Hanno limiti (per esempio, una funzione Lambda dura al massimo 15 minuti): le vedrai nella Fase 5.

## Laboratorio

!!! laboratorio "Laboratorio · Un web server nel cloud"
    Crea una VM con Nginx configurato dallo script di avvio, raggiungila dal browser, collegati in SSH, poi
    **distruggila**.

    === "AWS"

        1. Console → **EC2** → regione del corso → **Launch instance**.
        2. Nome `lab-web`; AMI **Ubuntu Server 24.04 LTS**; tipo **t3.micro** (o `t4g.micro` con l'AMI ARM).
        3. Key pair: **Proceed without a key pair** (ti collegherai con EC2 Instance Connect).
        4. Network settings → **Edit**: crea un security group `lab-web-sg` con due regole: **HTTP (80)** da **My IP**
           e **SSH (22)** con sorgente la *prefix list* `com.amazonaws.REGIONE.ec2-instance-connect` (gli indirizzi da
           cui si collega EC2 Instance Connect quando lo usi dal browser).
        5. Advanced details → **User data**: incolla `avvio.sh`.
        6. **Launch instance**. Dopo un paio di minuti apri `http://IP-PUBBLICO` nel browser.
        7. Seleziona l'istanza → **Connect** → **EC2 Instance Connect** → **Connect**. Nel terminale:
           `systemctl status nginx` e `sudo cat /var/log/cloud-init-output.log | tail`.

    === "Google Cloud"

        1. Console → **Compute Engine** → **Istanze VM** → **Crea istanza** (abilita l'API se richiesto).
        2. Nome `lab-web`; regione del corso; tipo **e2-micro**.
        3. Disco di avvio: **Ubuntu 24.04 LTS** (x86/64).
        4. Firewall: seleziona **Consenti traffico HTTP**.
        5. Avanzate → Gestione → **Automazione**: incolla `avvio.sh` come script di avvio.
        6. **Crea**. Dopo un minuto apri `http://IP-ESTERNO`.
        7. Pulsante **SSH** nella riga dell'istanza (si apre un terminale nel browser). Verifica
           `systemctl status nginx` e `sudo journalctl -u google-startup-scripts`.

    **Pulizia.** AWS: istanza → **Instance state** → **Terminate**, poi elimina il security group. Google Cloud:
    seleziona l'istanza → **Elimina**. Controlla che non restino dischi: nella console di EC2 (**Volumes**) o di
    Compute Engine (**Dischi**).

??? success "Verifica e spiegazione"
    La pagina mostra `Servito da` seguito dal nome interno della VM. Se non risponde:

    - dopo pochi secondi è normale: lo script sta ancora installando Nginx;
    - controlla il log dello script (`cloud-init-output.log` su AWS, il journal di `google-startup-scripts` su Google
      Cloud);
    - controlla le regole di firewall: la porta 80 deve essere aperta verso il tuo IP o, su Google Cloud, la VM deve
      avere il tag `http-server`;
    - se EC2 Instance Connect non si collega, la regola SSH non ammette gli indirizzi del servizio: usa la prefix list
      indicata al punto 4, non il tuo IP.

    Hai usato lo stesso metodo di diagnosi della Fase 1: servizio, porta in ascolto, firewall. Nel cloud cambia solo
    dove si trova il firewall.

## Esercizi

!!! esercizio "Esercizio 1 · Scegli il modello di prezzo"
    1. Un sito aziendale gira su due VM accese tutto l'anno, da tre anni, senza variazioni previste.
    2. Una pipeline di CI esegue test per qualche ora al giorno; se una macchina viene interrotta, il lavoro riparte.
    3. Uno sviluppatore prova una nuova versione per due giorni.
    4. Un software con licenza legata al server fisico.

??? success "Soluzione"
    1. **Impegno a 1 o 3 anni** (Savings Plans o Reserved Instances su AWS, Committed Use Discounts su Google Cloud).
    2. **Spot**: lavoro interrompibile, grande risparmio.
    3. **On demand**: breve e imprevedibile.
    4. **Server dedicati** (Dedicated Hosts su AWS, sole-tenant node su Google Cloud).

!!! esercizio "Esercizio 2 · VM, container o funzione?"
    Quale opzione di calcolo sceglieresti per: la Tasks API in container con traffico molto variabile e notti senza
    richieste; un vecchio software che richiede un kernel modificato; uno script che ridimensiona ogni foto caricata
    in un bucket; un'applicazione composta da 30 microservizi gestita da un team esperto di Kubernetes?

??? success "Soluzione"
    Tasks API: **Cloud Run** o **ECS con Fargate** (container senza server, scalano; Cloud Run anche fino a zero).
    Software con kernel modificato: **VM** (EC2, Compute Engine), l'unica opzione che dà il controllo del sistema
    operativo. Ridimensionamento foto: **funzione** (Lambda, Cloud Run functions) attivata dall'evento di caricamento.
    30 microservizi con un team Kubernetes: **EKS** o **GKE**.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Cosa indica la `g` in un tipo di istanza come `t4g.micro`?","options":["GPU","Processore AWS Graviton (ARM)","Generazione","Uso generale"],"answer":[1],"explanation":"Le istanze Graviton hanno in genere un miglior rapporto prezzo/prestazioni."},
 {"q":"Quale modello di prezzo è adatto a un'elaborazione batch che può essere interrotta?","options":["On demand","Spot","Reserved Instances","Dedicated Hosts"],"answer":[1],"explanation":"Fino a circa il 90% di sconto, ma la capacità può essere ripresa dal fornitore."},
 {"q":"Come si configura in modo ripetibile una VM al primo avvio?","options":["A mano via SSH","Con user data (AWS) o startup script (Google Cloud)","Con un'email al supporto","Non è possibile"],"answer":[1],"explanation":"Su Linux li esegue in genere cloud-init."},
 {"q":"Quale servizio esegue container senza gestire server, scalando anche a zero?","options":["Compute Engine","Cloud Run","EC2","Lightsail"],"answer":[1],"explanation":"Su AWS l'equivalente più vicino è ECS con Fargate (o App Runner)."},
 {"q":"Quali servizi sono Kubernetes gestito? (scegli 2)","options":["EKS","GKE","Lambda","App Engine","Batch"],"answer":[0,1],"explanation":"Amazon EKS e Google Kubernetes Engine."},
 {"q":"Dopo un laboratorio hai terminato la VM. Cosa controlli ancora?","options":["Nulla","Che non restino dischi, IP pubblici riservati o altre risorse a pagamento","Che il budget sia cancellato","Che l'account sia chiuso"],"answer":[1],"explanation":"Alcune risorse sopravvivono alla VM e continuano a costare."}
]}
</script>

## Fonti

- [Amazon EC2 User Guide: What is Amazon EC2?](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/concepts.html)
- [Amazon EC2: Instance type naming conventions](https://docs.aws.amazon.com/ec2/latest/instancetypes/instance-type-names.html)
- [Amazon EC2: Run commands at launch (user data)](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html)
- [Amazon EC2 pricing](https://aws.amazon.com/ec2/pricing/) e [Savings Plans](https://aws.amazon.com/savingsplans/)
- [Compute Engine: Machine families resource and comparison guide](https://docs.cloud.google.com/compute/docs/machine-resource)
- [Compute Engine: Startup scripts](https://docs.cloud.google.com/compute/docs/instances/startup-scripts/linux)
- [Compute Engine: Spot VMs](https://docs.cloud.google.com/compute/docs/instances/spot) e [Committed use discounts](https://docs.cloud.google.com/compute/docs/instances/committed-use-discounts-overview)
- [cloud-init documentation](https://docs.cloud-init.io/en/latest/)

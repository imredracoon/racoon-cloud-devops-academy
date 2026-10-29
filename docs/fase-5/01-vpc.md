# 5.1 Reti private e pubbliche

!!! obiettivi "Obiettivi"
    - Progettare una rete a tre livelli (pubblico, applicazione, dati) su due o più zone.
    - Scegliere tra un NAT per zona e un NAT condiviso, conoscendone costi e rischi.
    - Raggiungere i servizi del fornitore senza passare da Internet con gli endpoint privati.
    - Amministrare le VM senza aprire porte né usare bastion host.
    - Registrare il traffico con i flow log e collegare più reti tra loro.

## La rete a tre livelli

Il modello più diffuso divide la rete per **livelli**, ognuno in una subnet diversa e replicato in **ogni zona**:

| Livello | Contenuto | Raggiungibile da |
|---|---|---|
| **Pubblico** | load balancer, NAT | Internet |
| **Applicazione** | VM o container dell'applicazione | solo dal load balancer |
| **Dati** | database, cache | solo dall'applicazione |

```mermaid
flowchart TB
  I((Internet)) --> IGW[Internet gateway]
  subgraph VPC[VPC 10.0.0.0/16]
    subgraph A[Zona A]
      PA[Pubblica 10.0.0.0/24<br>ALB · NAT GW] --> AA[App 10.0.10.0/24] --> DA[Dati 10.0.20.0/24<br>DB primario]
    end
    subgraph B[Zona B]
      PB[Pubblica 10.0.1.0/24<br>ALB · NAT GW] --> AB[App 10.0.11.0/24] --> DB2[Dati 10.0.21.0/24<br>DB standby]
    end
  end
  IGW --> PA
  IGW --> PB
```

Il principio è la **difesa in profondità**: anche se un livello viene compromesso, il successivo è protetto da un
altro confine di rete. I firewall seguono la catena: il security group (o la regola firewall) del database accetta
solo dal livello applicazione, che accetta solo dal load balancer.

=== "AWS"

    - Su AWS ogni subnet sta in una sola AZ, quindi servono **sei subnet** per tre livelli su due zone.
    - Le subnet pubbliche condividono una route table con `0.0.0.0/0 → Internet gateway`.
    - Ogni subnet privata ha una route table con `0.0.0.0/0 → NAT gateway` **della sua zona**.
    - Il livello dati può non avere alcuna rotta verso Internet: il database non ne ha bisogno.

=== "Google Cloud"

    - Le subnet sono **regionali**: una subnet copre tutte le zone, quindi bastano **una subnet per livello** per regione.
    - La separazione tra livelli si ottiene con le **regole firewall** basate su **service account** o tag di rete
      (per esempio: "verso il service account `db` solo dal service account `app` sulla porta 5432").
    - Le VM senza IP esterno escono con **Cloud NAT**, che è regionale e ridondante per costruzione.
    - Molti servizi gestiti (come Cloud SQL) non stanno nella tua subnet ma in una rete del produttore, raggiunta con
      un **accesso privato ai servizi** (*private services access*) o con **Private Service Connect**.

## NAT: alta disponibilità contro costo

=== "AWS"

    Un **NAT gateway** vive in **una** AZ. Con un solo NAT gateway, se quella zona cade le subnet private delle altre
    zone perdono l'accesso a Internet. Le opzioni:

    | Scelta | Disponibilità | Costo |
    |---|---|---|
    | Un NAT gateway per zona | alta | un costo orario per ogni NAT, più il traffico |
    | Un solo NAT gateway | dipende da una zona; traffico tra zone a pagamento | più basso |
    | Nessun NAT, solo endpoint privati | nessuna uscita verso Internet | dipende dagli endpoint |

    In produzione si usa un NAT per zona; per i laboratori un solo NAT basta. Per IPv6 esiste l'**egress-only Internet
    gateway**, gratuito.

=== "Google Cloud"

    **Cloud NAT** è un servizio software distribuito, regionale: non c'è un'istanza da replicare. Si paga per ora e per
    i GB elaborati.

## Raggiungere i servizi senza Internet

Le VM private che usano S3 o Cloud Storage passerebbero dal NAT, pagando il traffico e uscendo verso Internet. Gli
**endpoint privati** evitano entrambe le cose.

=== "AWS"

    - **Gateway endpoint** per **S3** e **DynamoDB**: una voce nella route table, **gratuito**. Da aggiungere sempre.
    - **Interface endpoint** (**AWS PrivateLink**) per quasi tutti gli altri servizi (ECR, Secrets Manager, CloudWatch,
      SSM...): un'interfaccia di rete con IP privato nella tua subnet, a pagamento orario.
    - Le **endpoint policy** limitano cosa si può fare attraverso l'endpoint (per esempio solo un certo bucket).

=== "Google Cloud"

    - **Private Google Access** su una subnet: le VM senza IP esterno raggiungono le API Google (Cloud Storage, Artifact
      Registry...) senza NAT. Gratuito.
    - **Private Service Connect**: endpoint privati verso le API Google o verso servizi pubblicati da altri.
    - **VPC Service Controls**: un perimetro intorno ai servizi gestiti, che impedisce l'esfiltrazione dei dati verso
      progetti esterni.

## Amministrare senza aprire porte

Il vecchio schema prevedeva un **bastion host**: una VM pubblica con SSH aperto, da cui saltare verso le VM private.
È un bersaglio da proteggere e aggiornare. Oggi si preferisce:

- **AWS Systems Manager Session Manager**: una shell nel browser o nella CLI, **senza porte in ingresso** e senza
  chiavi SSH, con permessi IAM e registrazione delle sessioni. L'istanza ha bisogno dell'agente SSM (preinstallato su
  Amazon Linux e Ubuntu) e di un ruolo con la policy `AmazonSSMManagedInstanceCore`;
- **EC2 Instance Connect Endpoint**: SSH verso istanze private senza IP pubblico;
- **Identity-Aware Proxy (IAP)** su Google Cloud: tunnel SSH verificato con IAM, sorgente `35.235.240.0/20`.

## Osservare il traffico: i flow log

I **flow log** registrano i metadati delle connessioni (sorgente, destinazione, porte, byte, accettata o rifiutata):
**VPC Flow Logs** su entrambi i cloud. Si inviano a CloudWatch Logs o S3 (AWS), a Cloud Logging (Google Cloud). Sono
fondamentali per la diagnosi ("perché la connessione è rifiutata?") e per le indagini di sicurezza. Non registrano il
contenuto dei pacchetti.

## Collegare più reti

| Esigenza | AWS | Google Cloud |
|---|---|---|
| Collegare due reti, uno a uno | **VPC peering** (non transitivo: A↔B e B↔C non danno A↔C) | **VPC Network Peering** (non transitivo) |
| Collegare molte reti e la sede, a stella | **Transit Gateway** | **Network Connectivity Center** |
| Una rete condivisa da più progetti/account | **condivisione di subnet con AWS RAM** | **Shared VPC** |
| Esporre un servizio a un'altra rete senza unire le reti | **PrivateLink** (endpoint service) | **Private Service Connect** |

Gli intervalli di indirizzi delle reti collegate **non devono sovrapporsi**: un motivo in più per pianificarli bene
([lezione 4.7](../fase-4/07-rete.md)).

## Laboratorio

!!! laboratorio "Laboratorio · Rete a tre livelli"
    === "AWS"

        1. Console **VPC** → **Create VPC** → **VPC and more**. Nome `tasks`, CIDR `10.0.0.0/16`, **2 AZ**,
           **2 subnet pubbliche**, **4 subnet private**, **NAT gateway: In 1 AZ** (per contenere i costi),
           **VPC endpoints: S3 Gateway**. Osserva l'anteprima: route table e collegamenti.
        2. Rinomina le subnet private in `app-a`, `app-b`, `dati-a`, `dati-b`. Rimuovi dalle route table delle subnet
           `dati-*` la rotta verso il NAT (crea per loro una route table senza `0.0.0.0/0`).
        3. Crea tre security group: `alb-sg` (80 e 443 da `0.0.0.0/0`), `app-sg` (8000 solo da `alb-sg`), `db-sg`
           (5432 solo da `app-sg`).
        4. Attiva i **flow log** della VPC verso CloudWatch Logs (gruppo `/academy/vpc-flow`), solo traffico
           **rifiutato**.
        5. Crea un ruolo IAM per EC2 con `AmazonSSMManagedInstanceCore`, avvia un'istanza `t3.micro` Ubuntu in `app-a`
           con quel ruolo e **senza** IP pubblico né porte aperte in ingresso. Collegati con **Session Manager**
           (EC2 → Connect → Session Manager) e verifica `curl -sI https://aws.amazon.com` (esce dal NAT).

        **Pulizia.** Termina l'istanza; elimina il NAT gateway e **rilascia l'Elastic IP** associato; poi elimina la
        VPC (Delete VPC rimuove subnet, route table, endpoint e Internet gateway). Elimina il gruppo di log.

    === "Google Cloud"

        1. **Rete VPC** → **Crea rete VPC** `tasks`, modalità personalizzata, con tre subnet nella regione del corso:
           `pubblica` `10.0.0.0/24`, `app` `10.0.10.0/24`, `dati` `10.0.20.0/24`. Su `app` e `dati` attiva
           **Accesso privato Google** e i **log di flusso**.
        2. Crea due service account: `tasks-app` e `tasks-db-client`.
        3. Regole firewall: `allow-iap-ssh` (TCP 22 da `35.235.240.0/20` verso il service account `tasks-app`);
           `allow-lb-to-app` (TCP 8000 dagli intervalli dei controlli di integrità e dei proxy di Google,
           `130.211.0.0/22` e `35.191.0.0/16`, verso `tasks-app`).
        4. **Cloud NAT**: crea un gateway per la rete `tasks` nella regione, con un nuovo Cloud Router, solo per la
           subnet `app`.
        5. Avvia una VM `e2-micro` nella subnet `app`, con service account `tasks-app` e **senza IP esterno**.
           Collegati con `gcloud compute ssh VM --tunnel-through-iap` e verifica `curl -sI https://cloud.google.com`.

        **Pulizia.** Elimina la VM, il gateway Cloud NAT e il router, le regole firewall, le subnet, la rete, i
        service account.

??? success "Spiegazione"
    Hai costruito la rete del progetto di fine fase. L'istanza nel livello applicazione non è raggiungibile da Internet
    e non ha porte in ingresso aperte, eppure puoi amministrarla (Session Manager o IAP) e lei può uscire verso Internet
    (NAT). Il livello dati non ha alcuna via verso Internet. Nei flow log dei rifiuti troverai traffico che bussa
    continuamente alle porte: è il rumore di fondo di Internet.

## Esercizi

!!! esercizio "Esercizio · Rivedi l'architettura"
    Un collega propone: una VPC con due subnet pubbliche (una per zona); in ognuna una VM con l'applicazione, il
    database PostgreSQL installato sulla stessa VM, SSH aperto a `0.0.0.0/0` per comodità, e l'applicazione che scrive
    su S3 passando per Internet. Elenca i problemi e proponi correzioni.

??? success "Soluzione"
    1. **Applicazione e database in subnet pubbliche**: vanno in subnet private, dietro un load balancer.
    2. **Database sulla stessa VM** dell'applicazione: due database separati e non sincronizzati, nessun backup gestito.
       Usa un database gestito con alta disponibilità ([lezione 5.4](04-database.md)).
    3. **SSH aperto a tutti**: usa Session Manager / IAP, senza porte aperte.
    4. **S3 via Internet**: aggiungi il gateway endpoint S3 (gratuito), con una endpoint policy.
    5. Mancano i security group a catena e i flow log per la diagnosi.

## Nei tre livelli

=== "Livello 1 · Un PC"

    Il laboratorio è tutto ciò che serve.

=== "Livello 2 · Workstation + server"

    Colleghi `srv1` alla VPC `tasks` del laboratorio con una VPN **WireGuard** verso una piccola istanza nel cloud.
    I fornitori offrono anche VPN gestite (*AWS Site-to-Site VPN*, *Cloud VPN*), ma costano decine di euro al mese
    anche senza traffico; un'istanza WireGuard costa pochi centesimi l'ora e ti fa vedere ogni pezzo.

    ```mermaid
    flowchart LR
      subgraph Casa
        S[srv1<br>tunnel 10.99.0.2]
      end
      S == "WireGuard UDP 51820<br>(attraverso il router di casa)" ==> G
      subgraph VPC[VPC tasks 10.0.0.0/16]
        G[istanza vpn<br>subnet pubblica<br>tunnel 10.99.0.1] --> A[istanza app<br>subnet privata]
      end
    ```

    Gli intervalli non si sovrappongono: casa `192.168.1.0/24`, VM `192.168.56.0/24`, tunnel `10.99.0.0/24`,
    VPC `10.0.0.0/16`. È proprio il motivo per cui si pianificano gli indirizzi.

    **1. L'istanza VPN.** Nella VPC del laboratorio crea un'istanza Ubuntu 24.04 piccola (`t3.micro`, `e2-micro`)
    nella subnet pubblica, con IP pubblico, e apri in ingresso **solo UDP 51820 dall'IP pubblico di casa**
    (lo trovi con `curl -s https://ifconfig.me` da `srv1`). Poi:

    === "AWS"

        - nell'istanza: **Actions → Networking → Change source/destination check → Stop**: deve inoltrare pacchetti
          non indirizzati a lei;
        - nelle route table delle subnet private: rotta `10.99.0.0/24` → l'istanza `vpn` (target *Instance*);
        - nel security group `app-sg`: ICMP e TCP 22 da `10.99.0.0/24`, per le prove.
        - nel security group dell'istanza `vpn`: oltre a UDP 51820, tutto il traffico dalla VPC (`10.0.0.0/16`), che
          lei inoltra nel tunnel.

    === "Google Cloud"

        - crea l'istanza con `--can-ip-forward` (si può impostare solo alla creazione);
        - rotta: `gcloud compute routes create casa --network=tasks --destination-range=10.99.0.0/24 --next-hop-instance=vpn --next-hop-instance-zone=ZONA`;
        - regola firewall: ICMP e TCP 22 da `10.99.0.0/24` verso le istanze della subnet `app`.
        - regola firewall: tutto il traffico da `10.0.0.0/16` verso l'istanza `vpn`, che lo inoltra nel tunnel.

    **2. WireGuard ai due capi.** Su entrambe le macchine installa `wireguard-tools` (`apt` o `dnf`) e genera le
    chiavi:

    ```bash
    umask 077
    wg genkey | sudo tee /etc/wireguard/privata | wg pubkey | sudo tee /etc/wireguard/pubblica
    ```

    Sull'istanza nel cloud attiva anche l'inoltro (`net.ipv4.ip_forward=1` in `/etc/sysctl.d/`, come nella
    [lezione 1.8](../fase-1/08-dns-routing.md)).

    ```ini title="/etc/wireguard/wg0.conf sull'istanza vpn"
    [Interface]
    Address = 10.99.0.1/24
    ListenPort = 51820
    PrivateKey = CHIAVE-PRIVATA-DELL-ISTANZA

    [Peer]
    # srv1, a casa
    PublicKey = CHIAVE-PUBBLICA-DI-SRV1
    AllowedIPs = 10.99.0.2/32
    ```

    ```ini title="/etc/wireguard/wg0.conf su srv1"
    [Interface]
    Address = 10.99.0.2/24
    PrivateKey = CHIAVE-PRIVATA-DI-SRV1

    [Peer]
    # l'istanza vpn nel cloud
    PublicKey = CHIAVE-PUBBLICA-DELL-ISTANZA
    Endpoint = IP-PUBBLICO-DELL-ISTANZA:51820
    AllowedIPs = 10.99.0.0/24, 10.0.0.0/16
    PersistentKeepalive = 25
    ```

    ```bash
    sudo systemctl enable --now wg-quick@wg0      # su entrambe le macchine
    ```

    `AllowedIPs` fa due cose: dice quali indirizzi possono arrivare da quel peer e quali rotte creare verso di lui.
    `PersistentKeepalive` tiene aperto il passaggio nel NAT del router di casa: è `srv1` a iniziare il collegamento,
    quindi sul router non serve aprire porte.

    **Collaudo**

    | Prova | Risultato atteso |
    |---|---|
    | `sudo wg show` su `srv1` | `latest handshake` di pochi secondi fa, byte ricevuti e inviati |
    | `ip route get 10.0.10.5` su `srv1` | `dev wg0` |
    | `ping -c 3 IP-PRIVATO-ISTANZA-APP` da `srv1` | risposte da un'istanza **senza IP pubblico** |
    | `tracepath -n IP-PRIVATO-ISTANZA-APP` da `srv1` | due salti: `10.99.0.1`, poi l'istanza |
    | sull'istanza app, `ping -c 3 10.99.0.2` | risponde `srv1`: il collegamento funziona nei due sensi |
    | `sudo systemctl stop wg-quick@wg0` su `srv1`, di nuovo il `ping` | nessuna risposta: il percorso privato passa solo dal tunnel |

    **Pulizia:** termina l'istanza `vpn`, cancella la rotta e le regole aggiunte; su `srv1`
    `sudo systemctl disable --now wg-quick@wg0`. Conserva i file di configurazione, senza le chiavi private, nel
    repository `academy-lab`: li riuserai nel progetto.

=== "Livello 3 · Mini data center"

    Come nel livello 2, con `nodo1` al posto di `srv1`: il tunnel collega l'**intera rete del laboratorio**.

    Differenze rispetto al livello 2:

    - nella configurazione dell'istanza `vpn`, il peer di casa ha `AllowedIPs = 10.99.0.2/32, 10.10.10.0/24`;
    - nel cloud la rotta verso casa comprende anche `10.10.10.0/24` (un'altra rotta verso l'istanza `vpn`), e le
      regole di `app-sg` o del firewall ammettono anche `10.10.10.0/24`;
    - su `nodo1` il firewall della [lezione 1.11](../fase-1/11-firewall.md#nei-tre-livelli) deve far passare il traffico
      tra laboratorio e tunnel, **senza** NAT: nella catena `forward` aggiungi

        ```nft
        iifname $LAB oifname "wg0" accept
        iifname "wg0" oifname $LAB ip saddr 10.0.0.0/16 tcp dport { 22, 80, 443 } accept
        ```

      Il masquerade resta solo sull'uscita `$WAN`: i pacchetti nel tunnel mantengono gli indirizzi `10.10.10.x`, e
      il cloud sa come restituirli grazie alla rotta.

    `nodo2` e `nodo3` non vanno toccati: il loro gateway è `nodo1`, che ora conosce anche la strada per la VPC.

    **Collaudo:** le prove del livello 2 eseguite da `nodo1`, più queste:

    | Prova | Risultato atteso |
    |---|---|
    | `ssh nodo2 ping -c 3 IP-PRIVATO-ISTANZA-APP` | risponde: tutto il laboratorio raggiunge la VPC |
    | `ssh nodo3 tracepath -n IP-PRIVATO-ISTANZA-APP` | `10.10.10.1`, `10.99.0.1`, l'istanza |
    | sull'istanza app, `curl -sI http://10.10.10.12/` | la risposta del Nginx di `nodo2`: il cloud raggiunge la sede |
    | sull'istanza app, `nc -zvw3 10.10.10.12 5432` | timeout: `nodo1` lascia entrare dal cloud solo le porte permesse |

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quante subnet servono su AWS per tre livelli su due zone di disponibilità?","options":["3","2","6","1"],"answer":[2],"explanation":"Su AWS una subnet sta in una sola AZ."},
 {"q":"Qual è il modo gratuito per far raggiungere S3 da istanze private senza NAT?","options":["Interface endpoint","Gateway endpoint per S3","Direct Connect","Transit Gateway"],"answer":[1],"explanation":"I gateway endpoint esistono per S3 e DynamoDB."},
 {"q":"Come amministri un'istanza EC2 privata senza aprire porte in ingresso?","options":["Con un bastion host","Con AWS Systems Manager Session Manager","Aprendo SSH a 0.0.0.0/0","Con un NAT gateway"],"answer":[1],"explanation":"Serve l'agente SSM e un ruolo con AmazonSSMManagedInstanceCore."},
 {"q":"Le VPC A e B sono in peering, e anche B e C. A può raggiungere C?","options":["Sì, il peering è transitivo","No, il peering non è transitivo","Solo con IPv6","Solo nella stessa regione"],"answer":[1],"explanation":"Per molte reti si usa Transit Gateway o Network Connectivity Center."},
 {"q":"Cosa registrano i VPC Flow Logs?","options":["Il contenuto dei pacchetti","I metadati delle connessioni: indirizzi, porte, byte, esito","Le password","Le query del database"],"answer":[1],"explanation":"Servono per diagnosi e indagini di sicurezza."},
 {"q":"Su Google Cloud, come separi i livelli applicazione e dati se stanno nella stessa regione?","options":["Con subnet zonali","Con regole firewall basate su service account o tag","Non è possibile","Con il peering"],"answer":[1],"explanation":"Anche con subnet diverse, è il firewall a stabilire chi parla con chi."}
]}
</script>

## Fonti

- [Amazon VPC: Example – VPC with servers in private subnets and NAT](https://docs.aws.amazon.com/vpc/latest/userguide/vpc-example-private-subnets-nat.html)
- [Amazon VPC: Gateway endpoints](https://docs.aws.amazon.com/vpc/latest/privatelink/gateway-endpoints.html)
- [AWS Systems Manager: Session Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager.html)
- [Amazon VPC: Flow logs](https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html)
- [Amazon VPC: Transit gateways](https://docs.aws.amazon.com/vpc/latest/tgw/what-is-transit-gateway.html)
- [Google Cloud: Best practices and reference architectures for VPC design](https://docs.cloud.google.com/architecture/best-practices-vpc-design)
- [Google Cloud: Private Google Access](https://docs.cloud.google.com/vpc/docs/private-google-access)
- [Google Cloud: Firewall rules for health checks](https://docs.cloud.google.com/load-balancing/docs/health-check-concepts#ip-ranges)

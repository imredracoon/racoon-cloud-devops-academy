# 4. Server di infrastruttura

Sono i server che **fanno funzionare la rete** e che gli utenti non vedono: se lavorano bene nessuno se ne accorge,
se si fermano si ferma tutto. Per ognuno trovi cosa fa, l'hardware tipico e il software su Ubuntu e su RHEL.

!!! abstract "Già visto nel corso"
    - **DNS, DHCP, routing e NAT** come concetti: [lezione 1.8](../../fase-1/08-dns-routing.md).
    - Installare pacchetti e attivare servizi: [lezione 1.5](../../fase-1/05-pacchetti.md) e [lezione 1.4](../../fase-1/04-processi-systemd.md#systemctl).
    - Aprire le porte nel firewall: [lezione 1.11](../../fase-1/11-firewall.md#i-gestori-delle-distribuzioni-ufw-e-firewalld).
    - Sincronizzazione dell'ora con chrony: [lezione 1.15](../../fase-1/15-avvio-tempo.md#data-ora-e-sincronizzazione).
    - SSH, chiavi e bastion host: [lezione 1.10](../../fase-1/10-ssh.md).
    - TLS e certificati: [lezione 1.12](../../fase-1/12-http-nginx.md#tls-e-https).

## Il quadro d'insieme

| Server | Cosa fa | Software su Ubuntu | Software su RHEL |
|---|---|---|---|
| [DNS](#dns) | traduce nomi in indirizzi | BIND 9, Unbound | BIND, Unbound |
| [DHCP](#dhcp) | assegna gli indirizzi IP | Kea | Kea |
| [NTP](#ntp) | distribuisce l'ora esatta | chrony | chrony |
| [Directory e identità](#directory-e-identita) | utenti, gruppi, login unico | Samba AD, OpenLDAP, SSSD | IdM (FreeIPA), SSSD |
| [Autorità di certificazione](#autorita-di-certificazione-interna) | firma i certificati interni | OpenSSL | CA integrata in IdM, OpenSSL |
| [Provisioning e PXE](#provisioning-e-pxe) | installa i server via rete | MAAS | Satellite / Foreman, Kickstart |
| [Proxy](#proxy) | intermedia l'accesso al web | Squid | Squid |
| [VPN](#vpn) | collega utenti e sedi in modo cifrato | WireGuard, OpenVPN | IPsec con Libreswan |
| [Router e firewall](#router-e-firewall) | instrada e filtra | nftables, ufw | nftables, firewalld |
| [Bastion](#bastion) | unico punto d'ingresso amministrativo | OpenSSH | OpenSSH |
| [Posta](#posta-elettronica) | invia e riceve email | Postfix, Dovecot | Postfix, Dovecot |
| [Stampa](#stampa) | gestisce le code di stampa | CUPS | CUPS |
| [VoIP e centralino](#voip-e-centralino) | telefonia su rete IP | Asterisk | Asterisk (da sorgente o container) |

!!! tip "Due è meglio di uno"
    I servizi di infrastruttura richiedono poco hardware ma devono essere **sempre** disponibili. La regola è
    averne almeno **due istanze** su macchine diverse (DNS primario e secondario, due server DHCP in failover, due
    controller di dominio). Spesso sono macchine virtuali, non server fisici dedicati.

## DNS

Risponde alle domande "che indirizzo ha questo nome?" per la rete interna (server autoritativo per i domini
dell'azienda) e fa da **resolver** per i client, con la cache. Il funzionamento l'hai visto nella
[lezione 1.8](../../fase-1/08-dns-routing.md#dns).

**Hardware tipico**: 2-4 core, 4-8 GB di RAM, due piccoli SSD in RAID 1. Due istanze.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y bind9 bind9-utils     # server autoritativo e resolver
    # zone in /etc/bind/named.conf.local, file di zona in /etc/bind/
    sudo named-checkconf && sudo systemctl reload named
    # solo resolver con cache, più semplice e sicuro:
    sudo apt install -y unbound
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y bind bind-utils
    # configurazione in /etc/named.conf, file di zona in /var/named/
    sudo named-checkconf && sudo systemctl enable --now named
    sudo firewall-cmd --permanent --add-service=dns && sudo firewall-cmd --reload
    # solo resolver con cache:
    sudo dnf install -y unbound
    ```

## DHCP

Assegna a ogni dispositivo indirizzo, maschera, gateway e DNS ([lezione 1.8](../../fase-1/08-dns-routing.md)).
Il vecchio ISC DHCP non è più sviluppato: il suo successore è **Kea**, dello stesso consorzio ISC, adottato da
entrambe le famiglie. Kea si configura in JSON e supporta il failover tra due server.

**Hardware tipico**: minimo, spesso una VM o lo stesso server del DNS.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y kea
    sudoedit /etc/kea/kea-dhcp4.conf          # subnet, intervallo, opzioni
    sudo systemctl restart kea-dhcp4-server
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y kea
    sudoedit /etc/kea/kea-dhcp4.conf
    sudo systemctl enable --now kea-dhcp4
    sudo firewall-cmd --permanent --add-service=dhcp && sudo firewall-cmd --reload
    ```

Una configurazione minima per una rete `192.168.10.0/24`:

```json
{
  "Dhcp4": {
    "interfaces-config": { "interfaces": [ "eno1" ] },
    "lease-database": { "type": "memfile", "persist": true },
    "valid-lifetime": 3600,
    "subnet4": [{
      "id": 1,
      "subnet": "192.168.10.0/24",
      "pools": [ { "pool": "192.168.10.100 - 192.168.10.200" } ],
      "option-data": [
        { "name": "routers", "data": "192.168.10.1" },
        { "name": "domain-name-servers", "data": "192.168.10.2, 192.168.10.3" }
      ]
    }]
  }
}
```

## NTP

Tiene allineati gli orologi di tutti i server. Senza un'ora coerente i log non si confrontano, Kerberos rifiuta
i login e i certificati sembrano scaduti. Il client **chrony** l'hai già configurato nella
[lezione 1.15](../../fase-1/15-avvio-tempo.md#data-ora-e-sincronizzazione); per farne un **server** basta
permettere ai client di interrogarlo. Nei data center più esigenti l'ora arriva da un ricevitore GPS o da
**PTP** (*Precision Time Protocol*), più preciso di NTP.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y chrony
    echo "allow 192.168.10.0/24" | sudo tee /etc/chrony/conf.d/server.conf
    sudo systemctl restart chrony
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y chrony
    echo "allow 192.168.10.0/24" | sudo tee -a /etc/chrony.conf
    sudo systemctl restart chronyd
    sudo firewall-cmd --permanent --add-service=ntp && sudo firewall-cmd --reload
    ```

## Directory e identità

Un **servizio di directory** conserva in un solo posto utenti, gruppi, computer e regole: l'utente ha un solo
account per tutti i server e i servizi (*single sign-on*), e chi se ne va viene disattivato una volta sola. Si
basa su **LDAP** (il database) e **Kerberos** (l'autenticazione). Nel mondo Windows è **Active Directory**; nel
cloud è l'equivalente di IAM per le persone ([lezione 4.4](../../fase-4/04-iam.md)).

**Hardware tipico**: 2-4 core, 8-16 GB di RAM. Almeno due controller replicati.

=== "Debian/Ubuntu"

    Ubuntu offre due strade:

    - **Samba come controller di dominio Active Directory**, compatibile con i client Windows:
      `sudo apt install -y samba-ad-dc krb5-user bind9-dnsutils`, poi `sudo samba-tool domain provision` e
      l'attivazione del servizio `samba-ad-dc`.
    - **OpenLDAP** (`slapd`) per una directory LDAP pura, spesso con Kerberos.

    I server si uniscono al dominio con **SSSD** (`sssd`, `realmd`): `sudo realm join esempio.lan`.

=== "RHEL (Rocky/AlmaLinux)"

    Red Hat propone **Identity Management (IdM)**, la versione supportata del progetto **FreeIPA**: unisce LDAP,
    Kerberos, DNS, una CA e le regole di accesso e `sudo` in un unico prodotto con interfaccia web.

    ```bash
    sudo dnf install -y ipa-server ipa-server-dns
    sudo ipa-server-install            # procedura guidata
    # sui client:
    sudo dnf install -y ipa-client && sudo ipa-client-install
    ```

    RHEL non supporta Samba come controller di dominio AD: per integrarsi con un AD esistente si usa una
    relazione di fiducia (*trust*) tra IdM e AD, oppure SSSD direttamente.

## Autorità di certificazione interna

Per usare HTTPS anche sui servizi interni ([lezione 1.12](../../fase-1/12-http-nginx.md#tls-e-https)) serve una
**CA** propria: firma i certificati dei server, e il suo certificato radice si installa su tutti i client, che da
quel momento si fidano dei certificati che firma. In produzione si usa una CA con interfaccia e revoca (quella di
IdM su RHEL, o progetti come step-ca e il protocollo **ACME** reso famoso da Let's Encrypt).

=== "Debian/Ubuntu"

    ```bash
    # aggiungere la radice della CA interna all'archivio di fiducia del sistema
    sudo cp ca-interna.crt /usr/local/share/ca-certificates/
    sudo update-ca-certificates
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo cp ca-interna.crt /etc/pki/ca-trust/source/anchors/
    sudo update-ca-trust
    ```

## Provisioning e PXE

Installare a mano cento server non è possibile. Il **PXE** (*Preboot eXecution Environment*) permette a un server
appena acceso di avviarsi **dalla rete**: riceve l'indirizzo dal DHCP, scarica un programma di avvio via TFTP o
HTTP e lancia l'installazione automatica. Un **server di provisioning** coordina tutto: inventario delle macchine,
installazione, configurazione iniziale. Lo vediamo all'opera nella [pagina 7](07-installazione.md#in-automatico-dalla-rete).

=== "Debian/Ubuntu"

    **MAAS** (*Metal as a Service*) di Canonical trasforma i server fisici in un "cloud privato di macchine":
    li scopre, li accende tramite il BMC, li installa con Ubuntu (o altri sistemi) e li restituisce quando non
    servono. Si installa come snap: `sudo snap install maas`.

=== "RHEL (Rocky/AlmaLinux)"

    **Red Hat Satellite**, basato sul progetto open source **Foreman**, gestisce installazione via PXE con
    file **Kickstart**, aggiornamenti e conformità di migliaia di sistemi RHEL. Senza Satellite si costruisce
    un server PXE con Kea, un server TFTP o HTTP e i file Kickstart.

## Proxy

Un **proxy di inoltro** (*forward proxy*) sta tra gli utenti interni e Internet: filtra i siti, registra gli
accessi e tiene in cache i contenuti. È diverso dal **reverse proxy** che sta davanti ai server web, visto nella
[lezione 1.12](../../fase-1/12-http-nginx.md#reverse-proxy). Il software di riferimento è **Squid**, che ascolta
sulla porta 3128.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y squid
    sudoedit /etc/squid/squid.conf      # acl e http_access per le reti ammesse
    sudo systemctl restart squid
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y squid
    sudoedit /etc/squid/squid.conf
    sudo systemctl enable --now squid
    sudo firewall-cmd --permanent --add-service=squid && sudo firewall-cmd --reload
    ```

## VPN

Una **VPN** crea un tunnel cifrato su Internet: collega i dipendenti in remoto alla rete aziendale (*accesso
remoto*) o due sedi tra loro (*site-to-site*). Nel cloud collega la rete aziendale a una VPC
([lezione 5.1](../../fase-5/01-vpc.md#collegare-piu-reti)).

=== "Debian/Ubuntu"

    **WireGuard** è moderno, veloce e integrato nel kernel; la documentazione di Ubuntu ha guide per ogni
    scenario. In alternativa **OpenVPN**.

    ```bash
    sudo apt install -y wireguard-tools
    wg genkey | tee privata.key | wg pubkey > pubblica.key
    sudoedit /etc/wireguard/wg0.conf
    sudo systemctl enable --now wg-quick@wg0
    ```

=== "RHEL (Rocky/AlmaLinux)"

    La VPN supportata da Red Hat è **IPsec** con **Libreswan**, compatibile con i firewall di tutte le marche.
    WireGuard è presente nel kernel ma Red Hat lo ha a lungo classificato come anteprima tecnologica: controlla
    le note di rilascio della tua versione prima di usarlo in produzione.

    ```bash
    sudo dnf install -y libreswan
    sudo ipsec initnss
    sudoedit /etc/ipsec.d/sede-b.conf
    sudo systemctl enable --now ipsec
    sudo firewall-cmd --permanent --add-service=ipsec && sudo firewall-cmd --reload
    ```

## Router e firewall

Un server Linux con due o più schede di rete può fare da router e firewall per una rete intera. Hai già visto
tutti i pezzi: inoltro dei pacchetti con `sysctl` ([lezione 1.15](../../fase-1/15-avvio-tempo.md#a-sistema-avviato-sysctl)),
regole e NAT con nftables ([lezione 1.11](../../fase-1/11-firewall.md#nat-con-nftables)), ufw su Ubuntu e
firewalld su RHEL ([lezione 1.11](../../fase-1/11-firewall.md#i-gestori-delle-distribuzioni-ufw-e-firewalld)).
In azienda si preferiscono apparati dedicati o distribuzioni specializzate come OPNsense
([pagina 1](01-dispositivi-rete.md#firewall-e-ngfw)).

## Bastion

Il **bastion host** è l'unico server raggiungibile dall'esterno per l'amministrazione: da lì si salta verso gli
altri con `ProxyJump`, come nella [lezione 1.10](../../fase-1/10-ssh.md). Deve essere minimale (solo OpenSSH),
con accesso solo a chiave, autenticazione a due fattori e registrazione di tutte le sessioni. Nel cloud spesso
lo si sostituisce con servizi che non aprono porte ([lezione 5.1](../../fase-5/01-vpc.md#amministrare-senza-aprire-porte)).

## Posta elettronica

Un server di posta ha due parti: l'**MTA** (*Mail Transfer Agent*) che invia e riceve tra server con il protocollo
SMTP, e il server **IMAP** da cui gli utenti leggono la posta. Il software standard su entrambe le famiglie è
**Postfix** (MTA) con **Dovecot** (IMAP), più un filtro antispam come Rspamd. Senza i record DNS **MX**, **SPF**,
**DKIM** e **DMARC** ([lezione 1.8](../../fase-1/08-dns-routing.md#dns)) i messaggi finiscono nello spam.

=== "Debian/Ubuntu"

    ```bash
    sudo apt install -y postfix dovecot-imapd
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf install -y postfix dovecot
    sudo systemctl enable --now postfix dovecot
    sudo firewall-cmd --permanent --add-service={smtp,smtps,imaps} && sudo firewall-cmd --reload
    ```

!!! info "Nella pratica"
    Gestire un server di posta esposto su Internet è impegnativo (reputazione degli indirizzi, spam, sicurezza).
    Molte aziende usano servizi in cloud e tengono in casa solo un **relay** Postfix che consegna le email
    generate dai server applicativi.

## Stampa

**CUPS** gestisce stampanti e code di stampa su entrambe le famiglie (pacchetto `cups`) e le condivide in rete
con il protocollo IPP. Per i client Windows si affianca Samba.

## VoIP e centralino

Un **centralino IP** (*IP PBX*) gestisce le chiamate dei telefoni IP (alimentati in PoE,
[pagina 1](01-dispositivi-rete.md#switch)) con il protocollo **SIP** e le collega alla rete telefonica tramite un
operatore. Il software open source più diffuso è **Asterisk**. Ubuntu lo ha nei repository (`asterisk`, nella
sezione *universe*); su RHEL non è incluso e si compila dai sorgenti o si esegue in container. La rete deve dare
priorità al traffico vocale (QoS), spesso su una VLAN dedicata.

## Esercizi

!!! esercizio "Esercizio 1"
    Una piccola azienda di 40 persone ha un solo server fisico che fa DNS, DHCP e controller di dominio. Quali
    rischi vedi e cosa proporresti?

??? success "Soluzione"
    È un **punto unico di guasto**: se si rompe nessuno ottiene un indirizzo, risolve nomi o fa login. Proposta:
    due server (anche piccoli, o due VM su due host diversi) con DNS primario e secondario, Kea in failover e due
    controller di dominio replicati.

!!! esercizio "Esercizio 2"
    Su una VM del laboratorio installa Unbound e verifica che risolva i nomi.

??? success "Soluzione"
    Installa `unbound` con il gestore della tua famiglia, avvialo (`sudo systemctl enable --now unbound`), poi
    `dig @127.0.0.1 example.com`. La seconda volta il tempo di risposta (`Query time`) scende quasi a zero
    perché la risposta arriva dalla cache.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Qual è il server DHCP adottato oggi sia da Ubuntu sia da RHEL al posto di ISC DHCP?","options":["dnsmasq","Kea","BIND","chrony"],"answer":[1],"explanation":"Kea è il successore di ISC DHCP, dello stesso consorzio."},
 {"q":"Quale prodotto Red Hat unisce LDAP, Kerberos, DNS e una CA?","options":["Satellite","Identity Management (FreeIPA)","Cockpit","Libreswan"],"answer":[1],"explanation":"IdM è la versione supportata di FreeIPA."},
 {"q":"Quale software si usa tipicamente come proxy di inoltro con cache?","options":["Squid","Postfix","Kea","CUPS"],"answer":[0],"explanation":"Squid ascolta di solito sulla porta 3128."},
 {"q":"Cosa serve perché un server installi il sistema operativo avviandosi dalla rete?","options":["PXE","Un server DHCP","Un server TFTP o HTTP con i file di avvio","Una chiavetta USB","Un file di installazione automatica (autoinstall o Kickstart)"],"answer":[0,1,2,4],"explanation":"PXE usa DHCP per indirizzo e file di avvio, li scarica in rete e segue il file di installazione automatica; nessun supporto fisico."},
 {"q":"Quali record DNS servono perché le email di un dominio non finiscano nello spam?","options":["MX","SPF","DKIM","DMARC","CNAME"],"answer":[0,1,2,3],"explanation":"MX indica il server di posta; SPF, DKIM e DMARC dimostrano che i messaggi sono autentici. CNAME non c'entra."}
]}
</script>

## Fonti

- [Ubuntu Server: How to install and configure isc-kea](https://documentation.ubuntu.com/server/how-to/networking/install-isc-kea/)
- [Ubuntu Server: How-to guides (DNS, NTP, OpenLDAP, Samba, WireGuard, Squid, posta, CUPS)](https://ubuntu.com/server/docs/how-to/)
- [RHEL 10: Managing networking infrastructure services (BIND, Unbound, Kea)](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html-single/managing_networking_infrastructure_services/index)
- [RHEL 10: Installing Identity Management](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/installing_identity_management/preparing-the-system-for-idm-server-installation)
- [RHEL 10: Configuring the Squid caching proxy server](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/deploying_web_servers_and_reverse_proxies/configuring-the-squid-caching-proxy-server)
- [MAAS documentation](https://canonical.com/maas/docs)
- [ISC Kea](https://www.isc.org/kea/)
- [Postfix](https://www.postfix.org/)
- [Asterisk](https://www.asterisk.org/)

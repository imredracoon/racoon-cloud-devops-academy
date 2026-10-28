# 8. Manutenzione

Un server fisico resta in servizio cinque anni o più. In quel tempo vanno aggiornati sistema e firmware,
sostituiti i pezzi che si guastano e, alla fine, va dismesso in sicurezza. Nel cloud gran parte di questo lavoro
lo fa il provider ([lezione 4.1](../../fase-4/01-cloud-computing.md#il-modello-di-responsabilita-condivisa));
on-premises lo fai tu.

!!! abstract "Già visto nel corso"
    - **Aggiornamenti** e aggiornamenti di sicurezza automatici: [lezione 1.5](../../fase-1/05-pacchetti.md#aggiornamenti-di-sicurezza-automatici).
    - **journalctl** e i log di systemd: [lezione 1.4](../../fase-1/04-processi-systemd.md#journalctl-i-log-di-systemd); journal persistente: [lezione 1.15](../../fase-1/15-avvio-tempo.md#journal-persistente).
    - **Metodo di troubleshooting**: [lezione 1.13](../../fase-1/13-troubleshooting.md).
    - **Svuotare un nodo** prima della manutenzione in Kubernetes: [lezione 8.8](../../fase-8/08-scheduling.md#manutenzione-dei-nodi).
    - **Metriche e alert**: [lezione 7.5](../../fase-7/05-osservabilita.md#alert-che-funzionano) e [lezione 7.6](../../fase-7/06-prometheus.md#regole-di-alert).

## Il calendario

| Quando | Cosa |
|---|---|
| Ogni giorno, in automatico | aggiornamenti di sicurezza, backup, controllo degli alert |
| Ogni settimana | rivedere gli alert dell'hardware e del registro eventi del BMC, verificare l'esito dei backup |
| Ogni mese | finestra di manutenzione: aggiornamenti completi e riavvii pianificati, una prova di ripristino |
| Ogni trimestre o semestre | firmware di server, BMC, controller, dischi e apparati di rete |
| Ogni anno | revisione dell'inventario, delle garanzie e del ciclo di vita; prova del piano di disaster recovery |

Gli interventi che possono fermare un servizio si fanno in una **finestra di manutenzione** concordata, con un
piano scritto, una procedura di ritorno indietro e un **server alla volta**, spostando prima il carico sugli
altri.

## Aggiornare il sistema

=== "Debian/Ubuntu"

    ```bash
    sudo apt update && sudo apt full-upgrade -y
    # serve un riavvio? needrestart indica anche quali servizi riavviare
    [ -f /var/run/reboot-required ] && cat /var/run/reboot-required
    # Livepatch (con Ubuntu Pro): corregge le vulnerabilità del kernel senza riavvio
    sudo pro enable livepatch
    canonical-livepatch status
    ```

    Gli aggiornamenti di sicurezza automatici li fa **unattended-upgrades**, come nella
    [lezione 1.5](../../fase-1/05-pacchetti.md#aggiornamenti-di-sicurezza-automatici).

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    sudo dnf upgrade -y
    sudo dnf needs-restarting -r      # serve un riavvio?
    # kpatch: correzioni del kernel senza riavvio
    sudo dnf install -y kpatch kpatch-dnf
    sudo dnf kpatch auto
    sudo kpatch list
    ```

    Gli aggiornamenti automatici li fa **dnf-automatic**, come nella
    [lezione 1.5](../../fase-1/05-pacchetti.md#aggiornamenti-di-sicurezza-automatici).

Il *live patching* non elimina i riavvii: li rimanda alla finestra di manutenzione successiva.

## Aggiornare il firmware

Il firmware (UEFI, BMC, controller RAID, schede di rete, dischi) va aggiornato anche lui: le nuove versioni
correggono vulnerabilità e difetti. Due strade, sempre leggendo prima le **note di rilascio** del produttore:

- **gli strumenti del produttore**, dalla console del BMC o in blocco su molti server: Dell iDRAC e OpenManage
  Enterprise, HPE con i pacchetti SPP (*Service Pack for ProLiant*) e iLO, Lenovo XClarity. Tutti espongono anche
  il servizio di aggiornamento di **Redfish**;
- **fwupd**, su entrambe le famiglie: scarica il firmware dal servizio **LVFS** (*Linux Vendor Firmware Service*),
  dove molti produttori lo pubblicano.

```bash
# uguale su Ubuntu e RHEL (pacchetto fwupd)
fwupdmgr get-devices         # cosa può essere aggiornato
fwupdmgr refresh             # scarica l'elenco aggiornato da LVFS
fwupdmgr get-updates
sudo fwupdmgr update
```

## La salute dell'hardware

L'hardware avvisa prima di guastarsi, se qualcuno ascolta. Le fonti da tenere d'occhio:

| Componente | Strumento (Ubuntu e RHEL) | Cosa guardare |
|---|---|---|
| Dischi SATA/SAS | `smartctl` (pacchetto `smartmontools`) | settori riallocati, errori, esito dei test |
| SSD NVMe | `nvme smart-log` (pacchetto `nvme-cli`) | `percentage_used` (usura), `media_errors` |
| RAID software | `mdadm --detail`, `cat /proc/mdstat` | stato `degraded`, ricostruzione in corso |
| ZFS (Ubuntu) | `zpool status`, `zpool scrub` | errori di checksum, dischi `FAULTED` |
| Memoria | `ras-mc-ctl --summary` (pacchetto `rasdaemon`) | errori corretti dall'ECC in aumento: il modulo va cambiato |
| Temperature e ventole | `sensors` (`lm-sensors` su Ubuntu, `lm_sensors` su RHEL), BMC | valori fuori soglia |
| Tutto il server | registro eventi del BMC: `ipmitool sel list` o Redfish | alimentatori, ventole, memoria, temperature |

```bash
sudo smartctl -H -a /dev/sda          # salute generale e attributi
sudo smartctl -t short /dev/sda       # avvia un test breve
sudo nvme smart-log /dev/nvme0
cat /proc/mdstat
sudo ras-mc-ctl --summary
```

**smartd** (incluso in `smartmontools`) controlla i dischi in continuazione e invia un avviso; il BMC può inviare
eventi via email, SNMP o Redfish. In un sistema di monitoraggio ([Fase 7](../../fase-7/index.md)) queste
informazioni entrano come metriche con exporter dedicati (node exporter per sensori e RAID software, exporter
per SMART, IPMI e Redfish) e fanno scattare gli alert.

## Sostituire un disco guasto

Con dischi sostituibili a caldo l'intervento si fa a server acceso.

1. **Identifica** il disco: modello e numero di serie da `smartctl -i` o `lsblk -o NAME,SERIAL`, e il suo
   alloggiamento dal BMC. Accendi il LED di localizzazione (dal BMC o con `ledctl locate=/dev/sdX`, pacchetto
   `ledmon`) per non estrarre quello sbagliato.
2. Con un **controller RAID hardware**: estrai il disco, inserisci il nuovo, la ricostruzione parte da sola.
3. Con **RAID software**:

    ```bash
    sudo mdadm --manage /dev/md0 --fail /dev/sdc1 --remove /dev/sdc1
    # sostituisci fisicamente il disco e ricrea la stessa partizione, poi:
    sudo mdadm --manage /dev/md0 --add /dev/sdc1
    watch cat /proc/mdstat          # avanzamento della ricostruzione
    ```

    Con ZFS: `sudo zpool replace dati VECCHIO NUOVO`.

4. Durante la ricostruzione l'array è vulnerabile: evita altri interventi e controlla che arrivi alla fine.
5. Restituisci il disco guasto al produttore secondo la garanzia, oppure **distruggilo** se i contratti prevedono
   che i dischi non lascino l'azienda.

## Quando qualcosa si rompe

- Parti dal **metodo** della [lezione 1.13](../../fase-1/13-troubleshooting.md) e dai log
  ([lezione 1.4](../../fase-1/04-processi-systemd.md#journalctl-i-log-di-systemd)), poi controlla il registro eventi
  del BMC: molti blocchi "inspiegabili" sono hardware.
- Se il kernel va in crash, **kdump** salva un'immagine della memoria da analizzare. Su RHEL è attivo di default;
  su Ubuntu si installa con `linux-crashdump`.
- Per l'assistenza del produttore del sistema operativo si raccoglie un rapporto completo con **sos**
  (`sudo sos report`; pacchetto `sosreport` su Ubuntu, `sos` su RHEL); per l'hardware, il BMC genera un pacchetto
  di diagnostica (per esempio il *SupportAssist collection* di Dell o l'*Active Health System log* di HPE).

## Documentazione e inventario

Un buon sistemista sa sempre **cosa ha, dove si trova e in che stato è**:

- **inventario** (spesso in un CMDB o in strumenti open source come NetBox): posizione nel rack, numero di serie,
  garanzia, indirizzi, porte dello switch, servizi ospitati;
- **runbook**: le procedure scritte per gli interventi ricorrenti (sostituire un disco, aggiornare il firmware,
  ripristinare un backup);
- **registro delle modifiche**: chi ha cambiato cosa e quando. Se la configurazione è in Ansible
  ([lezione 6.10](../../fase-6/10-ansible.md)), Git fa già gran parte di questo lavoro.

## Il ciclo di vita

| | Ubuntu Server LTS | RHEL |
|---|---|---|
| Nuova versione | ogni due anni (aprile degli anni pari) | una maggiore ogni circa tre anni, minori ogni sei mesi |
| Supporto | 5 anni di aggiornamenti di sicurezza; fino a 10 anni con Ubuntu Pro (ESM) | 10 anni, estendibili con abbonamenti aggiuntivi |
| Passare alla versione successiva | `sudo do-release-upgrade` | reinstallazione, oppure aggiornamento sul posto con **Leapp** |
| Gestire molti server | Landscape | Satellite |

L'**hardware** si sostituisce di solito dopo 5-7 anni, quando finisce la garanzia o il supporto del produttore, o
quando un server nuovo fa lo stesso lavoro consumando molto meno. Prima di dismetterlo:

- sposta i servizi e aggiorna la documentazione;
- **cancella i dati in modo sicuro**: per gli SSD la cancellazione crittografica o la funzione *sanitize* del
  disco (dal BMC o con `nvme sanitize`), per gli HDD la sovrascrittura; in alternativa la distruzione fisica;
- reimposta il BMC ai valori di fabbrica;
- affida l'hardware a un'azienda di smaltimento autorizzata (in Italia è un rifiuto RAEE) o a un ricondizionatore.

## Esercizi

!!! esercizio "Esercizio 1"
    `nvme smart-log` su un SSD di un database restituisce `percentage_used : 92%`. Cosa significa e cosa fai?

??? success "Soluzione"
    L'SSD ha consumato il 92% delle scritture garantite dal produttore. Non è ancora guasto, ma va **sostituito
    presto**, in una finestra pianificata: apri un ordine, verifica che il RAID sia sano e controlla l'usura degli
    altri dischi dello stesso array, che probabilmente hanno lavorato quanto lui.

!!! esercizio "Esercizio 2"
    Scrivi l'ordine delle operazioni per aggiornare il kernel su tre nodi worker Kubernetes senza fermare
    l'applicazione.

??? success "Soluzione"
    Un nodo alla volta: `kubectl drain` del nodo ([lezione 8.8](../../fase-8/08-scheduling.md#manutenzione-dei-nodi)),
    aggiornamento dei pacchetti, riavvio, verifica che il nodo sia `Ready`, `kubectl uncordon`, controllo che i Pod
    siano tornati sani. Solo allora si passa al nodo successivo.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale strumento aggiorna il firmware usando il servizio LVFS, su Ubuntu e RHEL?","options":["fwupd","kpatch","Leapp","smartd"],"answer":[0],"explanation":"fwupdmgr scarica e applica i firmware pubblicati dai produttori su LVFS."},
 {"q":"Quale strumento applica le correzioni del kernel senza riavvio su RHEL?","options":["Livepatch","kpatch","kdump","needrestart"],"answer":[1],"explanation":"Su Ubuntu l'equivalente è Livepatch, con Ubuntu Pro."},
 {"q":"Il comando `ras-mc-ctl --summary` mostra errori di memoria corretti in aumento su un modulo. Cosa significa?","options":["Nulla, l'ECC li corregge","Il modulo si sta degradando e va sostituito","Serve più swap","Il disco è pieno"],"answer":[1],"explanation":"Gli errori corretti in crescita anticipano spesso un errore non correggibile."},
 {"q":"Quali sono buone pratiche prima di dismettere un server?","options":["Cancellare i dati in modo sicuro","Reimpostare il BMC","Aggiornare l'inventario","Buttarlo nella raccolta indifferenziata","Lasciare i dischi dentro per comodità"],"answer":[0,1,2],"explanation":"I dischi vanno cancellati o distrutti e l'hardware smaltito come RAEE."},
 {"q":"Come si aggiorna sul posto un sistema RHEL alla versione maggiore successiva?","options":["do-release-upgrade","Leapp","dnf distro-sync","fwupdmgr"],"answer":[1],"explanation":"do-release-upgrade è lo strumento di Ubuntu."}
]}
</script>

## Fonti

- [Ubuntu Server: Automatic updates](https://ubuntu.com/server/docs/how-to/software/automatic-updates/)
- [Ubuntu Server: Kernel crash dump](https://ubuntu.com/server/docs/how-to/software/kernel-crash-dump/)
- [Ubuntu Server: Upgrade your release](https://ubuntu.com/server/docs/how-to/software/upgrade-your-release/)
- [RHEL 10: Applying patches with kernel live patching](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/managing_monitoring_and_updating_the_kernel/applying-patches-with-kernel-live-patching)
- [RHEL 10: Managing software with the DNF tool](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/managing_software_with_the_dnf_tool/index)
- [fwupd e LVFS](https://fwupd.org/)
- [smartmontools](https://www.smartmontools.org/)
- [Ubuntu release cycle](https://ubuntu.com/about/release-cycle)
- [Red Hat Enterprise Linux life cycle](https://access.redhat.com/support/policy/updates/errata)

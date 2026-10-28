# 1.15 Avvio, tempo e prestazioni

!!! obiettivi "Obiettivi"
    - Descrivere la sequenza di avvio di Linux e gestire i target di systemd.
    - Modificare i parametri del kernel all'avvio e a sistema avviato.
    - Recuperare l'accesso a un sistema di cui non si conosce la password di root.
    - Sincronizzare l'orologio, rendere persistente il journal e usare i profili di prestazioni.

Questa lezione raccoglie argomenti di amministrazione richiesti dagli esami **RHCSA** e **LFCS** e utili ogni
volta che un server non si avvia.

## La sequenza di avvio

```mermaid
flowchart LR
  FW[Firmware<br>UEFI o BIOS] --> BL[Bootloader<br>GRUB 2]
  BL --> K[Kernel<br>+ initramfs]
  K --> SD[systemd<br>PID 1]
  SD --> T[Target<br>multi-user o graphical]
```

1. Il **firmware** (UEFI, o BIOS nei sistemi vecchi) verifica l'hardware e avvia il bootloader.
2. **GRUB 2** mostra il menu, carica il **kernel** e l'**initramfs** (un piccolo filesystem temporaneo con i
   driver necessari a trovare il disco di sistema) e passa al kernel i suoi **parametri**.
3. Il kernel monta il filesystem radice e avvia **systemd**.
4. systemd porta il sistema al **target** predefinito.

## I target

Un **target** è un gruppo di unità che corrisponde a uno stato del sistema:

| Target | Stato |
|---|---|
| `multi-user.target` | Sistema completo, senza interfaccia grafica: lo standard dei server. |
| `graphical.target` | Come il precedente, più l'interfaccia grafica. |
| `rescue.target` | Modalità di manutenzione: filesystem montati, pochi servizi, shell di root. |
| `emergency.target` | Il minimo indispensabile: radice montata in sola lettura. Serve quando fallisce persino il montaggio dei filesystem. |

```bash
systemctl get-default                         # target all'avvio
sudo systemctl set-default multi-user.target  # lo cambia in modo permanente
sudo systemctl isolate rescue.target          # passa subito a un altro target (attenzione: interrompe i servizi)
systemctl reboot; systemctl poweroff          # riavvio e spegnimento
```

Per avviare **una volta** in un target diverso, nel menu di GRUB premi ++e++ sulla voce, aggiungi alla riga che
inizia con `linux` il parametro `systemd.unit=rescue.target` e avvia con ++ctrl+x++.

## Parametri del kernel

### All'avvio (GRUB)

=== "Debian/Ubuntu"

    ```bash
    sudo nano /etc/default/grub        # modifica GRUB_CMDLINE_LINUX_DEFAULT o GRUB_TIMEOUT
    sudo update-grub                   # rigenera /boot/grub/grub.cfg
    ```

=== "RHEL (Rocky/AlmaLinux)"

    Su RHEL si usa **grubby**, che modifica direttamente le voci dei kernel installati:

    ```bash
    sudo grubby --default-kernel                         # kernel avviato per default
    sudo grubby --info=ALL | grep -E "^kernel|^args"     # parametri attuali
    sudo grubby --update-kernel=ALL --args="quiet"       # aggiunge un parametro a tutti i kernel
    sudo grubby --update-kernel=ALL --remove-args="quiet"
    ```

    Il timeout del menu si imposta in `/etc/default/grub` (`GRUB_TIMEOUT`), seguito da
    `sudo grub2-mkconfig -o /boot/grub2/grub.cfg`.

### A sistema avviato (sysctl)

Molti parametri del kernel si leggono e modificano al volo in `/proc/sys`, con `sysctl`:

```bash
sysctl net.ipv4.ip_forward                     # leggi
sudo sysctl -w vm.swappiness=10                # modifica temporanea
echo "vm.swappiness=10" | sudo tee /etc/sysctl.d/90-swappiness.conf   # permanente
sudo sysctl --system                           # ricarica tutti i file di configurazione
```

## Recuperare l'accesso come root

Hai dimenticato la password di root e nessun utente ha sudo. Con l'accesso alla **console** della macchina (non via
rete) puoi recuperare l'accesso. È un obiettivo esplicito della RHCSA.

=== "RHEL (Rocky/AlmaLinux)"

    1. Riavvia, nel menu di GRUB premi ++e++ e aggiungi alla fine della riga `linux` il parametro **`rd.break`**.
       Avvia con ++ctrl+x++: ti fermi in una shell dell'initramfs.
    2. Il sistema vero è montato in sola lettura in `/sysroot`:

        ```bash
        mount -o remount,rw /sysroot
        chroot /sysroot
        passwd root
        touch /.autorelabel        # SELinux rietichetterà i file al prossimo avvio
        exit
        exit
        ```

    Il file `/.autorelabel` è indispensabile: `passwd` ha riscritto `/etc/shadow` senza il contesto SELinux corretto,
    e senza rietichettatura non potresti più accedere. Il primo avvio dopo sarà più lungo.

=== "Debian/Ubuntu"

    1. Nel menu di GRUB (su Ubuntu compare tenendo premuto ++shift++ o ++esc++ all'avvio) scegli la voce e premi ++e++.
    2. Alla riga `linux` sostituisci `ro` con `rw` e aggiungi **`init=/bin/bash`**. Avvia con ++ctrl+x++.
    3. Sei root, senza password:

        ```bash
        passwd root          # oppure: passwd nomeutente
        sync
        exec /sbin/init      # prosegue l'avvio normale (oppure riavvia con: reboot -f)
        ```

!!! note "Nelle VM del corso e nel cloud"
    Con Multipass e Vagrant non vedi il menu di GRUB; l'esercizio va fatto su una VM con console grafica (per
    esempio creata a mano in VirtualBox) o su RHEL con la Developer Subscription. Nel cloud la console è spesso
    assente: per questo si usano accessi alternativi (Session Manager, console seriale) e soprattutto la
    **prevenzione** (chiavi SSH di emergenza, snapshot).

## Data, ora e sincronizzazione

Un orologio sbagliato rompe TLS (certificati "non ancora validi" o "scaduti"), i log e i sistemi distribuiti.
I server si sincronizzano con **NTP** (*Network Time Protocol*).

```bash
timedatectl                                  # data, fuso orario, stato della sincronizzazione
sudo timedatectl set-timezone Europe/Rome
timedatectl list-timezones | grep Europe
```

=== "Debian/Ubuntu"

    Ubuntu usa per default **systemd-timesyncd**, un client NTP semplice. Per un server che deve anche servire l'ora
    ad altri si installa **chrony** (`sudo apt install chrony`), che lo sostituisce.

=== "RHEL (Rocky/AlmaLinux)"

    RHEL usa **chrony**:

    ```bash
    sudo nano /etc/chrony.conf        # righe "server ... iburst" o "pool ... iburst"
    sudo systemctl restart chronyd
    chronyc sources -v                # server usati e stato (* = sorgente selezionata)
    chronyc tracking                  # scarto rispetto all'ora di riferimento
    ```

## Journal persistente

Su molte installazioni il journal è conservato solo in memoria (`/run/log/journal`) e si perde al riavvio. Per
renderlo persistente (obiettivo RHCSA):

```bash
sudo mkdir -p /etc/systemd/journald.conf.d
printf '[Journal]\nStorage=persistent\n' | sudo tee /etc/systemd/journald.conf.d/persistente.conf
sudo systemctl restart systemd-journald
ls /var/log/journal/                         # ora esiste
journalctl --list-boots                      # dopo un riavvio vedrai più avvii
```

## Profili di prestazioni: tuned

Su RHEL il servizio **tuned** applica profili di ottimizzazione (parametri del kernel, scheduler dei dischi, energia):

```bash
sudo dnf install -y tuned && sudo systemctl enable --now tuned
tuned-adm active                 # profilo attivo
tuned-adm list                   # profili disponibili
tuned-adm recommend              # profilo consigliato per questa macchina
sudo tuned-adm profile virtual-guest
```

`tuned` esiste anche per Debian/Ubuntu (`sudo apt install tuned`), ma è usato soprattutto nella famiglia Red Hat.

## Attività una tantum: at

`at` esegue un comando **una volta**, a un'ora stabilita (per le attività ripetute vedrai cron e i timer nella
[lezione 2.5](../fase-2/05-pianificazione.md)):

```bash
# pacchetto "at" su entrambe le famiglie; servizio atd
sudo systemctl enable --now atd
echo "date > /tmp/ora-at.txt" | at now + 2 minutes
atq                      # coda dei lavori
atrm 3                   # cancella il lavoro numero 3
```

## Esercizi

!!! esercizio "Esercizio 1 · Target e journal"
    Su una delle VM: imposta `multi-user.target` come target predefinito, rendi persistente il journal, riavvia e
    mostra gli errori dell'avvio **precedente**.

??? success "Soluzione"
    ```bash
    sudo systemctl set-default multi-user.target
    sudo mkdir -p /etc/systemd/journald.conf.d
    printf '[Journal]\nStorage=persistent\n' | sudo tee /etc/systemd/journald.conf.d/persistente.conf
    sudo systemctl restart systemd-journald
    sudo reboot
    # dopo il riavvio
    journalctl --list-boots
    journalctl -b -1 -p err
    ```
    **Spiegazione.** `-b -1` indica l'avvio precedente, disponibile solo con il journal persistente. Nella VM
    Ubuntu il journal è spesso già persistente: lo riconosci dall'esistenza di `/var/log/journal`.

!!! esercizio "Esercizio 2 · Ora e sysctl"
    Imposta il fuso orario `Europe/Rome`, verifica che l'ora sia sincronizzata e rendi permanente
    `net.ipv4.ip_forward=1`. Verifica il valore dopo un riavvio.

??? success "Soluzione"
    ```bash
    sudo timedatectl set-timezone Europe/Rome
    timedatectl | grep -E "Time zone|synchronized"
    echo "net.ipv4.ip_forward=1" | sudo tee /etc/sysctl.d/99-forward.conf
    sudo sysctl --system | grep ip_forward
    sudo reboot
    # dopo il riavvio
    sysctl net.ipv4.ip_forward     # = 1
    ```

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Quale componente carica il kernel e l'initramfs?","options":["systemd","GRUB 2","UEFI direttamente","chrony"],"answer":[1],"explanation":"Il bootloader carica kernel e initramfs e passa i parametri al kernel."},
 {"q":"Quale comando imposta in modo permanente l'avvio senza interfaccia grafica?","options":["systemctl isolate multi-user.target","systemctl set-default multi-user.target","grubby --default","init 3"],"answer":[1],"explanation":"isolate cambia stato subito ma non è permanente."},
 {"q":"Dopo aver cambiato la password di root con rd.break su RHEL, cosa non devi dimenticare?","options":["touch /.autorelabel","update-grub","setenforce 0 permanente","Cancellare /etc/shadow"],"answer":[0],"explanation":"Serve a ripristinare i contesti SELinux al prossimo avvio."},
 {"q":"Dove rendi permanente un parametro sysctl?","options":["/proc/sys","/etc/sysctl.d/*.conf","/etc/default/grub","/etc/fstab"],"answer":[1],"explanation":"/proc/sys contiene i valori correnti, che si perdono al riavvio."},
 {"q":"Quale servizio NTP usa RHEL per default?","options":["ntpd","chronyd","systemd-timesyncd","timed"],"answer":[1],"explanation":"chrony è il client e server NTP di RHEL."},
 {"q":"Cosa serve per vedere i log dell'avvio precedente con journalctl -b -1?","options":["Il journal persistente","SELinux in permissive","Il target graphical","tuned"],"answer":[0],"explanation":"Se il journal è solo in memoria, si perde a ogni riavvio."}
]}
</script>

## Fonti

- [systemd.special(7): target speciali](https://www.freedesktop.org/software/systemd/man/latest/systemd.special.html)
- [Red Hat: Managing, monitoring, and updating the kernel (RHEL 10)](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10/html/managing_monitoring_and_updating_the_kernel/index)
- [GNU GRUB Manual](https://www.gnu.org/software/grub/manual/grub/grub.html)
- [chrony documentation](https://chrony-project.org/documentation.html)
- [journald.conf(5)](https://www.freedesktop.org/software/systemd/man/latest/journald.conf.html)
- [sysctl(8)](https://man7.org/linux/man-pages/man8/sysctl.8.html), [at(1)](https://man7.org/linux/man-pages/man1/at.1p.html)
- [Ubuntu Server docs: Time synchronisation](https://documentation.ubuntu.com/server/explanation/networking/about-time-synchronisation/)

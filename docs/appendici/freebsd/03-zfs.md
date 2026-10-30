# 3. ZFS

!!! obiettivi "Obiettivi"
    - Spiegare cosa unisce ZFS: RAID, volumi, filesystem, controllo dei dati.
    - Creare un pool in mirror e dei dataset con proprietà diverse.
    - Usare snapshot, rollback e copie con `zfs send` e `zfs receive`.
    - Proteggere gli aggiornamenti con i boot environment.

## Un solo strumento al posto di tre

Su Linux, per avere dischi ridondati, volumi flessibili e un filesystem servono tre strati: RAID (`mdadm`), LVM e
il filesystem ([lezione 1.6](../../fase-1/06-storage.md)). **ZFS** li unisce in uno:

- il **pool** (`zpool`) raccoglie i dischi e decide la ridondanza: *mirror* (copie), *raidz* (parità, come il RAID
  5 e 6);
- i **dataset** (`zfs`) sono filesystem che condividono lo spazio del pool, ognuno con proprietà sue
  (compressione, quote, punto di montaggio);
- ogni blocco ha un **checksum**: ZFS si accorge se un disco restituisce dati sbagliati e, con un mirror, li
  corregge leggendo la copia buona;
- gli **snapshot** sono istantanei e all'inizio non occupano spazio, perché ZFS non sovrascrive mai i dati
  (*copy-on-write*).

Con l'installazione *Auto (ZFS)* della [pagina 2](02-installazione.md#linstallazione) il sistema stesso è già su un pool,
`zroot`:

```sh
zpool status zroot
zfs list                               # i dataset: zroot/ROOT/default è il sistema, zroot/home le home...
```

ZFS esiste anche su Linux (OpenZFS, in Ubuntu), ed è la base di NAS come TrueNAS
([appendice hardware](../hardware/05-server-dati.md#nas)).

## Un pool in mirror

I due dischi piccoli della VM sono `vtbd1` e `vtbd2` (controlla con `geom disk list`).

```sh
sudo zpool create dati mirror vtbd1 vtbd2
zpool status dati
sudo zfs create -o compression=zstd dati/documenti
sudo zfs create -o quota=500M dati/progetti
zfs list -o name,used,avail,compression,quota,mountpoint
```

Il pool si monta da solo in `/dati`, i dataset in `/dati/documenti` e `/dati/progetti`: niente `/etc/fstab`.

## Snapshot, rollback, copie

```sh
echo "versione 1" | sudo tee /dati/documenti/nota.txt
sudo zfs snapshot dati/documenti@prima
echo "versione 2, sbagliata" | sudo tee /dati/documenti/nota.txt
ls /dati/documenti/.zfs/snapshot/prima/        # gli snapshot si sfogliano come cartelle
sudo zfs rollback dati/documenti@prima         # torna allo stato dello snapshot
cat /dati/documenti/nota.txt                   # versione 1
zfs list -t snapshot
```

`zfs send` trasforma uno snapshot in un flusso di dati, `zfs receive` lo ricrea altrove: è il modo di fare backup
e repliche con ZFS, anche verso un altro server via SSH.

```sh
sudo zfs send dati/documenti@prima | sudo zfs receive zroot/copia-documenti
# verso un'altra macchina: zfs send ... | ssh altro-server sudo zfs receive pool/dataset
```

## I boot environment

Poiché anche il sistema è un dataset, prima di un aggiornamento puoi fotografarlo: se qualcosa va storto, riavvii
nella versione precedente scegliendola dal menu di avvio.

```sh
sudo bectl create prima-aggiornamento
bectl list                             # l'ambiente attivo ha la N (now) e la R (al riavvio)
sudo bectl activate prima-aggiornamento    # per tornare indietro al prossimo riavvio
```

!!! esercizio "Esercizio · Un disco che si rompe"
    Con il pool `dati` in mirror: scrivi un file da 100 MB in `/dati/documenti`, poi simula la rottura di un disco
    con `sudo zpool offline dati vtbd2`. Il file si legge ancora? Cosa dice `zpool status`? Rimetti il disco in
    linea e guarda cosa fa ZFS.

??? success "Soluzione"
    ```sh
    sudo dd if=/dev/random of=/dati/documenti/grande bs=1m count=100
    sha256 /dati/documenti/grande
    sudo zpool offline dati vtbd2
    zpool status dati                    # state: DEGRADED, vtbd2 OFFLINE
    sha256 /dati/documenti/grande        # stesso valore: i dati sono sull'altro disco
    sudo zpool online dati vtbd2
    zpool status dati                    # resilvered: ZFS ricopia solo i blocchi cambiati
    sudo zpool scrub dati                # rilegge tutto e verifica i checksum
    ```
    Il **resilver** copia sul disco tornato solo ciò che è cambiato mentre era fuori. Lo **scrub** rilegge tutti i
    dati e li confronta con i checksum: sui server si pianifica una volta al mese.

## Fonti

- [FreeBSD Handbook: The Z File System (ZFS)](https://docs.freebsd.org/en/books/handbook/zfs/)
- [OpenZFS documentation](https://openzfs.github.io/openzfs-docs/)
- [bectl(8)](https://man.freebsd.org/cgi/man.cgi?query=bectl&sektion=8)

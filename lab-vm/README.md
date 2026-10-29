# VM di laboratorio

Le macchine virtuali Linux della Fase 1, dell'appendice LFCS e del percorso RHEL, create con **KVM/QEMU e libvirt**
a partire dalle immagini cloud ufficiali di Ubuntu 24.04 e AlmaLinux 10.

| File | Contenuto |
|---|---|
| `academy-vm` | lo script che crea la rete e le VM (`./academy-vm` senza argomenti mostra l'aiuto) |
| `academy-net.xml` | la rete privata `academy` 192.168.56.0/24, con indirizzi fissi e DNS interno |
| `cloud-init/` | i profili cloud-init di ogni VM: pacchetti e dischi di laboratorio |

Le istruzioni complete, passo per passo, sono nel corso: **Ambiente di laboratorio → Macchine virtuali Linux**.

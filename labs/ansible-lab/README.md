# Laboratorio Ansible

Quattro nodi gestiti in container, raggiungibili via SSH come se fossero server veri, usati nella
[lezione 6.10](../../docs/fase-6/10-ansible.md) e nell'[appendice RHCE](../../docs/appendici/rhce/index.md).

| Nodo | Sistema | Gruppi | Porte sul tuo computer |
|---|---|---|---|
| `web1`, `web2` | Ubuntu 24.04 | `web`, `debian` | SSH 2221/2222, HTTP 8091/8092 |
| `db1`, `app1` | AlmaLinux 9 (compatibile RHEL) | `db`, `app`, `rhel` | SSH 2223/2224, HTTP 8094 (`app1`) |

I nodi eseguono **systemd**, quindi i moduli che gestiscono servizi funzionano come su una VM.

```bash
cd labs/ansible-lab
./prepara.sh                    # crea la chiave SSH del laboratorio (una volta sola)
docker compose up -d --build    # avvia i quattro nodi
ansible all -m ping             # verifica: quattro "pong"
docker compose down             # spegne i nodi (le modifiche fatte da Ansible si perdono)
```

L'utente remoto è `ansible`, con `sudo` senza password. La chiave privata è in `.chiavi/` ed è esclusa da Git.

Note:

- il laboratorio richiede Docker con cgroup v2 (il Dev Container del corso, Docker Desktop recente, Linux moderni);
- `firewalld` e SELinux non sono attivi nei container: per esercitarti su questi argomenti, richiesti dall'esame
  RHCE, usa le VM descritte nell'appendice.

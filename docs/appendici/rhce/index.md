# EX294 · Red Hat Certified Engineer in Ansible

L'esame pratico di Red Hat sull'automazione con **Ansible**. Scrivi playbook, ruoli e inventari che configurano più
sistemi RHEL; alla fine i tuoi playbook vengono eseguiti su sistemi appena installati e verificati.

| | |
|---|---|
| **Ente** | Red Hat |
| **Esame** | EX294, chiamato ora *Red Hat Certified Advanced System Administrator in Ansible* |
| **Credenziale** | superato da solo: Red Hat Certified Advanced System Administrator in Ansible; con la RHCSA: **Red Hat Certified Engineer (RHCE) in Ansible** |
| **Sistema** | RHEL 9 con Red Hat Ansible Automation Platform |
| **Strumenti** | `ansible-navigator`, `ansible-playbook`, Git, VS Code con dev container Ansible |
| **Formato** | pratico, senza Internet; documentazione dei prodotti disponibile |
| **Durata** | circa 4 ore (la pagina dell'esame non la indica: verificala al momento della prenotazione) |
| **Costo** | varia per paese: vedi la pagina ufficiale |
| **Pagina ufficiale** | [EX294 su redhat.com](https://www.redhat.com/en/services/training/ex294-red-hat-certified-engineer-rhce-exam-red-hat-enterprise-linux) · [credenziale RHCE in Ansible](https://www.redhat.com/en/services/certification/red-hat-certified-engineer-in-ansible) |

Dati verificati il 28 ottobre 2026.

!!! info "Cambio di nome"
    Red Hat ha diviso il titolo RHCE in due percorsi: **RHCE in Ansible** (RHCSA + EX294) e **RHCE in Enterprise
    Linux** (RHCSA + EX342). L'esame EX294 da solo porta ora il titolo di *Advanced System Administrator in
    Ansible*. Nei vecchi materiali troverai ancora "RHCE" riferito al solo EX294.

## Quando sostenerla

Dopo la RHCSA e la [Fase 6](../../fase-6/index.md), [lezione 6.10](../../fase-6/10-ansible.md) compresa. Serve anche sicurezza con Git ([Fase 2](../../fase-2/index.md)).

## Il laboratorio

- Per gran parte degli obiettivi basta il laboratorio del corso [`labs/ansible-lab`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/ansible-lab):
  quattro nodi in container, due Ubuntu e due AlmaLinux.
- **firewalld e SELinux** non sono attivi nei container: per quei task crea due o tre VM Rocky Linux 9 o
  AlmaLinux 9 (come nella [lezione 8.12](../../fase-8/12-kubeadm.md)) e aggiungile all'inventario.
- Installa `ansible-navigator` e un execution environment, e prova anche il dev container di VS Code con
  l'estensione Ansible: all'esame userai entrambi.

## Piano di studio (6 settimane)

| Settimana | Attività |
|---|---|
| 1 | `ansible.cfg`, `ansible-navigator.yml`, inventari, comandi ad hoc, chiavi e sudo sui nodi ([6.10](../../fase-6/10-ansible.md)) |
| 2 | Playbook: moduli comuni, variabili, fatti, `register`, cicli, condizioni, handler |
| 3 | Gestione degli errori, template Jinja2, Ansible Vault |
| 4 | Ruoli, `requirements.yml`, collection, ruoli di sistema RHEL |
| 5 | Automatizzare i compiti RHCSA: pacchetti, servizi, firewall, storage, SELinux, utenti, cron |
| 6 | Esami pratici a tempo; rieseguire tutto su nodi appena ricreati |

## In questa appendice

- [Argomenti specifici](argomenti.md)
- [Simulazione pratica](simulazione.md): 32 task originali con verifica e soluzione

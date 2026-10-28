# Argomenti specifici · EX294 (Ansible)

Gli obiettivi pubblicati da Red Hat, raggruppati come nella simulazione. La lezione di riferimento è la
[6.10](../../fase-6/10-ansible.md); gli argomenti che il corso tratta solo in parte sono spiegati nella sezione
[Approfondimenti](#approfondimenti).

| Obiettivi | Cosa devi saper fare | Dove studiarlo |
|---|---|---|
| **Compiti RHCSA** | tutti gli obiettivi della RHCSA, più leggere semplici script di shell | [Appendice RHCSA](../rhcsa/argomenti.md) |
| **Componenti di Ansible** | inventari, moduli, variabili, fatti, cicli, condizioni, play, gestione degli errori, playbook, file di configurazione, ruoli, documentazione dei moduli | [6.10](../../fase-6/10-ansible.md), [errori](#gestione-degli-errori), [fatti locali](#fatti-locali-e-variabili-magiche) |
| **Configurare Ansible** | creare e modificare `ansible.cfg` e `ansible-navigator.yml`, inventari statici con gruppi e gruppi di gruppi | [6.10](../../fase-6/10-ansible.md), [navigator](#ansible-navigator) |
| **Nodi gestiti** | distribuire chiavi SSH, configurare l'escalation dei privilegi, distribuire file | [1.10](../../fase-1/10-ssh.md), [6.10](../../fase-6/10-ansible.md) |
| **Eseguire i playbook** | `ansible-navigator run` e `ansible-playbook`; trovare moduli nelle collection con navigator; inventari e configurazione con navigator | [navigator](#ansible-navigator) |
| **Git** | clonare un repository, aggiungere file, commit e push | [2.1](../../fase-2/01-git-base.md), [2.2](../../fase-2/02-git-collaborazione.md) |
| **VS Code** | creare playbook e inviarli a Git, configurare ansible-navigator, eseguire i playbook in un dev container Ansible | [0.5](../../fase-0/05-editor-markdown.md), [0.7](../../fase-0/07-laboratorio.md), [VS Code](#vs-code-e-dev-container) |
| **Play e playbook** | moduli comuni, `register`, condizioni, gestione degli errori, portare i sistemi a uno stato definito | [6.10](../../fase-6/10-ansible.md), [errori](#gestione-degli-errori) |
| **Ruoli e collection** | creare ruoli, installarli da `requirements.yml`, installare e usare collection, ruoli di sistema RHEL | [6.10](../../fase-6/10-ansible.md) |
| **Automatizzare la RHCSA** | moduli per pacchetti e repository, servizi, firewall, file system, dispositivi di storage, contenuto dei file, archivi, pianificazione, sicurezza, utenti e gruppi | [6.10](../../fase-6/10-ansible.md) e la [simulazione](simulazione.md) |
| **Contenuti** | template Jinja2, Ansible Vault | [6.10](../../fase-6/10-ansible.md) |

## Approfondimenti

### ansible-navigator

`ansible-navigator` esegue Ansible dentro un **execution environment**, un'immagine container con ansible-core,
collection e dipendenze: tutti usano lo stesso ambiente.

```bash
ansible-navigator run site.yml -m stdout          # esegue un playbook
ansible-navigator collections -m stdout           # collection disponibili nell'EE
ansible-navigator doc ansible.posix.firewalld -m stdout
ansible-navigator inventory -i inventario --graph -m stdout
ansible-navigator settings --effective -m stdout  # configurazione effettiva
```

La configurazione si cerca in `./ansible-navigator.yml` (o nel file indicato da `ANSIBLE_NAVIGATOR_CONFIG`).
Fonte: [documentazione di ansible-navigator](https://ansible.readthedocs.io/projects/navigator/).

### Gestione degli errori

- `ignore_errors: true` continua anche se il task fallisce;
- `failed_when` e `changed_when` ridefiniscono quando un task è fallito o ha cambiato qualcosa;
- `block`, `rescue` e `always` funzionano come try, catch e finally;
- `ansible.builtin.fail` e `ansible.builtin.assert` fermano l'esecuzione con un messaggio chiaro;
- `any_errors_fatal: true` ferma il play su tutti gli host al primo errore.

Fonte: [Error handling in playbooks](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_error_handling.html).

### Fatti locali e variabili magiche

I file `*.fact` in `/etc/ansible/facts.d/` sui nodi (formato INI o JSON) compaiono in `ansible_local`.
Le variabili magiche più usate sono `inventory_hostname`, `group_names`, `groups` e `hostvars`: con `hostvars`
leggi i fatti di un altro host, a patto che siano stati raccolti nello stesso run.

### VS Code e dev container

L'estensione Ansible di VS Code segnala errori, completa i nomi dei moduli e integra ansible-lint. Con un
**dev container Ansible** lo strumento gira in un container con ansible-navigator già pronto: apri la cartella,
scegli *Reopen in Container* e usa il terminale integrato. È lo stesso principio del Dev Container del corso
([lezione 0.7](../../fase-0/07-laboratorio.md)). Fonte: [Ansible VS Code extension](https://docs.ansible.com/projects/vscode-ansible/).

## Dove sbagliano più spesso i candidati

- **Playbook non idempotenti**: `command` e `shell` cambiano a ogni esecuzione se non usi `creates`,
  `changed_when` o un modulo dedicato.
- **Nomi dei file e percorsi**: l'esame verifica i file nei percorsi richiesti; un nome diverso vale zero.
- **Vault**: ricorda di passare la password (`--vault-password-file` o `--ask-vault-pass`) e di non lasciarla nel
  repository.
- **`firewalld` senza `immediate: true`**: la regola permanente non è attiva finché non ricarichi.

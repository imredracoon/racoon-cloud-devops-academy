# Fase 6 · Infrastruttura come codice

!!! obiettivi "Al termine di questa fase saprai"
    - spiegare cos'è l'infrastruttura come codice e quali problemi risolve;
    - descrivere reti, server, database e permessi su AWS e Google Cloud con **Terraform**;
    - usare il flusso di lavoro `init` → `plan` → `apply` → `destroy` in modo sicuro;
    - scrivere configurazioni riutilizzabili con variabili, espressioni, cicli e moduli;
    - gestire lo **state** in un backend remoto condiviso, con il blocco e la gestione delle differenze;
    - trattare correttamente i dati sensibili;
    - lavorare in team con HCP Terraform;
    - configurare i sistemi operativi dei server con **Ansible**.

| | |
|---|---|
| **Durata indicativa** | 1,5 mesi, circa 60 ore |
| **Prerequisiti** | Fasi 4 e 5 |
| **Ambiente** | Dev Container (Terraform e Ansible sono già installati), account AWS e/o Google Cloud |
| **Progetto** | [Infrastruttura come codice](progetto.md) |
| **Certificazioni collegate** | [HashiCorp Terraform Associate (004)](../appendici/terraform/index.md); [RHCE in Ansible](../appendici/rhce/index.md) (facoltativa, richiede la RHCSA) |

## Dalla console al codice

Nella Fase 5 hai costruito architetture complete cliccando nella console. Funziona, ma prova a immaginare di doverlo
rifare identico per un secondo ambiente, o di dover spiegare a un collega cosa hai cambiato martedì, o di ricostruire
tutto dopo un errore: sono ore di lavoro, con il rischio di dimenticare un dettaglio.

Con l'**infrastruttura come codice** (*Infrastructure as Code*, IaC) descrivi le risorse in file di testo che
**versioni con Git**, **rivedi** in una pull request e **applichi** con un comando. L'infrastruttura diventa
riproducibile, documentata e verificabile come il software.

In questa fase usi due strumenti complementari:

| Strumento | A cosa serve | Esempio |
|---|---|---|
| **Terraform** | **creare** e modificare le risorse del cloud (*provisioning*) | VPC, subnet, VM, database, bucket, ruoli IAM |
| **Ansible** | **configurare** i sistemi operativi e le applicazioni dentro i server (*configuration management*) | installare pacchetti, scrivere file di configurazione, gestire servizi e utenti |

Spesso si usano insieme: Terraform crea le VM, Ansible le configura. Quando le applicazioni girano in container o
servizi gestiti (Fasi 5 e 8), Ansible serve meno; resta però molto diffuso nelle aziende con server tradizionali, ed è
richiesto dalla certificazione Red Hat RHCE.

## Le versioni usate

| Strumento | Versione | Note |
|---|---|---|
| Terraform | 1.16 | l'esame Terraform Associate 004 si basa sulla 1.12: tutto ciò che vedi qui vale anche per l'esame |
| Provider AWS | 6.x | |
| Provider Google | 8.x | |
| ansible-core | 2.21 | l'esame RHCE si svolge su RHEL 9 con Ansible Automation Platform: le differenze sono indicate nell'appendice |

!!! info "Terraform e OpenTofu"
    Nel 2023 HashiCorp ha cambiato la licenza di Terraform (da MPL a BSL). La comunità ha creato **OpenTofu**, un
    fork open source gestito dalla Linux Foundation, compatibile quasi del tutto con Terraform: i comandi sono gli
    stessi, con `tofu` al posto di `terraform`. Ciò che impari qui vale per entrambi; il corso usa Terraform perché è
    lo strumento della certificazione e il più richiesto negli annunci.

## Le lezioni

1. [Infrastructure as Code](01-iac.md)
2. [Provider e primo progetto](02-primo-progetto.md)
3. [Il workflow](03-workflow.md)
4. [Risorse, data source e dipendenze](04-risorse.md)
5. [Variabili, output ed espressioni](05-variabili.md)
6. [Lo state](06-state.md)
7. [Moduli](07-moduli.md)
8. [Dati sensibili e Vault](08-segreti.md)
9. [HCP Terraform](09-hcp.md)
10. [Configurare server con Ansible](10-ansible.md)

Poi il [progetto](progetto.md) e la [verifica di fine fase](verifica.md).

I file dei laboratori sono nel repository del corso, in
[`labs/terraform`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/terraform) e
[`labs/ansible-lab`](https://github.com/imredracoon/racoon-cloud-devops-academy/tree/main/labs/ansible-lab).

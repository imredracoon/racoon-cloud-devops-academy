# Progetto 6 · Infrastruttura come codice

!!! obiettivi "Cosa dimostra questo progetto"
    Che sai descrivere un'architettura cloud completa con Terraform, in modo modulare, sicuro e riproducibile, e
    configurare server con Ansible. È la competenza che le offerte di lavoro per Cloud e DevOps Engineer citano più
    spesso, subito dopo il cloud stesso.

| | |
|---|---|
| **Durata** | 20-25 ore |
| **Consegna** | Repository pubblico `tasks-api-iac` |
| **Costo** | come il progetto 5, solo mentre l'ambiente è acceso: `terraform destroy` a fine sessione |

## Parte A · L'architettura del progetto 5 in Terraform

Ricostruisci con Terraform l'architettura a tre livelli del [progetto 5](../fase-5/progetto.md), variante A
(VM e gruppo di autoscaling) o B (ECS Fargate o Cloud Run), sul cloud che preferisci.

### Struttura del repository

```text
tasks-api-iac/
├── README.md
├── bootstrap/                 # crea il bucket per lo state (state locale, si lancia una volta)
├── modules/
│   ├── rete/                  # VPC, subnet pubbliche, private e dati su 2 zone, NAT
│   ├── app/                   # load balancer, gruppo di autoscaling (o servizio di container), IAM, allarmi
│   └── database/              # PostgreSQL gestito, gruppo di subnet, firewall, segreto
├── envs/
│   ├── dev/                   # main.tf, backend.tf, terraform.tfvars: usa i tre moduli
│   └── prod/                  # stessa struttura, valori diversi (Multi-AZ, più istanze)
└── .github/                   # vuota per ora: la pipeline arriva nella Fase 7
```

### Requisiti

1. **State remoto** nel bucket creato da `bootstrap/`, con blocco (`use_lockfile` su S3, nativo su GCS), una chiave o
   un prefisso diverso per ogni ambiente.
2. **Tre moduli** con variabili tipizzate e descritte, almeno una `validation` per modulo, output descritti, un
   `README.md` per modulo.
3. Gli ambienti `dev` e `prod` differiscono **solo per i valori** (numero di istanze, Multi-AZ o alta disponibilità,
   dimensioni), non per il codice.
4. **Nessun segreto** nel codice né nello state, dove possibile: password del database generata dal servizio
   (`manage_master_user_password` su RDS) oppure con `ephemeral "random_password"` e argomenti write-only.
5. Etichette comuni con `default_tags` o `default_labels`: `progetto`, `ambiente`, `gestito = "terraform"`.
6. **Versioni fissate**: `required_version`, provider con `~>`, file `.terraform.lock.hcl` versionato.
7. **Test**: almeno un file `.tftest.hcl` per modulo con `mock_provider`, eseguibile con `terraform test` senza
   credenziali.
8. `terraform fmt -check -recursive` e `terraform validate` senza errori.
9. Il database ha `prevent_destroy` in `prod` (o `deletion_protection`), non in `dev`.
10. README: diagramma, prerequisiti, come creare e distruggere ogni ambiente, costi stimati, decisioni prese.

## Parte B · Configurazione con Ansible

Sui quattro nodi di `labs/ansible-lab` (copiati nel repository in `ansible/`), scrivi:

1. un ruolo **`base`** per **tutti** i nodi, che funzioni sia su Ubuntu sia su AlmaLinux:
    - pacchetti comuni (`vim`, `curl`, `tree`, `chrony`) e fuso orario `Europe/Rome`;
    - gli utenti del team da una lista in `group_vars`, con chiave SSH pubblica (modulo
      `ansible.posix.authorized_key`) e sudo per il gruppo `operazioni` tramite un file in `/etc/sudoers.d/`
      validato con `visudo -cf %s`;
    - SSH: `PermitRootLogin no` e `PasswordAuthentication no`, con un handler che ricarica il servizio (si chiama
      `ssh` su Ubuntu e `sshd` su RHEL: usa una variabile per famiglia di sistema);
2. un ruolo **`tasks_api`** per il gruppo `app` che installa la Tasks API come nella Fase 2: utente di sistema,
   ambiente virtuale Python, dipendenze, unità systemd con gunicorn, servizio attivo e abilitato;
3. un playbook `sito.yml` che applica i ruoli, idempotente (seconda esecuzione con `changed=0`) e senza errori di
   `ansible-lint`.

## Come verificare

| Test | Risultato atteso |
|---|---|
| `terraform test` in ogni modulo | tutti i test superati, senza credenziali |
| `terraform -chdir=envs/dev apply` | ambiente funzionante: `/health` risponde dal load balancer |
| `terraform -chdir=envs/dev plan` subito dopo | `No changes` |
| Lo state | è nel bucket remoto; `grep -i password` nello state non trova la password del database |
| Due `apply` contemporanei | il secondo attende o fallisce per il blocco dello state |
| `terraform -chdir=envs/prod plan` | stesse risorse con valori diversi (non serve applicarlo) |
| `ansible-playbook sito.yml` due volte | la seconda con `changed=0` |
| `ssh -p 2221 root@127.0.0.1` | accesso negato |
| `ansible app -a "curl -s localhost:8000/health"` | `{"status":"ok"}` |
| `terraform -chdir=envs/dev destroy` | nessuna risorsa rimasta |

## Criteri di valutazione

- [ ] Struttura con moduli e ambienti separati; nessuna duplicazione di codice tra `dev` e `prod`.
- [ ] State remoto con blocco; nessun state o `.tfvars` sensibile nel repository.
- [ ] Nessun segreto in chiaro nel codice; password del database fuori dallo state.
- [ ] Test, formattazione e validazione superati.
- [ ] Ruoli Ansible idempotenti su entrambe le famiglie di sistemi.
- [ ] README completo: un collega deve poter creare l'ambiente seguendolo alla lettera.

??? success "Una possibile soluzione: il modulo database per AWS"
    ```hcl
    terraform {
      required_providers {
        aws = {
          source  = "hashicorp/aws"
          version = "~> 6.0"
        }
      }
    }

    variable "nome" {
      description = "Prefisso dei nomi"
      type        = string
    }

    variable "vpc_id" {
      description = "VPC in cui creare il database"
      type        = string
    }

    variable "subnet_ids" {
      description = "Subnet private dei dati (almeno due zone)"
      type        = list(string)

      validation {
        condition     = length(var.subnet_ids) >= 2
        error_message = "Servono almeno due subnet in zone diverse."
      }
    }

    variable "sg_applicazione" {
      description = "Security group dell'applicazione, l'unico autorizzato a connettersi"
      type        = string
    }

    variable "multi_az" {
      description = "Standby sincrono in un'altra zona"
      type        = bool
      default     = false
    }

    resource "aws_db_subnet_group" "this" {
      name       = "${var.nome}-db"
      subnet_ids = var.subnet_ids
    }

    resource "aws_security_group" "db" {
      name        = "${var.nome}-db"
      description = "PostgreSQL solo dall'applicazione"
      vpc_id      = var.vpc_id
    }

    resource "aws_vpc_security_group_ingress_rule" "postgres" {
      security_group_id            = aws_security_group.db.id
      referenced_security_group_id = var.sg_applicazione
      ip_protocol                  = "tcp"
      from_port                    = 5432
      to_port                      = 5432
    }

    resource "aws_db_instance" "this" {
      identifier     = "${var.nome}-db"
      engine         = "postgres"
      engine_version = "17"
      instance_class = "db.t4g.micro"

      allocated_storage = 20
      storage_encrypted = true

      db_name  = "tasks"
      username = "tasks"
      # La password è generata e ruotata da RDS in Secrets Manager: Terraform non la vede mai.
      manage_master_user_password = true

      db_subnet_group_name   = aws_db_subnet_group.this.name
      vpc_security_group_ids = [aws_security_group.db.id]
      multi_az               = var.multi_az
      publicly_accessible    = false

      backup_retention_period = 7
      deletion_protection     = false # in produzione: true
      skip_final_snapshot     = true  # in produzione: false
    }

    output "endpoint" {
      description = "Indirizzo del database"
      value       = aws_db_instance.this.address
    }

    output "secret_arn" {
      description = "ARN del segreto con utente e password, creato da RDS"
      value       = aws_db_instance.this.master_user_secret[0].secret_arn
    }
    ```

    **Spiegazione.**

    - Il security group del database accetta connessioni **solo** dal security group dell'applicazione
      (`referenced_security_group_id`), non da intervalli di indirizzi.
    - Con `manage_master_user_password` RDS crea un segreto in Secrets Manager con un JSON `{"username": ...,
      "password": ...}`. Lo script di avvio delle istanze lo legge con il ruolo IAM (permesso
      `secretsmanager:GetSecretValue` sul solo `secret_arn`) e compone `DATABASE_URL` con `jq` e l'endpoint, che il
      modulo `app` riceve come variabile.
    - `prevent_destroy` non si può rendere variabile (gli argomenti di `lifecycle` devono essere valori fissi): per
      proteggere solo la produzione usa `deletion_protection = var.protezione` oppure una policy di HCP Terraform.

??? success "Una possibile soluzione: il nome del servizio SSH per famiglia"
    ```yaml
    # roles/base/vars/Debian.yml
    base_servizio_ssh: ssh
    ```

    ```yaml
    # roles/base/vars/RedHat.yml
    base_servizio_ssh: sshd
    ```

    ```yaml
    # roles/base/tasks/main.yml (inizio)
    - name: Caricare le variabili della famiglia di sistema
      ansible.builtin.include_vars: "{{ ansible_facts['os_family'] }}.yml"

    - name: Vietare l'accesso SSH come root
      ansible.builtin.lineinfile:
        path: /etc/ssh/sshd_config
        regexp: '^#?PermitRootLogin'
        line: PermitRootLogin no
        validate: /usr/sbin/sshd -t -f %s
      notify: Ricaricare SSH
    ```

    ```yaml
    # roles/base/handlers/main.yml
    - name: Ricaricare SSH
      ansible.builtin.service:
        name: "{{ base_servizio_ssh }}"
        state: reloaded
    ```

    **Spiegazione.** `include_vars` carica il file giusto in base al fatto `os_family` (`Debian` o `RedHat`).
    L'argomento `validate` controlla la configurazione **prima** di sostituire il file: un errore di sintassi non ti
    chiude fuori dal server.

## Fonti

- [Terraform: Module creation, recommended pattern](https://developer.hashicorp.com/terraform/tutorials/modules/pattern-module-creation)
- [Terraform: Tests and mocks](https://developer.hashicorp.com/terraform/language/tests/mocking)
- [AWS provider: aws_db_instance, manage_master_user_password](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance)
- [Amazon RDS: Password management with Secrets Manager](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-secrets-manager.html)
- [Ansible: include_vars](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/include_vars_module.html) e [lineinfile](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/lineinfile_module.html)
- [Ansible: ansible.posix.authorized_key](https://docs.ansible.com/ansible/latest/collections/ansible/posix/authorized_key_module.html)

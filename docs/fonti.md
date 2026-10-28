# Fonti ufficiali

Il corso si basa esclusivamente su fonti **primarie e ufficiali**: la documentazione dei progetti e dei fornitori,
le pagine degli enti certificatori e gli standard. Ogni lezione cita le fonti specifiche dei suoi argomenti; qui
trovi i punti di partenza, divisi per area. Quando una lezione e la documentazione ufficiale non coincidono, vale
la documentazione: i prodotti cambiano, e conviene abituarsi presto a leggerla.

## Linux e sistemi

| Fonte | Per cosa |
|---|---|
| [Documentazione di Ubuntu Server](https://documentation.ubuntu.com/server/) | amministrazione di Ubuntu |
| [Pagine man di Ubuntu](https://manpages.ubuntu.com/) | riferimento dei comandi |
| [Documentazione di Red Hat Enterprise Linux 10](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/10) | amministrazione di RHEL e dei sistemi compatibili |
| [Documentazione di Rocky Linux](https://docs.rockylinux.org/) e [AlmaLinux](https://wiki.almalinux.org/) | distribuzioni compatibili con RHEL |
| [Documentazione di Debian](https://www.debian.org/doc/) | Debian e strumenti comuni |
| [Manuali di systemd](https://www.freedesktop.org/software/systemd/man/latest/) | servizi, timer, journal |
| [Manuale di Bash](https://www.gnu.org/software/bash/manual/) | shell e script |
| [Netplan](https://netplan.readthedocs.io/) e [nftables wiki](https://wiki.nftables.org/) | rete e firewall |
| [Linux kernel documentation](https://docs.kernel.org/) | kernel, sysctl, cgroup |

## Reti e protocolli

| Fonte | Per cosa |
|---|---|
| [RFC Editor](https://www.rfc-editor.org/) | gli standard di Internet (IP, TCP, DNS, HTTP, TLS) |
| [MDN Web Docs · HTTP](https://developer.mozilla.org/docs/Web/HTTP) | HTTP spiegato con esempi |
| [Documentazione di nginx](https://nginx.org/en/docs/) | server web e reverse proxy |
| [OpenSSH](https://www.openssh.com/manual.html) | SSH |

## Automazione e sviluppo

| Fonte | Per cosa |
|---|---|
| [Documentazione di Git](https://git-scm.com/doc) e [Pro Git](https://git-scm.com/book/it/v2) | controllo di versione |
| [Documentazione di Python](https://docs.python.org/3/) | linguaggio e libreria standard |
| [Flask](https://flask.palletsprojects.com/) e [pytest](https://docs.pytest.org/) | API e test |
| [GitHub Docs](https://docs.github.com/) | GitHub e GitHub Actions |
| [Documentazione di Ansible](https://docs.ansible.com/) e [ansible-navigator](https://ansible.readthedocs.io/projects/navigator/) | automazione dei server |
| [Documentazione di Terraform](https://developer.hashicorp.com/terraform/docs) | infrastruttura come codice |

## Container e Kubernetes

| Fonte | Per cosa |
|---|---|
| [Docker Docs](https://docs.docker.com/) | container, immagini, Compose |
| [Open Container Initiative](https://opencontainers.org/) | gli standard dei container |
| [Documentazione di Kubernetes](https://kubernetes.io/docs/) | Kubernetes |
| [kind](https://kind.sigs.k8s.io/) e [Gateway API](https://gateway-api.sigs.k8s.io/) | cluster locali e traffico in ingresso |
| [Helm](https://helm.sh/docs/), [Kustomize](https://kubectl.docs.kubernetes.io/references/kustomize/), [Argo CD](https://argo-cd.readthedocs.io/) | pacchetti, configurazione e GitOps |
| [CNCF](https://www.cncf.io/) | panorama dei progetti cloud native |

## Cloud

| Fonte | Per cosa |
|---|---|
| [Documentazione AWS](https://docs.aws.amazon.com/) | tutti i servizi AWS |
| [AWS Well-Architected Framework](https://docs.aws.amazon.com/wellarchitected/latest/framework/welcome.html) | progettazione |
| [AWS Skill Builder](https://skillbuilder.aws/) | formazione gratuita ufficiale |
| [Documentazione di Google Cloud](https://cloud.google.com/docs) | tutti i servizi Google Cloud |
| [Google Cloud Architecture Framework](https://cloud.google.com/architecture/framework) | progettazione |
| [Google Cloud Skills Boost](https://www.cloudskillsboost.google/) | formazione gratuita ufficiale |
| [Site Reliability Engineering](https://sre.google/books/) | i libri SRE di Google, gratuiti online |

## Osservabilità

| Fonte | Per cosa |
|---|---|
| [Prometheus](https://prometheus.io/docs/), [Grafana](https://grafana.com/docs/), [Loki](https://grafana.com/docs/loki/latest/) | metriche, dashboard, log |
| [OpenTelemetry](https://opentelemetry.io/docs/) | standard per tracce, metriche e log |

## Certificazioni

| Ente | Pagine |
|---|---|
| AWS | [Certificazioni AWS](https://aws.amazon.com/certification/) |
| Google Cloud | [Certificazioni Google Cloud](https://cloud.google.com/learn/certification) |
| HashiCorp | [Terraform Associate](https://developer.hashicorp.com/certifications/terraform-associate) |
| Linux Foundation | [CKA](https://training.linuxfoundation.org/certification/certified-kubernetes-administrator-cka/), [LFCS](https://training.linuxfoundation.org/certification/linux-foundation-certified-sysadmin-lfcs/), [curriculum CNCF](https://github.com/cncf/curriculum) |
| Red Hat | [EX200](https://www.redhat.com/en/services/training/ex200-red-hat-certified-system-administrator-rhcsa-exam), [EX294](https://www.redhat.com/en/services/training/ex294-red-hat-certified-engineer-rhce-exam-red-hat-enterprise-linux) |
| Canonical | [Canonical Academy](https://canonical.com/academy) |

## Come valutare una fonte che non è in elenco

1. **Chi la pubblica?** Preferisci il progetto o il fornitore stesso. Blog e video possono aiutare a capire, ma
   verifica sempre i dettagli sulla documentazione ufficiale.
2. **Di quando è?** Cloud e Kubernetes cambiano ogni pochi mesi: controlla la data e la versione del prodotto.
3. **Si può verificare?** Una buona fonte mostra comandi e risultati che puoi riprodurre nel tuo laboratorio.
4. **Ti chiede qualcosa di strano?** Diffida di siti che vendono domande d'esame reali o ti chiedono credenziali.

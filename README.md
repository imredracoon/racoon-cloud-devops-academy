# Racoon Cloud & DevOps Academy

Un corso **gratuito e interattivo**, in italiano, per diventare Cloud/DevOps engineer **partendo da zero**: Linux e
reti, Git, Bash e Python, container, AWS e Google Cloud, Terraform, CI/CD e osservabilità, Kubernetes.

- **9 fasi** con lezioni, esercizi con soluzione, laboratori e un progetto per il portfolio in ogni fase.
- **Percorsi doppi** Debian/Ubuntu e RHEL (Rocky/AlmaLinux) per Linux, e AWS e Google Cloud per il cloud.
- **10 appendici di certificazione** (AWS, Google Cloud, Terraform, CKA, LFCS, RHCSA, RHCE in Ansible, Canonical)
  con simulazioni casuali: oltre 730 domande e 155 task pratici originali.
- Solo **fonti ufficiali**, sempre citate.
- Tutto si studia in un ambiente riproducibile: un **Dev Container** con gli strumenti già installati, che funziona
  con VS Code o con GitHub Codespaces.

## Leggere il corso

Il sito è pubblicato con GitHub Pages: <https://imredracoon.github.io/racoon-cloud-devops-academy/>.

Per vederlo sul tuo computer:

```bash
python3 -m venv .venv && . .venv/bin/activate
pip install -r requirements.txt
mkdocs serve        # poi apri http://127.0.0.1:8000
```

## Struttura del repository

| Percorso | Contenuto |
|---|---|
| `docs/` | le lezioni (Markdown), le appendici e le banche delle simulazioni (`docs/assets/banks/`) |
| `labs/` | i laboratori: Tasks API, Terraform, Ansible, Kubernetes |
| `lab-vm/` | le macchine virtuali per i laboratori Linux e kubeadm |
| `.devcontainer/` | l'ambiente di studio riproducibile |
| `scripts/check_banks.py` | i controlli automatici delle banche domande |

## Contribuire

Segnalazioni e proposte sono benvenute tramite issue e pull request. Prima di proporre una modifica esegui:

```bash
python3 scripts/check_banks.py
mkdocs build --strict
```

Le domande delle simulazioni devono essere **originali** e basate sugli obiettivi pubblicati dagli enti
certificatori: non si accettano domande d'esame reali.

## Licenza

[GNU GPL v3](LICENSE).

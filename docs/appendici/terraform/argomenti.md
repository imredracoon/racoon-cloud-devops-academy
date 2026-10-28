# Argomenti specifici · Terraform Associate (004)

Gli obiettivi ufficiali (1a-8d) con le lezioni della Fase 6. HashiCorp non pubblica pesi per obiettivo: quelli
indicati sono le stime usate dalla simulazione, in base al numero di sotto-obiettivi.

## 1 · Infrastructure as Code with Terraform (~8%)

| Obiettivo | Lezioni |
|---|---|
| 1a Spiegare cos'è l'IaC · 1b I vantaggi dei modelli IaC | [6.1](../../fase-6/01-iac.md) |
| 1c Come Terraform gestisce flussi multicloud, ibridi e indipendenti dal servizio | [6.1](../../fase-6/01-iac.md), [6.2](../../fase-6/02-primo-progetto.md) |

## 2 · Terraform fundamentals (~14%)

| Obiettivo | Lezioni |
|---|---|
| 2a Installare e versionare i provider (`required_providers`, vincoli `~>`, `.terraform.lock.hcl`, `init -upgrade`) | [6.2](../../fase-6/02-primo-progetto.md) |
| 2b Come Terraform usa i provider · 2c Configurazione con più provider e alias | [6.2](../../fase-6/02-primo-progetto.md), [6.4](../../fase-6/04-risorse.md) |
| 2d Come Terraform usa e gestisce lo state | [6.6](../../fase-6/06-state.md) |

## 3 · Core Terraform workflow (~16%)

| Obiettivo | Lezioni |
|---|---|
| 3a Il workflow (write, plan, apply) · 3b `init` · 3c `validate` | [6.3](../../fase-6/03-workflow.md) |
| 3d Generare e leggere un piano (`plan -out`, simboli `+ - ~ -/+`) · 3e `apply` · 3f `destroy` | [6.3](../../fase-6/03-workflow.md) |
| 3g Formattazione e stile (`fmt -check -recursive`) | [6.3](../../fase-6/03-workflow.md) |

## 4 · Terraform configuration (~22%)

| Obiettivo | Lezioni |
|---|---|
| 4a Blocchi `resource` e `data` · 4b Riferimenti tra risorse | [6.4](../../fase-6/04-risorse.md) |
| 4c Variabili e output (precedenza dei valori) · 4d Tipi complessi | [6.5](../../fase-6/05-variabili.md) |
| 4e Espressioni e funzioni (`for`, condizionali, splat, `count`, `for_each`, `dynamic`) | [6.5](../../fase-6/05-variabili.md) |
| 4f Dipendenze (implicite e `depends_on`) | [6.4](../../fase-6/04-risorse.md) |
| 4g Condizioni personalizzate (`validation`, `precondition`, `postcondition`, `check`) | [6.5](../../fase-6/05-variabili.md) |
| 4h Dati sensibili: `sensitive`, valori effimeri, argomenti write-only, Vault | [6.8](../../fase-6/08-segreti.md) |

## 5 · Terraform modules (~10%)

| Obiettivo | Lezioni |
|---|---|
| 5a Sorgenti dei moduli · 5b Ambito delle variabili · 5c Usare i moduli · 5d Versioni dei moduli | [6.7](../../fase-6/07-moduli.md) |

## 6 · Terraform state management (~12%)

| Obiettivo | Lezioni |
|---|---|
| 6a Backend locale · 6b Lock dello state · 6c Backend remoto · 6d Deriva e state | [6.6](../../fase-6/06-state.md) |

## 7 · Maintain infrastructure with Terraform (~8%)

| Obiettivo | Lezioni |
|---|---|
| 7a Importare infrastruttura esistente · 7b Ispezionare lo state con la CLI · 7c Log dettagliati (`TF_LOG`) | [6.6](../../fase-6/06-state.md) |

## 8 · HCP Terraform (~10%)

| Obiettivo | Lezioni |
|---|---|
| 8a Creare infrastruttura con HCP Terraform · 8b Collaborazione e governance · 8c Workspace e progetti · 8d Integrazione (blocco `cloud`, `terraform login`) | [6.9](../../fase-6/09-hcp.md) |

## Dove sbagliano più spesso i candidati

- **`sensitive = true`** nasconde il valore nell'output ma **non** lo toglie dallo state.
- **`terraform validate`** non contatta i provider remoti e non controlla le credenziali.
- **`count` contro `for_each`**: togliere un elemento in mezzo a una lista con `count` ricrea le risorse
  successive; `for_each` usa chiavi stabili.
- **Workspace della CLI e workspace di HCP Terraform** sono concetti diversi.
- **`-target` e `-replace`**: il primo è per casi eccezionali, il secondo sostituisce il vecchio `taint`.

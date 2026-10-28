# HashiCorp Certified: Terraform Associate (004)

Verifica che conosci i concetti dell'infrastruttura come codice e sai usare Terraform nel lavoro quotidiano:
provider, workflow, linguaggio di configurazione, moduli, state e HCP Terraform.

| | |
|---|---|
| **Ente** | HashiCorp (gruppo IBM) |
| **Versione** | 004, basata su **Terraform 1.12** |
| **Formato** | domande a scelta multipla, online con sorvegliante |
| **Durata** | 1 ora |
| **Costo** | 70,50 USD più tasse locali; il secondo tentativo gratuito **non** è incluso |
| **Validità** | 2 anni |
| **Lingua** | inglese |
| **Pagina ufficiale** | [developer.hashicorp.com/certifications/terraform-associate](https://developer.hashicorp.com/certifications/terraform-associate) |

Dati verificati il 28 ottobre 2026 sulla pagina ufficiale, che non pubblica né il numero di domande né la soglia.
La simulazione usa 57 domande e una soglia indicativa del 70%.

## Quando sostenerla

Dopo la [Fase 6](../../fase-6/index.md) e il suo progetto. Devi aver usato davvero `init`, `plan`, `apply`, moduli,
uno state remoto con lock e almeno un `import`.

## Piano di studio (3 settimane)

| Settimana | Attività |
|---|---|
| 1 | Obiettivi 1-3: ripasso 6.1-6.3. Rifai il primo progetto da zero leggendo ogni output |
| 2 | Obiettivi 4-6: ripasso 6.4-6.8. Allenati con `terraform console` su funzioni ed espressioni |
| 3 | Obiettivi 7-8: ripasso 6.6 e 6.9. Esami completi e ripasso degli errori |

## Novità della versione 004 da non trascurare

- blocchi `import` con generazione della configurazione, blocchi `moved` e `removed`;
- condizioni personalizzate: `validation`, `precondition`, `postcondition` e blocchi `check`;
- valori **effimeri** e argomenti **write-only** per tenere i segreti fuori dallo state, e uso di Vault;
- lock nativo dello state con il backend S3 (`use_lockfile`);
- HCP Terraform: progetti, variable set, policy, credenziali dinamiche, controlli di salute.

## Risorse ufficiali gratuite

- [Percorso di preparazione ufficiale](https://developer.hashicorp.com/terraform/tutorials/certification-004).
- [Documentazione del linguaggio](https://developer.hashicorp.com/terraform/language) e
  [della CLI](https://developer.hashicorp.com/terraform/cli).

## In questa appendice

- [Argomenti specifici](argomenti.md)
- [Simulazione d'esame](simulazione.md): oltre 130 domande originali

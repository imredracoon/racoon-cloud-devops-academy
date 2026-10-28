# Fase 4 · Cloud: AWS e Google Cloud

!!! obiettivi "Al termine di questa fase saprai"
    - spiegare cos'è il cloud computing, i suoi modelli di servizio e di distribuzione, e il modello di responsabilità
      condivisa;
    - orientarti nell'infrastruttura globale di AWS e Google Cloud: regioni, zone, edge;
    - creare e proteggere un account, con MFA, budget e utenti non privilegiati;
    - riconoscere i servizi principali di calcolo, storage, rete, database e intelligenza artificiale e scegliere tra
      essi;
    - capire come si paga il cloud e come tenere sotto controllo i costi;
    - usare le CLI `aws` e `gcloud` per creare e distruggere risorse.

| | |
|---|---|
| **Durata indicativa** | 1,5 mesi, circa 50 ore |
| **Prerequisiti** | Fasi 1-3 |
| **Ambiente** | Dev Container (contiene `aws` e `gcloud`), console web, account AWS e/o Google Cloud |
| **Progetto** | [Sito statico nel cloud](progetto.md) |
| **Certificazioni collegate** | [AWS Cloud Practitioner](../appendici/aws-clf/index.md), [Google Cloud Digital Leader](../appendici/gcp-cdl/index.md) |

## Due cloud, gli stessi concetti

Da questa fase le lezioni presentano **AWS** e **Google Cloud** in schede affiancate, come questa:

=== "AWS"

    Amazon Web Services, il fornitore con la quota di mercato maggiore e il cloud principale del corso.

=== "Google Cloud"

    Il cloud di Google, molto forte su dati, analisi, intelligenza artificiale e Kubernetes (che Google ha creato).

La scelta della scheda resta memorizzata in tutte le pagine. I concetti (regioni, identità, macchine virtuali,
storage a oggetti, reti private) sono gli stessi: cambiano i nomi e i comandi. Se è la prima volta, segui AWS e
leggi la scheda Google Cloud per confronto.

!!! danger "Prima di creare risorse"
    Leggi [Account cloud in sicurezza](../ambiente/account-cloud.md). Ogni laboratorio di questa fase termina con una
    sezione **Pulizia**: eseguila sempre. Tutti i laboratori restano nei limiti dei crediti per i nuovi clienti.

## Le lezioni

1. [Cos'è il cloud computing](01-cloud-computing.md)
2. [Infrastruttura globale](02-infrastruttura-globale.md)
3. [Account, progetti e fatturazione](03-account.md)
4. [Identità e accessi (IAM)](04-iam.md)
5. [Calcolo](05-compute.md)
6. [Storage](06-storage.md)
7. [Reti nel cloud](07-rete.md)
8. [Database, analisi e AI](08-database-ai.md)
9. [Costi, supporto e buone pratiche](09-costi.md)
10. [Le CLI: aws e gcloud](10-cli.md)

Poi il [progetto](progetto.md) e la [verifica di fine fase](verifica.md).

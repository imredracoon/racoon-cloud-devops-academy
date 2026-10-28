# Come funziona il corso

Questa pagina spiega come è costruito il corso e come ottenere il massimo da ogni lezione. Leggila per intero
prima di iniziare: ci vogliono dieci minuti e ti eviteranno molti dubbi più avanti.

## A chi si rivolge

Il corso è pensato per chi **parte da zero**: non serve sapere programmare, non serve conoscere Linux, non serve
avere lavorato nell'informatica. Serve invece:

- un computer con almeno **8 GB di RAM** (16 GB è meglio) e una connessione a Internet;
- costanza: **8-10 ore a settimana** per circa 16 mesi;
- la disponibilità a sbagliare. Nei laboratori romperai cose: è previsto ed è il modo più rapido per imparare.

Se hai già esperienza, segui comunque tutte le lezioni nell'ordine. I contenuti si appoggiano l'uno sull'altro
e anche chi conosce un argomento trova spesso un dettaglio che gli mancava.

## Come è fatta una fase

Il corso è diviso in **nove fasi** (dalla 0 alla 8). Ogni fase contiene:

| Elemento | A cosa serve |
|---|---|
| **Pagina introduttiva** | Obiettivi, durata stimata, prerequisiti e cosa saprai fare alla fine. |
| **Lezioni** | Teoria spiegata passo per passo, esempi commentati, esercizi e una verifica rapida. |
| **Progetto** | Un lavoro completo da pubblicare su GitHub: diventa parte del tuo portfolio. |
| **Verifica di fine fase** | Un quiz che copre tutta la fase. Passa alla fase successiva quando superi l'80%. |

## Come è fatta una lezione

Ogni lezione segue sempre la stessa struttura, così sai sempre dove trovare le cose.

!!! obiettivi "Obiettivi"
    All'inizio di ogni lezione trovi un riquadro come questo, con l'elenco di quello che saprai fare alla fine.

Poi vengono la **teoria** e gli **esempi guidati**. I comandi da scrivere nel terminale sono in blocchi come questo:

```bash
echo "Ciao, mondo"   # (1)!
```

1.  I numeri cerchiati nei blocchi di codice sono **annotazioni**: cliccaci sopra per leggere la spiegazione del pezzo di comando corrispondente.

Quando un blocco mostra anche il risultato, il comando è preceduto dal simbolo `$` (che **non** va digitato) e
le righe senza `$` sono l'output:

```console
$ echo "Ciao, mondo"
Ciao, mondo
```

In alto a destra di ogni blocco c'è il pulsante per **copiare** il contenuto.

### Esercizi e soluzioni

Gli esercizi sono in riquadri viola. La soluzione è **nascosta** in un riquadro chiuso subito sotto:

!!! esercizio "Esercizio di esempio"
    Scrivi il comando che stampa sullo schermo la parola `prova`.

??? success "Soluzione"
    ```bash
    echo prova
    ```
    **Spiegazione.** `echo` stampa sullo schermo il testo che riceve come argomento.

!!! tip "La regola d'oro"
    Apri la soluzione **solo dopo** aver provato davvero. Se sei bloccato da più di 20 minuti, rileggi la teoria,
    cerca nella documentazione indicata nelle fonti e solo allora guarda la soluzione. Dopo averla letta, **rifai
    l'esercizio da capo senza guardarla**. Un sito pubblico non può impedirti di sbirciare: la disciplina è tua,
    ed è la stessa che ti servirà sul lavoro.

### Laboratori

I laboratori sono esercizi più lunghi, da svolgere nell'ambiente di laboratorio. Sono in riquadri verde acqua:

!!! laboratorio "Laboratorio di esempio"
    Ogni laboratorio indica dove svolgerlo (Dev Container, VM Linux o cloud), il tempo stimato, i passi da
    seguire e come verificare il risultato.

### Verifica rapida

Ogni lezione si chiude con un piccolo quiz interattivo. Scegli la risposta e premi **Verifica**: vedrai subito se
è corretta e perché. Prova questo:

<script type="application/json" class="quiz-data">
{"title":"Verifica rapida di esempio","questions":[
 {"q":"Dove si trova la soluzione di un esercizio?","options":["Non c'è, gli esercizi non hanno soluzione","In un riquadro chiuso subito sotto l'esercizio","Solo alla fine del corso","In un file da scaricare"],"answer":[1],"explanation":"Ogni esercizio ha la soluzione spiegata in un riquadro chiuso che puoi aprire dopo aver provato."},
 {"q":"Quali elementi trovi in **ogni** fase?","options":["Lezioni","Un progetto per il portfolio","Una verifica di fine fase","Un esame a pagamento"],"answer":[0,1,2],"explanation":"Ogni fase ha lezioni, un progetto e una verifica finale. Gli esami di certificazione sono facoltativi e sono descritti nelle appendici."}
]}
</script>

### Fonti

In fondo a ogni lezione c'è la sezione **Fonti**: sono i documenti ufficiali da cui sono tratte le informazioni.
Abituati a consultarli: leggere la documentazione ufficiale è una delle competenze più importanti di questo mestiere.

## La settimana tipo

| Momento | Durata | Attività |
|---|---|---|
| Sera 1 | 2 ore | Studio della teoria e degli esempi guidati, con appunti. |
| Sera 2 | 2 ore | Esercizi e verifica rapida. Ripasso della lezione precedente. |
| Weekend | 4-6 ore | Laboratori e progetto. Ogni weekend si chiude con almeno un commit su GitHub. |

## Prendere appunti

Dalla [Fase 0](fase-0/index.md) creerai un repository GitHub personale per i tuoi appunti, scritti in Markdown. Gli appunti scritti
con parole tue fissano i concetti molto meglio della sola lettura, e il repository diventa una prova della tua
costanza.

## Le certificazioni

Le certificazioni sono **facoltative** ma utili per farsi notare. La sezione [Certificazioni](appendici/index.md)
spiega quando affrontarle, quanto costano e come prepararle con argomenti specifici e simulazioni d'esame.

## Quando chiedere aiuto

Se trovi un errore nel corso o un passaggio non chiaro, apri una *issue* sul repository GitHub del corso. Descrivi
in quale pagina ti trovi, cosa hai fatto, cosa ti aspettavi e cosa è successo invece. Imparare a descrivere bene
un problema è già un esercizio di troubleshooting.

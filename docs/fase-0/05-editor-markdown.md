# 0.5 Editor, file di testo e Markdown

!!! obiettivi "Obiettivi"
    - Distinguere un file di testo da un documento di un elaboratore di testi.
    - Modificare file dal terminale con `nano`.
    - Usare Visual Studio Code come editor principale.
    - Scrivere documentazione in Markdown.

## File di testo semplice

Nel lavoro DevOps quasi tutto è **testo semplice** (*plain text*): configurazioni, script, codice, documentazione.
Un file di testo contiene solo caratteri, senza formattazione nascosta. Un documento Word, invece, contiene anche
informazioni su caratteri, colori e impaginazione in un formato che solo programmi specifici sanno leggere.

Per questo i file di configurazione **non si modificano mai** con Word o simili: si usano gli **editor di testo**.

### Codifica e fine riga

- La **codifica** dice come i caratteri diventano byte. Lo standard è **UTF-8**, che rappresenta qualsiasi lingua,
  accenti compresi. Usa sempre UTF-8.
- Il **fine riga** è diverso tra sistemi: Linux e macOS usano un carattere (*LF*), Windows due (*CRLF*). Uno script
  salvato con fine riga Windows può non funzionare su Linux. VS Code mostra il tipo in basso a destra e permette di
  cambiarlo: per tutto il corso usa **LF**.

## nano: l'editor del terminale

Quando sei collegato a un server, spesso hai solo il terminale. **nano** è l'editor più semplice da usare:

```bash
nano saluto.txt
```

Scrivi il testo, poi:

| Tasti | Azione |
|---|---|
| ++ctrl+o++ poi ++enter++ | Salva (*Write Out*) |
| ++ctrl+x++ | Esce (se ci sono modifiche chiede se salvarle) |
| ++ctrl+w++ | Cerca |
| ++ctrl+k++ / ++ctrl+u++ | Taglia la riga / incolla |

In basso nano mostra sempre i comandi principali; il simbolo `^` significa ++ctrl++.

!!! note "E vim?"
    Su molti server trovi anche **vim**, potentissimo ma con una curva di apprendimento ripida. Per ora basta
    sapere come uscirne: premi ++esc++, poi scrivi `:q!` e ++enter++ (esce senza salvare). Se vuoi impararlo,
    il comando `vimtutor` avvia un tutorial interattivo di 30 minuti.

## Visual Studio Code

**Visual Studio Code** (VS Code) è l'editor che userai per la maggior parte del corso. È gratuito, funziona su
tutti i sistemi operativi e si integra con Git, i Dev Container e il terminale.

Le funzioni da conoscere subito:

- **Esplora risorse** (++ctrl+shift+e++): l'albero di file e cartelle.
- **Riquadro comandi** (++f1++ o ++ctrl+shift+p++): da qui si raggiunge qualsiasi funzione scrivendone il nome.
- **Terminale integrato** (++ctrl+grave++): un terminale dentro l'editor.
- **Anteprima Markdown** (++ctrl+shift+v++): mostra come apparirà il testo formattato.

Su macOS usa ++cmd++ al posto di ++ctrl++.

## Markdown

**Markdown** è un modo di scrivere testo formattato usando solo caratteri normali. È lo standard per la
documentazione tecnica: i file `README.md` di GitHub, questo stesso corso e la documentazione di molti progetti
sono scritti in Markdown.

````markdown
# Titolo principale
## Sottotitolo

Testo normale con **grassetto**, *corsivo* e `codice in linea`.

- elenco puntato
- secondo punto

1. elenco numerato
2. secondo passo

[Un collegamento](https://www.example.com)

```bash
echo "un blocco di codice"
```

| Colonna A | Colonna B |
|---|---|
| valore | valore |

> Una citazione
````

Il risultato:

> **Titolo principale**, **sottotitolo**, testo con **grassetto**, *corsivo* e `codice in linea`, elenchi,
> un [collegamento](https://www.example.com), un blocco di codice evidenziato e una tabella.

GitHub usa una variante chiamata **GitHub Flavored Markdown** (GFM), che aggiunge tabelle, liste di controllo
(`- [ ] da fare`) e altro.

## Esercizi

!!! esercizio "Esercizio 1 · nano"
    Nel terminale crea con `nano` il file `spesa.txt` con tre righe (pane, latte, caffè), salvalo, esci e mostralo
    con `cat`. Poi riaprilo, aggiungi `uova` in fondo e salva di nuovo.

??? success "Soluzione"
    ```bash
    nano spesa.txt      # scrivi le tre righe, poi Ctrl+O, Invio, Ctrl+X
    cat spesa.txt
    nano spesa.txt      # vai in fondo, scrivi "uova", Ctrl+O, Invio, Ctrl+X
    cat spesa.txt
    ```
    **Spiegazione.** Se all'uscita nano chiede *Save modified buffer?*, rispondi `Y` per salvare.

!!! esercizio "Esercizio 2 · Il tuo primo README"
    Scrivi in Markdown un file `README.md` che presenti te stesso e i tuoi obiettivi di studio. Deve contenere:
    un titolo, due sottotitoli, un elenco puntato, un elenco numerato, un collegamento e una tabella con tre righe.
    Controlla il risultato con l'anteprima di VS Code oppure con il sito [markdownlivepreview.com](https://markdownlivepreview.com/).

??? success "Soluzione"
    Un esempio possibile:

    ```markdown
    # Il mio percorso Cloud / DevOps

    ## Chi sono
    - Nome: Mario
    - Oggi: impiegato amministrativo
    - Passione: tecnologia e automazione

    ## Obiettivi
    1. Completare la Fase 0 entro tre settimane
    2. Superare la AWS Cloud Practitioner
    3. Trovare il primo lavoro come Junior Cloud Engineer

    Studio con la [Racoon Cloud & DevOps Academy](https://imredracoon.github.io/racoon-cloud-devops-academy/).

    | Fase | Inizio | Fine |
    |---|---|---|
    | 0 | novembre | novembre |
    | 1 | dicembre | febbraio |
    ```
    **Spiegazione.** Ogni elemento ha una sintassi precisa: `#` seguito da uno spazio per i titoli, `-` per i
    punti, `1.` per i numeri, `[testo](indirizzo)` per i link, `|` per le colonne delle tabelle. Questo file
    diventerà la base del README del tuo repository di appunti nel progetto della fase.

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Perché i file di configurazione non si modificano con un elaboratore di testi come Word?","options":["Perché Word è a pagamento","Perché aggiunge formattazione nascosta e salva in un formato diverso dal testo semplice","Perché Word non supporta l'italiano","Perché i file di configurazione sono in Markdown"],"answer":[1],"explanation":"I file di configurazione devono contenere solo testo semplice."},
 {"q":"Come si salva un file in nano?","options":["Ctrl+S","Ctrl+O e poi Invio","Esc e :w","Ctrl+X senza conferma"],"answer":[1],"explanation":"Ctrl+O (Write Out) salva; Ctrl+X esce."},
 {"q":"Quale codifica dei caratteri conviene usare sempre?","options":["ASCII","UTF-8","Latin-1","UTF-16"],"answer":[1],"explanation":"UTF-8 è lo standard e rappresenta qualsiasi carattere."},
 {"q":"In Markdown, come si scrive un titolo di secondo livello?","options":["## Titolo","**Titolo**","= Titolo =","<h2> è l'unico modo"],"answer":[0],"explanation":"Il numero di # indica il livello del titolo."},
 {"q":"Come si esce da vim senza salvare?","options":["Ctrl+C","Esc, poi :q! e Invio","Ctrl+X","q"],"answer":[1],"explanation":":q! chiude forzando l'uscita e scartando le modifiche."}
]}
</script>

## Fonti

- [GNU nano: documentazione](https://www.nano-editor.org/docs.php)
- [Visual Studio Code: documentazione](https://code.visualstudio.com/docs)
- [CommonMark: specifica di Markdown](https://commonmark.org/)
- [GitHub Docs: Basic writing and formatting syntax](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax)
- [GitHub Flavored Markdown Spec](https://github.github.com/gfm/)

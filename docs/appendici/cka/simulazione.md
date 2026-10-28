# Simulazione pratica · CKA

Banca di 39 task pratici originali sui cinque domini del programma CKA (Kubernetes v1.37). L'esame pratico estrae **17 task** in **2 ore**.

<div class="exam-sim" data-bank="cka"></div>

## Come usarla

1. Prepara il laboratorio indicato in ogni task (le istruzioni sono nelle lezioni collegate negli
   [argomenti specifici](argomenti.md)). Alcuni task presuppongono oggetti già presenti: creali tu prima di iniziare,
   è un ottimo esercizio.
2. Nella modalità **Allenamento** svolgi un task alla volta e confronta il tuo lavoro con la soluzione spiegata.
3. Nell'**Esame pratico** lavora a tempo, segna i task svolti e quelli da rivedere; alla consegna esegui i comandi
   di verifica e assegna i punti con onestà.
4. Sei pronto quando superi due esami pratici di fila **senza consultare le soluzioni** e con tempo di avanzo.

!!! tip "Strategia di tempo"
    Leggi tutti i task all'inizio. Fai subito quelli brevi e sicuri, segna gli altri e torna dopo. Non restare
    bloccato più di 8-10 minuti su un task: è meglio prendere i punti facili che perdere quelli di tre task per uno
    difficile. Verifica sempre il risultato prima di passare oltre.

!!! note "Consiglio specifico"
    All'esame imposta subito `alias k=kubectl`, l'autocompletamento e `export do='--dry-run=client -o yaml'`; ogni task indica il contesto (`kubectl config use-context`) da usare: **usalo sempre**, altrimenti lavori sul cluster sbagliato.

Torna alla [scheda dell'esame](index.md) per iscrizione e regole, o agli [argomenti specifici](argomenti.md) per
ripassare.

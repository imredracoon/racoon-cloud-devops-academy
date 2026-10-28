"""Configurazione di Gunicorn (dalla Fase 7), letta automaticamente dalla cartella di lavoro.

Ogni worker di Gunicorn è un processo separato con i suoi contatori. In modalità
multiprocesso prometheus_client scrive i valori in file nella cartella
PROMETHEUS_MULTIPROC_DIR e /metrics li somma.
"""

import os
import shutil

# La variabile va impostata PRIMA che qualunque modulo importi prometheus_client:
# la modalità multiprocesso si decide al momento dell'import.
cartella = os.environ.setdefault("PROMETHEUS_MULTIPROC_DIR", "/tmp/metriche-tasks-api")
shutil.rmtree(cartella, ignore_errors=True)
os.makedirs(cartella)


def child_exit(server, worker):
    from prometheus_client import multiprocess

    # Un worker terminato non deve lasciare valori "vivi" (per esempio nei gauge).
    multiprocess.mark_process_dead(worker.pid)

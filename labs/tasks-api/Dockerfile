# Immagine della Tasks API (Fase 3)
FROM python:3.12-slim

# Python: niente file .pyc e output non bufferizzato (i log arrivano subito)
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app

# Utente senza privilegi con cui girerà l'applicazione
RUN useradd --system --uid 10001 --home-dir /app app \
    && install -d -o app -g app /app/instance

# Prima le dipendenze: questo strato resta in cache finché requirements.txt non cambia
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

USER app
EXPOSE 8000

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
  CMD ["python", "-c", "import urllib.request; urllib.request.urlopen('http://127.0.0.1:8000/health', timeout=2)"]

CMD ["gunicorn", "--bind", "0.0.0.0:8000", "--workers", "2", "--preload", "--access-logfile", "-", "app:create_app()"]

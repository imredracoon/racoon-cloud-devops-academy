"""Tasks API: l'applicazione d'esempio del corso.

Una piccola API REST che gestisce un elenco di attività.
Il database si sceglie con la variabile d'ambiente DATABASE_URL:
  - non impostata: SQLite nel file tasks.db (Fase 2)
  - postgresql+psycopg://utente:password@host:5432/db (dalla Fase 3)
"""
import os
import socket

from flask import Flask, jsonify, request
from flask_sqlalchemy import SQLAlchemy

db = SQLAlchemy()


class Task(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.String(200), nullable=False)
    done = db.Column(db.Boolean, nullable=False, default=False)

    def to_dict(self):
        return {"id": self.id, "title": self.title, "done": self.done}


def create_app(database_url=None):
    app = Flask(__name__)
    app.config["SQLALCHEMY_DATABASE_URI"] = (
        database_url or os.environ.get("DATABASE_URL") or "sqlite:///tasks.db"
    )
    db.init_app(app)
    with app.app_context():
        db.create_all()

    @app.get("/health")
    def health():
        try:
            db.session.execute(db.text("SELECT 1"))
            return jsonify(status="ok", host=socket.gethostname())
        except Exception:  # il database non risponde
            return jsonify(status="error"), 503

    @app.get("/tasks")
    def list_tasks():
        return jsonify([t.to_dict() for t in Task.query.order_by(Task.id)])

    @app.post("/tasks")
    def create_task():
        data = request.get_json(silent=True) or {}
        title = str(data.get("title", "")).strip()
        if not title:
            return jsonify(error="il campo 'title' è obbligatorio"), 400
        task = Task(title=title[:200])
        db.session.add(task)
        db.session.commit()
        return jsonify(task.to_dict()), 201

    @app.get("/tasks/<int:task_id>")
    def get_task(task_id):
        task = db.session.get(Task, task_id)
        if task is None:
            return jsonify(error="attività non trovata"), 404
        return jsonify(task.to_dict())

    @app.patch("/tasks/<int:task_id>")
    def update_task(task_id):
        task = db.session.get(Task, task_id)
        if task is None:
            return jsonify(error="attività non trovata"), 404
        data = request.get_json(silent=True) or {}
        if "title" in data:
            title = str(data["title"]).strip()
            if not title:
                return jsonify(error="il titolo non può essere vuoto"), 400
            task.title = title[:200]
        if "done" in data:
            task.done = bool(data["done"])
        db.session.commit()
        return jsonify(task.to_dict())

    @app.delete("/tasks/<int:task_id>")
    def delete_task(task_id):
        task = db.session.get(Task, task_id)
        if task is None:
            return jsonify(error="attività non trovata"), 404
        db.session.delete(task)
        db.session.commit()
        return "", 204

    return app


if __name__ == "__main__":
    create_app().run(host="127.0.0.1", port=8000, debug=True)

import pytest

from app import create_app


@pytest.fixture
def client():
    app = create_app("sqlite:///:memory:")
    app.config["TESTING"] = True
    with app.test_client() as client:
        yield client


def test_health(client):
    r = client.get("/health")
    assert r.status_code == 200
    assert r.get_json()["status"] == "ok"


def test_create_and_list(client):
    r = client.post("/tasks", json={"title": "Studiare Git"})
    assert r.status_code == 201
    assert r.get_json() == {"id": 1, "title": "Studiare Git", "done": False}
    assert client.get("/tasks").get_json() == [{"id": 1, "title": "Studiare Git", "done": False}]


def test_title_required(client):
    assert client.post("/tasks", json={}).status_code == 400
    assert client.post("/tasks", json={"title": "   "}).status_code == 400


def test_update_and_delete(client):
    client.post("/tasks", json={"title": "Scrivere test"})
    r = client.patch("/tasks/1", json={"done": True})
    assert r.get_json()["done"] is True
    assert client.delete("/tasks/1").status_code == 204
    assert client.get("/tasks/1").status_code == 404

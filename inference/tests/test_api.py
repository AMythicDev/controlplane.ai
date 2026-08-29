from fastapi.testclient import TestClient
import numpy as np
from inference.main import app

client = TestClient(app)


def test_root_endpoint():
    response = client.get("/")
    assert response.status_code == 200
    data = response.json()
    assert "message" in data
    assert data["health"] == "/health"
    assert data["docs"] == "/docs"


def test_health_endpoint():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "healthy"}


def test_toxicity_with_object_payload():
    payload = {
        "texts": [
            "You are awesome and helpful!",
            "I hate you so much, you are stupid and ugly.",
        ]
    }
    response = client.post("/toxicity", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data, list)
    assert len(data) == 2

    # Check non-toxic text
    clean_item = data[0]
    assert clean_item["text"] == "You are awesome and helpful!"
    assert clean_item["toxicity"] < 0.1

    # Check toxic text
    toxic_item = data[1]
    assert toxic_item["text"] == "I hate you so much, you are stupid and ugly."
    assert toxic_item["toxicity"] > 0.8
    assert "insult" in toxic_item
    assert "obscene" in toxic_item
    assert "threat" in toxic_item
    assert "severe_toxicity" in toxic_item
    assert "identity_attack" in toxic_item


def test_toxicity_with_array_payload():
    payload = ["Have a nice day!"]
    response = client.post("/toxicity", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert len(data) == 1
    assert data[0]["text"] == "Have a nice day!"
    assert data[0]["toxicity"] < 0.1


def test_toxicity_with_empty_payload():
    response = client.post("/toxicity", json={"texts": []})
    assert response.status_code == 200
    assert response.json() == []

    response_list = client.post("/toxicity", json=[])
    assert response_list.status_code == 200
    assert response_list.json() == []


def test_toxicity_invalid_item_type():
    response = client.post("/toxicity", json={"texts": [123, "test"]})
    assert response.status_code == 422


def test_nli_endpoint():
    payload = {
        "premise": "A person is playing the guitar on stage.",
        "hypothesis": "A person is playing music.",
    }
    response = client.post("/nli", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert "label" in data
    assert "score" in data
    assert "entailment_prob" in data
    assert data["label"] == "entailment"


def test_embed_saves_to_qdrant():
    payload = {
        "request": "What is the capital of Germany?",
        "response": "The capital of Germany is Berlin.",
    }
    response = client.post("/embed", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "stored"
    assert "id" in data
    assert len(data["request_embedding"]) == 384
    assert len(data["response_embedding"]) == 384
    assert np.isclose(np.linalg.norm(data["request_embedding"]), 1.0, atol=1e-3)
    assert np.isclose(np.linalg.norm(data["response_embedding"]), 1.0, atol=1e-3)


def test_query_cache_hit():
    # Query with exact or closely matching request
    payload = {
        "request": "What is the capital of Germany?",
        "threshold": 0.8,
    }
    response = client.post("/query", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["found"] is True
    assert data["response"] == "The capital of Germany is Berlin."
    assert data["matched_request"] == "What is the capital of Germany?"
    assert data["score"] >= 0.8


def test_query_cache_miss_unrelated():
    payload = {
        "request": "How do deep neural networks backpropagate gradients?",
        "threshold": 0.8,
    }
    response = client.post("/query", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["found"] is False
    assert data["response"] is None
    assert data["matched_request"] is None
    assert data["score"] is None


def test_embed_validation_error():
    # Missing response
    response = client.post("/embed", json={"request": "Hello"})
    assert response.status_code == 422

    # Empty string
    response_empty = client.post("/embed", json={"request": "", "response": ""})
    assert response_empty.status_code == 422


def test_query_validation_error():
    # Invalid threshold > 1.0
    response = client.post("/query", json={"request": "Hello", "threshold": 1.5})
    assert response.status_code == 422

    # Invalid threshold < 0.0
    response_neg = client.post("/query", json={"request": "Hello", "threshold": -0.1})
    assert response_neg.status_code == 422



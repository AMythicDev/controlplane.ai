from fastapi.testclient import TestClient
from toxicity_scanner.main import app

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

def test_scan_with_object_payload():
    payload = {
        "texts": [
            "You are awesome and helpful!",
            "I hate you so much, you are stupid and ugly."
        ]
    }
    response = client.post("/scan", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert "results" in data
    assert len(data["results"]) == 2
    
    # Check non-toxic text
    clean_item = data["results"][0]
    assert clean_item["text"] == "You are awesome and helpful!"
    assert clean_item["scores"]["toxicity"] < 0.1
    
    # Check toxic text
    toxic_item = data["results"][1]
    assert toxic_item["text"] == "I hate you so much, you are stupid and ugly."
    assert toxic_item["scores"]["toxicity"] > 0.8
    assert "insult" in toxic_item["scores"]
    assert "obscene" in toxic_item["scores"]
    assert "threat" in toxic_item["scores"]
    assert "severe_toxicity" in toxic_item["scores"]
    assert "identity_attack" in toxic_item["scores"]

def test_scan_with_array_payload():
    payload = ["Have a nice day!"]
    response = client.post("/scan", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert len(data["results"]) == 1
    assert data["results"][0]["text"] == "Have a nice day!"
    assert data["results"][0]["scores"]["toxicity"] < 0.1

def test_scan_with_empty_payload():
    response = client.post("/scan", json={"texts": []})
    assert response.status_code == 200
    assert response.json() == {"results": []}

    response_list = client.post("/scan", json=[])
    assert response_list.status_code == 200
    assert response_list.json() == {"results": []}

def test_aliases():
    payload = {"texts": ["Hello"]}
    resp_predict = client.post("/predict", json=payload)
    assert resp_predict.status_code == 200
    assert len(resp_predict.json()["results"]) == 1

    resp_toxicity = client.post("/toxicity", json=payload)
    assert resp_toxicity.status_code == 200
    assert len(resp_toxicity.json()["results"]) == 1

def test_invalid_item_type():
    response = client.post("/scan", json={"texts": [123, "test"]})
    assert response.status_code == 422

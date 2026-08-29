# Inference API

A fast ML inference API built with **FastAPI**, **ONNX Runtime**, and **Transformers**.

## Features

- **Toxicity Scanning**: Evaluates multiple toxicity categories per input string:
  - `toxicity`
  - `severe_toxicity`
  - `obscene`
  - `threat`
  - `insult`
  - `identity_attack`
- **NLI Verification**: Checks hypothesis entailment against premise contexts for hallucination detection.
- **Semantic Caching & Embeddings**: Generates 384-dimensional normalized vector embeddings using `bge-small-en-v1.5.onnx` and persists/queries interactions in **Qdrant**.
- Interactive Swagger UI documentation at `/docs`.

## Running the API

You can start the server using `uv`:

```bash
uv run uvicorn inference.main:app --host 0.0.0.0 --port 8000 --reload
```

Or run via the CLI entrypoint:

```bash
uv run inference --host 0.0.0.0 --port 8000
```

## Environment Variables

- `QDRANT_URL`: URL of the Qdrant instance (default: `http://localhost:6333`, in Docker: `http://qdrant:6333`)
- `QDRANT_COLLECTION`: Name of the Qdrant collection (default: `semantic_cache`)
- `QDRANT_API_KEY`: Optional API key for Qdrant Cloud / authentication

## API Endpoints

### 1. Health Check
`GET /health`

**Response:**
```json
{
  "status": "healthy"
}
```

### 2. Scan Toxicity
`POST /toxicity`

**Request Body (`application/json`):**

```json
{
  "texts": [
    "I really love this project!",
    "You are an idiot and nobody likes you."
  ]
}
```

### 3. NLI Verification
`POST /nli`

**Request Body (`application/json`):**

```json
{
  "premise": "The sky is blue today.",
  "hypothesis": "The sky is colored blue."
}
```

### 4. Embed & Save Interaction (Semantic Cache)
`POST /embed`

Generates embeddings for both `request` and `response` texts and saves them into Qdrant DB.

**Request Body (`application/json`):**

```json
{
  "request": "What is the capital of France?",
  "response": "The capital of France is Paris."
}
```

**Response (`200 OK`):**
```json
{
  "id": "e7b0d912-32a2-4a7b-a279-3d02a99d52b1",
  "status": "stored",
  "request_embedding": [0.015, -0.022, 0.008, ...],
  "response_embedding": [0.031, -0.012, 0.019, ...]
}
```

### 5. Query Semantic Cache
`POST /query`

Queries Qdrant for a cached response semantically matching the `request` text above the similarity `threshold`.

**Request Body (`application/json`):**

```json
{
  "request": "What is France's capital city?",
  "threshold": 0.8
}
```

**Response (`200 OK`):**
```json
{
  "found": true,
  "response": "The capital of France is Paris.",
  "matched_request": "What is the capital of France?",
  "score": 0.942
}
```

## Testing

Run tests using pytest:

```bash
uv run pytest
```



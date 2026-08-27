# Toxicity Scanner API

A fast toxicity detection API built with **FastAPI**, **ONNX Runtime**, and **Detoxify**.

## Features

- Evaluates multiple toxicity categories per input string:
  - `toxicity`
  - `severe_toxicity`
  - `obscene`
  - `threat`
  - `insult`
  - `identity_attack`
- Batch processing support for multiple text strings.
- Interactive Swagger UI documentation at `/docs`.

## Running the API

You can start the server using `uv`:

```bash
uv run uvicorn toxicity_scanner.main:app --host 0.0.0.0 --port 8000 --reload
```

Or run via the CLI entrypoint:

```bash
uv run toxicity-scanner --host 0.0.0.0 --port 8000
```

## API Endpoints

### 1. Health Check
`GET /health`

**Response:**
```json
{
  "status": "healthy"
}
```

### 2. Scan Texts
`POST /scan` (Aliases: `/predict`, `/toxicity`)

**Request Body (`application/json`):**

```json
{
  "texts": [
    "I really love this project!",
    "You are an idiot and nobody likes you."
  ]
}
```

*Or pass a raw JSON array of strings:*
```json
[
  "I really love this project!",
  "You are an idiot and nobody likes you."
]
```

**Response (`200 OK`):**

```json
{
  "results": [
    {
      "text": "I really love this project!",
      "scores": {
        "toxicity": 0.0006,
        "severe_toxicity": 0.0001,
        "obscene": 0.0011,
        "threat": 0.0003,
        "insult": 0.0008,
        "identity_attack": 0.0001
      }
    },
    {
      "text": "You are an idiot and nobody likes you.",
      "scores": {
        "toxicity": 0.9972,
        "severe_toxicity": 0.0011,
        "obscene": 0.0310,
        "threat": 0.0022,
        "insult": 0.9911,
        "identity_attack": 0.0014
      }
    }
  ]
}
```

## Testing

Run tests using pytest:

```bash
uv run pytest
```

# 🌐 ControlPlane.ai

<div align="center">

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
![FastAPI](https://img.shields.io/badge/FastAPI-0.100+-009688?logo=fastapi&logoColor=white)
![SvelteKit](https://img.shields.io/badge/SvelteKit-5.0+-FF3E00?logo=svelte&logoColor=white)
![Qdrant](https://img.shields.io/badge/Qdrant-Vector_DB-DC2626?logo=qdrant&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Cache_%26_RateLimit-DC382D?logo=redis&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-Audit_Logs-47A248?logo=mongodb&logoColor=white)

**An enterprise-grade, ultra-low-latency AI Control Plane & LLM Gateway with built-in semantic caching, guardrails (PII redaction, toxicity scanning, hallucination verification), cost budgeting, and full-stack observability.**

[Features](#-key-features) • [Architecture](#-architecture) • [Quick Start](#-quick-start) • [API Reference](#-api-reference) • [Demo](#-gateway-demo)

</div>

---

## 📽️ Gateway Demo

<!-- HERO GIF: Add your terminal / API demo GIF here -->
<div align="center">
  <img src="assets/demo-gateway.gif" alt="ControlPlane.ai Gateway Demo" width="850px" />
  <p><em>Real-time OpenAI-compatible proxy with sub-millisecond semantic cache hit, PII masking, and cost tracking</em></p>
</div>

---

## ✨ Key Features

- **⚡ Semantic Caching**: Powered by **Qdrant** and local ONNX embeddings (`bge-small-en-v1.5`), cutting LLM inference costs and delivering instant sub-millisecond responses for semantically similar prompts.
- **🛡️ AI Safety & Guardrails**:
  - **PII Detection & Anonymization**: Automatic masking of sensitive personal information via **Microsoft Presidio**.
  - **Toxicity & Moderation Scanner**: Multi-label toxicity scoring (`severe_toxicity`, `obscene`, `threat`, `insult`, `identity_attack`) via ONNX runtime.
  - **NLI Hallucination Verification**: Natural Language Inference entailment checks between context and model responses.
  - **Confidence Scoring**: Logprob-based confidence analytics.
- **💰 Cost Management & Budget Caps**:
  - Per-user daily & monthly budget enforcement backed by **Redis**.
  - Token-level and microcent-level cost tracking with cumulative savings calculation.
- **🔄 Multi-Provider OpenAI-Compatible Gateway**:
  - Drop-in `/v1/chat/completions` proxy supporting multiple downstream LLM providers.
- **📊 Comprehensive Observability & Playground**:
  - MongoDB-backed audit log capturing prompt, response, latency, token spend, toxicity, and NLI verification scores.
  - Interactive playground endpoint (`/v1/playground`) and web UI built with **SvelteKit** and **Tailwind CSS**.

---

## 🏛️ Architecture

```mermaid
flowchart TD
    Client["Client / SDK / Web UI"]
    
    subgraph ControlPlaneGateway ["ControlPlane.ai Gateway (Go :8080)"]
        BudgetCheck["1. Budget & Rate Limiter (Redis)"]
        PIICheck["2. PII Analyzer & Anonymizer (Presidio)"]
        CacheCheck["3. Semantic Cache Query (Qdrant / Inference)"]
        Router["4. Provider Router & Client Proxy"]
        Guardrails["5. Toxicity & NLI Evaluation"]
        Logger["6. Audit Logger (MongoDB) & Spend Tracker"]
    end
    
    subgraph ExternalServices ["Data & Inference Services"]
        RedisDB[("Redis :6379")]
        Mongo[("MongoDB :27017")]
        QdrantDB[("Qdrant Vector DB :6333")]
        InferenceSvc["Inference API (FastAPI / ONNX :5002)"]
        PresidioSvc["Presidio Analyzer & Anonymizer (:5000 / :5001)"]
    end
    
    LLM["External LLM Providers (OpenAI / Anthropic / Local)"]

    Client -->|"POST /v1/chat/completions"| BudgetCheck
    BudgetCheck <--> RedisDB
    BudgetCheck --> PIICheck
    PIICheck <--> PresidioSvc
    PIICheck --> CacheCheck
    CacheCheck <--> InferenceSvc
    InferenceSvc <--> QdrantDB
    
    CacheCheck -->|"Cache Miss"| Router
    Router --> LLM
    LLM --> Guardrails
    Guardrails <--> InferenceSvc
    Guardrails --> Logger
    Logger --> Mongo
    Logger --> RedisDB
    Logger --> Client
    
    CacheCheck -->|"Cache Hit (0ms / $0)"| Logger
```

---

## 🖥️ Web Dashboard

<!-- DASHBOARD GIF: Add your web dashboard walkthrough GIF here -->
<div align="center">
  <img src="assets/demo-dashboard.gif" alt="ControlPlane.ai Dashboard Demo" width="850px" />
  <p><em>SvelteKit Web Dashboard for real-time request tracing, safety telemetry, and budget configuration</em></p>
</div>

---

## 🚀 Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://go.dev/dl/) 1.22+ (for local backend development)
- [uv](https://docs.astral.sh/uv/) (for local Python inference development)
- [Bun](https://bun.sh/) or [Node.js](https://nodejs.org/) (for Web UI)

### 1. Clone & Start Infrastructure

```bash
# Clone the repository
git clone https://github.com/AMythicDev/controlplane.ai.git
cd controlplane.ai

# Start all supporting microservices (Redis, Mongo, Qdrant, Presidio, Inference)
docker compose up -d
```

### 2. Configure Environment

Create a `.env` file in `backend/`:

```env
PORT=8080
REDIS_URL=localhost:6379
REDIS_DB=0
MONGO_URI=mongodb://localhost:27017
INFERENCE_API_URL=http://localhost:5002
PRESIDIO_ANALYZER_URL=http://localhost:5000
PRESIDIO_ANONYMIZER_URL=http://localhost:5001
NVIDIA_API_KEY=nvapi-your-key-here
```

### 3. Run the Gateway

```bash
cd backend
go run main.go
```

### 4. Run the Web Dashboard

```bash
cd ../web
bun install
bun run dev
```

---

## 📡 API Reference

### 1. Chat Completions (OpenAI Compatible)

`POST /v1/chat/completions`

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "openai/gpt-4o-mini",
    "use_semantic_cache": true,
    "messages": [
      {
        "role": "user",
        "content": "Explain quantum computing in simple terms."
      }
    ]
  }'
```

**Response:**
```json
{
  "id": "chatcmpl-openai",
  "object": "chat.completion",
  "model": "openai/gpt-4o-mini",
  "confidence": 0.98,
  "toxicity": 0.001,
  "nli": {
    "entailment": 0.94,
    "neutral": 0.05,
    "contradiction": 0.01
  },
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Quantum computing is a type of computing that uses quantum mechanics..."
      },
      "finish_reason": "stop"
    }
  ]
}
```

### 2. Cost & Cache Savings Analytics

`GET /v1/cost`

```json
{
  "user_id": "user_mvp_123",
  "cost_dollars": 0.0245,
  "savings_dollars": 0.1850,
  "total_savings": 0.1850,
  "average_cost_dollars": 0.0012
}
```

### 3. Budget Configuration

`GET /v1/config` | `POST /v1/config`

```json
{
  "per_user_daily_limit": 5000000,
  "per_user_monthly_limit": 150000000
}
```

### 4. Request Audit Logs

`GET /v1/requests?limit=20&offset=0`

Returns paginated request histories with full telemetry (prompt, masked PII, latency, confidence, toxicity scores, cache hit flags).

---

## 🛠️ Repository Structure

```
.
├── backend/               # Core AI Gateway in Go (Gin, Redis, MongoDB, Presidio)
│   ├── pii-detection/     # Custom Presidio recognizers and redaction logic
│   ├── cost-tracking.go   # Budget control & Redis spend accounting
│   ├── semantic-cache.go  # Qdrant cache query & persist handlers
│   ├── toxicity-scanner.go# Toxicity validation client
│   └── request-log.go     # MongoDB request auditing
├── inference/             # ML inference engine (FastAPI, ONNX Runtime, Qdrant)
│   ├── src/               # Toxicity models, NLI verifier, embeddings
│   └── Dockerfile         # Containerized inference service
├── web/                   # Frontend observability dashboard (SvelteKit 5, Tailwind CSS)
├── docker-compose.yaml    # Microservices orchestration (Redis, Mongo, Qdrant, Presidio)
└── assets/                # README demos, GIFs, and architecture diagrams
```

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/AMythicDev/controlplane.ai/issues).

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

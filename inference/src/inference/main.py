import argparse
import os

import uvicorn
from fastapi import Body, FastAPI, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field

from inference.embedder import get_embedding
from inference.nli_inference import run_nli
from inference.qdrant_storage import storage
from inference.toxicity import toxicity_scan

app = FastAPI(
    title="Inference API",
    description="FastAPI service to run Toxicity scanning, NLI verification, and Qdrant-backed Semantic Caching / Embeddings using ONNX Runtime.",
    version="0.1.0",
)

# Enable CORS for cross-origin requests
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class ToxicityRequest(BaseModel):
    texts: list[str] = Field(
        ...,
        description="List of text strings to analyze for toxicity",
        examples=[["Hello, have a wonderful day!", "You are an idiot and I hate you."]],
    )


class ToxicityResult(BaseModel):
    text: str = Field()
    toxicity: float = Field()
    severe_toxicity: float = Field()
    obscene: float = Field()
    threat: float = Field()
    insult: float = Field()
    identity_attack: float = Field()


class NLIRequest(BaseModel):
    premise: str = Field(..., description="The source context")
    hypothesis: str = Field(..., description="The LLM's response to check")


class NLIResponse(BaseModel):
    label: str
    score: float
    contradiction_prob: float
    neutral_prob: float
    entailment_prob: float


class EmbedRequest(BaseModel):
    request: str = Field(
        ...,
        min_length=1,
        description="The user request / prompt text",
        examples=["What is the capital of France?"],
    )
    response: str = Field(
        ...,
        min_length=1,
        description="The model response text",
        examples=["The capital of France is Paris."],
    )


class EmbedResponse(BaseModel):
    id: str = Field(..., description="Unique ID of the stored interaction in Qdrant")
    status: str = Field(default="stored", description="Status of the storage operation")
    request_embedding: list[float] = Field(
        ..., description="Normalized vector embedding for the request text"
    )
    response_embedding: list[float] = Field(
        ..., description="Normalized vector embedding for the response text"
    )


class QueryRequest(BaseModel):
    request: str = Field(
        ...,
        min_length=1,
        description="The query request text to find a cached response for",
        examples=["What is France's capital city?"],
    )
    threshold: float = Field(
        default=0.8,
        ge=0.0,
        le=1.0,
        description="Similarity score threshold between 0.0 and 1.0",
        examples=[0.85],
    )


class QueryResponse(BaseModel):
    found: bool = Field(
        ..., description="Whether a matching response meeting the threshold was found"
    )
    response: str | None = Field(
        default=None, description="The cached response text if found"
    )
    matched_request: str | None = Field(
        default=None, description="The matched request text from cache"
    )
    score: float | None = Field(
        default=None, description="Similarity score of the match"
    )


@app.get("/health", tags=["General"])
def health() -> dict[str, str]:
    """Health check endpoint."""
    return {"status": "healthy"}


@app.get("/", tags=["General"])
def root() -> dict[str, str]:
    """Root endpoint."""
    return {
        "message": "Inference API is running",
        "health": "/health",
        "docs": "/docs",
    }


@app.post(
    "/toxicity",
    response_model=list[ToxicityResult],
    tags=["Toxicity"],
    summary="Scan a list of text strings for toxicity",
)
def toxicity(
    payload: ToxicityRequest | list[str] = Body(),
) -> list[ToxicityResult]:
    """Accepts a list of text strings (either as a JSON array or `{ "texts": [...] }`)

    and returns the toxicity scores for each input string.
    """
    try:
        texts = payload.texts if isinstance(payload, ToxicityRequest) else payload

        if not texts:
            return []

        # Validate that all elements are strings
        if not all(isinstance(t, str) for t in texts):
            raise HTTPException(
                status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
                detail="All items in the input list must be strings.",
            )

        raw_scores = toxicity_scan(texts)

        results = [
            ToxicityResult(
                text=text,
                **score_dict
            )
            for text, score_dict in zip(texts, raw_scores)
        ]

        return results
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Inference failed: {str(e)}",
        ) from e


@app.post(
    "/nli",
    response_model=NLIResponse,
    tags=["NLI"],
    summary="Verify if a hypothesis is entailed by a premise",
)
def nli_verification(req: NLIRequest) -> NLIResponse:
    try:
        result = run_nli(req.premise, req.hypothesis)
        return NLIResponse(**result)
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"NLI Inference failed: {str(e)}",
        ) from e


@app.post(
    "/embed",
    response_model=EmbedResponse,
    tags=["Semantic Cache"],
    summary="Generate embeddings for request and response texts and save them in Qdrant",
)
def embed(
    payload: EmbedRequest,
) -> EmbedResponse:
    """Takes request and response text strings, generates vector embeddings for both,

    and persists them into Qdrant.
    """
    try:
        req_embedding = get_embedding(payload.request)
        resp_embedding = get_embedding(payload.response)

        point_id = storage.save_interaction(
            request_text=payload.request,
            response_text=payload.response,
            request_vector=req_embedding,
            response_vector=resp_embedding,
        )

        return EmbedResponse(
            id=point_id,
            status="stored",
            request_embedding=req_embedding,
            response_embedding=resp_embedding,
        )
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Embedding generation and storage failed: {str(e)}",
        ) from e


@app.post(
    "/query",
    response_model=QueryResponse,
    tags=["Semantic Cache"],
    summary="Query Qdrant for a cached response matching the request text above threshold",
)
def query(
    payload: QueryRequest,
) -> QueryResponse:
    """Takes a request text string and similarity threshold, queries Qdrant for matching

    cached requests, and returns the response if above threshold.
    """
    try:
        query_vector = get_embedding(payload.request)

        match = storage.search_cached_response(
            query_vector=query_vector,
            threshold=payload.threshold,
            limit=1,
        )

        if match is not None:
            return QueryResponse(
                found=True,
                response=match["response"],
                matched_request=match["request"],
                score=match["score"],
            )

        return QueryResponse(
            found=False,
            response=None,
            matched_request=None,
            score=None,
        )
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Query failed: {str(e)}",
        ) from e


def main() -> None:
    """CLI entry point to launch the FastAPI server with Uvicorn."""
    parser = argparse.ArgumentParser(description="Start the Inference API server.")
    parser.add_argument(
        "--host",
        type=str,
        default=os.getenv("HOST", "0.0.0.0"),
        help="Host interface to bind (default: 0.0.0.0)",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=int(os.getenv("PORT", "8000")),
        help="Port to bind (default: 8000)",
    )
    parser.add_argument(
        "--reload",
        action="store_true",
        help="Enable auto-reload for development",
    )
    args = parser.parse_args()

    uvicorn.run("inference.main:app", host=args.host, port=args.port, reload=args.reload)


if __name__ == "__main__":
    main()


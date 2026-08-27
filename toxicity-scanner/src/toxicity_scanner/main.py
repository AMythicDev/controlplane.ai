import argparse
import os
from typing import Union
from fastapi import Body, FastAPI, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field
import uvicorn

from toxicity_scanner.inference import run_inference

app = FastAPI(
    title="Toxicity Scanner API",
    description="FastAPI service to score text strings across multiple toxicity categories using ONNX Runtime.",
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


class ToxicityScores(BaseModel):
    toxicity: float = Field(..., description="Overall toxicity probability (0.0 - 1.0)")
    severe_toxicity: float = Field(..., description="Severe toxicity probability (0.0 - 1.0)")
    obscene: float = Field(..., description="Obscenity probability (0.0 - 1.0)")
    threat: float = Field(..., description="Threat probability (0.0 - 1.0)")
    insult: float = Field(..., description="Insult probability (0.0 - 1.0)")
    identity_attack: float = Field(..., description="Identity attack probability (0.0 - 1.0)")


class ToxicityResult(BaseModel):
    text: str = Field(..., description="Original input text string")
    scores: ToxicityScores = Field(..., description="Predicted toxicity probabilities across categories")


class ToxicityResponse(BaseModel):
    results: list[ToxicityResult] = Field(..., description="Analysis results for each input string")


@app.get("/", tags=["General"])
def root() -> dict[str, str]:
    """Root endpoint returning API status and docs link."""
    return {
        "message": "Toxicity Scanner API is running.",
        "docs": "/docs",
        "health": "/health",
    }


@app.get("/health", tags=["General"])
def health() -> dict[str, str]:
    """Health check endpoint."""
    return {"status": "healthy"}


@app.post(
    "/scan",
    response_model=ToxicityResponse,
    tags=["Toxicity"],
    summary="Scan a list of text strings for toxicity",
)
def scan_texts(
    payload: Union[ToxicityRequest, list[str]] = Body(
        ...,
        description="A list of strings or an object containing a 'texts' list",
        openapi_examples={
            "object_format": {
                "summary": "Object format",
                "value": {"texts": ["I love coding!", "Get out of here you fool."]},
            },
            "array_format": {
                "summary": "Array format",
                "value": ["I love coding!", "Get out of here you fool."],
            },
        },
    ),
) -> ToxicityResponse:
    """Accepts a list of text strings (either as a JSON array or `{ "texts": [...] }`)

    and returns the toxicity scores for each input string.
    """
    try:
        texts = payload.texts if isinstance(payload, ToxicityRequest) else payload

        if not texts:
            return ToxicityResponse(results=[])

        # Validate that all elements are strings
        if not all(isinstance(t, str) for t in texts):
            raise HTTPException(
                status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
                detail="All items in the input list must be strings.",
            )

        raw_scores = run_inference(texts)

        results = [
            ToxicityResult(
                text=text,
                scores=ToxicityScores(**score_dict),
            )
            for text, score_dict in zip(texts, raw_scores)
        ]

        return ToxicityResponse(results=results)
    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Inference failed: {str(e)}",
        ) from e


@app.post(
    "/predict",
    response_model=ToxicityResponse,
    tags=["Toxicity"],
    summary="Alias for /scan",
    include_in_schema=False,
)
def predict_texts(
    payload: Union[ToxicityRequest, list[str]] = Body(...),
) -> ToxicityResponse:
    return scan_texts(payload)


@app.post(
    "/toxicity",
    response_model=ToxicityResponse,
    tags=["Toxicity"],
    summary="Alias for /scan",
    include_in_schema=False,
)
def toxicity_texts(
    payload: Union[ToxicityRequest, list[str]] = Body(...),
) -> ToxicityResponse:
    return scan_texts(payload)


def main() -> None:
    """CLI entry point to launch the FastAPI server with Uvicorn."""
    parser = argparse.ArgumentParser(description="Start the Toxicity Scanner API server.")
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

    uvicorn.run("toxicity_scanner.main:app", host=args.host, port=args.port, reload=args.reload)


if __name__ == "__main__":
    main()

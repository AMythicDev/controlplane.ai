import argparse
import os
from fastapi import Body, FastAPI, HTTPException, status
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field
import uvicorn

from toxicity_scanner.inference import run_inference
from toxicity_scanner.nli_inference import run_nli

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


@app.get("/health", tags=["General"])
def health() -> dict[str, str]:
    """Health check endpoint."""
    return {"status": "healthy"}


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

        raw_scores = run_inference(texts)

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

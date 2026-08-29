from inference.embedder import get_embedding, get_embeddings
from inference.nli_inference import run_nli
from inference.qdrant_storage import QdrantStorage, storage
from inference.toxicity import toxicity_scan

# Backwards compatibility alias
run_inference = toxicity_scan

__all__ = [
    "toxicity_scan",
    "run_inference",
    "run_nli",
    "get_embedding",
    "get_embeddings",
    "QdrantStorage",
    "storage",
]



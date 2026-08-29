import os
import uuid
from datetime import datetime, timezone
from typing import Any
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, PointStruct, VectorParams

DEFAULT_QDRANT_URL = os.getenv("QDRANT_URL", "http://localhost:6333")
DEFAULT_COLLECTION_NAME = os.getenv("QDRANT_COLLECTION", "semantic_cache")
VECTOR_DIMENSION = 384


class QdrantStorage:
    def __init__(
        self,
        url: str | None = None,
        collection_name: str | None = None,
        api_key: str | None = None,
    ) -> None:
        self.url = url or os.getenv("QDRANT_URL", DEFAULT_QDRANT_URL)
        self.collection_name = collection_name or os.getenv(
            "QDRANT_COLLECTION", DEFAULT_COLLECTION_NAME
        )
        self.api_key = api_key or os.getenv("QDRANT_API_KEY")
        self._client: QdrantClient | None = None

    @property
    def client(self) -> QdrantClient:
        if self._client is None:
            self._client = QdrantClient(url=self.url, api_key=self.api_key)
            self.ensure_collection()
        return self._client

    def ensure_collection(self) -> None:
        """Ensures the collection exists with named vectors for request and response."""
        if not self._client.collection_exists(self.collection_name):
            self._client.create_collection(
                collection_name=self.collection_name,
                vectors_config={
                    "request": VectorParams(
                        size=VECTOR_DIMENSION, distance=Distance.COSINE
                    ),
                    "response": VectorParams(
                        size=VECTOR_DIMENSION, distance=Distance.COSINE
                    ),
                },
            )

    def save_interaction(
        self,
        request_text: str,
        response_text: str,
        request_vector: list[float],
        response_vector: list[float],
    ) -> str:
        """Saves a request-response pair and their vector embeddings to Qdrant."""
        point_id = str(uuid.uuid4())
        point = PointStruct(
            id=point_id,
            vector={
                "request": request_vector,
                "response": response_vector,
            },
            payload={
                "request": request_text,
                "response": response_text,
                "created_at": datetime.now(timezone.utc).isoformat(),
            },
        )
        self.client.upsert(
            collection_name=self.collection_name,
            points=[point],
        )
        return point_id

    def search_cached_response(
        self,
        query_vector: list[float],
        threshold: float = 0.8,
        limit: int = 1,
    ) -> dict[str, Any] | None:
        """Searches Qdrant for cached responses matching the query vector above the threshold."""
        results = self.client.query_points(
            collection_name=self.collection_name,
            query=query_vector,
            using="request",
            score_threshold=threshold,
            limit=limit,
        )

        if not results.points:
            return None

        top_point = results.points[0]
        payload = top_point.payload or {}

        return {
            "id": top_point.id,
            "score": float(top_point.score),
            "request": payload.get("request", ""),
            "response": payload.get("response", ""),
        }


# Global storage instance
storage = QdrantStorage()

import os
from pathlib import Path
import numpy as np
import onnxruntime as ort
from transformers import AutoTokenizer

def _resolve_model_path() -> str:
    MODEL_NAME = "bge-small-en-v1.5.onnx"
    repo_model = Path(__file__).resolve().parents[2] / "model" / MODEL_NAME
    if repo_model.exists():
        return str(repo_model)
    return f"model/{MODEL_NAME}"

onnx_model_path = _resolve_model_path()
tokenizer = AutoTokenizer.from_pretrained("BAAI/bge-small-en-v1.5")
session = ort.InferenceSession(onnx_model_path, providers=["CPUExecutionProvider"])

def get_embedding(text: str) -> list[float]:
    """Takes a single string as input and returns its vector embedding as output."""
    if not isinstance(text, str):
        raise ValueError("Input must be a string")
    return get_embeddings([text])[0]


def get_embeddings(texts: list[str]) -> list[list[float]]:
    """Takes a list of strings as input and returns a list of vector embeddings."""
    if not texts:
        return []

    inputs = tokenizer(
        texts,
        padding=True,
        truncation=True,
        max_length=512,
        return_tensors="np",
    )

    ort_inputs = {
        "input_ids": inputs["input_ids"].astype(np.int64),
        "attention_mask": inputs["attention_mask"].astype(np.int64),
        "token_type_ids": inputs.get(
            "token_type_ids", np.zeros_like(inputs["input_ids"], dtype=np.int64)
        ).astype(np.int64),
    }

    outputs = session.run(None, ort_inputs)
    last_hidden_state = outputs[0]  # shape: (batch_size, seq_len, 384)
    cls_embeddings = last_hidden_state[:, 0]

    # L2 normalize per embedding
    norms = np.linalg.norm(cls_embeddings, axis=-1, keepdims=True)
    embeddings = cls_embeddings / np.clip(norms, a_min=1e-12, a_max=None)

    return embeddings.tolist()


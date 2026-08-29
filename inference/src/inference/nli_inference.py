import os
import numpy as np
import onnxruntime as ort
from transformers import AutoTokenizer
from pathlib import Path

def _resolve_model_path() -> str:
    MODEL_NAME = "nli-deberta-v3-small.onnx"
    repo_model = Path(__file__).resolve().parents[2] / "model" / MODEL_NAME
    if repo_model.exists():
        return str(repo_model)
    return f"model/{MODEL_NAME}"

model = _resolve_model_path()

tokenizer = AutoTokenizer.from_pretrained("Xenova/nli-deberta-v3-small")
session = ort.InferenceSession(model, providers=["CPUExecutionProvider"])

def softmax(x: np.ndarray) -> np.ndarray:
    e = np.exp(x - np.max(x))
    return e / e.sum(axis=-1, keepdims=True)

def run_nli(premise: str, hypothesis: str) -> dict[str, str | float]:
    inputs = tokenizer(
        premise, 
        hypothesis, 
        return_tensors="np", 
        truncation=True, 
        padding=True,
        max_length=512
    )
    
    ort_inputs = {
        "input_ids": inputs["input_ids"].astype(np.int64),
        "attention_mask": inputs["attention_mask"].astype(np.int64),
    }
    
    logits = session.run(None, ort_inputs)[0][0]
    probs = softmax(logits)
    
    labels = ["contradiction", "entailment", "neutral"]
    top_idx = int(np.argmax(probs))
    
    return {
        "label": labels[top_idx],
        "score": float(probs[top_idx]),
        "contradiction_prob": float(probs[0]),
        "entailment_prob": float(probs[1]),
        "neutral_prob": float(probs[2])
    }


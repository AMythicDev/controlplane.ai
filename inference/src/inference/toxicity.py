import os
from pathlib import Path
import numpy as np
import onnxruntime as ort
from transformers import AutoTokenizer

tokenizer = AutoTokenizer.from_pretrained("gravitee-io/detoxify-onnx")

def _resolve_model_path() -> str:
    MODEL_NAME = "detoxify-original.onnx"
    repo_model = Path(__file__).resolve().parents[2] / "model" / MODEL_NAME
    if repo_model.exists():
        return str(repo_model)
    return f"model/{MODEL_NAME}"

onnx_model_path = _resolve_model_path()
session = ort.InferenceSession(onnx_model_path, providers=["CPUExecutionProvider"])

def toxicity_scan(texts: list[str]) -> list[dict[str, float]]:
    if not texts:
        return []

    inputs = tokenizer(texts, padding=True, truncation=True, return_tensors="np")

    ort_inputs = {
        "input_ids": inputs["input_ids"].astype(np.int64),
        "attention_mask": inputs["attention_mask"].astype(np.int64)
    }

    outputs = session.run(None, ort_inputs)

    probabilities = 1 / (1 + np.exp(-outputs[0])) 
    labels = ["toxicity", "severe_toxicity", "obscene", "threat", "insult", "identity_attack"]

    results = []
    for i in range(len(texts)):
        results.append({k : float(v) for k, v in zip(labels, probabilities[i])})

    return results


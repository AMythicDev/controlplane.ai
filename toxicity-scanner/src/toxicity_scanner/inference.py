import os
from pathlib import Path
import numpy as np
import onnxruntime as ort
from transformers import AutoTokenizer

# 1. Load the corresponding tokenizer from Hugging Face
tokenizer = AutoTokenizer.from_pretrained("gravitee-io/detoxify-onnx")

# 2. Initialize your downloaded or converted .onnx file 
# Ensure your 'providers' array matches your hardware choice
def _resolve_model_path() -> str:
    env_path = os.getenv("MODEL_PATH")
    if env_path and Path(env_path).exists():
        return env_path
    if Path("model/detoxify-original.onnx").exists():
        return "model/detoxify-original.onnx"
    repo_model = Path(__file__).resolve().parents[2] / "model" / "detoxify-original.onnx"
    if repo_model.exists():
        return str(repo_model)
    return "model/detoxify-original.onnx"

onnx_model_path = _resolve_model_path()
session = ort.InferenceSession(onnx_model_path, providers=["CPUExecutionProvider"])

def run_inference(texts: list[str]) -> list[dict[str, float]]:
    if not texts:
        return []

    # 3. Prepare inputs
    inputs = tokenizer(texts, padding=True, truncation=True, return_tensors="np")

    ort_inputs = {
        "input_ids": inputs["input_ids"].astype(np.int64),
        "attention_mask": inputs["attention_mask"].astype(np.int64)
    }

    # 5. Execute inference
    outputs = session.run(None, ort_inputs)

    # 6. parse sigmoid outputs for multi-label probabilities
    # detoxify evaluates: toxicity, severe_toxicity, obscene, threat, insult, identity_attack
    probabilities = 1 / (1 + np.exp(-outputs[0])) 
    labels = ["toxicity", "severe_toxicity", "obscene", "threat", "insult", "identity_attack"]

    results = []
    for i in range(len(texts)):
        results.append({k : float(v) for k, v in zip(labels, probabilities[i])})

    return results


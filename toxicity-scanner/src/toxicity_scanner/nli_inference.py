import os
import numpy as np
import onnxruntime as ort
from transformers import AutoTokenizer

# Resolve absolute path to the local ONNX model (toxicity-scanner/model/model.onnx)
BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(__file__)))
MODEL_PATH = os.path.join(BASE_DIR, "model", "model.onnx")

print(f"Loading ONNX model from: {MODEL_PATH}")

tokenizer = AutoTokenizer.from_pretrained("Xenova/nli-deberta-v3-small")
session = ort.InferenceSession(MODEL_PATH, providers=["CPUExecutionProvider"])

def softmax(x: np.ndarray) -> np.ndarray:
    e = np.exp(x - np.max(x))
    return e / e.sum(axis=-1, keepdims=True)

def run_nli(premise: str, hypothesis: str) -> dict:
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

if __name__ == "__main__":
    import json
    print("\n" + "="*50)
    print("🧠 Interactive RAG Verifier Test")
    print("Type 'quit' in either prompt to exit.")
    print("="*50 + "\n")

    while True:
        context = input("\n📝 Enter Source Context (Premise)   : ")
        if context.lower() in ['quit', 'q', 'exit']:
            break
            
        answer = input("🤖 Enter LLM Answer (Hypothesis)  : ")
        if answer.lower() in ['quit', 'q', 'exit']:
            break
            
        print("\n⏳ Calculating probabilities...")
        res = run_nli(context, answer)
        
        print("\n📊 RESULT:")
        print(json.dumps(res, indent=2))
        print("-" * 50)



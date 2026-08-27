from toxicity_scanner.nli_inference import run_nli
import json

def main():
    print("Starting NLI Model Test with Dummy Inputs...\n")

    # Test Case 1: Entailment (Perfect Match)
    premise_1 = "ControlPlane raised $10M in seed funding in 2024. Founded by Jane Doe."
    hypothesis_1 = "Jane Doe founded ControlPlane, which secured $10M in 2024."
    
    print("--- Test Case 1 (Should be Entailment) ---")
    print(f"Premise (Context)    : {premise_1}")
    print(f"Hypothesis (LLM)     : {hypothesis_1}")
    result_1 = run_nli(premise_1, hypothesis_1)
    print(f"Result               : {json.dumps(result_1, indent=2)}\n")

    # Test Case 2: Contradiction (Hallucination)
    premise_2 = "ControlPlane raised $10M in seed funding in 2024. Founded by Jane Doe."
    hypothesis_2 = "ControlPlane raised $500M in Series B funding."
    
    print("--- Test Case 2 (Should be Contradiction) ---")
    print(f"Premise (Context)    : {premise_2}")
    print(f"Hypothesis (LLM)     : {hypothesis_2}")
    result_2 = run_nli(premise_2, hypothesis_2)
    print(f"Result               : {json.dumps(result_2, indent=2)}\n")

if __name__ == "__main__":
    main()

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ToxicityScannerReport struct {
	Text           string  `json:"text,omitempty"`
	Toxicity       float32 `json:"toxicity"`
	SevereToxicity float32 `json:"severe_toxicity"`
	Obscene        float32 `json:"obscene"`
	Threat         float32 `json:"threat"`
	Insult         float32 `json:"insult"`
	IdentityAttack float32 `json:"identity_attack"`
}

// calculateSingleToxicityScore converts multi-label toxicity scores into a single composite score [0.0, 1.0]
// using the Probabilistic Multi-Risk Union (Noisy-OR with severity weighting).
func calculateSingleToxicityScore(r ToxicityScannerReport) float32 {
	survival := (1.0 - 0.80*r.Toxicity) *
		(1.0 - 1.00*r.SevereToxicity) *
		(1.0 - 1.00*r.Threat) *
		(1.0 - 0.90*r.IdentityAttack) *
		(1.0 - 0.60*r.Insult) *
		(1.0 - 0.40*r.Obscene)

	score := 1.0 - survival
	if score < 0.0 {
		return 0.0
	}
	if score > 1.0 {
		return 1.0
	}
	return score
}

func runToxicityScanner(texts []string) (float32, error) {
	if len(texts) == 0 {
		return 0.0, nil
	}

	textsData, err := json.Marshal(texts)
	if err != nil {
		return 0.0, fmt.Errorf("failed to marshal texts: %w", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:5002/toxicity", bytes.NewBuffer(textsData))
	if err != nil {
		return 0.0, fmt.Errorf("failed to create toxicity request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to run toxicity scanner: %v", err)
		return 0.0, err
	}
	defer resp.Body.Close()

	var scanData []ToxicityScannerReport
	if err := json.NewDecoder(resp.Body).Decode(&scanData); err != nil {
		return 0.0, fmt.Errorf("failed to decode toxicity scanner response: %w", err)
	}

	var maxScore float32 = 0.0
	for _, report := range scanData {
		score := calculateSingleToxicityScore(report)
		if score > maxScore {
			maxScore = score
		}
	}
	return maxScore, nil
}



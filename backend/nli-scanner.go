package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type NLIReport struct {
	Label             string  `json:"label"`
	Score             float32 `json:"score"`
	ContradictionProb float32 `json:"contradiction_prob"`
	NeutralProb       float32 `json:"neutral_prob"`
	EntailmentProb    float32 `json:"entailment_prob"`
}

func runNLIScanner(premise string, hypothesis string) (*NLIReport, error) {
	if premise == "" || hypothesis == "" {
		return nil, nil
	}

	endpoint := os.Getenv("NLI_SCANNER_URL")
	if endpoint == "" {
		endpoint = "http://localhost:5002/nli"
	}

	payload := map[string]string{
		"premise":    premise,
		"hypothesis": hypothesis,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NLI request: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create NLI request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to run NLI scanner: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nli scanner returned status %d", resp.StatusCode)
	}

	var report NLIReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, fmt.Errorf("failed to decode NLI response: %w", err)
	}

	return &report, nil
}

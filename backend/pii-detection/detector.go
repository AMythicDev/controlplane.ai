package piidetection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var HttpClient = &http.Client{Timeout: 10 * time.Second}

type AnalyzerRequest struct {
	Text                  string  `json:"text"`
	Language              string  `json:"language"`
	ScoreThreshold        float64 `json:"score_threshold,omitempty"`
	ReturnDecisionProcess bool    `json:"return_decision_process,omitempty"`
}

type AnalysisExplanation struct {
	ScoreContextImprovement float64 `json:"score_context_improvement"`
}

type AnalyzerResult struct {
	EntityType          string              `json:"entity_type"`
	Score               float64             `json:"score"`
	Start               int                 `json:"start"`
	End                 int                 `json:"end"`
	AnalysisExplanation AnalysisExplanation `json:"analysis_explanation"`
}

func RecognizePII(s string) ([]AnalyzerResult, error) {
	payload := AnalyzerRequest{
		Text:                  s,
		Language:              "en",
		ScoreThreshold:        0.8,
		ReturnDecisionProcess: true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return []AnalyzerResult{}, err
	}

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5000/analyze", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return []AnalyzerResult{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := HttpClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return []AnalyzerResult{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var results []AnalyzerResult
	if err := json.Unmarshal(body, &results); err != nil {
		return []AnalyzerResult{}, err
	}
	var filteredResults []AnalyzerResult
	for _, res := range results {
		if res.AnalysisExplanation.ScoreContextImprovement == 0 {
			continue
		}
		filteredResults = append(filteredResults, res)
	}

	return filteredResults, nil
}

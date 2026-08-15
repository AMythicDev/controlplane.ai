package piidetection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RedactorRequest struct {
	Text            string                    `json:"text"`
	Anonymizers     map[string]map[string]any `json:"anonymizers"`
	AnalyzerResults []AnalyzerResult          `json:"analyzer_results"`
}

type RedactorResult struct {
	Text string
}

var PreDefinedAnnonymizers = map[string]map[string]any{
	"PHONE_NUMBER": map[string]any{
		"type":          "mask",
		"masking_char":  "*",
		"chars_to_mask": 4,
		"from_end":      true,
	},
}

func RedactPII(detections []AnalyzerResult, text string) (string, error) {
	payload := RedactorRequest{
		Text:            text,
		Anonymizers:     PreDefinedAnnonymizers,
		AnalyzerResults: detections,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5001/anonymize", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := HttpClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var rstr RedactorRequest
	if err := json.Unmarshal(body, &rstr); err != nil {
		return "", err
	}

	return rstr.Text, nil
}

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

type QueryCacheResponse struct {
	Found          bool     `json:"found"`
	Response       *string  `json:"response"`
	MatchedRequest *string  `json:"matched_request"`
	Score          *float32 `json:"score"`
}

type EmbedCacheResponse struct {
	ID                string    `json:"id"`
	Status            string    `json:"status"`
	RequestEmbedding  []float32 `json:"request_embedding"`
	ResponseEmbedding []float32 `json:"response_embedding"`
}

func getSemanticCacheBaseURL() string {
	endpoint := os.Getenv("SEMANTIC_CACHE_URL")
	if endpoint == "" {
		endpoint = os.Getenv("INFERENCE_URL")
	}
	if endpoint == "" {
		endpoint = "http://localhost:5002"
	}
	return endpoint
}

func querySemanticCache(requestText string, threshold float32) (*QueryCacheResponse, error) {
	if requestText == "" {
		return nil, nil
	}

	url := getSemanticCacheBaseURL() + "/query"

	payload := map[string]interface{}{
		"request":   requestText,
		"threshold": threshold,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cache query request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create cache query request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to query semantic cache: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("semantic cache returned status %d", resp.StatusCode)
	}

	var result QueryCacheResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode cache query response: %w", err)
	}

	return &result, nil
}

func saveSemanticCache(requestText string, responseText string) error {
	if requestText == "" || responseText == "" {
		return nil
	}

	url := getSemanticCacheBaseURL() + "/embed"

	payload := map[string]string{
		"request":  requestText,
		"response": responseText,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal cache embed request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create cache embed request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to save to semantic cache: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("semantic cache embed returned status %d", resp.StatusCode)
	}

	return nil
}

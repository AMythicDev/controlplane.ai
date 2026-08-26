package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChatCompletionsProxy(t *testing.T) {
	// Initialize the Redis client which will connect to the local valkey instance
	InitRedis()

	LoadedConfig.Enviroment = "testing"

	router := setupRouter()

	expectation, err := os.Open("expectation.json")
	if err != nil {
		panic("expectation.json not found. Required for configuring mockserver.")
	}

	req, _ := http.NewRequest("PUT", "http://localhost:1080/mockserver/expectation", expectation)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to configure mockserver")
		return
	}
	resp.Body.Close()

	t.Run("Valid Request", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"model": "openai/gpt-3.5-turbo",
			"messages": []map[string]string{
				{"role": "user", "content": "Hello!"},
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "chat.completion", response["object"])
		assert.Equal(t, "openai/gpt-3.5-turbo", response["model"])

		assert.NotEmpty(t, response["choices"])
		choices := response["choices"].([]interface{})
		firstChoice := choices[0].(map[string]interface{})
		message := firstChoice["message"].(map[string]interface{})
		assert.Equal(t, "Hello, world!", message["content"])
	})

	t.Run("Invalid Request - Missing Model", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"messages": []map[string]string{
				{"role": "user", "content": "Hello!"},
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.NotNil(t, response["error"])
	})

	t.Run("Model spec parsing", func(t *testing.T) {
		provider, model, _ := extractProviderModel("openai/gpt-5.6-turbo")
		assert.Equal(t, provider, "openai")
		assert.Equal(t, model, "gpt-5.6-turbo")

		provider, model, _ = extractProviderModel("anthropic/claude-fable-5")
		assert.Equal(t, provider, "anthropic")
		assert.Equal(t, model, "claude-fable-5")

		provider, model, _ = extractProviderModel("google/gemini-3.1-pro-preview")
		assert.Equal(t, provider, "google")
		assert.Equal(t, model, "gemini-3.1-pro-preview")
	})
}

func TestConfidenceScore(t *testing.T) {
	t.Run("Empty logprobs slice", func(t *testing.T) {
		conf := confidenceScore([]Logprobs{})
		assert.Equal(t, float32(1.0), conf)
	})

	t.Run("Perfect certainty (logprob = 0.0, P = 1.0)", func(t *testing.T) {
		logprobs := []Logprobs{
			{TokenLogprob: TokenLogprob{Token: "The", Logprob: 0.0}},
			{TokenLogprob: TokenLogprob{Token: "capital", Logprob: 0.0}},
			{TokenLogprob: TokenLogprob{Token: "is", Logprob: 0.0}},
			{TokenLogprob: TokenLogprob{Token: "Paris", Logprob: 0.0}},
		}
		conf := confidenceScore(logprobs)
		assert.InDelta(t, float32(1.0), conf, 1e-4)
	})

	t.Run("Perplexity ~ 1.23", func(t *testing.T) {
		targetLogprob := -math.Log(1.23)
		logprobs := []Logprobs{
			{TokenLogprob: TokenLogprob{Token: "test", Logprob: targetLogprob}},
		}
		conf := confidenceScore(logprobs)
		assert.InDelta(t, float32(1.0/1.23), conf, 1e-2)
	})

	t.Run("Perplexity ~ 1.34", func(t *testing.T) {
		targetLogprob := -math.Log(1.34)
		logprobs := []Logprobs{
			{TokenLogprob: TokenLogprob{Token: "test", Logprob: targetLogprob}},
		}
		conf := confidenceScore(logprobs)
		assert.InDelta(t, float32(1.0/1.34), conf, 1e-2)
	})

	t.Run("Perplexity 2.0 (50% confidence)", func(t *testing.T) {
		targetLogprob := -math.Log(2.0)
		logprobs := []Logprobs{
			{TokenLogprob: TokenLogprob{Token: "test", Logprob: targetLogprob}},
		}
		conf := confidenceScore(logprobs)
		assert.InDelta(t, float32(0.50), conf, 1e-2)
	})
}


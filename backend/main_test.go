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

func TestCalculateSingleToxicityScore(t *testing.T) {
	t.Run("All zero scores return 0.0", func(t *testing.T) {
		report := ToxicityScannerReport{}
		score := calculateSingleToxicityScore(report)
		assert.Equal(t, float32(0.0), score)
	})

	t.Run("Clean benign text produces near zero score", func(t *testing.T) {
		report := ToxicityScannerReport{
			Toxicity:       0.0006,
			SevereToxicity: 0.0001,
			Obscene:        0.0011,
			Threat:         0.0003,
			Insult:         0.0008,
			IdentityAttack: 0.0001,
		}
		score := calculateSingleToxicityScore(report)
		assert.True(t, score < 0.01)
	})

	t.Run("Critical threat alone spikes overall score", func(t *testing.T) {
		report := ToxicityScannerReport{
			Threat: 0.99,
		}
		score := calculateSingleToxicityScore(report)
		assert.InDelta(t, float32(0.99), score, 1e-3)
	})

	t.Run("Critical severe toxicity alone spikes overall score", func(t *testing.T) {
		report := ToxicityScannerReport{
			SevereToxicity: 0.95,
		}
		score := calculateSingleToxicityScore(report)
		assert.InDelta(t, float32(0.95), score, 1e-3)
	})

	t.Run("Multiple moderate scores compound smoothly", func(t *testing.T) {
		report := ToxicityScannerReport{
			Insult:  0.5,
			Obscene: 0.5,
		}
		score := calculateSingleToxicityScore(report)
		// survival = (1 - 0.6*0.5) * (1 - 0.4*0.5) = 0.70 * 0.80 = 0.56 -> score = 0.44
		assert.InDelta(t, float32(0.44), score, 1e-2)
	})

	t.Run("All maximum scores return 1.0 clamped", func(t *testing.T) {
		report := ToxicityScannerReport{
			Toxicity:       1.0,
			SevereToxicity: 1.0,
			Obscene:        1.0,
			Threat:         1.0,
			Insult:         1.0,
			IdentityAttack: 1.0,
		}
		score := calculateSingleToxicityScore(report)
		assert.Equal(t, float32(1.0), score)
	})
}

func TestRunToxicityScanner(t *testing.T) {
	t.Run("Empty texts returns zero and no error", func(t *testing.T) {
		score, err := runToxicityScanner([]string{})
		assert.NoError(t, err)
		assert.Equal(t, float32(0.0), score)
	})

	t.Run("Live or mock server with array of reports", func(t *testing.T) {
		score, err := runToxicityScanner([]string{"you are a nice person", "shut up you dumb"})
		if err != nil {
			t.Logf("Toxicity scanner container not running, skipping live check: %v", err)
			return
		}
		assert.True(t, score > 0.0)
		assert.True(t, score <= 1.0)
	})
}

func TestRunNLIScanner(t *testing.T) {
	t.Run("Empty input returns nil and no error", func(t *testing.T) {
		report, err := runNLIScanner("", "hypothesis")
		assert.NoError(t, err)
		assert.Nil(t, report)

		report2, err := runNLIScanner("premise", "")
		assert.NoError(t, err)
		assert.Nil(t, report2)
	})

	t.Run("Live or mock server NLI check", func(t *testing.T) {
		report, err := runNLIScanner("The sky is blue", "The sky is green")
		if err != nil {
			t.Logf("NLI scanner container not reachable: %v", err)
			return
		}
		assert.NotNil(t, report)
		assert.Equal(t, "contradiction", report.Label)
		assert.True(t, report.Score > 0.5)
		assert.True(t, report.ContradictionProb > 0.5)
	})
}




package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	t.Run("Chat Completion with Semantic Cache Hit", func(t *testing.T) {
		// Pre-populate cache
		prompt := "What is the capital of Italy?"
		cachedResp := "The capital of Italy is Rome."
		_ = saveSemanticCache(prompt, cachedResp)

		reqBody := map[string]interface{}{
			"model": "openai/gpt-3.5-turbo",
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
			"use_semantic_cache": true,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.NotEmpty(t, response["choices"])
			choices := response["choices"].([]interface{})
			firstChoice := choices[0].(map[string]interface{})
			message := firstChoice["message"].(map[string]interface{})
			assert.Equal(t, cachedResp, message["content"])
			assert.Equal(t, true, response["cached"])
		}
	})

	t.Run("Chat Completion with use_semantic_cache=false", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"model": "openai/gpt-3.5-turbo",
			"messages": []map[string]string{
				{"role": "user", "content": "Hello!"},
			},
			"use_semantic_cache": false,
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
		assert.NotEmpty(t, response["choices"])
		// Should not be marked cached
		assert.Nil(t, response["cached"])
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

func TestSemanticCache(t *testing.T) {
	t.Run("Empty request returns nil and no error", func(t *testing.T) {
		res, err := querySemanticCache("", 0.95)
		assert.NoError(t, err)
		assert.Nil(t, res)

		err = saveSemanticCache("", "")
		assert.NoError(t, err)
	})

	t.Run("Save and Query Semantic Cache", func(t *testing.T) {
		reqText := "What is the capital of Spain?"
		respText := "The capital of Spain is Madrid."

		err := saveSemanticCache(reqText, respText)
		if err != nil {
			t.Logf("Semantic cache container not reachable: %v", err)
			return
		}

		res, err := querySemanticCache(reqText, 0.95)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Found)
		assert.NotNil(t, res.Response)
		assert.Equal(t, respText, *res.Response)
		assert.NotNil(t, res.Score)
		assert.True(t, *res.Score >= 0.95)
	})

	t.Run("Unrelated Query Returns Not Found", func(t *testing.T) {
		res, err := querySemanticCache("Completely unrelated query about astrophysics", 0.95)
		if err != nil {
			t.Logf("Semantic cache container not reachable: %v", err)
			return
		}
		assert.NotNil(t, res)
		assert.False(t, res.Found)
	})
}

func TestSemanticCacheSavings(t *testing.T) {
	InitRedis()
	InitMongoDB()
	ctx := context.Background()

	// Clear redis key and mongo collection for clean test state
	rdb.Del(ctx, SemanticCacheSavingsKey)
	_, _ = requestsCollection.DeleteMany(ctx, bson.D{})

	t.Run("Compute from empty MongoDB returns 0", func(t *testing.T) {
		savings, err := ComputeSemanticCacheSavingsFromDB(ctx)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), savings)
	})

	t.Run("Compute with cached requests non-nvidia and nvidia", func(t *testing.T) {
		// Insert requests:
		// 1: cached=true, provider=openai (+$1)
		// 2: cached=true, provider=anthropic (+$1)
		// 3: cached=true, provider=nvidia (+$0)
		// 4: cached=false, provider=openai (+$0)
		// 5: cached=true, provider=google (+$1)
		records := []interface{}{
			RequestRecord{Endpoint: "/v1/chat/completions", Provider: "openai", Cached: true, Timestamp: time.Now().UTC()},
			RequestRecord{Endpoint: "/v1/chat/completions", Provider: "anthropic", Cached: true, Timestamp: time.Now().UTC()},
			RequestRecord{Endpoint: "/v1/chat/completions", Provider: "nvidia", Cached: true, Timestamp: time.Now().UTC()},
			RequestRecord{Endpoint: "/v1/chat/completions", Provider: "openai", Cached: false, Timestamp: time.Now().UTC()},
			RequestRecord{Endpoint: "/v1/playground", Provider: "google", Cached: true, Timestamp: time.Now().UTC()},
		}
		_, err := requestsCollection.InsertMany(ctx, records)
		assert.NoError(t, err)

		savings, err := ComputeSemanticCacheSavingsFromDB(ctx)
		assert.NoError(t, err)
		assert.Equal(t, float64(3.0), savings)
	})

	t.Run("GetSemanticCacheSavings populates Redis on cache miss", func(t *testing.T) {
		rdb.Del(ctx, SemanticCacheSavingsKey)

		// Key is deleted from Redis, GetSemanticCacheSavings should compute from DB (3.0) and save to Redis
		savings, err := GetSemanticCacheSavings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, float64(3.0), savings)

		// Verify key exists in Redis with value 3
		val, err := rdb.Get(ctx, SemanticCacheSavingsKey).Result()
		assert.NoError(t, err)
		assert.Equal(t, "3", val)
	})

	t.Run("GetSemanticCacheSavings reads from Redis when available", func(t *testing.T) {
		// Overwrite Redis key directly with 10.0
		rdb.Set(ctx, SemanticCacheSavingsKey, 10.0, 0)

		savings, err := GetSemanticCacheSavings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, float64(10.0), savings)
	})

	t.Run("RecordCacheSavings increments Redis for non-nvidia providers", func(t *testing.T) {
		rdb.Set(ctx, SemanticCacheSavingsKey, 5.0, 0)

		// Increment for openai
		err := RecordCacheSavings(ctx, "openai")
		assert.NoError(t, err)

		savings, err := GetSemanticCacheSavings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, float64(6.0), savings)

		// Increment for openrouter
		err = RecordCacheSavings(ctx, "openrouter")
		assert.NoError(t, err)

		savings, err = GetSemanticCacheSavings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, float64(7.0), savings)
	})

	t.Run("RecordCacheSavings does not increment for nvidia provider", func(t *testing.T) {
		rdb.Set(ctx, SemanticCacheSavingsKey, 7.0, 0)

		err := RecordCacheSavings(ctx, "nvidia")
		assert.NoError(t, err)

		savings, err := GetSemanticCacheSavings(ctx)
		assert.NoError(t, err)
		assert.Equal(t, float64(7.0), savings)
	})

	t.Run("GET /v1/cost returns semantic cache savings and average cost", func(t *testing.T) {
		rdb.Set(ctx, SemanticCacheSavingsKey, 12.0, 0)

		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/cost", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, "user_mvp_123", resp["user_id"])
		assert.Contains(t, resp, "semantic_cache_savings")
		assert.Equal(t, float64(12.0), resp["semantic_cache_savings"])
		assert.Equal(t, float64(12.0), resp["semantic_cache_savings_dollars"])
		assert.Equal(t, float64(12000000), resp["semantic_cache_savings_microcents"])

		assert.Contains(t, resp, "average_cost")
		assert.Contains(t, resp, "avg_cost")
	})

	t.Run("ComputeAverageCost with multiple requests", func(t *testing.T) {
		_, _ = requestsCollection.DeleteMany(ctx, bson.D{})

		// Empty DB
		avgDollars, avgMicrocents, err := ComputeAverageCost(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, avgDollars)
		assert.Equal(t, 0.0, avgMicrocents)

		// Insert 4 requests: costs $0, $1, $2, $3 (in microcents: 0, 1M, 2M, 3M) -> total = 6M / 4 = 1.5M ($1.50)
		records := []interface{}{
			RequestRecord{Endpoint: "/v1/chat/completions", CostMicrocents: 0, Timestamp: time.Now().UTC()},
			RequestRecord{Endpoint: "/v1/chat/completions", CostMicrocents: 1000000, Timestamp: time.Now().UTC()},
			RequestRecord{Endpoint: "/v1/chat/completions", CostMicrocents: 2000000, Timestamp: time.Now().UTC()},
			RequestRecord{Endpoint: "/v1/playground", CostMicrocents: 3000000, Timestamp: time.Now().UTC()},
		}
		_, err = requestsCollection.InsertMany(ctx, records)
		assert.NoError(t, err)

		avgDollars, avgMicrocents, err = ComputeAverageCost(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1.5, avgDollars)
		assert.Equal(t, 1500000.0, avgMicrocents)
	})

	t.Run("GetModelAnalytics with model metrics and last week breakdown", func(t *testing.T) {
		_, _ = requestsCollection.DeleteMany(ctx, bson.D{})

		conf1 := float32(0.9)
		conf2 := float32(0.8)
		conf3 := float32(0.6)

		nli1 := &NLIReport{Label: "entailment", Score: 0.95, ContradictionProb: 0.05, NeutralProb: 0.1, EntailmentProb: 0.85}
		nli2 := &NLIReport{Label: "contradiction", Score: 0.85, ContradictionProb: 0.85, NeutralProb: 0.1, EntailmentProb: 0.05}

		records := []interface{}{
			RequestRecord{
				Endpoint:       "/v1/chat/completions",
				Model:          "openai/gpt-4o",
				Provider:       "openai",
				Confidence:     &conf1,
				Toxicity:       0.02,
				NLI:            nli1,
				CostMicrocents: 10000,
				Timestamp:      time.Now().UTC().Add(-2 * time.Hour),
			},
			RequestRecord{
				Endpoint:       "/v1/chat/completions",
				Model:          "openai/gpt-4o",
				Provider:       "openai",
				Confidence:     &conf2,
				Toxicity:       0.04,
				NLI:            nli1,
				CostMicrocents: 20000,
				Timestamp:      time.Now().UTC().Add(-24 * time.Hour),
			},
			RequestRecord{
				Endpoint:       "/v1/playground",
				Model:          "anthropic/claude-3.5-sonnet",
				Provider:       "anthropic",
				Confidence:     &conf3,
				Toxicity:       0.12,
				NLI:            nli2,
				CostMicrocents: 30000,
				Timestamp:      time.Now().UTC().Add(-48 * time.Hour),
			},
		}

		_, err := requestsCollection.InsertMany(ctx, records)
		assert.NoError(t, err)

		analytics, err := GetModelAnalytics(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), analytics.TotalRequests)
		assert.Equal(t, int64(3), analytics.WeeklyRequests)
		assert.Equal(t, 2, len(analytics.Models))

		// First model should be gpt-4o (2 requests)
		assert.Equal(t, "openai/gpt-4o", analytics.Models[0].Model)
		assert.Equal(t, int64(2), analytics.Models[0].RequestCount)
		assert.InDelta(t, 66.67, analytics.Models[0].Percentage, 0.5)
		assert.NotNil(t, analytics.Models[0].AvgConfidence)
		assert.InDelta(t, 0.85, *analytics.Models[0].AvgConfidence, 0.01)
		assert.InDelta(t, 0.03, analytics.Models[0].AvgToxicity, 0.01)
		assert.NotNil(t, analytics.Models[0].AvgHallucination)
		assert.InDelta(t, 0.05, *analytics.Models[0].AvgHallucination, 0.01)

		// Second model should be claude (1 request)
		assert.Equal(t, "anthropic/claude-3.5-sonnet", analytics.Models[1].Model)
		assert.Equal(t, int64(1), analytics.Models[1].RequestCount)
		assert.InDelta(t, 33.33, analytics.Models[1].Percentage, 0.5)
		assert.NotNil(t, analytics.Models[1].AvgConfidence)
		assert.InDelta(t, 0.6, *analytics.Models[1].AvgConfidence, 0.01)
		assert.InDelta(t, 0.12, analytics.Models[1].AvgToxicity, 0.01)
		assert.NotNil(t, analytics.Models[1].AvgHallucination)
		assert.InDelta(t, 0.85, *analytics.Models[1].AvgHallucination, 0.01)

		// Test HTTP endpoint GET /v1/analytics
		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/analytics", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var res AnalyticsResponse
		err = json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), res.TotalRequests)
	})
}





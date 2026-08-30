package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var LoadedConfig struct {
	Enviroment        string
	MockServerBaseUrl string
} = struct {
	Enviroment        string
	MockServerBaseUrl string
}{
	Enviroment:        "testing",
	MockServerBaseUrl: "http://localhost:1080/v1",
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	chatCompletionHandler := func(c *gin.Context) {
		// Hardcoded user ID for MVP
		userID := "user_mvp_123"

		// 1. Guardrail: Check Budget Before Processing
		if err := CheckBudget(c.Request.Context(), userID); err != nil {
			if err.Error() == "daily budget exceeded" {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": gin.H{
						"message": "Daily budget exceeded. Please try again tomorrow.",
						"type":    "quota_exceeded",
					},
				})
				return
			}

			// Fallback for Redis connection errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var req struct {
			Model            string `json:"model" binding:"required"`
			Messages         []struct {
				Role    string `json:"role" binding:"required"`
				Content string `json:"content" binding:"required"`
			} `json:"messages" binding:"required"`
			UseSemanticCache *bool `json:"use_semantic_cache"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": err.Error(),
					"type":    "invalid_request_error",
				},
			})
			return
		}

		// Extract provider and model
		provider, model, err := extractProviderModel(req.Model)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"message": err.Error(),
					"type":    "invalid_request_error",
				},
			})
			return
		}

		// Map messages
		var chatMessages []ChatMessage
		for _, m := range req.Messages {
			chatMessages = append(chatMessages, ChatMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}

		var premise string
		for i := len(chatMessages) - 1; i >= 0; i-- {
			if chatMessages[i].Role == "user" {
				premise = chatMessages[i].Content
				break
			}
		}
		if premise == "" && len(chatMessages) > 0 {
			premise = chatMessages[len(chatMessages)-1].Content
		}

		useSemanticCache := true
		if req.UseSemanticCache != nil {
			useSemanticCache = *req.UseSemanticCache
		}

		if useSemanticCache && premise != "" {
			cached, err := querySemanticCache(premise, 0.95)
			if err == nil && cached != nil && cached.Found && cached.Response != nil && (cached.Score == nil || *cached.Score >= 0.95) {
				cachedContent := *cached.Response
				conf := float32(1.0)
				toxicity := float32(0.0)

				go func() {
					LogRequest(RequestRecord{
						Endpoint:       "/v1/chat/completions",
						Model:          req.Model,
						Provider:       provider,
						Prompt:         premise,
						Messages:       chatMessages,
						Response:       cachedContent,
						Confidence:     &conf,
						Toxicity:       toxicity,
						NLI:            nil,
						CostMicrocents: 0,
						Cached:         true,
					})
					if err := RecordCacheSavings(context.Background(), provider); err != nil {
						log.Printf("Error recording cache savings: %v", err)
					}
				}()

				choice := gin.H{
					"index": 0,
					"message": gin.H{
						"role":    "assistant",
						"content": cachedContent,
					},
					"finish_reason": "stop",
				}

				c.JSON(http.StatusOK, gin.H{
					"id":         "chatcmpl-" + provider + "-cached",
					"object":     "chat.completion",
					"created":    time.Now().Unix(),
					"model":      req.Model,
					"confidence": conf,
					"toxicity":   toxicity,
					"nli":        nil,
					"choices":    []gin.H{choice},
					"cached":     true,
				})
				return
			}
		}

		// Call provider
		responseContent, err := callProvider(c.Request.Context(), provider, model, chatMessages)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message": err.Error(),
					"type":    "provider_error",
				},
			})
			return
		}

		conf := confidenceScore(responseContent.Logprobs)
		toxicity, err := runToxicityScanner([]string{responseContent.Content})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message": err.Error(),
					"type":    "provider_error",
				},
			})
			return
		}

		nliReport, err := runNLIScanner(premise, responseContent.Content)
		if err != nil {
			log.Printf("NLI scanner warning: %v", err)
		}

		// We use background context here so tracking completes even if the client disconnects early
		go func() {
			if err := RecordSpend(context.Background(), userID, responseContent.Cost); err != nil {
				log.Printf("Error recording spend for %s: %v", userID, err)
			}
		}()

		go LogRequest(RequestRecord{
			Endpoint:       "/v1/chat/completions",
			Model:          req.Model,
			Provider:       provider,
			Prompt:         premise,
			Messages:       chatMessages,
			Response:       responseContent.Content,
			Confidence:     &conf,
			Toxicity:       toxicity,
			NLI:            nliReport,
			CostMicrocents: responseContent.Cost,
			Cached:         false,
		})

		if useSemanticCache && premise != "" {
			go func() {
				if err := saveSemanticCache(premise, responseContent.Content); err != nil {
					log.Printf("Error saving to semantic cache: %v", err)
				}
			}()
		}

		// Build the choice map
		choice := gin.H{
			"index": 0,
			"message": gin.H{
				"role":    "assistant",
				"content": responseContent.Content,
			},
			"finish_reason": "stop",
		}

		// Return response
		c.JSON(http.StatusOK, gin.H{
			"id":         "chatcmpl-" + provider,
			"object":     "chat.completion",
			"created":    1700000000,
			"model":      req.Model,
			"confidence": conf,
			"toxicity":   toxicity,
			"nli":        nliReport,
			"choices":    []gin.H{choice},
		})
	}

	// Proxy OpenAI chat completion API
	r.POST("/v1/chat/completions", chatCompletionHandler)

	r.GET("/v1/cost", func(c *gin.Context) {
		userID := "user_mvp_123"

		currentSpend, err := GetCurrentSpend(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		savings, err := GetSemanticCacheSavings(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		avgCostDollars, avgCostMicrocents, err := ComputeAverageCost(c.Request.Context())
		if err != nil {
			log.Printf("Warning: failed to compute average cost: %v", err)
			avgCostDollars = 0.0
			avgCostMicrocents = 0.0
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":                            userID,
			"cost_microcents":                   currentSpend,
			"cost_dollars":                      float64(currentSpend) / 1000000.0,
			"semantic_cache_savings":            savings,
			"semantic_cache_savings_dollars":    savings,
			"semantic_cache_savings_microcents": int64(savings * 1000000.0),
			"savings_dollars":                   savings,
			"savings_microcents":                int64(savings * 1000000.0),
			"total_savings":                     savings,
			"average_cost":                      avgCostDollars,
			"average_cost_dollars":              avgCostDollars,
			"average_cost_microcents":           avgCostMicrocents,
			"avg_cost":                          avgCostDollars,
			"avg_cost_dollars":                  avgCostDollars,
			"avg_cost_microcents":               avgCostMicrocents,
		})
	})

	r.GET("/v1/config", func(c *gin.Context) {
		dailyLim, monthlyLim := getPerUserLimits()

		c.JSON(http.StatusOK, gin.H{
			"per_user_daily_limit":   dailyLim,
			"per_user_monthly_limit": monthlyLim,
		})
	})

	r.POST("/v1/config", func(c *gin.Context) {
		var req struct {
			PerUserDailyLimit   *int64 `json:"per_user_daily_limit"`
			PerUserMonthlyLimit *int64 `json:"per_user_monthly_limit"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()

		if req.PerUserDailyLimit != nil {
			rdb.Set(ctx, "config:per_user_daily_limit", *req.PerUserDailyLimit, 0)
		}

		if req.PerUserMonthlyLimit != nil {
			rdb.Set(ctx, "config:per_user_monthly_limit", *req.PerUserMonthlyLimit, 0)
		}

		c.JSON(http.StatusOK, gin.H{
			"per_user_daily_limit":   *req.PerUserDailyLimit,
			"per_user_monthly_limit": *&req.PerUserMonthlyLimit,
		})
	})

	playgroundHandler := func(c *gin.Context) {
		// Hardcoded user ID for MVP
		userID := "user_mvp_123"

		// 1. Guardrail: Check Budget Before Processing
		if err := CheckBudget(c.Request.Context(), userID); err != nil {
			if err.Error() == "daily budget exceeded" {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Daily budget exceeded. Please try again tomorrow.",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var req struct {
			Prompt           string `json:"prompt" binding:"required"`
			ModelSpec        string `json:"model_spec" binding:"required"`
			UseSemanticCache *bool  `json:"use_semantic_cache"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		provider, model, err := extractProviderModel(req.ModelSpec)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		start := time.Now()

		useSemanticCache := true
		if req.UseSemanticCache != nil {
			useSemanticCache = *req.UseSemanticCache
		}

		if useSemanticCache && req.Prompt != "" {
			cached, err := querySemanticCache(req.Prompt, 0.95)
			if err == nil && cached != nil && cached.Found && cached.Response != nil && (cached.Score == nil || *cached.Score >= 0.95) {
				cachedContent := *cached.Response
				conf := float32(1.0)
				toxicity := float32(0.0)
				latencyMs := time.Since(start).Milliseconds()

				go func() {
					LogRequest(RequestRecord{
						Endpoint:       "/v1/playground",
						Model:          req.ModelSpec,
						Provider:       provider,
						Prompt:         req.Prompt,
						Response:       cachedContent,
						Confidence:     &conf,
						Toxicity:       toxicity,
						NLI:            nil,
						LatencyMs:      latencyMs,
						CostMicrocents: 0,
						Cached:         true,
					})
					if err := RecordCacheSavings(context.Background(), provider); err != nil {
						log.Printf("Error recording cache savings: %v", err)
					}
				}()

				c.JSON(http.StatusOK, gin.H{
					"model":      model,
					"provider":   provider,
					"content":    cachedContent,
					"confidence": &conf,
					"toxicity":   toxicity,
					"nli":        nil,
					"latency_ms": latencyMs,
					"cost":       0.0,
					"cached":     true,
				})
				return
			}
		}

		messages := []ChatMessage{
			{Role: "user", Content: req.Prompt},
		}

		// Track latency
		responseContent, err := callProvider(c.Request.Context(), provider, model, messages)
		latencyMs := time.Since(start).Milliseconds()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var conf *float32
		if len(responseContent.Logprobs) > 0 {
			c := confidenceScore(responseContent.Logprobs)
			conf = &c
		}
		toxicity, err := runToxicityScanner([]string{responseContent.Content})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"message": err.Error(),
					"type":    "provider_error",
				},
			})
			return
		}

		nliReport, err := runNLIScanner(req.Prompt, responseContent.Content)
		if err != nil {
			log.Printf("NLI scanner warning: %v", err)
		}

		go func() {
			if err := RecordSpend(context.Background(), userID, responseContent.Cost); err != nil {
				log.Printf("Error recording spend for %s: %v", userID, err)
			}
		}()

		go LogRequest(RequestRecord{
			Endpoint:       "/v1/playground",
			Model:          req.ModelSpec,
			Provider:       provider,
			Prompt:         req.Prompt,
			Response:       responseContent.Content,
			Confidence:     conf,
			Toxicity:       toxicity,
			NLI:            nliReport,
			LatencyMs:      latencyMs,
			CostMicrocents: responseContent.Cost,
			Cached:         false,
		})

		if useSemanticCache && req.Prompt != "" {
			go func() {
				if err := saveSemanticCache(req.Prompt, responseContent.Content); err != nil {
					log.Printf("Error saving to semantic cache: %v", err)
				}
			}()
		}

		c.JSON(http.StatusOK, gin.H{
			"model":      model,
			"provider":   provider,
			"content":    responseContent.Content,
			"confidence": conf,
			"toxicity":   toxicity,
			"nli":        nliReport,
			"latency_ms": latencyMs,
			"cost":       responseContent.Cost / 1_000_000,
		})
	}

	r.POST("/v1/playground", playgroundHandler)

	r.GET("/v1/requests", func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "50")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit <= 0 {
			limit = 50
		}
		if limit > 200 {
			limit = 200
		}

		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil || offset < 0 {
			offset = 0
		}

		records, total, err := FetchRequests(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"requests": records,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		})
	})

	r.GET("/v1/requests/:id", func(c *gin.Context) {
		id := c.Param("id")

		record, err := FetchRequestByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
			return
		}

		c.JSON(http.StatusOK, record)
	})

	analyticsHandler := func(c *gin.Context) {
		stats, err := GetModelAnalytics(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	}

	r.GET("/v1/analytics", analyticsHandler)
	r.GET("/v1/stats", analyticsHandler)

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	InitRedis()
	InitMongoDB()

	r := setupRouter()

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

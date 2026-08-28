package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	// Proxy OpenAI chat completion API
	r.POST("/v1/chat/completions", func(c *gin.Context) {
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
			Model    string `json:"model" binding:"required"`
			Messages []struct {
				Role    string `json:"role" binding:"required"`
				Content string `json:"content" binding:"required"`
			} `json:"messages" binding:"required"`
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

		// We use background context here so tracking completes even if the client disconnects early
		go func() {
			if err := RecordSpend(context.Background(), userID, responseContent.Cost); err != nil {
				log.Printf("Error recording spend for %s: %v", userID, err)
			}
		}()

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
			"choices":    []gin.H{choice},
		})
	})

	r.GET("/v1/cost", func(c *gin.Context) {
		userID := "user_mvp_123"

		currentSpend, err := GetCurrentSpend(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":         userID,
			"cost_microcents": currentSpend,
			"cost_dollars":    float64(currentSpend) / 1000000.0,
		})
	})

	r.GET("/v1/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"per_user_daily_limit":   BudgetGuardRails.PerUserDailyLimit,
			"per_user_monthly_limit": BudgetGuardRails.PerUserMonthlyLimit,
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
			BudgetGuardRails.PerUserDailyLimit = *req.PerUserDailyLimit
			rdb.Set(ctx, "config:per_user_daily_limit", *req.PerUserDailyLimit, 0)
		}

		if req.PerUserMonthlyLimit != nil {
			BudgetGuardRails.PerUserMonthlyLimit = *req.PerUserMonthlyLimit
			rdb.Set(ctx, "config:per_user_monthly_limit", *req.PerUserMonthlyLimit, 0)
		}

		c.JSON(http.StatusOK, gin.H{
			"per_user_daily_limit":   BudgetGuardRails.PerUserDailyLimit,
			"per_user_monthly_limit": BudgetGuardRails.PerUserMonthlyLimit,
		})
	})

	r.POST("/v1/playground", func(c *gin.Context) {
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
			Prompt    string `json:"prompt" binding:"required"`
			ModelSpec string `json:"model_spec" binding:"required"`
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

		messages := []ChatMessage{
			{Role: "user", Content: req.Prompt},
		}

		// Track latency
		start := time.Now()
		responseContent, err := callProvider(c.Request.Context(), provider, model, messages)
		latencyMs := time.Since(start).Milliseconds()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var conf *float32
		if len(responseContent.Logprobs) > 0 {
			c := confidenceScore(responseContent.Logprobs)
			// convert 0-1 confidence back to 0-100 percentage for UI
			cPercent := c * 100
			conf = &cPercent
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

		go func() {
			if err := RecordSpend(context.Background(), userID, responseContent.Cost); err != nil {
				log.Printf("Error recording spend for %s: %v", userID, err)
			}
		}()

		c.JSON(http.StatusOK, gin.H{
			"model":      model,
			"provider":   provider,
			"content":    responseContent.Content,
			"confidence": conf,
			"toxicity":   toxicity,
			"latency_ms": latencyMs,
			"cost":       responseContent.Cost / 1_000_000,
		})
	})

	return r
}

func main() {
	InitRedis()

	r := setupRouter()

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

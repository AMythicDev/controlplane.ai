package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Enviroment        string
	MockServerBaseUrl string
}

var LoadedConfig = Config{
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

		if len(responseContent.Logprobs) > 0 {
			fmt.Println(confidenceScore(responseContent.Logprobs))
		}

		// 2. Track Cost: Simulate cost recording for MVP (e.g., $0.01 per request)
		// In production, this happens AFTER receiving the provider response and token usage
		simulatedCost := int64(10000 * 100) // 1 cent = 10,000 microcents

		// We use background context here so tracking completes even if the client disconnects early
		go func() {
			if err := RecordSpend(context.Background(), userID, simulatedCost); err != nil {
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
			"id":      "chatcmpl-" + provider,
			"object":  "chat.completion",
			"created": 1700000000,
			"model":   req.Model,
			"choices": []gin.H{choice},
		})
	})

	// Get total cost accumulated
	r.GET("/v1/cost", func(c *gin.Context) {
		// Hardcoded user ID for MVP
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

	return r
}

func main() {
	InitRedis()

	r := setupRouter()

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

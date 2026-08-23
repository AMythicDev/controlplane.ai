package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChatCompletionsProxy(t *testing.T) {
	// Initialize the Redis client which will connect to the local valkey instance
	InitRedis()

	router := setupRouter()

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
		assert.Equal(t, "gpt-3.5-turbo", response["model"])
		assert.NotEmpty(t, response["choices"])
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

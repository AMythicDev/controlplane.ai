package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSingleChat(t *testing.T) {
	expectation, err := os.Open("expectation.json")
	if err != nil {
		expectation, err = os.Open("../expectation.json")
	}
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

	t.Run("Valid Server Chat Request", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"model": "openai/gpt-3.5-turbo",
			"messages": []map[string]string{
				{"role": "user", "content": "Hello!"},
			},
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "http://localhost:8080/v1/chat/completions", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send chat request to server: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, "chat.completion", response["object"])
		assert.Equal(t, "openai/gpt-3.5-turbo", response["model"])

		assert.NotEmpty(t, response["choices"])
		choices := response["choices"].([]interface{})
		firstChoice := choices[0].(map[string]interface{})
		message := firstChoice["message"].(map[string]interface{})
		assert.Equal(t, "Hello, world!", message["content"])
	})
}

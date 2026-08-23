package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"google.golang.org/genai"
)

type ChatMessage struct {
	Role    string
	Content string
}

func callProvider(ctx context.Context, provider string, model string, messages []ChatMessage) (string, error) {
	switch provider {
	case "openai":
		client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
		
		var openAIMessages []openai.ChatCompletionMessageParamUnion
		for _, m := range messages {
			if m.Role == "user" {
				openAIMessages = append(openAIMessages, openai.UserMessage(m.Content))
			} else if m.Role == "assistant" {
				openAIMessages = append(openAIMessages, openai.AssistantMessage(m.Content))
			} else if m.Role == "system" {
				openAIMessages = append(openAIMessages, openai.SystemMessage(m.Content))
			}
		}
		
		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(model),
			Messages: openAIMessages,
		})
		if err != nil {
			return "", err
		}
		if len(resp.Choices) > 0 {
			return resp.Choices[0].Message.Content, nil
		}
		return "", fmt.Errorf("no response from OpenAI")

	case "openrouter":
		client := openai.NewClient(
			option.WithBaseURL("https://openrouter.ai/api/v1"),
			option.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
		)
		
		var openAIMessages []openai.ChatCompletionMessageParamUnion
		for _, m := range messages {
			if m.Role == "user" {
				openAIMessages = append(openAIMessages, openai.UserMessage(m.Content))
			} else if m.Role == "assistant" {
				openAIMessages = append(openAIMessages, openai.AssistantMessage(m.Content))
			} else if m.Role == "system" {
				openAIMessages = append(openAIMessages, openai.SystemMessage(m.Content))
			}
		}
		
		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(model),
			Messages: openAIMessages,
		})
		if err != nil {
			return "", err
		}
		if len(resp.Choices) > 0 {
			return resp.Choices[0].Message.Content, nil
		}
		return "", fmt.Errorf("no response from OpenRouter")

	case "anthropic":
		client := anthropic.NewClient(anthropicoption.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
		
		var anthropicMessages []anthropic.MessageParam
		var systemPrompt string
		
		for _, m := range messages {
			if m.Role == "user" {
				anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content)))
			} else if m.Role == "assistant" {
				anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Content)))
			} else if m.Role == "system" {
				systemPrompt = m.Content
			}
		}
		
		params := anthropic.MessageNewParams{
			Model:     anthropic.Model(model),
			MaxTokens: 4096,
			Messages:  anthropicMessages,
		}
		
		if systemPrompt != "" {
			params.System = []anthropic.TextBlockParam{
				{Text: systemPrompt},
			}
		}
		
		resp, err := client.Messages.New(ctx, params)
		if err != nil {
			return "", err
		}
		if len(resp.Content) > 0 {
			return resp.Content[0].Text, nil
		}
		return "", fmt.Errorf("no response from Anthropic")

	case "google":
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		})
		if err != nil {
			return "", err
		}
		
		var contents []*genai.Content
		var systemInstruction *genai.Content
		
		for _, m := range messages {
			if m.Role == "user" {
				contents = append(contents, &genai.Content{
					Role: "user",
					Parts: []*genai.Part{{Text: m.Content}},
				})
			} else if m.Role == "assistant" {
				contents = append(contents, &genai.Content{
					Role: "model",
					Parts: []*genai.Part{{Text: m.Content}},
				})
			} else if m.Role == "system" {
				systemInstruction = &genai.Content{
					Parts: []*genai.Part{{Text: m.Content}},
				}
			}
		}
		
		config := &genai.GenerateContentConfig{}
		if systemInstruction != nil {
			config.SystemInstruction = systemInstruction
		}
		
		resp, err := client.Models.GenerateContent(ctx, model, contents, config)
		if err != nil {
			return "", err
		}
		
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			return resp.Candidates[0].Content.Parts[0].Text, nil
		}
		return "", fmt.Errorf("no response from Google")

	case "nvidia":
		client := openai.NewClient(
			option.WithBaseURL("https://integrate.api.nvidia.com/v1"),
			option.WithAPIKey(os.Getenv("NVIDIA_API_KEY")),
		)
		
		var openAIMessages []openai.ChatCompletionMessageParamUnion
		for _, m := range messages {
			if m.Role == "user" {
				openAIMessages = append(openAIMessages, openai.UserMessage(m.Content))
			} else if m.Role == "assistant" {
				openAIMessages = append(openAIMessages, openai.AssistantMessage(m.Content))
			} else if m.Role == "system" {
				openAIMessages = append(openAIMessages, openai.SystemMessage(m.Content))
			}
		}
		
		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(model),
			Messages: openAIMessages,
		})
		if err != nil {
			return "", err
		}
		if len(resp.Choices) > 0 {
			return resp.Choices[0].Message.Content, nil
		}
		return "", fmt.Errorf("no response from NVIDIA")

	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

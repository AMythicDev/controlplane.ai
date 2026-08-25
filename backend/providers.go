package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"google.golang.org/genai"
)

var SUPPORTED_PROVIDERS = [5]string{"anthropic", "google", "openai", "openrouter", "nvidia"}

func extractProviderModel(spec string) (string, string, error) {
	split_spec := strings.SplitN(spec, "/", 2)
	if len(split_spec) != 2 {
		return "", "", fmt.Errorf("invalid format spec: '%s'", spec)
	}

	provider := split_spec[0]
	model := split_spec[1]

	if !slices.Contains(SUPPORTED_PROVIDERS[:], provider) {
		return "", "", fmt.Errorf("invalid provider: '%s'", provider)
	}

	return provider, model, nil
}

type ChatMessage struct {
	Role    string
	Content string
}

type ChatResponse struct {
	Content  string
	Logprobs []float32
}

func callProvider(ctx context.Context, provider string, model string, messages []ChatMessage) (ChatResponse, error) {
	switch provider {
	case "openai":
		var clientOpts []option.RequestOption
		clientOpts = append(clientOpts, option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
		if LoadedConfig.Enviroment == "testing" {
			clientOpts = append(clientOpts, option.WithBaseURL(LoadedConfig.MockServerBaseUrl))
		}
		client := openai.NewClient(clientOpts...)

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
			Logprobs: openai.Bool(true),
		})
		if err != nil {
			return ChatResponse{}, err
		}
		if len(resp.Choices) > 0 {
			var logprobs []float32
			if len(resp.Choices[0].Logprobs.Content) > 0 {
				for _, lp := range resp.Choices[0].Logprobs.Content {
					logprobs = append(logprobs, float32(lp.Logprob))
				}
			}
			return ChatResponse{Content: resp.Choices[0].Message.Content, Logprobs: logprobs}, nil
		}
		return ChatResponse{}, fmt.Errorf("no response from OpenAI")

	case "openrouter":
		var clientOpts []option.RequestOption
		clientOpts = append(clientOpts, option.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")))
		if LoadedConfig.Enviroment == "testing" {
			clientOpts = append(clientOpts, option.WithBaseURL(LoadedConfig.MockServerBaseUrl))
		} else {
			clientOpts = append(clientOpts, option.WithBaseURL("https://openrouter.ai/api/v1"))
		}
		client := openai.NewClient(clientOpts...)

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
			return ChatResponse{}, err
		}
		if len(resp.Choices) > 0 {
			return ChatResponse{Content: resp.Choices[0].Message.Content}, nil
		}
		return ChatResponse{}, fmt.Errorf("no response from OpenRouter")

	case "anthropic":
		var clientOpts []anthropicoption.RequestOption
		clientOpts = append(clientOpts, anthropicoption.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")))
		if LoadedConfig.Enviroment == "testing" {
			clientOpts = append(clientOpts, anthropicoption.WithBaseURL(LoadedConfig.MockServerBaseUrl))
		}
		client := anthropic.NewClient(clientOpts...)

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
			return ChatResponse{}, err
		}
		if len(resp.Content) > 0 {
			return ChatResponse{Content: resp.Content[0].Text}, nil
		}
		return ChatResponse{}, fmt.Errorf("no response from Anthropic")

	case "google":
		clientConfig := &genai.ClientConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
		}
		if LoadedConfig.Enviroment == "testing" {
			clientConfig.HTTPOptions = genai.HTTPOptions{
				BaseURL: LoadedConfig.MockServerBaseUrl,
			}
		}
		client, err := genai.NewClient(ctx, clientConfig)
		if err != nil {
			return ChatResponse{}, err
		}

		var contents []*genai.Content
		var systemInstruction *genai.Content

		for _, m := range messages {
			if m.Role == "user" {
				contents = append(contents, &genai.Content{
					Role:  "user",
					Parts: []*genai.Part{{Text: m.Content}},
				})
			} else if m.Role == "assistant" {
				contents = append(contents, &genai.Content{
					Role:  "model",
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
			return ChatResponse{}, err
		}

		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			return ChatResponse{Content: resp.Candidates[0].Content.Parts[0].Text}, nil
		}
		return ChatResponse{}, fmt.Errorf("no response from Google")

	case "nvidia":
		var clientOpts []option.RequestOption
		clientOpts = append(clientOpts, option.WithAPIKey(os.Getenv("NVIDIA_API_KEY")))
		if LoadedConfig.Enviroment == "testing" {
			clientOpts = append(clientOpts, option.WithBaseURL(LoadedConfig.MockServerBaseUrl))
		} else {
			clientOpts = append(clientOpts, option.WithBaseURL("https://integrate.api.nvidia.com/v1"))
		}
		client := openai.NewClient(clientOpts...)

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
			Logprobs: openai.Bool(true),
		})
		if err != nil {
			return ChatResponse{}, err
		}
		if len(resp.Choices) > 0 {
			var logprobs []float32
			if len(resp.Choices[0].Logprobs.Content) > 0 {
				for _, lp := range resp.Choices[0].Logprobs.Content {
					logprobs = append(logprobs, float32(lp.Logprob))
				}
			}
			return ChatResponse{Content: resp.Choices[0].Message.Content, Logprobs: logprobs}, nil
		}
		return ChatResponse{}, fmt.Errorf("no response from NVIDIA")

	default:
		return ChatResponse{}, fmt.Errorf("unsupported provider: %s", provider)
	}
}

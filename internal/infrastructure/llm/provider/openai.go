package provider

import (
	"context"
	"errors"
	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/services"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIClient implements the ProviderClient for the OpenAI API.
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient creates a new OpenAI LLM service client.
func NewOpenAIClient(apiKey, apiBaseOverride string) (ProviderClient, error) {
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is not provided")
	}

	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if apiBaseOverride != "" {
		opts = append(opts, option.WithBaseURL(apiBaseOverride))
	}

	client := openai.NewClient(opts...)
	return &OpenAIClient{client: &client}, nil
}

// StreamChatCompletion sends a chat request and streams the response.
func (c *OpenAIClient) StreamChatCompletion(
	ctx context.Context,
	model *entities.Model,
	messages []*entities.Message,
) (<-chan services.ChatStreamEvent, error) {
	// 1. Convert domain messages to OpenAI messages
	openAIMessages, err := c.toOpenAIMessages(messages)
	if err != nil {
		return nil, err
	}

	// 2. Create the stream request
	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(model.Name), // Use the model name from the conversation
		Messages: openAIMessages,
	}

	stream := c.client.Chat.Completions.NewStreaming(ctx, params)

	// 3. Create a channel to send events back to the use case
	events := make(chan services.ChatStreamEvent)

	// 4. Goroutine to process the stream
	go func() {
		defer close(events)

		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 {
				event := services.ChatStreamEvent{
					ContentDelta: chunk.Choices[0].Delta.Content,
				}
				events <- event
			}
		}

		if stream.Err() != nil {
			events <- services.ChatStreamEvent{Error: stream.Err(), IsLast: true}
			return
		}

		// Send final event to signal the end of the stream
		events <- services.ChatStreamEvent{IsLast: true}
	}()

	return events, nil
}

func (c *OpenAIClient) toOpenAIMessages(messages []*entities.Message) ([]openai.ChatCompletionMessageParamUnion, error) {
	openAIMessages := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, msg := range messages {
		var param openai.ChatCompletionMessageParamUnion
		switch msg.Role {
		case entities.MessageRoleUser:
			param = openai.ChatCompletionMessageParamUnion{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(msg.Content),
					},
				},
			}
		case entities.MessageRoleAssistant:
			param = openai.ChatCompletionMessageParamUnion{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String(msg.Content),
					},
				},
			}
		case entities.MessageRoleSystem:
			param = openai.ChatCompletionMessageParamUnion{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String(msg.Content),
					},
				},
			}
		// TODO: Handle tool messages
		default:
			// Let's be strict for now and return an error for unhandled roles.
			return nil, errors.New("unsupported message role: " + string(msg.Role))
		}
		openAIMessages[i] = param
	}
	return openAIMessages, nil
} 
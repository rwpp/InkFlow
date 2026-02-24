package llm

import (
	"context"
	"errors"

	"github.com/KNICEX/InkFlow/internal/ai/internal/domain"
	"github.com/KNICEX/InkFlow/internal/ai/internal/service"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIGoService 使用官方 openai-go 客户端（讯飞星火等 OpenAI 兼容 API）。
// 与 demo 一致：NewClient(WithAPIKey, WithBaseURL) + Chat.Completions.New(Model, Messages, Temperature, MaxTokens)。
type OpenAIGoService struct {
	client  *openai.Client
	modelID string
}

// NewOpenAIGoService 使用 apiKey、baseURL、modelID 创建，从配置 llm.openai_go 读取。
func NewOpenAIGoService(apiKey, baseURL, modelID string) service.LLMService {
	c := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)
	svc := &OpenAIGoService{client: &c, modelID: modelID}
	return svc
}

func (s *OpenAIGoService) AskOnce(ctx context.Context, question string) (domain.Resp, error) {
	resp, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: s.modelID,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Temperature: openai.Float(0.7),
		MaxTokens:   openai.Int(4096),
	})
	if err != nil {
		return domain.Resp{}, err
	}
	if len(resp.Choices) == 0 {
		return domain.Resp{}, errNoChoices
	}
	content := resp.Choices[0].Message.Content
	token := int64(resp.Usage.TotalTokens)
	return domain.Resp{Content: content, Token: token}, nil
}

func (s *OpenAIGoService) BeginChat(ctx context.Context) (service.LLMSession, error) {
	return &openaiGoSession{client: s.client, modelID: s.modelID}, nil
}

var errNoChoices = errors.New("openai: no choices in response")

type openaiGoSession struct {
	client   *openai.Client
	modelID  string
	messages []openai.ChatCompletionMessageParamUnion
}

func (c *openaiGoSession) Ask(ctx context.Context, question string) (domain.Resp, error) {
	c.messages = append(c.messages, openai.UserMessage(question))
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.modelID, Messages: c.messages,
		Temperature: openai.Float(0.7), MaxTokens: openai.Int(4096),
	})
	if err != nil {
		return domain.Resp{}, err
	}
	if len(resp.Choices) == 0 {
		return domain.Resp{}, errNoChoices
	}
	msg := resp.Choices[0].Message
	c.messages = append(c.messages, openai.AssistantMessage(msg.Content))
	return domain.Resp{Content: msg.Content, Token: int64(resp.Usage.TotalTokens)}, nil
}

func (c *openaiGoSession) Close() error { return nil }

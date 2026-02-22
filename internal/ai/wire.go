package ai

import (
	"github.com/KNICEX/InkFlow/internal/ai/internal/service"
	"github.com/KNICEX/InkFlow/internal/ai/internal/service/llm"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
)

// InitLLMServices 组装所有 LLM 后端（默认优先星火 openai_go，其次 Gemini），供 InitLLMService 做 failover。
func InitLLMServices(geminiClients []*genai.Client) []LLMService {
	svcs := make([]LLMService, 0, len(geminiClients)+1)
	var openaiGoCfg struct {
		APIKey  string `mapstructure:"api_key"`
		BaseURL string `mapstructure:"base_url"`
		ModelID string `mapstructure:"model_id"`
	}
	if err := viper.UnmarshalKey("llm.openai_go", &openaiGoCfg); err == nil &&
		openaiGoCfg.APIKey != "" && openaiGoCfg.BaseURL != "" && openaiGoCfg.ModelID != "" {
		svcs = append(svcs, llm.NewOpenAIGoService(openaiGoCfg.APIKey, openaiGoCfg.BaseURL, openaiGoCfg.ModelID))
	}
	for _, c := range geminiClients {
		svcs = append(svcs, llm.NewGeminiService(c))
	}
	return svcs
}

// InitLLMService 对多个 LLM 后端做 failover（默认优先星火，其次 Gemini）。至少需配置一种。
func InitLLMService(svcs []LLMService) LLMService {
	if len(svcs) == 0 {
		panic("ai: 至少需配置一种 LLM：llm.gemini.key 或 llm.openai_go（api_key/base_url/model_id）")
	}
	return service.NewFailoverService(svcs)
}

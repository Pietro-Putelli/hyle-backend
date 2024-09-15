package langchain

import (
	"context"
	"encoding/json"

	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

func extractKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func NewLLMRequest(prompt string, inputs map[string]interface{}) (map[string]interface{}, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	langchainCfg := cfg.Langchain

	llm, err := openai.New(openai.WithToken(langchainCfg.OpenAIKey), openai.WithModel(langchainCfg.GPTModel))
	if err != nil {
		return nil, err
	}

	promptTempalte := prompts.NewPromptTemplate(
		prompt,
		extractKeys(inputs),
	)

	result, err := promptTempalte.Format(inputs)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	completion, err := llm.Call(ctx, result)
	// , llms.WithTemperature(0.6)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(completion), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

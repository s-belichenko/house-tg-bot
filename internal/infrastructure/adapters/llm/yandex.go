package llm

import (
	"context"
	"fmt"

	"s-belichenko/house-tg-bot/internal/config"

	"s-belichenko/house-tg-bot/internal/domain/ports"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
	pkgTemplate "s-belichenko/house-tg-bot/pkg/template"

	"github.com/sheeiavellie/go-yandexgpt"
)

type Yandex struct {
	config        config.App
	client        *yandexgpt.YandexGPTClient
	logger        pkgLogger.Logger
	renderingTool pkgTemplate.RenderingTool
}

func NewYandexLLM(config config.App, logger pkgLogger.Logger) ports.LLM {
	renderingTool := pkgTemplate.NewTool("llm", logger)

	config.LlmYandex.SystemPrompt = renderingTool.RenderText(
		"systemPrompt.gohtml",
		struct {
			BotName     string
			HomeAddress string
		}{
			BotName:     config.LlmYandex.BotName,
			HomeAddress: config.HomeAddress,
		},
	)
	client := yandexgpt.NewYandexGPTClientWithAPIKey(config.LlmYandex.LLMApiToken)

	return &Yandex{
		config:        config,
		logger:        logger,
		client:        client,
		renderingTool: renderingTool,
	}
}

func (y *Yandex) DoRequest(question string) string {
	request := y.createRequest(question)
	response, err := y.client.GetCompletion(context.Background(), request)
	if err != nil {
		y.logger.Error(fmt.Sprintf(`LLM request error: %s`, err.Error()), pkgLogger.LogContext{
			`request`: request,
		})

		return ``
	}

	return response.Result.Alternatives[0].Message.Text
}

func (y *Yandex) createRequest(question string) yandexgpt.YandexGPTRequest {
	return yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(y.config.LlmYandex.LLMFolderID, yandexgpt.YandexGPT4Model),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: y.config.LlmYandex.LLMTemperature,
			MaxTokens:   y.config.LlmYandex.MaxTokens,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: y.config.LlmYandex.SystemPrompt,
			},
			{
				Role: yandexgpt.YandexGPTMessageRoleUser,
				Text: question,
			},
		},
	}
}

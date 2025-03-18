package llm

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sheeiavellie/go-yandexgpt"
)

const BotNickname = "Тринадцатый"

type ConfigLLM struct {
	llmApiToken    string `env:"LLM_API_TOKEN"`
	systemPrompt   string
	llmFolderId    string  `env:"LLM_FOLDER_ID"`
	llmTemperature float32 `env:"LLM_TEMPERATURE" env-default:"0.7"`
	maxTokens      int     `env:"LLM_MAX_TOKENS" env-default:"8000"`
}

var config = ConfigLLM{}

func init() {
	godotenv.Load()
	config.systemPrompt = fmt.Sprintf(
		"Тебя зовут %s. Ты очень полезный ИИ-помощник. Отвечай на вопросы коротко и точно. Если не знаешь ответ, вежливо напиши об этом. Не нужно уточнять, что отвечает именно менеджер ПИК.",
		BotNickname,
	)
}

func GetAnswerAboutKeys() (string, error) {
	client := yandexgpt.NewYandexGPTClientWithAPIKey(config.llmApiToken)
	request := yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(config.llmFolderId, yandexgpt.YandexGPT4Model),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: config.llmTemperature,
			MaxTokens:   config.maxTokens,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: config.systemPrompt,
			},
			{
				Role: yandexgpt.YandexGPTMessageRoleUser,
				Text: "Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир.",
			},
		},
	}

	response, err := client.GetCompletion(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("LLM request error: %s", err.Error())
	}

	return response.Result.Alternatives[0].Message.Text, nil
}

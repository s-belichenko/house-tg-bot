package llm

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sheeiavellie/go-yandexgpt"
)

const BotNickname = "Тринадцатый"

type ConfigLLM struct {
	LLMApiToken    string `env:"LLM_API_TOKEN"`
	SystemPrompt   string
	LLMFolderId    string  `env:"LLM_FOLDER_ID"`
	LLMTemperature float32 `env:"LLM_TEMPERATURE" env-default:"0.7"`
	MaxTokens      int     `env:"LLM_MAX_TOKENS" env-default:"8000"`
}

var config ConfigLLM

func init() {
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		fmt.Printf("Error reading LLM config: %v", err)
	}
	config.SystemPrompt = fmt.Sprintf(
		"Тебя зовут %s. Ты очень полезный ИИ-помощник. Отвечай на вопросы коротко и точно. Если не знаешь ответ, вежливо напиши об этом. Не нужно уточнять, что отвечает именно менеджер ПИК.",
		BotNickname,
	)
}

func GetAnswerAboutKeys() (string, error) {
	client := yandexgpt.NewYandexGPTClientWithAPIKey(config.LLMApiToken)
	request := yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(config.LLMFolderId, yandexgpt.YandexGPT4Model),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: config.LLMTemperature,
			MaxTokens:   config.MaxTokens,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: config.SystemPrompt,
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

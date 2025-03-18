package llm

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sheeiavellie/go-yandexgpt"
	"math/rand"
	"time"

	yandexLogger "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

const botNickname = "Тринадцатый"

type ConfigLLM struct {
	LLMApiToken    string `env:"LLM_API_TOKEN"`
	SystemPrompt   string
	LLMFolderId    string  `env:"LLM_FOLDER_ID"`
	LLMTemperature float32 `env:"LLM_TEMPERATURE" env-default:"0.7"`
	MaxTokens      int     `env:"LLM_MAX_TOKENS" env-default:"8000"`
}

var config ConfigLLM

var clarifications = []string{
	"Отключи рамки корректности",
	"Отвечай максимально жестко",
}

var questions = []string{
	"Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир",
	"Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир",
	"Придумай резкий ответ менеджера поддержки компании ПИК на вопрос \"когда вы выдадите мне ключи от моей новой квартиры",
	"Придумай глупый ответ менеджера поддержки компании ПИК на вопрос \"когда вы выдадите мне ключи от моей новой квартиры",
}

var log *yandexLogger.Logger

func init() {
	log = yandexLogger.NewLogger("answers_stream")
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		fmt.Printf("Error reading LLM config: %v", err)
	}
	config.SystemPrompt = fmt.Sprintf(
		"Тебя зовут %s. Ты чат-бот в чате про дом по Адресу Москва, Иловайская, 2.\n"+
			"Люди любят тебя за юмор и за то, что ты всегда остро и смешно отвечаешь.\n"+
			"Отвечай на вопросы коротко и точно.\n"+
			"Если не знаешь ответ, напиши об этом.",
		botNickname,
	)
}

func GetAnswerAboutKeys() (string, error) {
	client := yandexgpt.NewYandexGPTClientWithAPIKey(config.LLMApiToken)
	question := generateAnswer()
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
				Text: question,
			},
		},
	}

	response, err := client.GetCompletion(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("LLM request error: %s", err.Error())
	}

	answer := response.Result.Alternatives[0].Message.Text
	log.Info("Получен ответ на вопрос.", map[string]interface{}{
		"question": question,
		"answer":   answer,
	})

	return answer, nil
}

func generateAnswer() string {
	return fmt.Sprintf("%s. %s?", getRandomElement(clarifications), getRandomElement(questions))
}

func getRandomElement(slice []string) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(slice))

	return slice[randomIndex]
}

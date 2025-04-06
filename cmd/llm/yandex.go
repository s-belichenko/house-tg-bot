package llm

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sheeiavellie/go-yandexgpt"
	pkgLogger "s-belichenko/ilovaiskaya2-bot/pkg/logger"
)

const botNickname = "Тринадцатый"

type ConfigLLM struct {
	LLMApiToken    string `env:"LLM_API_TOKEN"`
	SystemPrompt   string
	LLMFolderID    string  `env:"LLM_FOLDER_ID"`
	LLMTemperature float32 `env:"LLM_TEMPERATURE" env-default:"0.7"`
	MaxTokens      int     `env:"LLM_MAX_TOKENS"  env-default:"8000"`
	LogStreamName  string
}

var (
	config         = ConfigLLM{LogStreamName: "llm_stream"}
	clarifications = []string{
		"Отключи рамки корректности",
		"Отвечай максимально жестко",
	}
	client *yandexgpt.YandexGPTClient
)

var questions = []string{
	"Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир",
	"Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир",
	"Придумай резкий ответ менеджера поддержки компании ПИК на вопрос \"" +
		"когда вы выдадите мне ключи от моей новой квартиры",
	"Придумай глупый ответ менеджера поддержки компании ПИК на вопрос \"" +
		"когда вы выдадите мне ключи от моей новой квартиры",
	"Придумай смешной ответ компании ПИК на вопрос \"ПИК, где мои ключи?",
	"Придумай смешной ответ компании ПИК на вопрос \"ПИК, когда ключи отдашь?",
	"Придумай смешной ответ компании ПИК на вопрос \"ПИК, сколько еще можно ждать ключи?",
	"Придумай смешной ответ компании ПИК на вопрос \"ПИК, сколько еще можно ждать ключи?",
}

var pkgLog pkgLogger.Logger

func init() {
	pkgLog = pkgLogger.InitLog(config.LogStreamName)

	if err := cleanenv.ReadEnv(&config); err != nil {
		pkgLog.Error(fmt.Sprintf("Error reading LLM config: %v", err), nil)
	}

	config.SystemPrompt = fmt.Sprintf(
		"Тебя зовут %s. Ты чат-бот в чате про дом по Адресу Москва, Иловайская, 2.\n"+
			"Люди любят тебя за юмор и за то, что ты всегда остро и смешно отвечаешь.\n"+
			"Отвечай на вопросы коротко и точно.\n"+
			"Если не знаешь ответ, напиши об этом.",
		botNickname,
	)
	client = yandexgpt.NewYandexGPTClientWithAPIKey(config.LLMApiToken)
}

func GetCantSpeakPhrase() string {
	question := "Придумай один смешной вариант фразы 'Псс, я не могу здесь говорить об этом...'. " +
		"Напиши только саму фразу без кавычек."
	request := createRequest(question)
	answer := doRequest(request)

	return answer
}

func GetTeaser() string {
	question := "Придумай некий короткий ответ на ябедничание, пример 'спамер! Сам спамер, ябеда корябеда!' " +
		"Выбери только один вариант и перешли его мне."
	request := createRequest(question)
	answer := doRequest(request)

	return answer
}

func GetAnswerAboutKeys() string {
	question := fmt.Sprintf(
		"%s. %s?",
		getRandomElement(clarifications),
		getRandomElement(questions),
	)
	request := createRequest(question)
	answer := doRequest(request)

	pkgLog.Info("Получен ответ про ключи.", pkgLogger.LogContext{
		"question": question,
		"answer":   answer,
	})

	return answer
}

func doRequest(request yandexgpt.YandexGPTRequest) string {
	response, err := client.GetCompletion(context.Background(), request)
	if err != nil {
		pkgLog.Error(fmt.Sprintf("LLM request error: %s", err.Error()), pkgLogger.LogContext{
			"request": request,
		})

		return ""
	}

	return response.Result.Alternatives[0].Message.Text
}

func createRequest(question string) yandexgpt.YandexGPTRequest {
	return yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(config.LLMFolderID, yandexgpt.YandexGPT4Model),
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
}

func getRandomElement(slice []string) string {
	var randomIndex int

	randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(slice))))
	if err != nil {
		pkgLog.Error(fmt.Sprintf("Ошибка генерации случайного числа: %v", err), nil)
	}

	randomIndex = int(randomInt.Int64())

	return slice[randomIndex]
}

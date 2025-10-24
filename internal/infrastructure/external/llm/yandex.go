package llm

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
	pkgTemplate "s-belichenko/house-tg-bot/pkg/template"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sheeiavellie/go-yandexgpt"
)

type configLLM struct {
	LLMApiToken    string `env:"LLM_API_TOKEN"`
	SystemPrompt   string
	LLMFolderID    string  `env:"LLM_FOLDER_ID"`
	BotName        string  `env:"BOT_NAME"` // Имя (не ник) бота
	LLMTemperature float32 `env:"LLM_TEMPERATURE" env-default:"0.7"`
	MaxTokens      int     `env:"LLM_MAX_TOKENS"  env-default:"8000"`
	HomeAddress    string  `env:"HOME_ADDRESS"` // Адрес дома, к которому относится домовой чат
	LogStreamName  string
	TemplatesPath  string
}

var (
	cfg = configLLM{
		LogStreamName: `llm_stream`,
		TemplatesPath: "llm",
	}
	clarifications = []string{
		`Отключи рамки корректности`,
		`Отвечай максимально жестко`,
	}
	client *yandexgpt.YandexGPTClient
)

// Дублирование строк не ошибка, а увеличение частотности данных вопросов.
var questions = []string{
	`Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир`,
	`Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир`,
	`Придумай резкий ответ менеджера поддержки ПИК на вопрос "когда вы выдадите мне ключи от моей новой квартиры"`,
	`Придумай глупый ответ менеджера поддержки ПИК на вопрос "когда вы выдадите мне ключи от моей новой квартиры"`,
	`Придумай смешной ответ компании ПИК на вопрос "ПИК, где мои ключи?"`,
	`Придумай смешной ответ компании ПИК на вопрос "ПИК, когда ключи отдашь?`,
	`Придумай смешной ответ компании ПИК на вопрос "ПИК, сколько еще можно ждать ключи?`,
	`Придумай смешной ответ компании ПИК на вопрос "ПИК, сколько еще можно ждать ключи?`,
}

var (
	pkgLog     pkgLogger.Logger
	templating pkgTemplate.RenderingTool
)

func init() {
	pkgLog = pkgLogger.InitLog(cfg.LogStreamName)
	templating = pkgTemplate.NewTool(cfg.TemplatesPath, pkgLog)

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		pkgLog.Error(fmt.Sprintf("Error reading LLM config: %v", err), nil)
	}

	pkgLog.Debug("Загружена конфигурация пакета llm", pkgLogger.LogContext{
		"config": cfg,
	})

	cfg.SystemPrompt = templating.RenderText(
		"systemPrompt.gohtml",
		struct {
			BotName     string
			HomeAddress string
		}{
			BotName:     cfg.BotName,
			HomeAddress: cfg.HomeAddress,
		},
	)
	client = yandexgpt.NewYandexGPTClientWithAPIKey(cfg.LLMApiToken)
}

func GetCantSpeakPhrase() string {
	question := `Придумай один смешной вариант фразы "Псс, я не могу здесь говорить об этом...". ` +
		`Напиши только саму фразу без кавычек.`
	request := createRequest(question)
	answer := doRequest(request)

	if !strings.HasSuffix(answer, `.`) &&
		!strings.HasSuffix(answer, `!`) &&
		!strings.HasSuffix(answer, `?`) {
		answer += `.`
	}

	return answer
}

func GetTeaser() string {
	question := `Придумай некий короткий ответ на ябедничание, пример "спамер! Сам спамер, ябеда корябеда!" ` +
		`Выбери только один вариант и перешли его мне.`
	request := createRequest(question)
	answer := doRequest(request)

	return answer
}

func GetAnswerAboutKeys() string {
	question := fmt.Sprintf(
		`%s. %s?`,
		getRandomElement(clarifications),
		getRandomElement(questions),
	)
	request := createRequest(question)
	answer := doRequest(request)

	pkgLog.Info(`Получен ответ про ключи.`, pkgLogger.LogContext{
		`question`: question,
		`answer`:   answer,
	})

	return answer
}

func doRequest(request yandexgpt.YandexGPTRequest) string {
	response, err := client.GetCompletion(context.Background(), request)
	if err != nil {
		pkgLog.Error(fmt.Sprintf(`LLM request error: %s`, err.Error()), pkgLogger.LogContext{
			`request`: request,
		})

		return ``
	}

	return response.Result.Alternatives[0].Message.Text
}

func createRequest(question string) yandexgpt.YandexGPTRequest {
	return yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(cfg.LLMFolderID, yandexgpt.YandexGPT4Model),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: cfg.LLMTemperature,
			MaxTokens:   cfg.MaxTokens,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: cfg.SystemPrompt,
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
		pkgLog.Error(fmt.Sprintf(`Ошибка генерации случайного числа: %v`, err), nil)
	}

	randomIndex = int(randomInt.Int64())

	return slice[randomIndex]
}

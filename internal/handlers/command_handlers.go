package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
	"s-belichenko/ilovaiskaya2-bot/cmd/llm"
	yandexLogger "s-belichenko/ilovaiskaya2-bot/internal/logger"
	"strings"
)

type TeleContext interface {
	Chat() *tele.Chat
	Sender() *tele.User
	Message() *tele.Message
	Send(what interface{}, opts ...interface{}) error
}

type Config struct {
	HomeThreadBot int    `env:"HOME_THREAD_BOT"` // Тема в домашней группе, где нет ограничений для бота
	LogStreamName string // Имя потока в YC Logs
}

// Общие переменные пакета
var (
	config = Config{LogStreamName: "main_stream"}
	log    *yandexLogger.Logger
)

func init() {
	initConfig()
	log = yandexLogger.InitLog(config.LogStreamName)
}

func initConfig() {
	err := cleanenv.ReadEnv(&config)
	config.LogStreamName = "main_stream"
	if err != nil {
		fmt.Printf("Error reading Bot config: %v", err)
	}
}

func CommandStartHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("Привет, %s", getUsername(*c.Sender())))
}

func CommandKeysHandler(c tele.Context) error {
	if !isBotHouse(c) {
		cantSpeakPhrase := llm.GetCantSpeakPhrase()
		if "" != cantSpeakPhrase {
			if !strings.HasSuffix(cantSpeakPhrase, ".") &&
				!strings.HasSuffix(cantSpeakPhrase, "!") &&
				!strings.HasSuffix(cantSpeakPhrase, "?") {
				cantSpeakPhrase += "."
			}
			err := c.Reply(cantSpeakPhrase + " Попробуйте использовать команду в теме \"Оффтоп.\"")
			if err != nil {
				log.Error(fmt.Sprintf("Бот не смог рассказать об ограничениях команды /keys: %v", err), nil)
			}
		}
		return nil
	}

	return c.Send(llm.GetAnswerAboutKeys())
}

func CommandTestHandler(c tele.Context) error {
	commands, err := c.Bot().Commands()
	if err != nil {
		log.Error(fmt.Sprintf("Не удалось получить команды: %v", err), nil)
	} else {
		log.Debug("Команды бота", yandexLogger.LogContext{
			"commands": commands,
		})
	}

	// Преобразуем сообщение в JSON
	messageJSON, err := json.Marshal(c.Message())
	if err != nil {
		log.Error(fmt.Sprintf("Ошибка при преобразовании в JSON: %v", err), nil)
		return nil
	}

	// Форматируем JSON с отступами
	prettyJSON, err := json.MarshalIndent(json.RawMessage(messageJSON), "", "  ")
	if err != nil {
		log.Error(fmt.Sprintf("Ошибка при форматировании JSON: %v", err), nil)
		return nil
	}

	userID := c.Sender().ID
	// Отправляем отформатированный JSON
	formattedJSON := fmt.Sprintf("```json\n%s\n```", prettyJSON)
	if err := c.Send(formattedJSON, tele.ModeMarkdownV2); err != nil {
		log.Error(
			fmt.Sprintf("Не удалось отправить тестовый JSON пользователю %d: %v", userID, err),
			nil,
		)
	}

	return nil
}

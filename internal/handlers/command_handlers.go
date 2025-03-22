package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"s-belichenko/ilovaiskaya2-bot/cmd/llm"
	"strings"

	tele "gopkg.in/telebot.v4"
	yandexLogger "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

type TeleContext interface {
	Chat() *tele.Chat
	Sender() *tele.User
	Message() *tele.Message
	Send(what interface{}, opts ...interface{}) error
}

type Config struct {
	HomeThreadBot int `env:"HOME_THREAD_BOT"` // Тема в супергруппе, где нет ограничений для общения с ботом
	LogStreamName string
}

var config Config

var log *yandexLogger.Logger

func init() {
	initConfig()
	initLog()
}

func initLog() {
	log = yandexLogger.InitLog(config.LogStreamName)
}

func initConfig() {
	err := cleanenv.ReadEnv(&config)
	config.LogStreamName = "main_stream"
	if err != nil {
		fmt.Printf("Error reading Bot config: %v", err)
	}
}

func getUsername(u tele.User) string {
	username := ""
	if r := strings.TrimSpace(u.Username); r != "" {
		username = r
	}

	return username
}

func CommandStartHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("Привет, %s", getUsername(*c.Sender())))
}

func CommandKeysHandler(c tele.Context) error {
	if !isBotHouse(c) {
		cantSpeakPhrase := llm.GetCantSpeakPhrase()
		if "" != cantSpeakPhrase {
			err := c.Reply(cantSpeakPhrase)
			if err != nil {
				log.Error(fmt.Sprintf("Бот не смог рассказать об ограничениях команды /keys: %v", err), nil)
			}
		}
		return nil
	}
	answer := llm.GetAnswerAboutKeys()
	if !strings.HasSuffix(answer, ".") {
		answer += "."
	}
	return c.Send(answer + " Попробуйте использовать команду в теме \"Оффтоп.\"")
}

func isBotHouse(c TeleContext) bool {
	//if c.Message().ThreadID == config.HomeThreadBot || c.Message().ReplyTo != nil {
	if c.Message().ThreadID == config.HomeThreadBot {
		return true
	} else {
		err := c.Send("Псс, я не могу здесь говорить об этом...", tele.SendOptions{
			AllowWithoutReply: true,
		})
		if err != nil {
			log.Error(fmt.Sprintf(
				"Бот не смог рассказать об ограничениях команды /keys"),
				map[string]interface{}{
					"error": err.Error(),
				},
			)
		}
		return false
	}
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

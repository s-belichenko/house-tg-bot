package security

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
	yaLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

type TeleID tele.ChatID
type TeleIDList []TeleID
type TeleContext interface {
	Chat() *tele.Chat
	Sender() *tele.User
	Message() *tele.Message
}

type Config struct {
	BotAdminsIDs         TeleIDList `env:"CHAT_ADMINS"`
	AdministrationChatID TeleID     `env:"ADMINISTRATION_CHAT_ID"`
	HouseChatId          TeleID     `env:"HOUSE_CHAT_ID"`   // Домовой чат, управляемый ботом
	HomeThreadBot        int        `env:"HOME_THREAD_BOT"` // Тема в супергруппе, где нет ограничений для бота
	LogStreamName        string
}

var (
	config = Config{LogStreamName: "main_stream"}
	log    *yaLog.Logger
)

func init() {
	initConfig()
	log = yaLog.InitLog(config.LogStreamName)
}

func initConfig() {
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		fmt.Printf("Error reading Bot config: %v", err)
	}
}

// SetValue сеттер для загрузки в конфигурацию типа TeleIDList
func (f *TeleIDList) SetValue(s string) error {
	*f = getAllowedIDs(s)
	return nil
}

// SetValue сеттер для загрузки в конфигурацию типа TeleID
func (f *TeleID) SetValue(s string) error {
	r, err := parseChatID(s)
	if err != nil {
		return nil
	}
	*f = r
	return nil
}

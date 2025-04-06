package security

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"

	tele "gopkg.in/telebot.v4"
	pkgLog "s-belichenko/ilovaiskaya2-bot/pkg/logger"
)

type TeleID tele.ChatID
type TeleIDList []TeleID
type TeleContext interface {
	Chat() *tele.Chat
	Sender() *tele.User
	Message() *tele.Message
}

type Config struct {
	AdministrationChatID TeleID `env:"ADMINISTRATION_CHAT_ID"`
	HouseChatId          TeleID `env:"HOUSE_CHAT_ID"`   // Домовой чат, управляемый ботом
	HomeThreadBot        int    `env:"HOME_THREAD_BOT"` // Тема в супергруппе, где нет ограничений для бота
	LogStreamName        string
}

var (
	config = Config{LogStreamName: "main_stream"}
	log    pkgLog.Logger
)

func init() {
	initConfig()
	log = pkgLog.InitLog(config.LogStreamName)
}

func initConfig() {
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		fmt.Printf("Error reading Bot config: %v", err)
	}
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

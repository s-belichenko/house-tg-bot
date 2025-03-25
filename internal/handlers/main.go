package handlers

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
	intLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

type Config struct {
	HouseChatId          int64  `env:"HOUSE_CHAT_ID"`          // Домовой чат, управляемый ботом
	AdministrationChatID int64  `env:"ADMINISTRATION_CHAT_ID"` // Чат администраторов, куда поступают уведомления и тп
	BotID                int64  // Собственный идентификатор бота
	LogStreamName        string // Имя потока в YC Logs
}

// Общие переменные пакета
var (
	config = Config{LogStreamName: "main_stream"}
	log    intLog.Logger
)

type TeleContext interface {
	Chat() *tele.Chat
	Sender() *tele.User
	Message() *tele.Message
	Send(what interface{}, opts ...interface{}) error
}

func init() {
	initConfig()
	log = intLog.InitLog(config.LogStreamName)
}

func initConfig() {
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		fmt.Printf("Error reading Bot config: %v", err)
	}
}

func SetBotID(botID int64) {
	config.BotID = botID
}

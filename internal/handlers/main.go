package handlers

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
	pkgLogger "s-belichenko/ilovaiskaya2-bot/pkg/logger"
)

type Config struct {
	HouseChatID          int64  `env:"HOUSE_CHAT_ID"`          // Домовой чат, управляемый ботом
	AdministrationChatID int64  `env:"ADMINISTRATION_CHAT_ID"` // Чат администраторов, куда поступают уведомления и тп
	BotID                int64  // Собственный идентификатор бота
	LogStreamName        string // Имя потока в YC Logs
}

// Общие переменные пакета.
var (
	config = Config{LogStreamName: "main_stream"}
	pkgLog pkgLogger.Logger
)

type TeleContext interface {
	tele.Context
}

func init() {
	pkgLog = pkgLogger.InitLog(config.LogStreamName)

	initConfig()
}

func initConfig() {
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		pkgLog.Error(fmt.Sprintf("Error reading Bot config: %v", err), nil)
	}
}

func SetBotID(botID int64) {
	config.BotID = botID
}

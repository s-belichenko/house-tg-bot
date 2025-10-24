package middleware

import (
	"fmt"

	pkgLog "s-belichenko/house-tg-bot/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
)

type (
	TeleID      tele.ChatID
	TeleIDList  []TeleID
	TeleContext interface {
		Chat() *tele.Chat
		Sender() *tele.User
		Message() *tele.Message
	}
)

type config struct {
	AdministrationChatID TeleID `env:"ADMINISTRATION_CHAT_ID"`
	HouseChatID          TeleID `env:"HOUSE_CHAT_ID"`      // Домовой чат, управляемый ботом
	HomeThreadBot        int    `env:"HOME_THREAD_BOT"`    // Тема в супергруппе, где нет ограничений для бота
	HouseIsCompleted     bool   `env:"HOUSE_IS_COMPLETED"` // Признак, что дом уже сдан
	LogStreamName        string
}

var (
	cfg = config{LogStreamName: "main_stream"}
	log pkgLog.Logger
)

func init() {
	log = pkgLog.InitLog(cfg.LogStreamName)

	initConfig()
}

func initConfig() {
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Error(fmt.Sprintf("Error reading Bot config: %v", err), nil)
	}

	log.Debug("Загружена конфигурация пакета middleware", pkgLog.LogContext{
		"config": cfg,
	})
}

// SetValue сеттер для загрузки в конфигурацию типа TeleID.
func (f *TeleID) SetValue(s string) error {
	r, err := parseChatID(s)
	if err != nil {
		return err
	}

	*f = r

	return nil
}

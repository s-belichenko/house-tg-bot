package handlers

import (
	"fmt"
	"net/url"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
	pkgTemplate "s-belichenko/house-tg-bot/pkg/template"

	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
)

type config struct {
	HouseChatID          int64   `env:"HOUSE_CHAT_ID"`          // Домовой чат, управляемый ботом
	AdministrationChatID int64   `env:"ADMINISTRATION_CHAT_ID"` // Чат администраторов, куда поступают уведомления и тп
	RulesURL             url.URL `env:"RULES_URL"`              // Ссылка на правила чата
	OwnerNickname        string  `env:"OWNER_NICKNAME"`         // Никнейм владельца чата
	InviteURL            url.URL `env:"INVITE_URL"`             // Ссылка для вступления в домовой чат
	BotNickname          string  `env:"BOT_NICKNAME"`           // Ник бота
	HomeAddress          string  `env:"HOME_ADDRESS"`           // Адрес дома, к которому относится домовой чат
	VerifyRules          string  `env:"VERIFY_RULES"`           // Правила верификации
	HouseIsCompleted     bool    `env:"HOUSE_IS_COMPLETED"`     // Признак, что дом уже сдан
	BotID                int64   // Собственный идентификатор бота
	LogStreamName        string  // Имя потока в YC Logs
	TemplatesPath        string  // Путь к шаблонам
}

// Общие переменные пакета.
var (
	cfg           = config{LogStreamName: "main_stream", TemplatesPath: "handlers"}
	pkgLog        pkgLogger.Logger
	renderingTool pkgTemplate.RenderingTool
)

type TeleContext interface {
	tele.Context
}

func init() {
	pkgLog = pkgLogger.InitLog(cfg.LogStreamName)
	renderingTool = pkgTemplate.NewTool(cfg.TemplatesPath, pkgLog)

	initConfig()
}

func initConfig() {
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		pkgLog.Error(fmt.Sprintf("Error reading Bot config: %v", err), nil)
	}

	pkgLog.Debug("Загружена конфигурация пакета handlers", pkgLogger.LogContext{
		"config": cfg,
	})
}

func SetBotID(botID int64) {
	cfg.BotID = botID
}

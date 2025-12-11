package config

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"s-belichenko/house-tg-bot/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
)

type TeleID tele.ChatID

type App struct {
	BotToken         string  `env:"TELEGRAM_BOT_TOKEN"`
	AdminChatID      TeleID  `env:"ADMINISTRATION_CHAT_ID"` // Чат администраторов, куда поступают уведомления и тп
	HomeAddress      string  `env:"HOME_ADDRESS"`           // Адрес дома, к которому относится домовой чат
	HouseChatID      TeleID  `env:"HOUSE_CHAT_ID"`          // Домовой чат, управляемый ботом
	HomeThreadBot    int     `env:"HOME_THREAD_BOT"`        // Тема в супергруппе, где нет ограничений для бота
	HouseIsCompleted bool    `env:"HOUSE_IS_COMPLETED"`     // Признак, что дом уже сдан
	RulesURL         url.URL `env:"RULES_URL"`              // Ссылка на правила чата
	OwnerNickname    string  `env:"OWNER_NICKNAME"`         // Никнейм владельца чата
	InviteURL        url.URL `env:"INVITE_URL"`             // Ссылка для вступления в домовой чат
	BotNickname      string  `env:"BOT_NICKNAME"`           // Ник бота
	VerifyRules      string  `env:"VERIFY_RULES"`           // Правила верификации
	BotID            int64   // Собственный идентификатор бота

	LlmYandex struct {
		LLMApiToken    string `env:"LLM_API_TOKEN"`
		SystemPrompt   string
		LLMFolderID    string  `env:"LLM_FOLDER_ID"`
		BotName        string  `env:"BOT_NAME"` // Имя (не ник) бота
		LLMTemperature float32 `env:"LLM_TEMPERATURE" env-default:"0.7"`
		MaxTokens      int     `env:"LLM_MAX_TOKENS"  env-default:"8000"`
	}
}

func LoadConfig(loggerP logger.Logger) App {
	var cfg App
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		loggerP.Fatal("Ошибка загрузки конфигурации: %v", logger.LogContext{
			"error": err,
		})
		os.Exit(1)
	}

	loggerP.Debug("Загружена конфигурация.", logger.LogContext{
		"config": cfg,
	})

	return cfg
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

func parseChatID(s string) (TeleID, error) {
	idStr := strings.TrimSpace(s)
	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			return TeleID(id), nil
		}

		return 0, err
	}

	return 0, nil
}

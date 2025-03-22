package security

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"slices"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
	yandexLogger "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

type TeleID tele.ChatID
type TeleIDList []TeleID
type TeleContext interface {
	Chat() *tele.Chat
	Sender() *tele.User
}

type Config struct {
	BotAdminsIDs         TeleIDList `env:"CHAT_ADMINS"`
	AdministrationChatID TeleID     `env:"ADMINISTRATION_CHAT_ID"`
	AllowedChats         TeleIDList `env:"ALLOWED_CHATS"`
	LogStreamName        string
}

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
	if err != nil {
		fmt.Printf("Error reading Bot config: %v", err)
	}
}

func (f *TeleIDList) SetValue(s string) error {
	*f = getAllowedIDs(s)
	return nil
}

func (f *TeleID) SetValue(s string) error {
	r, err := parseChatID(s)
	if err != nil {
		return nil
	}
	*f = r
	return nil
}

// getAllowedIDs Получает из текстового списка идентификаторов валидные
func getAllowedIDs(IDs string) TeleIDList {
	var allowedIDs TeleIDList
	allowedIDs = make(TeleIDList, 0)
	if IDs != "" {
		userIDs := strings.Split(IDs, ",")
		for _, idStr := range userIDs {
			if id, err := parseChatID(idStr); err == nil {
				allowedIDs = append(allowedIDs, id)
			} else {
				log.Warn(fmt.Sprintf("Не удалось распознать идентфикатор %s", idStr), nil)
			}
		}
	}

	return allowedIDs
}

func parseChatID(s string) (TeleID, error) {
	idStr := strings.TrimSpace(s)
	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			return TeleID(id), nil
		} else {
			return 0, err
		}
	}

	return 0, nil
}

// isAllowed Проверяем, разрешен ли пользователь или группа
func isAllowed(c TeleContext) (bool, string) {
	var msg string
	r := true

	switch c.Chat().Type {
	case "private", "privatechannel":
		userID := TeleID(c.Sender().ID)

		if !slices.Contains(config.BotAdminsIDs, userID) {
			r = false
			msg = fmt.Sprintf("Извините, у вас нет доступа к этому боту, ваш идентификатор %d", userID)
		}
	case "group", "supergroup":
		chatID := TeleID(c.Chat().ID)

		if !slices.Contains(config.AllowedChats, chatID) && (config.AdministrationChatID != chatID) {
			r = false
			msg = fmt.Sprintf("Извините, бот не предназначен для группы с идентификатором %d", chatID)
		}
	case "channel":
		r = false
		channelID := TeleID(c.Chat().ID)
		msg = fmt.Sprintf("Извините, бот не предназначен для канала с идентификатором %d", channelID)
	}

	return r, msg
}

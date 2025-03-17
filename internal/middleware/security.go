package middleware

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"slices"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

type Config struct {
	ChatAdmins   []tele.ChatID
	AllowedChats []tele.ChatID
}

var config Config

func init() {
	// Читаем список разрешенных пользователей из переменной окружения
	allowedUsersEnv := os.Getenv("CHAT_ADMINS")
	config.ChatAdmins = getAllowedIDs(allowedUsersEnv)
	// Читаем список разрешенных групп из переменной окружения
	allowedChatsEnv := os.Getenv("ALLOWED_CHATS")
	config.AllowedChats = getAllowedIDs(allowedChatsEnv)
}

func SecurityMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if result, msg := isAllowed(c); result != true {
			if err := c.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			// Прерываем дальнейшую обработку
			return nil
		}
		return next(c)
	}
}

// getAllowedIDs Получает из текстового списка идентификаторов валидные
func getAllowedIDs(IDs string) []tele.ChatID {
	var allowedIDs []tele.ChatID
	allowedIDs = make([]tele.ChatID, 0)
	if IDs != "" {
		userIDs := strings.Split(IDs, ",")
		for _, idStr := range userIDs {
			idStr = strings.TrimSpace(idStr)
			if idStr != "" {
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err == nil {
					allowedIDs = append(allowedIDs, tele.ChatID(id))
				} else {
					log.Warn().Msg(fmt.Sprintf("Invalid allowed ID: %s", idStr))
				}
			}
		}
	}

	return allowedIDs
}

// isAllowed Проверяем, разрешен ли пользователь или группа
func isAllowed(c tele.Context) (bool, string) {
	var userID tele.ChatID
	var chatID tele.ChatID
	var msg string
	switch c.Chat().Type {
	case "private", "privatechannel":
		userID = tele.ChatID(c.Sender().ID)

		if slices.Contains(config.ChatAdmins, userID) {
			return true, msg
		} else {
			msg = fmt.Sprintf("Извините, у вас нет доступа к этому боту, ваш идентификатор %d", userID)
		}
	case "group", "supergroup":
		chatID = tele.ChatID(c.Chat().ID)

		if slices.Contains(config.AllowedChats, chatID) {
			return true, msg
		} else {
			msg = fmt.Sprintf("Извините, бот не предназначен для группы с идентификатором %d", chatID)
		}
	case "channel":
		chatID = tele.ChatID(c.Chat().ID)
		msg = fmt.Sprintf("Извините, бот не предназначен для канала с идентификатором %d", chatID)
	}

	return false, msg
}

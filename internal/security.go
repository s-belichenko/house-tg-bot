package internal

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"slices"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

// GetAllowedIDs Получает из текстового списка идентификаторов валидные
func GetAllowedIDs(IDs string) []tele.ChatID {
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

// IsAllowed Проверяем, разрешен ли пользователь или группа
func IsAllowed(c tele.Context, allowedUsers []tele.ChatID, allowedChats []tele.ChatID) bool {
	var userID tele.ChatID
	var chatID tele.ChatID
	if c.Sender() != nil {
		userID = tele.ChatID(c.Sender().ID)
	}
	if c.Chat() != nil {
		chatID = tele.ChatID(c.Chat().ID)
	}

	if slices.Contains(allowedUsers, userID) || slices.Contains(allowedChats, chatID) {
		return true
	}
	return false
}

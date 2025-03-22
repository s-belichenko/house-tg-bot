package handlers

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"
	yandexLogger "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

func isBotHouse(c TeleContext) bool {
	if c.Message().ThreadID == config.HomeThreadBot {
		return true
	} else {
		return false
	}
}

func getUsername(u tele.User) string {
	username := ""
	if r := strings.TrimSpace(u.Username); r != "" {
		username = r
	}

	return username
}

func generateMessageLink(chat *tele.Chat, messageID int) string {
	if chat.Type == tele.ChatChannel || chat.Type == tele.ChatSuperGroup || chat.Type == tele.ChatGroup {
		if chat.Username != "" { // Проверяем, есть ли у чата username
			// если есть username, формируем публичную ссылку
			return fmt.Sprintf("https://t.me/%s/%d", chat.Username, messageID)
		} else { // Если username нет, формируем приватную ссылку
			// удаляем -100 из начала chat.ID
			chatID := chat.ID
			if chatID < 0 {
				chatID = -chatID
			}
			if chatID > 1000000000000 {
				chatID = chatID - 1000000000000
			}
			return fmt.Sprintf("https://t.me/c/%d/%d", chatID, messageID)
		}
	} else {
		log.Error("Невозможно сформировать ссылку для этого типа чата", yandexLogger.LogContext{
			"chat":       chat,
			"message_id": messageID,
		})
		return ""
	}
}

package handlers

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"
	yandexLogger "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

func getUsername(u tele.User) string {
	username := ""
	if r := strings.TrimSpace(u.Username); r != "" {
		username = r
	}

	return username
}

func GenerateMessageLink(chat *tele.Chat, messageID int) string {
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

func setCommands(c tele.Context, commands []tele.Command, scope tele.CommandScope) {
	if err := c.Bot().SetCommands(commands, scope); err != nil {
		log.Fatal(fmt.Sprintf("Не удалось инициализировать команды бота: %v", err), yandexLogger.LogContext{
			"commands": commands,
			"scope":    scope,
		})
	} else {
		log.Info("Успешно установлены команды бота", yandexLogger.LogContext{
			"commands": commands,
			"scope":    scope,
		})
	}
}

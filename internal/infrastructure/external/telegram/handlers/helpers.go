package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

const (
	usernameRegex = `^(?:[a-z_0-9]){5,64}$`
	userIDRegex   = `^[0-9]+$`
)

func GetGreetingName(user *tele.User) string {
	var name string

	if user.Username != "" {
		return "@" + user.Username
	}

	firstname := strings.TrimSpace(user.FirstName)
	if firstname != "" {
		name = firstname
	}

	lastname := strings.TrimSpace(user.LastName)
	if lastname != "" {
		if firstname != "" {
			name += " "
		} else {
			name = ""
		}

		name += lastname
	}

	regExp, err := regexp.Compile(`[a-zA-Z0-9а-яА-Я.\-_—–!@#$%^&*()"'/?><,]+`)
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf(`Ошибка компиляции регулярного выражения для вычисления имени соседа: %v`, err),
			pkgLogger.LogContext{`name`: name},
		)

		return ""
	}

	matchString, err := regexp.MatchString(regExp.String(), name)
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf(`Не удалось применить регулярное выражение для вычисления имени соседа: %v`, err),
			pkgLogger.LogContext{`name`: name},
		)

		return ""
	}

	if matchString {
		return name
	}

	return "сосед"
}

func GenerateMessageLink(chat *tele.Chat, messageID int) string {
	if chat.Type == tele.ChatChannel || chat.Type == tele.ChatSuperGroup ||
		chat.Type == tele.ChatGroup {
		if chat.Username != "" { // Проверяем, есть ли у чата username
			// если есть username, формируем публичную ссылку
			return fmt.Sprintf("https://t.me/%s/%d", chat.Username, messageID)
		}

		// удаляем -100 из начала chat.ID
		chatID := chat.ID
		if chatID < 0 {
			chatID = -chatID
		}

		if chatID > 1000000000000 {
			chatID -= 1000000000000
		}

		return fmt.Sprintf("https://t.me/c/%d/%d", chatID, messageID)
	}

	pkgLog.Error("Невозможно сформировать ссылку для этого типа чата", pkgLogger.LogContext{
		"chat":       chat,
		"message_id": messageID,
	})

	return ""
}

func setCommands(c TeleContext, commands []tele.Command, scope tele.CommandScope) {
	err := c.Bot().SetCommands(commands, scope)
	if err != nil {
		pkgLog.Fatal(
			fmt.Sprintf("Не удалось установить команды бота: %v", err),
			pkgLogger.LogContext{
				"commands": commands,
				"scope":    scope,
			},
		)
	} else {
		pkgLog.Info("Успешно установлены команды бота", pkgLogger.LogContext{
			"commands": commands,
			"scope":    scope,
		})
	}
}

func deleteCommands(c TeleContext, scope tele.CommandScope) {
	err := c.Bot().DeleteCommands(scope)
	if err != nil {
		pkgLog.Fatal(fmt.Sprintf("Не удалось удалить команды бота: %v", err), pkgLogger.LogContext{
			"scope": scope,
		})
	} else {
		pkgLog.Info("Успешно удалены команды бота", pkgLogger.LogContext{
			"scope": scope,
		})
	}
}

func parseUsername(str string) string {
	re := regexp.MustCompile(usernameRegex)
	res := re.FindString(str)

	return res
}

func parseUserID(str string) int64 {
	re := regexp.MustCompile(userIDRegex)
	res := re.FindString(str)
	i, _ := strconv.ParseInt(res, 10, 64)

	return i
}

func parseDays(s string) int64 {
	days, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось распарсить days %q в int64 %v", s, err), nil)

		return 0
	}

	return days
}

func createUserViolator(s string) *tele.User {
	if userID := parseUserID(s); userID > 0 {
		return &tele.User{ID: userID}
	}

	return nil
}

func createUnixTimeFromDays(d string) int64 {
	r := parseDays(d)
	// Дни в секундах плюс один час для просмотра после бана в настройках
	return time.Now().Unix() + (r*86400 + 600)
}

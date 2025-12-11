package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

const (
	usernameRegex = `^(?:[a-z_0-9]){5,64}$`
	userIDRegex   = `^[0-9]+$`
)

func GetGreetingName(user *tele.User) (string, error) {
	const defaultName = "сосед"
	var name string

	if user.Username != "" {
		return "@" + user.Username, nil
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
		return defaultName, fmt.Errorf(
			`ошибка компиляции регулярного выражения для вычисления имени соседа "%s": %w`,
			name,
			err,
		)
	}

	matchString, err := regexp.MatchString(regExp.String(), name)
	if err != nil {
		return defaultName, fmt.Errorf(
			`не удалось применить регулярное выражение для вычисления имени соседа "%s": %w`,
			name,
			err,
		)
	}

	if matchString {
		return name, nil
	}

	return defaultName, nil
}

func GenerateMessageLink(chat *tele.Chat, messageID int) (string, error) {
	if chat.Type == tele.ChatChannel || chat.Type == tele.ChatSuperGroup ||
		chat.Type == tele.ChatGroup {
		if chat.Username != "" { // Проверяем, есть ли у чата username
			// если есть username, формируем публичную ссылку
			return fmt.Sprintf("https://t.me/%s/%d", chat.Username, messageID), nil
		}

		// удаляем -100 из начала chat.ID
		chatID := chat.ID
		if chatID < 0 {
			chatID = -chatID
		}

		if chatID > 1000000000000 {
			chatID -= 1000000000000
		}

		return fmt.Sprintf("https://t.me/c/%d/%d", chatID, messageID), nil
	}

	return "", nil
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

func createUserViolator(s string) *tele.User {
	if userID := parseUserID(s); userID > 0 {
		return &tele.User{ID: userID}
	}

	return nil
}

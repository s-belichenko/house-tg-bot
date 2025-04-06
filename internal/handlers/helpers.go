package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
	pkgLogger "s-belichenko/ilovaiskaya2-bot/pkg/logger"
)

const (
	usernameRegex = `^(?:[a-z_0-9]){5,64}$`
	userIDRegex   = `^[0-9]+$`
)

func GetGreetingName(user *tele.User) string {
	name := "сосед"

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

	return name
}

func GenerateMessageLink(chat *tele.Chat, messageID int) string {
	if chat.Type == tele.ChatChannel || chat.Type == tele.ChatSuperGroup ||
		chat.Type == tele.ChatGroup {
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
				chatID -= 1000000000000
			}

			return fmt.Sprintf("https://t.me/c/%d/%d", chatID, messageID)
		}
	} else {
		pkgLog.Error("Невозможно сформировать ссылку для этого типа чата", pkgLogger.LogContext{
			"chat":       chat,
			"message_id": messageID,
		})

		return ""
	}
}

func setCommands(c TeleContext, commands []tele.Command, scope tele.CommandScope) {
	if err := c.Bot().SetCommands(commands, scope); err != nil {
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
	if err := c.Bot().DeleteCommands(scope); err != nil {
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
	if days, err := strconv.ParseInt(s, 10, 64); err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось распарсить days %q в int64 %v", s, err), nil)

		return 0
	} else {
		return days
	}
}

func createUserViolator(s string) *tele.User {
	if userID := parseUserID(s); userID > 0 {
		return &tele.User{ID: userID}
	} else {
		// FIXME: Если не найдется способа получать user_id по username, удалить закомментированный код ниже.
		// username := parseUsername(s)
		// if username != "" {
		//	if chat, err := c.Bot().ChatByUsername(username); err != nil {
		//		pkgLog.Error(fmt.Sprintf("Не удалось получить чат для блокировки пользователя: %v", err), pkgLog.LogContext{
		//			"username": username,
		//		})
		//		return nil
		//	} else if chat != nil {
		//		return &tele.User{ID: chat.ID}
		//	} else {
		//		pkgLog.Warn(fmt.Sprintf("В команду /ban передан невалидный username или user_id"), pkgLog.LogContext{
		//			"username_or_user_id": s,
		//		})
		//		if err := c.Reply(fmt.Sprintf(
		//			"Не удалось распознать username нарушителя. Верный формат команды: %s",
		//			muteCommandFormat,
		//		), tele.ModeHTML); err != nil {
		//			pkgLog.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /ban: %v", err), pkgLog.LogContext{
		//				"message": c.Message(),
		//			})
		//		}
		//	}
		// } else {
		//	pkgLog.Error(fmt.Sprintf("Для администрирования не удалось определить user_id и username"), pkgLog.LogContext{
		//		"text": s,
		//	})
		//}
	}

	return nil
}

func createUnixTimeFromDays(d string) int64 {
	r := parseDays(d)
	// Дни в секундах плюс один час для просмотра после бана в настройках
	return time.Now().Unix() + (r*86400 + 600)
}

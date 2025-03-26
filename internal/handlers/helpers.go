package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
	intLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

const usernameRegex = `^(?:[a-z_0-9]){5,64}$`
const userIDRegex = `^[0-9]+$`

func GetGreetingName(u *tele.User) string {
	name := "сосед"
	if u.Username != "" {
		return "@" + u.Username
	}

	f := strings.TrimSpace(u.FirstName)
	l := strings.TrimSpace(u.LastName)
	if f != "" {
		name = f
	}
	if l != "" {
		if f != "" {
			name += " "
		} else {
			name = ""
		}
		name += l
	}

	return name
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
		log.Error("Невозможно сформировать ссылку для этого типа чата", intLog.LogContext{
			"chat":       chat,
			"message_id": messageID,
		})
		return ""
	}
}

func setCommands(c tele.Context, commands []tele.Command, scope tele.CommandScope) {
	if err := c.Bot().SetCommands(commands, scope); err != nil {
		log.Fatal(fmt.Sprintf("Не удалось установить команды бота: %v", err), intLog.LogContext{
			"commands": commands,
			"scope":    scope,
		})
	} else {
		log.Info("Успешно установлены команды бота", intLog.LogContext{
			"commands": commands,
			"scope":    scope,
		})
	}
}

func deleteCommands(c tele.Context, scope tele.CommandScope) {
	if err := c.Bot().DeleteCommands(scope); err != nil {
		log.Fatal(fmt.Sprintf("Не удалось удалить команды бота: %v", err), intLog.LogContext{
			"scope": scope,
		})
	} else {
		log.Info("Успешно удалены команды бота", intLog.LogContext{
			"scope": scope,
		})
	}
}

func parseUsername(s string) string {
	var res string
	if re, err := regexp.Compile(usernameRegex); err != nil {
		log.Error(fmt.Sprintf("Не удалось распарсить username %q: %v", s, err), nil)
		return ""
	} else {
		res = re.FindString(s)
	}
	return res
}

func parseUserID(s string) int64 {
	var res string
	if re, err := regexp.Compile(userIDRegex); err != nil {
		return 0
	} else {
		res = re.FindString(s)
	}
	i, _ := strconv.ParseInt(res, 10, 64)
	return i
}

func parseDays(s string) int64 {
	if days, err := strconv.ParseInt(s, 10, 64); err != nil {
		log.Error(fmt.Sprintf("Не удалось распарсить days %q в int64 %v", s, err), nil)
		return 0
	} else {
		return days
	}
}

func createUserViolator(c tele.Context, s string) *tele.User {
	if userID := parseUserID(s); userID > 0 {
		return &tele.User{ID: userID}
	} else {
		// FIXME: Если не найдется способа получать user_id по username, удалить закомментированный код ниже.
		//username := parseUsername(s)
		//if username != "" {
		//	if chat, err := c.Bot().ChatByUsername(username); err != nil {
		//		log.Error(fmt.Sprintf("Не удалось получить чат для блокировки пользователя: %v", err), intLog.LogContext{
		//			"username": username,
		//		})
		//		return nil
		//	} else if chat != nil {
		//		return &tele.User{ID: chat.ID}
		//	} else {
		//		log.Warn(fmt.Sprintf("В команду /ban передан невалидный username или user_id"), intLog.LogContext{
		//			"username_or_user_id": s,
		//		})
		//		if err := c.Reply(fmt.Sprintf(
		//			"Не удалось распознать username нарушителя. Верный формат команды: %s",
		//			muteCommandFormat,
		//		), tele.ModeHTML); err != nil {
		//			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /ban: %v", err), intLog.LogContext{
		//				"message": c.Message(),
		//			})
		//		}
		//	}
		//} else {
		//	log.Error(fmt.Sprintf("Для администрирования не удалось определить user_id и username"), intLog.LogContext{
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

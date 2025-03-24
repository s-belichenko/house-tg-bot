package handlers

import (
	"fmt"
	tele "gopkg.in/telebot.v4"
	yaLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

const hi = "Привет, вы подали заявку на вступление в чат дома по адресу Иловайская, 2 (бывшие 13-е корпуса) в ЖК Люблинский парк. Ожидайте, скоро с вами свяжутся."

func JoinRequestHandler(c tele.Context) error {
	log.Info("Получена заявка на вступление в чат", yaLog.LogContext{
		"user_id":   c.Sender().ID,
		"username":  c.Sender().Username,
		"firstname": c.Sender().FirstName,
		"lastname":  c.Sender().LastName,
	})

	if _, err := c.Bot().Send(c.Sender(), hi); err != nil {
		log.Error(fmt.Sprintf("Не удалось ответить на заявку: %v", err), yaLog.LogContext{
			"user_id":   c.Sender().ID,
			"username":  c.Sender().Username,
			"firstname": c.Sender().FirstName,
			"lastname":  c.Sender().LastName,
		})
		return err
	}
	adminChat := &tele.Chat{ID: config.AdministrationChatID}
	requestMsg := fmt.Sprintf(`
#JOIN_REQUEST
Новая заявка на вступление в чат.

user_id: %d
username: @%s
firstname: %s
lastname: %s
`, c.Sender().ID, c.Sender().Username, c.Sender().FirstName, c.Sender().LastName)

	if _, err := c.Bot().Send(adminChat, requestMsg); err != nil {
		log.Error(fmt.Sprintf("Не удалось ответить на заявку на вступление: %v", err), yaLog.LogContext{
			"user_id":   c.Sender().ID,
			"username":  c.Sender().Username,
			"firstname": c.Sender().FirstName,
			"lastname":  c.Sender().LastName,
		})
	}

	return nil
}

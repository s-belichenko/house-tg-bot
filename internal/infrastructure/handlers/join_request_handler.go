package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

const hiMsg = "Привет, вы подали заявку на вступление в чат дома по адресу Иловайская, 2 (бывшие 13-е корпуса) " +
	"в ЖК Люблинский парк. Ожидайте, скоро с вами свяжутся."

func JoinRequestHandler(ctx tele.Context) error {
	pkgLog.Info("Получена заявка на вступление в чат", pkgLogger.LogContext{
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	// FIXME: Не отправляется тем, кто не начал общение с ботом, то есть всем. Подсмотреть алгоритм в других домовых чатах.
	if _, err := ctx.Bot().Send(ctx.Sender(), hiMsg); err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось ответить на заявку: %v", err), pkgLogger.LogContext{
			"user_id":   ctx.Sender().ID,
			"username":  ctx.Sender().Username,
			"firstname": ctx.Sender().FirstName,
			"lastname":  ctx.Sender().LastName,
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
`, ctx.Sender().ID, ctx.Sender().Username, ctx.Sender().FirstName, ctx.Sender().LastName)

	if _, err := ctx.Bot().Send(adminChat, requestMsg); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось ответить на заявку на вступление: %v", err),
			pkgLogger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}

	return nil
}

package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

func JoinRequestHandler(ctx tele.Context) error {
	pkgLog.Info("Получена заявка на вступление в чат", pkgLogger.LogContext{
		"chat_id":   ctx.Chat().ID,
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	if _, err := ctx.Bot().Send(
		ctx.Sender(),
		templating.RenderText(`hi.txt`, config),
		tele.ModeHTML, tele.NoPreview,
	); err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось отправить правила вступления: %v", err), pkgLogger.LogContext{
			"user_id":   ctx.Sender().ID,
			"username":  ctx.Sender().Username,
			"firstname": ctx.Sender().FirstName,
			"lastname":  ctx.Sender().LastName,
		})
	}

	tmplInfo := `#JOIN_REQUEST
Новая заявка на вступление в чат.

chat_title: <a href="%s">%s</>
user_id: %d
username: @%s
firstname: %s
lastname: %s
`

	adminChat := &tele.Chat{ID: config.AdministrationChatID}
	requestMsg := fmt.Sprintf(
		tmplInfo,
		config.InviteURL.String(),
		ctx.Chat().Title,
		ctx.Sender().ID,
		ctx.Sender().Username,
		ctx.Sender().FirstName,
		ctx.Sender().LastName,
	)

	if _, err := ctx.Bot().Send(adminChat, requestMsg, tele.ModeHTML, tele.NoPreview); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось оповестить администраторов о заявке на вступление: %v", err),
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

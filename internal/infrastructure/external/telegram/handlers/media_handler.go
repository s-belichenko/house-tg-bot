package handlers

import (
	"fmt"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"

	tele "gopkg.in/telebot.v4"
)

func MediaHandler(ctx tele.Context) error {
	pkgLog.Info("Получено медиа в переписке с ботом", pkgLogger.LogContext{
		"chat_id":   ctx.Chat().ID,
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	var (
		menuInline = &tele.ReplyMarkup{
			ResizeKeyboard: true,
			Placeholder:    "Inline placeholder",
		}
		btnContactAdmin = menuInline.URL(
			"Написать администратору",
			"https://t.me/"+cfg.OwnerNickname,
		)
	)

	menuInline.Inline(menuInline.Row(btnContactAdmin))

	err := ctx.Reply(
		`Я не обрабатываю медиафайлы. Если вы хотели пройти верификацию, прочитайте правила вступления чуть внимательнее.`,
		menuInline,
	)
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось отправить сообщение в ответ на медиа: %v", err),
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

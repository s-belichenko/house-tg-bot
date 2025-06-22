package handlers

import (
	"fmt"
	"html/template"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"

	tele "gopkg.in/telebot.v4"
)

func JoinRequestHandler(ctx tele.Context) error {
	pkgLog.Info("Получена заявка на вступление в чат", pkgLogger.LogContext{
		"chat_id":   ctx.Chat().ID,
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	notifyAdmins(ctx)
	sendHi(ctx)

	return nil
}

func notifyAdmins(ctx tele.Context) {
	if _, err := ctx.Bot().Send(
		&tele.Chat{ID: config.AdministrationChatID},
		renderingTool.RenderText(`join_request.txt`, struct {
			ChatURL   template.URL
			ChatName  string
			UserID    int64
			Username  string
			Firstname string
			Lastname  string
		}{
			ChatURL:   template.URL(config.InviteURL.String()),
			ChatName:  ctx.Chat().Title,
			UserID:    ctx.Sender().ID,
			Username:  ctx.Sender().Username,
			Firstname: ctx.Sender().FirstName,
			Lastname:  ctx.Sender().LastName,
		}),
		tele.ModeHTML,
		tele.NoPreview,
	); err != nil {
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
}

func sendHi(ctx tele.Context) {
	var (
		menuInline = &tele.ReplyMarkup{
			ResizeKeyboard: true,
			Placeholder:    "Inline placeholder",
		}
		btnContactAdmin = menuInline.URL(
			"Написать администратору",
			"https://t.me/"+config.OwnerNickname,
		)
	)

	menuInline.Inline(menuInline.Row(btnContactAdmin))

	if _, err := ctx.Bot().Send(
		ctx.Sender(),
		renderingTool.RenderEscapedText(
			`hi.txt`,
			struct {
				InviteURL   template.URL
				HomeAddress template.HTML
				VerifyRules template.HTML
			}{
				InviteURL:   template.URL(config.InviteURL.String()),
				HomeAddress: template.HTML(config.HomeAddress),
				VerifyRules: template.HTML(config.VerifyRules),
			},
			[]string{"VerifyRules"},
		),
		menuInline,
		tele.ModeHTML, tele.NoPreview,
	); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось отправить правила вступления: %v", err),
			pkgLogger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}
}

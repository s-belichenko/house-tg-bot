package handlers

import (
	"fmt"
	"html/template"

	"s-belichenko/house-tg-bot/internal/config"

	template2 "s-belichenko/house-tg-bot/pkg/template"

	"s-belichenko/house-tg-bot/pkg/logger"

	tele "gopkg.in/telebot.v4"
)

type joinRequestHandlers struct {
	renderingTool template2.RenderingTool
	config        config.App
	logger        logger.Logger
}

type JoinRequestHandlers interface {
	JoinRequestHandler(ctx tele.Context) error
}

func NewJoinRequestHandlersHandlers(cfg config.App, logger logger.Logger) JoinRequestHandlers {
	renderingTool := template2.NewTool("handlers", logger)

	return &joinRequestHandlers{
		renderingTool: renderingTool,
		config:        cfg,
		logger:        logger,
	}
}

func (h *joinRequestHandlers) JoinRequestHandler(ctx tele.Context) error {
	h.logger.Info("Получена заявка на вступление в чат", logger.LogContext{
		"chat_id":   ctx.Chat().ID,
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	h.notifyAdmins(ctx)
	h.sendHi(ctx)

	return nil
}

func (h *joinRequestHandlers) notifyAdmins(ctx tele.Context) {
	if _, err := ctx.Bot().Send(
		&tele.Chat{ID: int64(h.config.AdminChatID)},
		h.renderingTool.RenderText(`join_request.gohtml`, struct {
			ChatURL   template.URL
			ChatName  string
			UserID    int64
			Username  string
			Firstname string
			Lastname  string
		}{
			ChatURL:   template.URL(h.config.InviteURL.String()),
			ChatName:  ctx.Chat().Title,
			UserID:    ctx.Sender().ID,
			Username:  ctx.Sender().Username,
			Firstname: ctx.Sender().FirstName,
			Lastname:  ctx.Sender().LastName,
		}),
		tele.ModeHTML,
		tele.NoPreview,
	); err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось оповестить администраторов о заявке на вступление: %v", err),
			logger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}
}

func (h *joinRequestHandlers) sendHi(ctx tele.Context) {
	var (
		menuInline = &tele.ReplyMarkup{
			ResizeKeyboard: true,
			Placeholder:    "Inline placeholder",
		}
		btnContactAdmin = menuInline.URL(
			"Написать администратору",
			"https://t.me/"+h.config.OwnerNickname,
		)
	)

	menuInline.Inline(menuInline.Row(btnContactAdmin))

	if _, err := ctx.Bot().Send(
		ctx.Sender(),
		h.renderingTool.RenderEscapedText(
			`hi.gohtml`,
			struct {
				InviteURL   template.URL
				HomeAddress template.HTML
				VerifyRules template.HTML
			}{
				InviteURL:   template.URL(h.config.InviteURL.String()),
				HomeAddress: template.HTML(h.config.HomeAddress),
				VerifyRules: template.HTML(h.config.VerifyRules),
			},
			[]string{"VerifyRules"},
		),
		menuInline,
		tele.ModeHTML, tele.NoPreview,
	); err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось отправить правила вступления: %v", err),
			logger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}
}

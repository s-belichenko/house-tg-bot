package handlers

import (
	"fmt"
	"html/template"

	"s-belichenko/house-tg-bot/internal/config"

	template2 "s-belichenko/house-tg-bot/pkg/template"

	"s-belichenko/house-tg-bot/pkg/logger"

	tele "gopkg.in/telebot.v4"
)

type JoinRequestHandlers struct {
	renderingTool template2.RenderingTool
	config        config.App
	logger        logger.Logger
}

func NewJoinRequestHandlersHandlers(cfg config.App, logger logger.Logger) *JoinRequestHandlers {
	renderingTool := template2.NewTool("handlers", logger)

	return &JoinRequestHandlers{
		renderingTool: renderingTool,
		config:        cfg,
		logger:        logger,
	}
}

func (h *JoinRequestHandlers) JoinRequestHandler(ctx tele.Context) error {
	h.logger.Info("Получена заявка на вступление в чат", logger.LogContext{
		"chat_id":   ctx.Chat().ID,
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	h.notifyAdminsAboutJoinRequest(ctx)
	h.sendJoinRules(ctx)

	return nil
}

// UserLeftHandler TODO: Точно ли не работает именно в закрытых группах?
func (h *JoinRequestHandlers) UserLeftHandler(ctx tele.Context) error {
	h.logger.Info("Пользователь покинул чат", logger.LogContext{
		"chat_id":   ctx.Chat().ID,
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	h.sendYouLeftMessage(ctx)
	h.notifyAdminsAboutUserLeft(ctx)

	return nil
}

func (h *JoinRequestHandlers) UserJoinedHandler(ctx tele.Context) error {
	h.logger.Info("Пользователь успешно добавлен в чат", logger.LogContext{
		"chat_id":   ctx.Chat().ID,
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	h.sendHiMessage(ctx)
	h.notifyAdminsAboutUserJoined(ctx)

	return nil
}

func (h *JoinRequestHandlers) notifyAdminsAboutJoinRequest(ctx tele.Context) {
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

func (h *JoinRequestHandlers) sendJoinRules(ctx tele.Context) {
	var (
		menuInline = &tele.ReplyMarkup{
			ResizeKeyboard: true,
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
			`join_rules.gohtml`,
			struct {
				InviteURL   template.URL
				HomeAddress template.HTML
				JoinRules   template.HTML
			}{
				InviteURL:   template.URL(h.config.InviteURL.String()),
				HomeAddress: template.HTML(h.config.HomeAddress),
				JoinRules:   template.HTML(h.config.JoinRules),
			},
			[]string{"JoinRules"},
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

func (h *JoinRequestHandlers) sendHiMessage(ctx tele.Context) {
	if _, err := ctx.Bot().Send(
		ctx.Sender(),
		h.renderingTool.RenderEscapedText(
			`hi.gohtml`,
			struct {
				InviteURL   template.URL
				HomeAddress template.HTML
				HiMessage   template.HTML
				ChatSiteURL template.URL
			}{
				InviteURL:   template.URL(h.config.InviteURL.String()),
				HomeAddress: template.HTML(h.config.HomeAddress),
				HiMessage:   template.HTML(h.config.HiMessage),
				ChatSiteURL: template.URL(h.config.ChatSiteURL.String()),
			},
			[]string{"HiMessage"},
		),
		tele.ModeHTML,
		tele.NoPreview,
	); err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось отправить приветственное сообщение: %v", err),
			logger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}
}

func (h *JoinRequestHandlers) notifyAdminsAboutUserJoined(ctx tele.Context) {
	if _, err := ctx.Bot().Send(
		&tele.Chat{ID: int64(h.config.AdminChatID)},
		h.renderingTool.RenderText(`user_joined.gohtml`, struct {
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
			fmt.Sprintf("Не удалось оповестить администраторов о вступлении в чат: %v", err),
			logger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}
}

func (h *JoinRequestHandlers) sendYouLeftMessage(ctx tele.Context) {
	if _, err := ctx.Bot().Send(
		ctx.Sender(),
		h.renderingTool.RenderText(
			`you_left.gohtml`,
			struct {
				InviteURL template.URL
			}{
				InviteURL: template.URL(h.config.InviteURL.String()),
			},
		),
		tele.ModeHTML,
		tele.NoPreview,
	); err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось отправить сообщение что пользователь покинул чат: %v", err),
			logger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}
}

func (h *JoinRequestHandlers) notifyAdminsAboutUserLeft(ctx tele.Context) {
	if _, err := ctx.Bot().Send(
		&tele.Chat{ID: int64(h.config.AdminChatID)},
		h.renderingTool.RenderText(`user_left.gohtml`, struct {
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
			fmt.Sprintf("Не удалось оповестить администраторов о покидании пользователем чата: %v", err),
			logger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}
}

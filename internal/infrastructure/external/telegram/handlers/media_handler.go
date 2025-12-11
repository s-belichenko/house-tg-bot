package handlers

import (
	"fmt"

	"s-belichenko/house-tg-bot/internal/config"

	"s-belichenko/house-tg-bot/pkg/logger"

	tele "gopkg.in/telebot.v4"
)

type commandMediaHandlers struct {
	config config.App
	logger logger.Logger
}

type CommandMediaHandlers interface {
	MediaHandler(ctx tele.Context) error
}

func NewCommandMediaHandlers(cfg config.App, logger logger.Logger) CommandMediaHandlers {
	return &commandMediaHandlers{
		config: cfg,
		logger: logger,
	}
}

func (h *commandMediaHandlers) MediaHandler(ctx tele.Context) error {
	h.logger.Info("Получено медиа в переписке с ботом", logger.LogContext{
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
			"https://t.me/"+h.config.OwnerNickname,
		)
	)

	menuInline.Inline(menuInline.Row(btnContactAdmin))

	err := ctx.Reply(
		`Я не обрабатываю медиафайлы. Если вы хотели пройти верификацию, прочитайте правила вступления чуть внимательнее.`,
		menuInline,
	)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось отправить сообщение в ответ на медиа: %v", err),
			logger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}

	return nil
}

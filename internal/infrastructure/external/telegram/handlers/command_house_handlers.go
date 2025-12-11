package handlers

import (
	"fmt"
	"html/template"

	"s-belichenko/house-tg-bot/internal/config"

	template2 "s-belichenko/house-tg-bot/pkg/template"

	"s-belichenko/house-tg-bot/internal/domain/models"

	tele "gopkg.in/telebot.v4"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

type CommandHouseHandlers struct {
	config        config.App
	ai            models.AI
	renderingTool template2.RenderingTool
	logger        pkgLogger.Logger
}

// Команды бота для домового чата.
var (
	// KeysCommand TODO: Добавить возможность изменять название застройщика.
	KeysCommand   = tele.Command{Text: "keys", Description: "ПИК, где ключи?"}
	ReportCommand = tele.Command{
		Text:        "report",
		Description: "Сообщить о нарушении правил, формат: /report [уточнение]",
	}
	RulesCommand = tele.Command{
		Text:        "rules",
		Description: "Посмотреть правила чата",
	}
)

func NewCommandHouseHandlers(cfg config.App, ai models.AI, logger pkgLogger.Logger) *CommandHouseHandlers {
	renderingTool := template2.NewTool("handlers", logger)

	return &CommandHouseHandlers{
		config:        cfg,
		ai:            ai,
		renderingTool: renderingTool,
		logger:        logger,
	}
}

func (h *CommandHouseHandlers) CommandKeysHandler(c tele.Context) error {
	if !h.config.HouseIsCompleted {
		err := c.Send(h.ai.GetAnswerAboutKeys())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *CommandHouseHandlers) CommandReportHandler(ctx tele.Context) error {
	msg := ctx.Message()
	reporter := msg.Sender

	var (
		violatorID       int64
		violatorUsername = "username неизвестен"
		violationText    = "Сообщение не содержало текста."
		clarification    = "Не оставлено."
	)

	if h.incorrectUseReportCommand(ctx, msg) {
		h.cleanUpReport(ctx, msg, reporter)

		return nil
	}

	if msg.ReplyTo.Text != "" {
		violationText = msg.ReplyTo.Text
	}

	if ctx.Data() != "" {
		clarification = ctx.Data()
	}

	violator := msg.ReplyTo.Sender
	violatorID = violator.ID

	if violator.Username != "" {
		violatorUsername = "@" + violator.Username
	}

	violationMessageID := msg.ReplyTo.ID
	messageLink, err := GenerateMessageLink(msg.ReplyTo.Chat, violationMessageID)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Не удалось получить ссылку на сообщение с нарушением: %v", err), nil)
	}

	greetingName, err := GetGreetingName(reporter)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Не удалось сформировать обращение к пользователю %d: %v", reporter.ID, err), nil)
	}
	h.logger.Info(
		fmt.Sprintf("Получен отчет о новом нарушении правил от %s", greetingName),
		pkgLogger.LogContext{
			"reporter_username": reporter.Username,
			"reporter_id":       reporter.ID,
			"violator":          violatorUsername,
			"violator_id":       violatorID,
			"violation_text":    violationText,
			"clarification":     clarification,
			"message_link":      messageLink,
		},
	)

	if h.reportAboutBot(ctx, violatorID, reporter) {
		h.cleanUpReport(ctx, msg, reporter)

		return nil
	}

	h.sendNotification(ctx, violationText, violator, reporter, clarification, messageLink)
	h.cleanUpReport(ctx, msg, reporter)
	h.thxForReport(ctx, violationText, clarification, reporter)

	return nil
}

func (h *CommandHouseHandlers) CommandRulesHandler(ctx tele.Context) error {
	var (
		targetUser    *tele.User
		targetMessage *tele.Message
	)

	if ctx.Message().ReplyTo == nil {
		targetMessage = ctx.Message()
		targetUser = ctx.Message().Sender
	} else {
		targetMessage = ctx.Message().ReplyTo
		targetUser = ctx.Message().ReplyTo.Sender
	}

	h.logger.Info(
		`Получен запрос правил`,
		pkgLogger.LogContext{
			"message": ctx.Message(),
		},
	)

	greetingName, err := GetGreetingName(targetUser)
	if err != nil {
		h.logger.Warn(
			fmt.Sprintf("Не удалось сформировать обращение к пользователю %d: %v", targetUser.ID, err),
			nil,
		)
	}
	if targetMessage == nil {
		err := ctx.Send(
			fmt.Sprintf(
				`Привет, %s! Вот <a href="%s">правила чата</a>, ознакомься.`,
				greetingName,
				h.config.RulesURL.String(),
			),
			tele.ModeHTML,
			tele.NoPreview,
		)
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось отправить правила чата по команде /rules: %v", err),
				pkgLogger.LogContext{
					"ctx": ctx,
				},
			)
		}
	} else {
		if _, err := ctx.Bot().Reply(
			targetMessage,
			fmt.Sprintf(
				`Привет, %s! Вот <a href="%s">правила чата</a>, ознакомься.`,
				greetingName,
				h.config.RulesURL.String(),
			),
			tele.ModeHTML,
			tele.NoPreview,
		); err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось отправить правила чата по команде /rules: %v", err),
				pkgLogger.LogContext{
					"ctx": ctx,
				},
			)
		}
	}

	err = ctx.Bot().Delete(ctx.Message())
	if err != nil {
		h.logger.Error(
			fmt.Sprintf(
				"Не удалось удалить сообщение с командой /rules от %d: %v",
				ctx.Message().Sender.ID,
				err,
			),
			pkgLogger.LogContext{
				"message_id":   ctx.Message().ID,
				"message_text": ctx.Message().Text,
			},
		)
	}

	return nil
}

func (h *CommandHouseHandlers) thxForReport(
	ctx tele.Context,
	violationText string,
	clarification string,
	reporter *tele.User,
) {
	if _, err := ctx.Bot().Send(
		reporter,
		h.renderingTool.RenderText(`report_thx.gohtml`, struct {
			Text          string
			Clarification string
		}{
			Text:          violationText,
			Clarification: clarification,
		}),
		tele.ModeHTML,
		tele.NoPreview,
	); err != nil {
		h.logger.Error(fmt.Sprintf(
			"Не удалось послать благодарность за жалобу %d: %v",
			reporter.ID,
			err,
		), nil)
	}
}

func (h *CommandHouseHandlers) cleanUpReport(ctx tele.Context, msg *tele.Message, reporter *tele.User) {
	err := ctx.Bot().Delete(msg)
	if err != nil {
		h.logger.Error(fmt.Sprintf(
			"Не удалось удалить сообщение с жалобой от %d: %v", reporter.ID, err),
			pkgLogger.LogContext{
				"message_id":   msg.ID,
				"message_text": msg.Text,
			},
		)
	}
}

func (h *CommandHouseHandlers) reportAboutBot(ctx tele.Context, violatorID int64, reporter *tele.User) bool {
	if violatorID == h.config.BotID {
		greetingName, err := GetGreetingName(reporter)
		if err != nil {
			h.logger.Warn(fmt.Sprintf("Не удалось сформировать обращение к пользователю %d", reporter.ID), nil)
		}
		err = ctx.Reply(fmt.Sprintf("%s, ай-яй-яй! %s", greetingName, h.ai.GetTeaser()))
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось пообзываться в ответ на репорт на бота: %v", err),
				pkgLogger.LogContext{
					"reporter": reporter,
					"message":  ctx.Message(),
				},
			)
		}

		return true
	}

	return false
}

func (h *CommandHouseHandlers) incorrectUseReportCommand(ctx tele.Context, msg *tele.Message) bool {
	if msg.ReplyTo == nil || msg.ReplyTo.Sender.ID == msg.Sender.ID {
		if _, err := ctx.Bot().Send(
			msg.Sender,
			"Пожалуйста, используйте команду /report в ответе на сообщение с нарушением. "+
				"Подробнее: /help.",
		); err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось отправить уточнение про команду /report: %v", err),
				pkgLogger.LogContext{
					"report_message_object": msg,
				},
			)
		}

		return true
	}

	return false
}

func (h *CommandHouseHandlers) sendNotification(
	ctx tele.Context,
	violationText string,
	violator *tele.User,
	reporter *tele.User,
	clarification string,
	messageLink string,
) {
	var (
		reporterUsername = "username неизвестен"
		violatorUsername = "username неизвестен"
	)

	if reporter.Username != "" {
		reporterUsername = "@" + reporter.Username
	}

	if violator.Username != "" {
		violatorUsername = "@" + violator.Username
	}

	if _, err := ctx.Bot().Send(
		&tele.Chat{ID: int64(h.config.AdminChatID)},
		h.renderingTool.RenderText(`report_notice.gohtml`, struct {
			ChatTitle        string
			ChatURL          template.URL
			ReporterUsername string
			ReporterID       int64
			Clarification    string
			MessageLink      string
			ViolatorUsername string
			ViolatorID       int64
			Text             string
		}{
			ChatTitle:        ctx.Chat().Title,
			ChatURL:          template.URL(h.config.InviteURL.String()),
			ReporterUsername: reporterUsername,
			ReporterID:       reporter.ID,
			Clarification:    clarification,
			MessageLink:      messageLink,
			ViolatorUsername: violatorUsername,
			ViolatorID:       violator.ID,
			Text:             violationText,
		}),
		tele.ModeHTML,
		tele.NoPreview,
	); err != nil {
		h.logger.Error(fmt.Sprintf(
			"Не удалось послать в административный чат жалобу от %d на %d: %v",
			reporter.ID,
			violator.ID,
			err,
		), nil)
	}
}

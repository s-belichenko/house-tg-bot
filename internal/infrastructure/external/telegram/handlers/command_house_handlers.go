package handlers

import (
	"fmt"
	"s-belichenko/house-tg-bot/internal/infrastructure/external/llm"

	tele "gopkg.in/telebot.v4"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

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

func CommandKeysHandler(c tele.Context) error {
	if config.HouseIsCompleted {
		return c.Send(llm.GetAnswerAboutKeys())
	}

	return nil
}

func CommandReportHandler(ctx tele.Context) error {
	msg := ctx.Message()
	reporter := msg.Sender

	var (
		violatorID       int64
		violatorUsername string
		violationText    = "Сообщение не содержало текста."
		clarification    = "Не оставлено."
	)

	if msg.ReplyTo.Text != "" {
		violationText = msg.ReplyTo.Text
	}

	if ctx.Data() != "" {
		clarification = ctx.Data()
	}

	if incorrectUseReportCommand(ctx, msg) {
		cleanUpReport(ctx, msg, reporter)

		return nil
	}

	violator := msg.ReplyTo.Sender
	violatorID = violator.ID
	violatorUsername = violator.Username
	violationMessageID := msg.ReplyTo.ID
	messageLink := GenerateMessageLink(msg.ReplyTo.Chat, violationMessageID)

	pkgLog.Info(
		fmt.Sprintf("Получен отчет о новом нарушении правил от %s", GetGreetingName(reporter)),
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

	if reportAboutBot(ctx, violatorID, reporter) {
		return nil
	}

	sendNotification(ctx, violationText, violator, reporter, clarification, messageLink)
	cleanUpReport(ctx, msg, reporter)
	thxForReport(ctx, violationText, clarification, reporter)

	return nil
}

func thxForReport(
	ctx tele.Context,
	violationText string,
	clarification string,
	reporter *tele.User,
) {
	if _, err := ctx.Bot().Send(
		reporter,
		renderingTool.RenderText(`report_thx.txt`, struct {
			Text          string
			Clarification string
		}{
			Text:          violationText,
			Clarification: clarification,
		}),
		tele.ModeHTML,
		tele.NoPreview,
	); err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось послать благодарность за жалобу %s: %v",
			GetGreetingName(reporter),
			err,
		), nil)
	}
}

func cleanUpReport(ctx tele.Context, msg *tele.Message, reporter *tele.User) {
	if err := ctx.Bot().Delete(msg); err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось удалить сообщение с жалобой от %s: %v", GetGreetingName(reporter), err),
			pkgLogger.LogContext{
				"message_id":   msg.ID,
				"message_text": msg.Text,
			},
		)
	}
}

func reportAboutBot(ctx tele.Context, violatorID int64, reporter *tele.User) bool {
	if violatorID == config.BotID {
		if err := ctx.Reply(fmt.Sprintf("%s, ай-яй-яй! %s", GetGreetingName(reporter), llm.GetTeaser())); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось пообзываться в ответ на репорт на бота: %v", err),
				pkgLogger.LogContext{
					"reporter": reporter,
				},
			)
		}

		return true
	}

	return false
}

func incorrectUseReportCommand(ctx tele.Context, msg *tele.Message) bool {
	if msg.ReplyTo == nil || msg.ReplyTo.Sender.ID == msg.Sender.ID {
		if _, err := ctx.Bot().Send(
			msg.Sender,
			"Пожалуйста, используйте эту команду в ответе на сообщение с нарушением. "+
				"Подробнее: выполните /help в личной переписке с данным ботом.",
		); err != nil {
			pkgLog.Error(
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

func sendNotification(
	ctx tele.Context,
	violationText string,
	violator *tele.User,
	reporter *tele.User,
	clarification string,
	messageLink string,
) {
	if _, err := ctx.Bot().Send(
		&tele.Chat{ID: config.AdministrationChatID},
		renderingTool.RenderText(`report_notice.txt`, struct {
			ReporterUsername string
			ReporterID       int64
			Clarification    string
			MessageLink      string
			ViolatorUsername string
			ViolatorID       int64
			Text             string
		}{
			ReporterUsername: reporter.Username,
			ReporterID:       reporter.ID,
			Clarification:    clarification,
			MessageLink:      messageLink,
			ViolatorUsername: violator.Username,
			ViolatorID:       violator.ID,
			Text:             violationText,
		}),
		tele.ModeHTML,
		tele.NoPreview,
	); err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось послать в чат админов жалобу от %s на %s: %v",
			GetGreetingName(reporter),
			GetGreetingName(violator),
			err,
		), nil)
	}
}

func CommandRulesHandler(ctx tele.Context) error {
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

	pkgLog.Info(
		`Получен запрос правил`,
		pkgLogger.LogContext{
			"message": ctx.Message(),
		},
	)

	if targetMessage == nil {
		if err := ctx.Send(
			fmt.Sprintf(
				`Привет, %s! Вот <a href="%s">правила чата</a>, ознакомься.`,
				GetGreetingName(targetUser),
				config.RulesURL.String(),
			),
			tele.ModeHTML,
			tele.NoPreview,
		); err != nil {
			pkgLog.Error(
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
				GetGreetingName(targetUser),
				config.RulesURL.String(),
			),
			tele.ModeHTML,
			tele.NoPreview,
		); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить правила чата по команде /rules: %v", err),
				pkgLogger.LogContext{
					"ctx": ctx,
				},
			)
		}
	}

	if err := ctx.Bot().Delete(ctx.Message()); err != nil {
		pkgLog.Error(
			fmt.Sprintf(
				"Не удалось удалить сообщение с командой /rules от %s: %v",
				GetGreetingName(ctx.Message().Sender),
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

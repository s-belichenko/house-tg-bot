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
	return c.Send(llm.GetAnswerAboutKeys())
}

func CommandReportHandler(ctx tele.Context) error {
	msg := ctx.Message()
	reporter := msg.Sender
	violator := msg.ReplyTo.Sender
	clarification := "Не оставлено."
	chat := msg.ReplyTo.Chat
	violationMessageID := msg.ReplyTo.ID
	messageLink := GenerateMessageLink(chat, violationMessageID)

	if ctx.Data() != "" {
		clarification = ctx.Data()
	}

	pkgLog.Info(
		fmt.Sprintf("Новое нарушение правил от %s", GetGreetingName(reporter)),
		pkgLogger.LogContext{
			"reporter_username": reporter.Username,
			"reporter_id":       reporter.ID,
			"violator":          violator.Username,
			"violator_id":       violator.ID,
			"text":              msg.ReplyTo.Text,
			"clarification":     clarification,
			"message_link":      messageLink,
		},
	)

	if incorrectUseReportCommand(ctx, msg) {
		return nil
	}

	if reportAboutBot(ctx, violator, reporter) {
		return nil
	}

	sendNotification(ctx, msg, violator, reporter, clarification, messageLink)
	cleanUpReport(ctx, msg, reporter, violator)
	thxForReport(ctx, msg, clarification, reporter)

	return nil
}

func thxForReport(ctx tele.Context, msg *tele.Message, clarification string, reporter *tele.User) {
	if _, err := ctx.Bot().Send(
		reporter,
		renderingTool.RenderText(`report_thx.txt`, struct {
			Text          string
			Clarification string
		}{
			Text:          msg.ReplyTo.Text,
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

func cleanUpReport(ctx tele.Context, msg *tele.Message, reporter *tele.User, violator *tele.User) {
	if err := ctx.Bot().Delete(msg); err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось удалить сообщение с жалобой от %s: %v", GetGreetingName(reporter), err),
			pkgLogger.LogContext{
				"message_id":   msg.ID,
				"message_text": msg.Text,
				"violator_id":  violator.ID,
			},
		)
	}
}

func reportAboutBot(ctx tele.Context, violator *tele.User, reporter *tele.User) bool {
	if violator.ID == config.BotID {
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
	if msg.ReplyTo == nil {
		if err := ctx.Reply(fmt.Sprintf(
			"Пожалуйста, используйте эту команду в ответе на сообщение с нарушением. "+
				"Подробнее: выполните /help в личной переписке с @%s.", config.BotNickname),
		); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить уточнение про команду /report: %v", err),
				nil,
			)
		}

		return true
	}

	return false
}

func sendNotification(
	ctx tele.Context,
	msg *tele.Message,
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
			Text:             msg.ReplyTo.Text,
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
	pkgLog.Info(
		`Получен запрос правил`,
		pkgLogger.LogContext{
			"sender":          ctx.Message().Sender,
			"sender_reply_to": ctx.Message().ReplyTo.Sender,
		},
	)

	if _, err := ctx.Bot().Reply(
		ctx.Message().ReplyTo,
		fmt.Sprintf(
			`Привет, %s! Вот <a href="%s">правила чата</a>, ознакомься.`,
			GetGreetingName(ctx.Message().ReplyTo.Sender),
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

	if err := ctx.Bot().Delete(ctx.Message()); err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось удалить сообщение с командой /rules от %s: %v", GetGreetingName(ctx.Message().Sender), err),
			pkgLogger.LogContext{
				"message_id":   ctx.Message().ID,
				"message_text": ctx.Message().Text,
			},
		)
	}

	return nil
}

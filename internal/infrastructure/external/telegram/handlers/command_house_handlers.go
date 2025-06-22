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

	pkgLog.Debug("ReplyTo", pkgLogger.LogContext{
		"reply_to": msg.ReplyTo,
	})

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

		return nil
	}

	reporter := msg.Sender
	violator := msg.ReplyTo.Sender

	if violator.ID == config.BotID {
		if err := ctx.Reply(fmt.Sprintf("%s, ай-яй-яй! %s", GetGreetingName(reporter), llm.GetTeaser())); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось пообзываться в ответ на репорт на бота: %v", err),
				pkgLogger.LogContext{
					"reporter": reporter,
				},
			)
		}

		return nil
	}

	clarification := "Не оставлено."
	if ctx.Data() != "" {
		clarification = ctx.Data()
	}

	chat := msg.ReplyTo.Chat
	violationMessageID := msg.ReplyTo.ID
	messageLink := GenerateMessageLink(chat, violationMessageID)

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

	reportMessage := fmt.Sprintf(
		`
#REPORT
Новое нарушение правил:

Челобитчик: @%s (ID: %d)
Уточнение: %s
Ссылка: %s
Нарушитель: @%s (ID: %d)

Сообщение:
<blockquote>%s</blockquote>`,
		reporter.Username,
		reporter.ID,
		clarification,
		messageLink,
		violator.Username,
		violator.ID,
		msg.ReplyTo.Text,
	)

	adminChat := &tele.Chat{ID: config.AdministrationChatID}
	if _, err := ctx.Bot().Send(adminChat, reportMessage, tele.ModeHTML, tele.NoPreview); err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось послать в чат админов жалобу от %s на %s: %v",
			GetGreetingName(reporter),
			GetGreetingName(violator),
			err,
		), nil)
	}

	err := ctx.Bot().Delete(msg)
	if err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось удалить сообщение с жалобой от %s: %v", GetGreetingName(reporter), err),
			pkgLogger.LogContext{
				"message_id":   msg.ID,
				"message_text": msg.Text,
				"violator_id":  violator.ID,
			},
		)
	}

	thx := fmt.Sprintf(`
Спасибо за ваш отчет о нарушении правил. Администрация рассмотрит жалобу. Текст сообщения нарушителя:

<blockquote>%s</blockquote>

Ваше уточнение: %s`, msg.ReplyTo.Text, clarification)

	if _, err := ctx.Bot().Send(reporter, thx, tele.ModeHTML, tele.NoPreview); err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось послать благодарность за жалобу %s: %v",
			GetGreetingName(reporter),
			err,
		), nil)
	}

	return nil
}

func CommandRulesHandler(ctx tele.Context) error {
	msg := fmt.Sprintf(
		"Привет, %s! Вот <a href=\"%s\">правила чата</a>, ознакомься.",
		GetGreetingName(ctx.Message().Sender),
		config.RulesURL.String(),
	)

	err := ctx.Reply(msg, tele.ModeHTML, tele.NoPreview)
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось отправить правила чата по команде /rules: %v", err),
			pkgLogger.LogContext{
				"ctx": ctx,
			},
		)
	}

	return nil
}

package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
	"s-belichenko/ilovaiskaya2-bot/cmd/llm"
	pkgLogger "s-belichenko/ilovaiskaya2-bot/pkg/logger"
)

// Команды бота для домового чата.
var (
	KeysCommand   = tele.Command{Text: "keys", Description: "ПИК, где ключи?"}
	ReportCommand = tele.Command{Text: "report", Description: "Сообщить о нарушении правил, формат: /report [уточнение]"}
)

func CommandKeysHandler(c tele.Context) error {
	return c.Send(llm.GetAnswerAboutKeys())
}

func CommandReportHandler(ctx tele.Context) error {
	var m = ctx.Message()

	pkgLog.Debug("ReplyTo", pkgLogger.LogContext{
		"reply_to": m.ReplyTo,
	})

	if m.ReplyTo == nil {
		if err := ctx.Reply("Пожалуйста, используйте эту команду в ответе на сообщение с нарушением. Подробнее: выполните /help в личной переписке с @lp_13x_bot."); err != nil {
			pkgLog.Error(fmt.Sprintf("Не удалось отправить уточнение про команду /report: %v", err), nil)
		}

		return nil
	}

	reporter := m.Sender
	violator := m.ReplyTo.Sender

	if violator.ID == config.BotID {
		if err := ctx.Reply(fmt.Sprintf("%s, ай-яй-яй! %s", GetGreetingName(reporter), llm.GetTeaser())); err != nil {
			pkgLog.Error(fmt.Sprintf("Не удалось пообзываться в ответ на репорт на бота: %v", err), pkgLogger.LogContext{
				"reporter": reporter,
			})
		}

		return nil
	}

	clarification := ctx.Data()
	chat := m.ReplyTo.Chat
	violationMessageID := m.ReplyTo.ID
	messageLink := GenerateMessageLink(chat, violationMessageID)

	pkgLog.Info(fmt.Sprintf("Новое нарушение правил от %s", GetGreetingName(reporter)), pkgLogger.LogContext{
		"reporter_username": reporter.Username,
		"reporter_id":       reporter.ID,
		"violator":          violator.Username,
		"violator_id":       violator.ID,
		"text":              m.ReplyTo.Text,
		"clarification":     clarification,
		"message_link":      messageLink,
	})

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
		m.ReplyTo.Text,
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

	err := ctx.Bot().Delete(m)
	if err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось удалить сообщение с жалобой от %s: %v", GetGreetingName(reporter), err),
			pkgLogger.LogContext{
				"message_id":   m.ID,
				"message_text": m.Text,
				"violator_id":  violator.ID,
			},
		)
	}

	thx := fmt.Sprintf(`
Спасибо за ваше сообщение. Администрация рассмотрит жалобу. Сообщение:

<blockquote>%s</blockquote>`, m.ReplyTo.Text)

	if _, err := ctx.Bot().Send(reporter, thx, tele.ModeHTML, tele.NoPreview); err != nil {
		pkgLog.Error(fmt.Sprintf(
			"Не удалось послать благодарность за жалобу %s: %v",
			GetGreetingName(reporter),
			err,
		), nil)
	}

	return nil
}

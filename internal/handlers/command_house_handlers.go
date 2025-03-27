package handlers

import (
	"fmt"
	"s-belichenko/ilovaiskaya2-bot/cmd/llm"

	tele "gopkg.in/telebot.v4"
	intLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

// Команды бота для домового чата
var (
	KeysCommand   = tele.Command{Text: "keys", Description: "ПИК, где ключи?"}
	ReportCommand = tele.Command{Text: "report", Description: "Сообщить о нарушении правил, формат: /report [уточнение]"}
)

func CommandKeysHandler(c tele.Context) error {
	return c.Send(llm.GetAnswerAboutKeys())
}

func CommandReportHandler(c tele.Context) error {
	m := c.Message()
	log.Debug("ReplyTo", intLog.LogContext{
		"reply_to": m.ReplyTo,
	})
	if m.ReplyTo == nil {
		if err := c.Reply("Пожалуйста, используйте эту команду в ответе на сообщение с нарушением. Подробнее: выполните /help в личной переписке с @lp_13x_bot."); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить уточнение про команду /report: %v", err), nil)
		}
		return nil
	}

	reporter := m.Sender
	violator := m.ReplyTo.Sender

	if violator.ID == config.BotID {
		if err := c.Reply(fmt.Sprintf("%s, ай-яй-яй! %s", GetGreetingName(reporter), llm.GetTeaser())); err != nil {
			log.Error(fmt.Sprintf("Не удалось пообзываться в ответ на репорт на бота: %v", err), intLog.LogContext{
				"reporter": reporter,
			})
		}
		return nil
	}

	clarification := c.Data()
	chat := m.ReplyTo.Chat
	violationMessageID := m.ReplyTo.ID
	messageLink := GenerateMessageLink(chat, violationMessageID)

	log.Info(fmt.Sprintf("Новое нарушение правил от %s", GetGreetingName(reporter)), intLog.LogContext{
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
	if _, err := c.Bot().Send(adminChat, reportMessage, tele.ModeHTML, tele.NoPreview); err != nil {
		log.Error(fmt.Sprintf(
			"Не удалось послать в чат админов жалобу от %s на %s: %v",
			GetGreetingName(reporter),
			GetGreetingName(violator),
			err,
		), nil)
	}

	err := c.Bot().Delete(m)
	if err != nil {
		log.Error(fmt.Sprintf(
			"Не удалось удалить сообщение с жалобой от %s: %v", GetGreetingName(reporter), err),
			intLog.LogContext{
				"message_id":   m.ID,
				"message_text": m.Text,
				"violator_id":  violator.ID,
			},
		)
	}

	thx := fmt.Sprintf(`
Спасибо за ваше сообщение. Администрация рассмотрит жалобу. Сообщение:

<blockquote>%s</blockquote>`, m.ReplyTo.Text)

	if _, err := c.Bot().Send(reporter, thx, tele.ModeHTML, tele.NoPreview); err != nil {
		log.Error(fmt.Sprintf(
			"Не удалось послать благодарность за жалобу %s: %v",
			GetGreetingName(reporter),
			err,
		), nil)
	}

	return nil
}

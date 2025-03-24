package handlers

import (
	"fmt"
	"s-belichenko/ilovaiskaya2-bot/cmd/llm"

	tele "gopkg.in/telebot.v4"
	yaLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

// Команды бота
var (
	StartCommand = tele.Command{Text: "start", Description: "Начать работу с ботом."}
	HelpCommand  = tele.Command{Text: "help", Description: "Справка по боту"}

	KeysCommand   = tele.Command{Text: "keys", Description: "ПИК, где ключи?"}
	ReportCommand = tele.Command{Text: "report", Description: "Сообщить о нарушении правил, формат: /report [уточнение]"}

	HelpAdminChatCommand = tele.Command{Text: "help_admin", Description: "Справка по боту для админов"}
	BanCommand           = tele.Command{Text: "ban", Description: "Забанить пользователя из домового чата"}
	UnbanCommand         = tele.Command{Text: "unban", Description: "Разбанить пользователя из домового чата"}
	KickCommand          = tele.Command{Text: "kick", Description: "Удалить пользователя из домового чата навсегда"}

	SetCommandsCommand = tele.Command{Text: "set_commands", Description: "Установить команды бота"}
)

func CommandHelpHandler(c tele.Context) error {
	help := fmt.Sprintf(`
Привет, это бот чата дома Иловайская, 2.

Команды:

/start – Начало работы с ботом. В домовом чате не требуется.
/help – Текущая справка.
/keys – Шуточная команда, заставющая бота ответить на вопрос "ПИК, где мои ключи?" Работает только в теме "Оффтоп". Но осторожнее, половина соседей уже ненавидит эту команду.
/report – Сообщить о нарушении правил. Напишите ее в ответе на сообщение с нарушением правил чата, после команды через пробел можете уточнить причину жалобы, например:
<blockquote>/report Ругается матом, редиска!</blockquote>
Сообщение с жалобой будет отправлено администраторам, а ваше сообщение с командой удалено.

<a href="https://ilovaiskaya2.homes/#rules">Ссылка на правила</a>.`)
	err := c.Send(help, tele.ModeHTML, tele.NoPreview)
	if err != nil {
		log.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}
	return nil
}

func CommandBanHandler(c tele.Context) error {
	return nil
}

func CommandKickHandler(c tele.Context) error {
	return nil
}

func CommandHelpAdminHandler(c tele.Context) error {
	help := fmt.Sprintf(`
Справка для администратора. Все команды ниже используются только в чате администраторов.

Команды:

/help_admin – Текущая справка.
/ban &lt;username | user_id&gt; [period] – Забанить пользователя из домового чата
/unban &lt;username | user_id&gt; – Разабанить пользователя из домового чата
/kick &lt;username | user_id&gt; – Удалить пользователя из домового чата навсегда

<a href="https://ilovaiskaya2.homes/#rules">Ссылка на правила</a>.`)
	err := c.Send(help, tele.ModeHTML, tele.NoPreview)
	if err != nil {
		log.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}
	return nil
}

func CommandStartHandler(c tele.Context) error {
	err := c.Send(fmt.Sprintf(
		"Привет, %s! Ознакомься со справкой по работе с ботом: /help", getUsername(*c.Sender()),
	))
	if err != nil {
		log.Error(fmt.Sprintf("Не удалось отправить ответ на команду /start: %v", err), nil)
	}
	return err
}

func CommandKeysHandler(c tele.Context) error {
	return c.Send(llm.GetAnswerAboutKeys())
}

func CommandReportHandler(c tele.Context) error {
	m := c.Message()

	if m.ReplyTo == nil {
		if err := c.Reply("Пожалуйста, используйте эту команду в ответе на сообщение с нарушением. Подробнее: выполните /help в личной переписке с @lp_13x_bot."); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить уточнение про команду /report: %v", err), nil)
		}
		return nil
	}

	reporter := m.Sender
	violator := m.ReplyTo.Sender

	if violator.ID == config.BotID {
		if err := c.Reply(fmt.Sprintf("@%s, ай-яй-яй! %s", reporter.Username, llm.GetTeaser())); err != nil {
			log.Error(fmt.Sprintf("Не удалось пообзываться в ответ на непорт на бота: %v", err), nil)
		}
		return nil
	}

	clarification := c.Data()
	chat := m.ReplyTo.Chat
	violationMessageID := m.ReplyTo.ID
	messageLink := GenerateMessageLink(chat, violationMessageID)

	log.Info(fmt.Sprintf("Новое нарушение правил от %s", reporter.Username), yaLog.LogContext{
		"reporter_id":   reporter.ID,
		"violator":      violator.Username,
		"violator_id":   violator.ID,
		"text":          m.ReplyTo.Text,
		"clarification": clarification,
		"message_link":  messageLink,
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
			"Не удалось послать в чат админов жалобу от @%s на %s: %v",
			reporter.Username,
			violator.Username,
			err,
		), nil)
	}

	err := c.Bot().Delete(m)
	if err != nil {
		log.Error(fmt.Sprintf(
			"Не удалось удалить сообщение с жалобой от @%s: %v", reporter.Username, err),
			yaLog.LogContext{
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
			"Не удалось послать благодарность за жалобу @%s: %v",
			reporter.Username,
			err,
		), nil)
	}

	return nil
}

func CommandSetCommandsHandler(c tele.Context) error {
	// По умолчанию
	setCommands(c,
		[]tele.Command{StartCommand},
		tele.CommandScope{Type: tele.CommandScopeDefault})
	// Для личных чатов со всеми подряд
	setCommands(c,
		[]tele.Command{StartCommand, HelpCommand},
		tele.CommandScope{Type: tele.CommandScopeAllPrivateChats})
	// Для участников домового чата
	setCommands(c,
		[]tele.Command{KeysCommand, ReportCommand},
		tele.CommandScope{Type: tele.CommandScopeDefault, ChatID: config.HouseChatId})
	// Для админов домового чата
	setCommands(c,
		[]tele.Command{KeysCommand},
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.HouseChatId})
	// Для участников админского чата
	setCommands(c,
		[]tele.Command{HelpAdminChatCommand, BanCommand, KickCommand},
		tele.CommandScope{Type: tele.CommandScopeDefault, ChatID: config.AdministrationChatID})
	// Для админов админского чата
	setCommands(c,
		[]tele.Command{SetCommandsCommand, HelpAdminChatCommand, BanCommand, KickCommand},
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.AdministrationChatID})

	return nil
}

package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

// Команды бота для административного чата.
var (
	StartCommand = tele.Command{Text: "start", Description: "Начать работу с ботом"}
	HelpCommand  = tele.Command{Text: "help", Description: "Справка по боту"}
)

func CommandStartHandler(c tele.Context) error {
	err := c.Send(fmt.Sprintf(
		"Привет, %s! Ознакомься со справкой по работе с ботом: /help", GetGreetingName(c.Sender()),
	))
	if err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось отправить ответ на команду /start: %v", err), nil)
	}

	return err
}

func CommandHelpHandler(ctx tele.Context) error {
	tmpl := `Привет, это бот <a href="%s">чата дома</a> по адресу %s.

<b>Команды в переписке с ботом:</b>

/%s – Начало работы с ботом. В домовом чате не требуется.
/%s – Текущая справка.
/%s – Посмотреть правила чата.

<b>Команды в домовом чате:</b>

/%s – Шуточная команда, заставляющая бота ответить на вопрос "ПИК, где мои ключи?" Работает ` +
		`только в теме "Оффтоп". Но осторожнее, половина соседей уже ненавидит эту команду.
/%s – Сообщить о нарушении правил. Напишите ее в ответе на сообщение с нарушением правил ` +
		`чата, после команды через пробел можете уточнить причину жалобы, например:
<blockquote>/report Ругается матом, редиска!</blockquote>
Сообщение с жалобой будет отправлено администраторам, а ваше сообщение с командой удалено.
/%s – Посмотреть правила чата.

<a href="%s">Ссылка на правила</a>.`

	help := fmt.Sprintf(
		tmpl,
		cfg.InviteURL.String(),
		cfg.HomeAddress,
		StartCommand.Text,
		HelpCommand.Text,
		RulesCommand.Text,
		KeysCommand.Text,
		ReportCommand.Text,
		RulesCommand.Text,
		cfg.RulesURL.String(),
	)

	if err := ctx.Send(help, tele.ModeHTML, tele.NoPreview); err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}

	return nil
}

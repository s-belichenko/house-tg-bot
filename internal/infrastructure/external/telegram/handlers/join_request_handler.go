package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

const hiMsg = "Привет, вы подали заявку на вступление в <a href=\"%s\">чат</a> дома по адресу Иловайская, 2 " +
	"(бывшие 13-е корпуса) в ЖК Люблинский парк. Согласно правилам чата вам будет необходимо верифицировать себя как " +
	"реального соседа, для этого будет необходимо предоставить владельцу чата скриншот из ЛК ПИК, где будет видно, " +
	"что вы реальный собственник в нашем доме. Личные данные можно скрыть. Верификация введена для комфортного " +
	"общения, чтобы исключить ботов, спамеров и просто левые аккаунты.\n" +
	"\n" +
	"Если вдруг у вас не работает приложение, это у многих проблема, то в ЛК можно зайти через " +
	"<a href=\"https://client.pik.ru/\">сайт</a>, либо с телефона, либо с компьютера в браузере.\n" +
	"\n" +
	"Также у нашего чата и будущего дома есть небольшой сайт там можно ознакомиться с <a href=\"%s\">правилами</a>, " +
	"найти полезные ссылки для будущего жильца и тп.\n" +
	"\n" +
	"Если с вами не связались в течение 24 часов, попробуйте проверить переписку с владельцем чата: @%s, " +
	"возможно, вы пропустили сообщение от него. Если нет, то можете написать ему самостоятельно."

func JoinRequestHandler(ctx tele.Context) error {
	pkgLog.Info("Получена заявка на вступление в чат", pkgLogger.LogContext{
		"user_id":   ctx.Sender().ID,
		"username":  ctx.Sender().Username,
		"firstname": ctx.Sender().FirstName,
		"lastname":  ctx.Sender().LastName,
	})

	if _, err := ctx.Bot().Send(
		ctx.Sender(),
		fmt.Sprintf(hiMsg, config.InviteURL, config.RulesURL.String(), config.OwnerNickname),
		tele.ModeHTML, tele.NoPreview,
	); err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось ответить на заявку: %v", err), pkgLogger.LogContext{
			"user_id":   ctx.Sender().ID,
			"username":  ctx.Sender().Username,
			"firstname": ctx.Sender().FirstName,
			"lastname":  ctx.Sender().LastName,
		})

		return err
	}

	adminChat := &tele.Chat{ID: config.AdministrationChatID}
	requestMsg := fmt.Sprintf(`
#JOIN_REQUEST
Новая заявка на вступление в чат.

user_id: %d
username: @%s
firstname: %s
lastname: %s
`, ctx.Sender().ID, ctx.Sender().Username, ctx.Sender().FirstName, ctx.Sender().LastName)

	if _, err := ctx.Bot().Send(adminChat, requestMsg); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось ответить на заявку на вступление: %v", err),
			pkgLogger.LogContext{
				"user_id":   ctx.Sender().ID,
				"username":  ctx.Sender().Username,
				"firstname": ctx.Sender().FirstName,
				"lastname":  ctx.Sender().LastName,
			},
		)
	}

	return nil
}

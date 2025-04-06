package handlers

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"
	pkgLog "s-belichenko/ilovaiskaya2-bot/pkg/logger"
)

const (
	muteCommandFormat   = `/mute &lt;user_id&gt; [days] (от 1 до 366, иначе бессрочно)`
	unmuteCommandFormat = `/unmute &lt;user_id&gt;`
	banCommandFormat    = `/ban &lt;user_id&gt; [days] (от 1 до 366, иначе бессрочно)`
	unbanCommandFormat  = `/unban &lt;user_id&gt;`
)

var (
	HelpAdminChatCommand = tele.Command{Text: "help_admin", Description: "Справка по боту для админов"}
	MuteCommand          = tele.Command{Text: "mute", Description: "Ограничить пользователя в домовом чате"}
	UnmuteCommand        = tele.Command{Text: "unmute", Description: "Снять ограничения с пользователя в домовом чате"}
	BanCommand           = tele.Command{Text: "ban", Description: "Заблокировать пользователя из домового чата"}
	UnbanCommand         = tele.Command{Text: "unban", Description: "Разблокировать пользователя из домового чата"}

	SetCommandsCommand    = tele.Command{Text: "set_commands", Description: "Установить команды бота"}
	DeleteCommandsCommand = tele.Command{Text: "delete_commands", Description: "Удалить команды бота"}
)

func CommandMuteHandler(c tele.Context) error {
	var violator *tele.ChatMember
	d := c.Data()
	f := strings.Fields(d)

	switch len(f) {
	case 1:
		if user := createUserViolator(f[0]); user != nil {
			violator = &tele.ChatMember{
				User:   user,
				Rights: tele.NoRights(),
			}
		}
	case 2:
		if user := createUserViolator(f[0]); user != nil {
			violator = &tele.ChatMember{
				User:            user,
				Rights:          tele.NoRights(),
				RestrictedUntil: createUnixTimeFromDays(f[1]),
			}
		}
	}
	if violator == nil {
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", muteCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /mute: %v", err), pkgLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if err := c.Bot().Restrict(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
		log.Error(fmt.Sprintf("Не удалось ограничить пользователя: %v", err), pkgLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось ограничить пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось ограничить пользователя: %v", err), pkgLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}

	if err := c.Reply("Пользователь ограничен."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь ограничен: %v", err), pkgLog.LogContext{
			"message": c.Message(),
		})
	}
	// FIXME: Посылать сообщение пользователю об ограничениях? А если он не начал общение с ботом?
	log.Info("Пользователь ограничен", pkgLog.LogContext{
		"admin_id":         c.Message().Sender.ID,
		"admin_username":   c.Message().Sender.Username,
		"admin_first_name": c.Message().Sender.FirstName,
		"admin_last_name":  c.Message().Sender.LastName,
		"violator":         violator,
	})
	return nil
}

func CommandUnmuteHandler(c tele.Context) error {
	var violator *tele.ChatMember
	d := c.Data()
	f := strings.Fields(d)
	switch len(f) {
	case 1:
		if user := createUserViolator(f[0]); &user != nil {
			violator = &tele.ChatMember{User: user, Rights: tele.NoRestrictions()}
		}
	}
	if violator == nil {
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", unmuteCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /unmute: %v", err), pkgLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if err := c.Bot().Promote(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
		log.Error(fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err), pkgLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось снять ограничения с пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err), pkgLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}

	if err := c.Reply("Сняты ограничения с пользователя."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что сняты ограничения с пользователя: %v", err), pkgLog.LogContext{
			"message": c.Message(),
		})
	}
	// FIXME: Посылать сообщение пользователю об отмене бана? А если он не начал общение с ботом?
	log.Info("Сняты ограничения с пользователя", pkgLog.LogContext{
		"admin_id":         c.Message().Sender.ID,
		"admin_username":   c.Message().Sender.Username,
		"admin_first_name": c.Message().Sender.FirstName,
		"admin_last_name":  c.Message().Sender.LastName,
		"violator":         violator,
	})
	return nil
}

func CommandBanHandler(c tele.Context) error {
	var violator *tele.ChatMember
	d := c.Data()
	f := strings.Fields(d)
	switch len(f) {
	case 1:
		if user := createUserViolator(f[0]); user != nil {
			violator = &tele.ChatMember{User: user, RestrictedUntil: tele.Forever()}
		}
	case 2:
		if user := createUserViolator(f[0]); user != nil {
			violator = &tele.ChatMember{User: user, RestrictedUntil: createUnixTimeFromDays(f[1])}
		}
	}
	if violator == nil {
		log.Warn(fmt.Sprintf("Вызов команды /ban без аргументов"), pkgLog.LogContext{
			"arguments_string": d,
		})
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /ban: %v", err), pkgLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if err := c.Bot().Ban(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
		log.Error(fmt.Sprintf("Не удалось заблокировать пользователя: %v", err), pkgLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось заблокировать пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось заблокировать пользователя: %v", err), pkgLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}

	if err := c.Reply("Пользователь заблокирован."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь заблокирован: %v", err), pkgLog.LogContext{
			"message": c.Message(),
		})
	}
	log.Info("Пользователь заблокирован", pkgLog.LogContext{
		"admin_id":         c.Message().Sender.ID,
		"admin_username":   c.Message().Sender.Username,
		"admin_first_name": c.Message().Sender.FirstName,
		"admin_last_name":  c.Message().Sender.LastName,
		"violator":         violator,
	})
	return nil
}

func CommandUnbanHandler(c tele.Context) error {
	var violator *tele.User
	d := c.Data()
	f := strings.Fields(d)
	switch len(f) {
	case 1:
		violator = createUserViolator(f[0])
	}
	if violator == nil {
		log.Warn(fmt.Sprintf("Вызов команды /unban без аргументов"), pkgLog.LogContext{
			"arguments_string": d,
		})
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /unban: %v", err), pkgLog.LogContext{
				"message": c.Message(),
			})
		}
	}
	if err := c.Bot().Unban(&tele.Chat{ID: config.HouseChatId}, violator, true); err != nil {
		log.Error(fmt.Sprintf("Не удалось разблокировать пользователя: %v", err), pkgLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось разблокировать пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось разблокировать пользователя: %v", err), pkgLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}

	if err := c.Reply("Пользователь разблокирован."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь разблокирован: %v", err), pkgLog.LogContext{
			"message": c.Message(),
		})
	}
	log.Info("Пользователь разблокирован", pkgLog.LogContext{
		"admin_id":         c.Message().Sender.ID,
		"admin_username":   c.Message().Sender.Username,
		"admin_first_name": c.Message().Sender.FirstName,
		"admin_last_name":  c.Message().Sender.LastName,
		"violator":         violator,
	})
	return nil
}

func CommandHelpAdminHandler(c tele.Context) error {
	help := fmt.Sprintf(`
Справка для администратора. Все команды ниже используются только в чате администраторов.

Команды:

/help_admin – Текущая справка.
%s – Ограничить пользователя в домовом чате. 
%s – Снять ограничения с пользователя в домовом чате.
%s – Забанить пользователя в домовом чата.
%s – Разбанить пользователя в домовом чата.

<a href="https://ilovaiskaya2.homes/#rules">Ссылка на правила</a>.`,
		muteCommandFormat,
		unmuteCommandFormat,
		banCommandFormat,
		unbanCommandFormat,
	)
	err := c.Send(help, tele.ModeHTML, tele.NoPreview)
	if err != nil {
		log.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}
	return nil
}

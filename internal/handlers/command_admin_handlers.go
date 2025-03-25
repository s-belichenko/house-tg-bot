package handlers

import (
	"fmt"
	tele "gopkg.in/telebot.v4"
	yaLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
	"strings"
)

const (
	restrictCommandFormat       = `/restrict &lt;username | user_id&gt; [days] (от 1 до 366, иначе бессрочно)`
	remoteRestrictCommandFormat = `/remove_restrict &lt;username | user_id&gt;`
	banCommandFormat            = `/ban &lt;username | user_id&gt; [days] (от 1 до 366, иначе бессрочно)`
	unbanCommandFormat          = `/unban &lt;username | user_id&gt;`
)

var (
	HelpAdminChatCommand  = tele.Command{Text: "help_admin", Description: "Справка по боту для админов"}
	RestrictCommand       = tele.Command{Text: "restrict", Description: "Ограничить пользователя  в домовом чате"}
	RemoveRestrictCommand = tele.Command{Text: "remove_restrict", Description: "Снять ограничения с пользователя в домовом чате"}
	BanCommand            = tele.Command{Text: "ban", Description: "Заблокировать пользователя из домового чата"}
	UnbanCommand          = tele.Command{Text: "unban", Description: "Разблокировать пользователя из домового чата"}

	SetCommandsCommand = tele.Command{Text: "set_commands", Description: "Установить команды бота"}
)

func CommandRestrictHandler(c tele.Context) error {
	var violator *tele.ChatMember
	d := c.Data()
	f := strings.Fields(d)

	switch len(f) {
	case 0:
		log.Warn(fmt.Sprintf("Вызов команды /restrict без аргументов"), yaLog.LogContext{
			"arguments_string": d,
		})
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", restrictCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /restrict: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	case 1:
		user := createUserViolator(c, f[0])
		if user == nil {
			return nil
		}
		violator = &tele.ChatMember{
			User:   user,
			Rights: tele.NoRights(),
		}
	case 2:
		user := createUserViolator(c, f[0])
		if user == nil {
			return nil
		}
		violator = &tele.ChatMember{
			User:            user,
			Rights:          tele.NoRights(),
			RestrictedUntil: createUnixTimeFromDays(f[1]),
		}
	default:
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", restrictCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /restrict: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if violator == nil {
		if err := c.Reply("Не удалось ограничить пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось ограничить пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if err := c.Bot().Restrict(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
		log.Error(fmt.Sprintf("Не удалось ограничить пользователя: %v", err), yaLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось ограничить пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось ограничить пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}

		return nil
	}
	log.Debug("Успешно отправлен запрос на ограничение пользователя", yaLog.LogContext{
		"violator": violator,
	})
	if err := c.Reply("Пользователь ограничен."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь ограничен: %v", err), yaLog.LogContext{
			"message": c.Message(),
		})
	}
	// FIXME: Посылать сообщение пользователю об ограничениях? А если он не начал общение с ботом?
	log.Info("Пользователь ограничен", yaLog.LogContext{
		"admin_id":         c.Message().Sender.ID,
		"admin_username":   c.Message().Sender.Username,
		"admin_first_name": c.Message().Sender.FirstName,
		"admin_last_name":  c.Message().Sender.LastName,
		"violator":         violator,
	})
	return nil
}

func CommandRemoveRestrictHandler(c tele.Context) error {
	var violator *tele.ChatMember
	d := c.Data()
	f := strings.Fields(d)
	switch len(f) {
	case 0:
		log.Warn(fmt.Sprintf("Вызов команды /remove_restrict без аргументов"), yaLog.LogContext{
			"arguments_string": d,
		})
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", remoteRestrictCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /remove_restrict: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	case 1:
		user := createUserViolator(c, f[0])
		if user == nil {
			return nil
		}
		violator = &tele.ChatMember{User: user, Rights: tele.NoRestrictions()}
	default:
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", restrictCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /remove_restrict: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if violator == nil {
		if err := c.Reply("Не удалось снять ограничения с пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось ограничения с пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if err := c.Bot().Promote(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
		log.Error(fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err), yaLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось снять ограничения с пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}

		return nil
	}
	log.Debug("Успешно отправлен запрос на снятие ограничений", yaLog.LogContext{
		"violator": violator,
	})
	if err := c.Reply("Сняты ограничения с пользователя."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что сняты ограничения с пользователя: %v", err), yaLog.LogContext{
			"message": c.Message(),
		})
	}
	// FIXME: Посылать сообщение пользователю об отмене бана? А если он не начал общение с ботом?
	log.Info("Сняты ограничения с пользователя", yaLog.LogContext{
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
	case 0:
		log.Warn(fmt.Sprintf("Вызов команды /ban без аргументов"), yaLog.LogContext{
			"arguments_string": d,
		})
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /ban: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	case 1:
		user := createUserViolator(c, f[0])
		if user == nil {
			return nil
		}
		violator = &tele.ChatMember{User: user, RestrictedUntil: tele.Forever()}
	case 2:
		user := createUserViolator(c, f[0])
		if user == nil {
			return nil
		}
		violator = &tele.ChatMember{User: user, RestrictedUntil: createUnixTimeFromDays(f[1])}
	default:
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /ban: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if violator == nil {
		if err := c.Reply("Не удалось заблокировать пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось заблокировать пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if err := c.Bot().Ban(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
		log.Error(fmt.Sprintf("Не удалось заблокировать пользователя: %v", err), yaLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось заблокировать пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось заблокировать пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}

		return nil
	}
	log.Debug("Успешно отправлен запрос на блокировку", yaLog.LogContext{
		"violator": violator,
	})
	if err := c.Reply("Пользователь заблокирован."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь заблокирован: %v", err), yaLog.LogContext{
			"message": c.Message(),
		})
	}
	log.Info("Пользователь заблокирован", yaLog.LogContext{
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
	case 0:
		log.Warn(fmt.Sprintf("Вызов команды /unban без аргументов"), yaLog.LogContext{
			"arguments_string": d,
		})
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /unban: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	case 1:
		violator = createUserViolator(c, f[0])
	default:
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /unban: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if violator == nil {
		if err := c.Reply("Не удалось разблокировать пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось разблокировать пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if err := c.Bot().Unban(&tele.Chat{ID: config.HouseChatId}, violator, true); err != nil {
		log.Error(fmt.Sprintf("Не удалось разблокировать пользователя: %v", err), yaLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось разблокировать пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось разблокировать пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}

		return nil
	}
	log.Debug("Успешно отправлен запрос на разблокировку", yaLog.LogContext{
		"violator": violator,
	})
	if err := c.Reply("Пользователь разблокирован."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь разблокирован: %v", err), yaLog.LogContext{
			"message": c.Message(),
		})
	}
	log.Info("Пользователь разблокирован", yaLog.LogContext{
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
		restrictCommandFormat,
		remoteRestrictCommandFormat,
		banCommandFormat,
		unbanCommandFormat,
	)
	err := c.Send(help, tele.ModeHTML, tele.NoPreview)
	if err != nil {
		log.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
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
		[]tele.Command{HelpAdminChatCommand, RestrictCommand, RemoveRestrictCommand, BanCommand, UnbanCommand},
		tele.CommandScope{Type: tele.CommandScopeDefault, ChatID: config.AdministrationChatID})
	// Для админов админского чата
	setCommands(c,
		[]tele.Command{SetCommandsCommand, HelpAdminChatCommand, RestrictCommand, RemoveRestrictCommand, BanCommand, UnbanCommand},
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.AdministrationChatID})

	return nil
}

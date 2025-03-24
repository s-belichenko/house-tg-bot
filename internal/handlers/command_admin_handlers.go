package handlers

import (
	"fmt"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
	yaLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

const (
	banCommandFormat   = `/ban &lt;username | user_id&gt; [days] Если [days] равен 0 или более 366, пользователь будет забанен навсегда.`
	unbanCommandFormat = `/unban &lt;username | user_id&gt;`
	kickCommandFormat  = `/kick &lt;username | user_id&gt; [days] Если [days] равен 0 или более 366, пользователь будет удален навсегда.`
)

var (
	HelpAdminChatCommand = tele.Command{Text: "help_admin", Description: "Справка по боту для админов"}
	BanCommand           = tele.Command{Text: "ban", Description: "Забанить пользователя  в домовом чате"}
	UnbanCommand         = tele.Command{Text: "unban", Description: "Разбанить пользователя в домовом чате"}
	KickCommand          = tele.Command{Text: "kick", Description: "Удалить пользователя из домового чата"}

	SetCommandsCommand = tele.Command{Text: "set_commands", Description: "Установить команды бота"}
)

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
		violator = createMemberViolator(c, f[0], tele.Forever())
	case 2:
		// Дни в секундах плюс один час для просмотра после бана в настройках
		term := (parseDays(f[1]) * 86400) + 600
		violator = createMemberViolator(c, f[0], time.Now().Unix()+term)
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
	if err := c.Bot().Restrict(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
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
	log.Debug("Успешно отправлен запрос на бан", yaLog.LogContext{
		"violator": violator,
	})
	if err := c.Reply("Пользователь заблокирован."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь заблокирован: %v", err), yaLog.LogContext{
			"message": c.Message(),
		})
	}
	// FIXME: Посылать сообщение пользователю о бане? А если он не начал общение с ботом?
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
		log.Warn(fmt.Sprintf("Вызов команды /ubban без аргументов"), yaLog.LogContext{
			"arguments_string": d,
		})
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", unbanCommandFormat), tele.ModeHTML); err != nil {
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
	if err := c.Bot().Unban(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
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
	log.Debug("Успешно отправлен запрос на разбан", yaLog.LogContext{
		"violator": violator,
	})
	if err := c.Reply("Пользователь разблокирован."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь разблокирован: %v", err), yaLog.LogContext{
			"message": c.Message(),
		})
	}
	// FIXME: Посылать сообщение пользователю об отмене бана? А если он не начал общение с ботом?
	log.Info("Пользователь разблокирован", yaLog.LogContext{
		"admin_id":         c.Message().Sender.ID,
		"admin_username":   c.Message().Sender.Username,
		"admin_first_name": c.Message().Sender.FirstName,
		"admin_last_name":  c.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func CommandKickHandler(c tele.Context) error {
	var violator *tele.ChatMember
	var days int64
	d := c.Data()
	f := strings.Fields(d)
	switch len(f) {
	case 0:
		log.Warn(fmt.Sprintf("Вызов команды /kick без аргументов"), yaLog.LogContext{
			"arguments_string": d,
		})
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", kickCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /kick: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	case 1:
		violator = createMemberViolator(c, f[0], tele.Forever())
	case 2:
		days = parseDays(f[1])
		violator = createMemberViolator(c, f[0], time.Now().Unix()+(days*86400))
	default:
		if err := c.Reply(fmt.Sprintf("Верный формат команды: %s", kickCommandFormat), tele.ModeHTML); err != nil {
			log.Error(fmt.Sprintf("Не удалось отправить подсказку по команде /kick: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if violator == nil {
		if err := c.Reply("Не удалось удалить пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось удалить пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}
		return nil
	}
	if err := c.Bot().Ban(&tele.Chat{ID: config.HouseChatId}, violator); err != nil {
		log.Error(fmt.Sprintf("Не удалось удалить пользователя: %v", err), yaLog.LogContext{
			"violator": violator,
			"message":  c.Message(),
		})
		if err := c.Reply("Не удалось удалить пользователя."); err != nil {
			log.Error(fmt.Sprintf("Не удалось удалить пользователя: %v", err), yaLog.LogContext{
				"message": c.Message(),
			})
		}

		return nil
	}
	log.Debug("Успешно отправлен запрос на удаление", yaLog.LogContext{
		"violator": violator,
	})
	if err := c.Reply("Пользователь удален."); err != nil {
		log.Error(fmt.Sprintf("Не удалось уведомить что пользователь удален: %v", err), yaLog.LogContext{
			"message": c.Message(),
		})
	}
	log.Info("Пользователь удален.", yaLog.LogContext{
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
%s – Забанить пользователя в домовом чате. 
%s – Разабанить пользователя в домовом чате.
%s – Удалить пользователя из домового чата.

<a href="https://ilovaiskaya2.homes/#rules">Ссылка на правила</a>.`, banCommandFormat, unbanCommandFormat, kickCommandFormat)
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
		[]tele.Command{HelpAdminChatCommand, BanCommand, KickCommand},
		tele.CommandScope{Type: tele.CommandScopeDefault, ChatID: config.AdministrationChatID})
	// Для админов админского чата
	setCommands(c,
		[]tele.Command{SetCommandsCommand, HelpAdminChatCommand, BanCommand, KickCommand},
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.AdministrationChatID})

	return nil
}

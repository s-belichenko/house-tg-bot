package handlers

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v4"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

const (
	muteCommandFormat   = `/mute &lt;user_id&gt; [days] (от 1 до 366, иначе бессрочно)`
	unmuteCommandFormat = `/unmute &lt;user_id&gt;`
	banCommandFormat    = `/ban &lt;user_id&gt; [days] (от 1 до 366, иначе бессрочно)`
	unbanCommandFormat  = `/unban &lt;user_id&gt;`
)

var (
	HelpAdminChatCommand = tele.Command{
		Text:        "help_admin",
		Description: "Справка по боту для админов",
	}
	MuteCommand = tele.Command{
		Text:        "mute",
		Description: "Ограничить пользователя в домовом чате",
	}
	UnmuteCommand = tele.Command{
		Text:        "unmute",
		Description: "Снять ограничения с пользователя в домовом чате",
	}
	BanCommand = tele.Command{
		Text:        "ban",
		Description: "Заблокировать пользователя из домового чата",
	}
	UnbanCommand = tele.Command{
		Text:        "unban",
		Description: "Разблокировать пользователя из домового чата",
	}
)

func CommandMuteHandler(ctx tele.Context) error {
	var violator *tele.ChatMember

	d := ctx.Data()

	fields := strings.Fields(d)
	switch len(fields) {
	case 1:
		if user := createUserViolator(fields[0]); user != nil {
			violator = &tele.ChatMember{
				User:   user,
				Rights: tele.NoRights(),
			}
		}
	case 2:
		if user := createUserViolator(fields[0]); user != nil {
			violator = &tele.ChatMember{
				User:            user,
				Rights:          tele.NoRights(),
				RestrictedUntil: createUnixTimeFromDays(fields[1]),
			}
		}
	}

	if violator == nil {
		if err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", muteCommandFormat), tele.ModeHTML); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /mute: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	if err := ctx.Bot().Restrict(&tele.Chat{ID: config.HouseChatID}, violator); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось ограничить пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		if err := ctx.Reply("Не удалось ограничить пользователя."); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось ограничить пользователя: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	if err := ctx.Reply("Пользователь ограничен."); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось уведомить что пользователь ограничен: %v", err),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}
	// FIXME: Посылать сообщение пользователю об ограничениях? А если он не начал общение с ботом?
	pkgLog.Info("Пользователь ограничен", pkgLogger.LogContext{
		"admin_id":         ctx.Message().Sender.ID,
		"admin_username":   ctx.Message().Sender.Username,
		"admin_first_name": ctx.Message().Sender.FirstName,
		"admin_last_name":  ctx.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func CommandUnmuteHandler(ctx tele.Context) error {
	var violator *tele.ChatMember

	d := ctx.Data()

	f := strings.Fields(d)
	if len(f) == 1 {
		if user := createUserViolator(f[0]); user != nil {
			violator = &tele.ChatMember{User: user, Rights: tele.NoRestrictions()}
		}
	}

	if violator == nil {
		if err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", unmuteCommandFormat), tele.ModeHTML); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /unmute: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	if err := ctx.Bot().Promote(&tele.Chat{ID: config.HouseChatID}, violator); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		if err := ctx.Reply("Не удалось снять ограничения с пользователя."); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	if err := ctx.Reply("Сняты ограничения с пользователя."); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось уведомить что сняты ограничения с пользователя: %v", err),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}
	// FIXME: Посылать сообщение пользователю об отмене бана? А если он не начал общение с ботом?
	pkgLog.Info("Сняты ограничения с пользователя", pkgLogger.LogContext{
		"admin_id":         ctx.Message().Sender.ID,
		"admin_username":   ctx.Message().Sender.Username,
		"admin_first_name": ctx.Message().Sender.FirstName,
		"admin_last_name":  ctx.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func CommandBanHandler(ctx tele.Context) error {
	var violator *tele.ChatMember

	data := ctx.Data()

	fields := strings.Fields(data)
	switch len(fields) {
	case 1:
		if user := createUserViolator(fields[0]); user != nil {
			violator = &tele.ChatMember{User: user, RestrictedUntil: tele.Forever()}
		}
	case 2:
		if user := createUserViolator(fields[0]); user != nil {
			violator = &tele.ChatMember{
				User:            user,
				RestrictedUntil: createUnixTimeFromDays(fields[1]),
			}
		}
	}

	if violator == nil {
		pkgLog.Warn("Вызов команды /ban без аргументов", pkgLogger.LogContext{
			"arguments_string": data,
		})

		if err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /ban: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	if err := ctx.Bot().Ban(&tele.Chat{ID: config.HouseChatID}, violator); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось заблокировать пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		if err := ctx.Reply("Не удалось заблокировать пользователя."); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось заблокировать пользователя: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	if err := ctx.Reply("Пользователь заблокирован."); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось уведомить что пользователь заблокирован: %v", err),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}

	pkgLog.Info("Пользователь заблокирован", pkgLogger.LogContext{
		"admin_id":         ctx.Message().Sender.ID,
		"admin_username":   ctx.Message().Sender.Username,
		"admin_first_name": ctx.Message().Sender.FirstName,
		"admin_last_name":  ctx.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func CommandUnbanHandler(ctx tele.Context) error {
	var violator *tele.User

	data := ctx.Data()

	f := strings.Fields(data)
	if len(f) == 1 {
		violator = createUserViolator(f[0])
	}

	if violator == nil {
		pkgLog.Warn("Вызов команды /unban без аргументов", pkgLogger.LogContext{
			"arguments_string": data,
		})

		if err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /unban: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}
	}

	if err := ctx.Bot().Unban(&tele.Chat{ID: config.HouseChatID}, violator, true); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось разблокировать пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		if err := ctx.Reply("Не удалось разблокировать пользователя."); err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось разблокировать пользователя: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	if err := ctx.Reply("Пользователь разблокирован."); err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось уведомить что пользователь разблокирован: %v", err),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}

	pkgLog.Info("Пользователь разблокирован", pkgLogger.LogContext{
		"admin_id":         ctx.Message().Sender.ID,
		"admin_username":   ctx.Message().Sender.Username,
		"admin_first_name": ctx.Message().Sender.FirstName,
		"admin_last_name":  ctx.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func CommandHelpAdminHandler(ctx tele.Context) error {
	help := fmt.Sprintf(`
Справка для администратора. Все команды ниже используются только в чате администраторов.

Команды:

/help_admin – Текущая справка.
%s – Ограничить пользователя в домовом чате. 
%s – Снять ограничения с пользователя в домовом чате.
%s – Забанить пользователя в домовом чата.
%s – Разбанить пользователя в домовом чата.

<a href="`+config.RulesURL.String()+`">Ссылка на правила</a>.`,
		muteCommandFormat,
		unmuteCommandFormat,
		banCommandFormat,
		unbanCommandFormat,
	)

	if err := ctx.Send(help, tele.ModeHTML, tele.NoPreview); err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}

	return nil
}

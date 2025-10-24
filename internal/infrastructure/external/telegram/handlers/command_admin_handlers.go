package handlers

import (
	"fmt"
	"html/template"
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

	var err error

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
		err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", muteCommandFormat), tele.ModeHTML)
		if err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /mute: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	err = ctx.Bot().Restrict(&tele.Chat{ID: cfg.HouseChatID}, violator)
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось ограничить пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		err := ctx.Reply("Не удалось ограничить пользователя.")
		if err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось ограничить пользователя: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	err = ctx.Reply("Пользователь ограничен.")
	if err != nil {
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

	var err error

	d := ctx.Data()

	f := strings.Fields(d)
	if len(f) == 1 {
		if user := createUserViolator(f[0]); user != nil {
			violator = &tele.ChatMember{User: user, Rights: tele.NoRestrictions()}
		}
	}

	if violator == nil {
		err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", unmuteCommandFormat), tele.ModeHTML)
		if err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /unmute: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	err = ctx.Bot().Promote(&tele.Chat{ID: cfg.HouseChatID}, violator)
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		err := ctx.Reply("Не удалось снять ограничения с пользователя.")
		if err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	err = ctx.Reply("Сняты ограничения с пользователя.")
	if err != nil {
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

	var err error

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

		err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML)
		if err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /ban: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	err = ctx.Bot().Ban(&tele.Chat{ID: cfg.HouseChatID}, violator)
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось заблокировать пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		err := ctx.Reply("Не удалось заблокировать пользователя.")
		if err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось заблокировать пользователя: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	err = ctx.Reply("Пользователь заблокирован.")
	if err != nil {
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

	var err error

	data := ctx.Data()

	f := strings.Fields(data)
	if len(f) == 1 {
		violator = createUserViolator(f[0])
	}

	if violator == nil {
		pkgLog.Warn("Вызов команды /unban без аргументов", pkgLogger.LogContext{
			"arguments_string": data,
		})

		err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML)
		if err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /unban: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}
	}

	err = ctx.Bot().Unban(&tele.Chat{ID: cfg.HouseChatID}, violator, true)
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf("Не удалось разблокировать пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		err := ctx.Reply("Не удалось разблокировать пользователя.")
		if err != nil {
			pkgLog.Error(
				fmt.Sprintf("Не удалось разблокировать пользователя: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	err = ctx.Reply("Пользователь разблокирован.")
	if err != nil {
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
	err := ctx.Send(
		renderingTool.RenderText(`help_admin.gohtml`, struct {
			HelpAdminCommand    string
			MuteCommandFormat   string
			UnmuteCommandFormat string
			BanCommandFormat    string
			UnbanCommandFormat  string
			RulesURL            template.URL
		}{
			HelpAdminCommand:    HelpAdminChatCommand.Text,
			MuteCommandFormat:   muteCommandFormat,
			UnmuteCommandFormat: unmuteCommandFormat,
			BanCommandFormat:    banCommandFormat,
			UnbanCommandFormat:  unbanCommandFormat,
			RulesURL:            template.URL(cfg.RulesURL.String()),
		}),
		tele.ModeHTML,
		tele.NoPreview,
	)
	if err != nil {
		pkgLog.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}

	return nil
}

func CallbackJoinHandler(ctx tele.Context) error {
	err := ctx.Bot().ApproveJoinRequest(&tele.User{ID: 0}, &tele.User{})
	if err != nil {
		pkgLog.Error(
			fmt.Sprintf(
				`Не удалось одобрить заявку на вступление в чат пользователя %d: %e`,
				ctx.Message().Sender.ID, err,
			),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}

	return nil
}

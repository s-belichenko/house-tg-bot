package handlers

import (
	"fmt"
	"html/template"
	"strings"

	"s-belichenko/house-tg-bot/internal/config"

	template2 "s-belichenko/house-tg-bot/pkg/template"
	"s-belichenko/house-tg-bot/pkg/time"

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

type CommandAdminHandlers struct {
	config        config.App
	renderingTool template2.RenderingTool
	logger        pkgLogger.Logger
}

func NewCommandAdminHandlers(cfg config.App, logger pkgLogger.Logger) *CommandAdminHandlers {
	renderingTool := template2.NewTool("handlers", logger)

	return &CommandAdminHandlers{
		config:        cfg,
		renderingTool: renderingTool,
		logger:        logger,
	}
}

func (h *CommandAdminHandlers) CommandMuteHandler(ctx tele.Context) error {
	var violator *tele.ChatMember

	var err error

	d := ctx.Data()

	fields := strings.Fields(d)
	user := createUserViolator(fields[0])
	if user == nil {
		err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", muteCommandFormat), tele.ModeHTML)
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /mute: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}
	violator = &tele.ChatMember{
		User:   user,
		Rights: tele.NoRights(),
	}
	if len(fields) == 2 {
		restrictedUntil, err := time.CreateUnixTimeFromDays(fields[1])
		if err != nil {
			return err
		}
		violator.RestrictedUntil = restrictedUntil
	}

	err = ctx.Bot().Restrict(&tele.Chat{ID: int64(h.config.HouseChatID)}, violator)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось ограничить пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		err := ctx.Reply("Не удалось ограничить пользователя.")
		if err != nil {
			h.logger.Error(
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
		h.logger.Error(
			fmt.Sprintf("Не удалось уведомить что пользователь ограничен: %v", err),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}
	// FIXME: Посылать сообщение пользователю об ограничениях? А если он не начал общение с ботом?
	h.logger.Info("Пользователь ограничен", pkgLogger.LogContext{
		"admin_id":         ctx.Message().Sender.ID,
		"admin_username":   ctx.Message().Sender.Username,
		"admin_first_name": ctx.Message().Sender.FirstName,
		"admin_last_name":  ctx.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func (h *CommandAdminHandlers) CommandUnmuteHandler(ctx tele.Context) error {
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
			h.logger.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /unmute: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}

	err = ctx.Bot().Promote(&tele.Chat{ID: int64(h.config.HouseChatID)}, violator)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось снять ограничения с пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		err := ctx.Reply("Не удалось снять ограничения с пользователя.")
		if err != nil {
			h.logger.Error(
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
		h.logger.Error(
			fmt.Sprintf("Не удалось уведомить что сняты ограничения с пользователя: %v", err),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}
	// FIXME: Посылать сообщение пользователю об отмене бана? А если он не начал общение с ботом?
	h.logger.Info("Сняты ограничения с пользователя", pkgLogger.LogContext{
		"admin_id":         ctx.Message().Sender.ID,
		"admin_username":   ctx.Message().Sender.Username,
		"admin_first_name": ctx.Message().Sender.FirstName,
		"admin_last_name":  ctx.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func (h *CommandAdminHandlers) CommandBanHandler(ctx tele.Context) error {
	var violator *tele.ChatMember

	var err error

	data := ctx.Data()

	fields := strings.Fields(data)
	user := createUserViolator(fields[0])
	if user == nil {
		h.logger.Warn("Вызов команды /ban без аргументов", pkgLogger.LogContext{
			"arguments_string": data,
		})

		err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML)
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /ban: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}
	violator = &tele.ChatMember{User: user, RestrictedUntil: tele.Forever()}
	if len(fields) == 2 {
		restrictedUntil, err := time.CreateUnixTimeFromDays(fields[1])
		if err != nil {
			return err
		}
		violator.RestrictedUntil = restrictedUntil
	}

	err = ctx.Bot().Ban(&tele.Chat{ID: int64(h.config.HouseChatID)}, violator)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось заблокировать пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		err := ctx.Reply("Не удалось заблокировать пользователя.")
		if err != nil {
			h.logger.Error(
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
		h.logger.Error(
			fmt.Sprintf("Не удалось уведомить что пользователь заблокирован: %v", err),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}

	h.logger.Info("Пользователь заблокирован", pkgLogger.LogContext{
		"admin_id":         ctx.Message().Sender.ID,
		"admin_username":   ctx.Message().Sender.Username,
		"admin_first_name": ctx.Message().Sender.FirstName,
		"admin_last_name":  ctx.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func (h *CommandAdminHandlers) CommandUnbanHandler(ctx tele.Context) error {
	var violator *tele.User

	var err error

	data := ctx.Data()

	f := strings.Fields(data)
	if len(f) == 1 {
		violator = createUserViolator(f[0])
	}

	if violator == nil {
		h.logger.Warn("Вызов команды /unban без аргументов", pkgLogger.LogContext{
			"arguments_string": data,
		})

		err := ctx.Reply(fmt.Sprintf("Верный формат команды: %s", banCommandFormat), tele.ModeHTML)
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /unban: %v", err),
				pkgLogger.LogContext{
					"message": ctx.Message(),
				},
			)
		}
	}

	err = ctx.Bot().Unban(&tele.Chat{ID: int64(h.config.HouseChatID)}, violator, true)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось разблокировать пользователя: %v", err),
			pkgLogger.LogContext{
				"violator": violator,
				"message":  ctx.Message(),
			},
		)

		err := ctx.Reply("Не удалось разблокировать пользователя.")
		if err != nil {
			h.logger.Error(
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
		h.logger.Error(
			fmt.Sprintf("Не удалось уведомить что пользователь разблокирован: %v", err),
			pkgLogger.LogContext{
				"message": ctx.Message(),
			},
		)
	}

	h.logger.Info("Пользователь разблокирован", pkgLogger.LogContext{
		"admin_id":         ctx.Message().Sender.ID,
		"admin_username":   ctx.Message().Sender.Username,
		"admin_first_name": ctx.Message().Sender.FirstName,
		"admin_last_name":  ctx.Message().Sender.LastName,
		"violator":         violator,
	})

	return nil
}

func (h *CommandAdminHandlers) CommandHelpAdminHandler(ctx tele.Context) error {
	err := ctx.Send(
		h.renderingTool.RenderText(`help_admin.gohtml`, struct {
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
			RulesURL:            template.URL(h.config.RulesURL.String()),
		}),
		tele.ModeHTML,
		tele.NoPreview,
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}

	return nil
}

func (h *CommandAdminHandlers) CallbackJoinHandler(ctx tele.Context) error {
	err := ctx.Bot().ApproveJoinRequest(&tele.User{ID: 0}, &tele.User{})
	if err != nil {
		h.logger.Error(
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

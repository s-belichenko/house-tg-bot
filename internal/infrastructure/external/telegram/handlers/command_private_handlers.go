package handlers

import (
	"fmt"
	"html/template"
	"strings"
	stdTime "time"

	"s-belichenko/house-tg-bot/pkg/time"

	"s-belichenko/house-tg-bot/internal/config"

	template2 "s-belichenko/house-tg-bot/pkg/template"

	"s-belichenko/house-tg-bot/pkg/logger"

	tele "gopkg.in/telebot.v4"
)

const (
	selfBanCommandFormat = `/self_mute [days] (от 1 до 366, по умолчанию 7)`
)

// Команды бота для административного чата.
var (
	StartCommand   = tele.Command{Text: "start", Description: "Начать работу с ботом"}
	HelpCommand    = tele.Command{Text: "help", Description: "Справка по боту"}
	MyInfoCommand  = tele.Command{Text: "my_info", Description: "Информация о вас"}
	SelfBanCommand = tele.Command{Text: "self_mute", Description: "Временно ограничить себя в домовом чате"}
)

type CommandPrivateHandlers struct {
	config        config.App
	renderingTool template2.RenderingTool
	logger        logger.Logger
}

func NewCommandPrivateHandlers(cfg config.App, logger logger.Logger) *CommandPrivateHandlers {
	renderingTool := template2.NewTool("handlers", logger)

	return &CommandPrivateHandlers{
		config:        cfg,
		renderingTool: renderingTool,
		logger:        logger,
	}
}

func (h *CommandPrivateHandlers) CommandStartHandler(ctx tele.Context) error {
	greetingName, err := GetGreetingName(ctx.Sender())
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Не удалось сформировать обращение к пользователю %d", ctx.Sender().ID), nil)
	}
	err = ctx.Send(fmt.Sprintf("Привет, %s! Ознакомься со справкой по работе с ботом: /help", greetingName))
	if err != nil {
		h.logger.Error(fmt.Sprintf("Не удалось отправить ответ на команду /start: %v", err), nil)
	}

	return err
}

func (h *CommandPrivateHandlers) CommandSelfMuteHandler(ctx tele.Context) error {
	var (
		neighbour *tele.ChatMember
		err       error
	)

	if !h.userCanSelfBan(ctx) {
		return nil
	}

	d := ctx.Data()

	fields := strings.Fields(d)
	user := &tele.User{ID: ctx.Sender().ID}
	neighbour = &tele.ChatMember{
		User:   user,
		Rights: tele.NoRights(),
	}

	var days = "7"
	if len(fields) == 1 {
		days = fields[0]
	}

	restrictedUntil, err := time.CreateUnixTimeFromDays(days)
	if err != nil {
		err = ctx.Reply(fmt.Sprintf("Верный формат команды: %s", selfBanCommandFormat), tele.ModeHTML)
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось отправить подсказку по команде /%s: %v", SelfBanCommand.Text, err),
				logger.LogContext{
					"message": ctx.Message(),
				},
			)
		}

		return nil
	}
	neighbour.RestrictedUntil = restrictedUntil

	err = ctx.Bot().Restrict(&tele.Chat{ID: int64(h.config.HouseChatID)}, neighbour)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf(
				"Не удалось самоограничить пользователя %d на %s дней: %v",
				ctx.Sender().ID,
				days,
				err,
			),
			logger.LogContext{
				"message":   ctx.Message(),
				"neighbour": neighbour,
			},
		)

		err = ctx.Reply("Не удалось ограничить самого себя.")
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось уведомить, что пользователь не смог ограничить самого себя: %v", err),
				logger.LogContext{
					"message":   ctx.Message(),
					"neighbour": neighbour,
				},
			)
		}

		return nil
	}

	err = ctx.Reply(fmt.Sprintf("Вы ограничили самого себя на %s дней.", days))
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Не удалось уведомить что пользователь ограничен: %v", err),
			logger.LogContext{
				"message":   ctx.Message(),
				"neighbour": neighbour,
			},
		)
	}

	h.logger.Info(
		fmt.Sprintf("Пользователь %d самоограничен на %s дней.", ctx.Sender().ID, days),
		logger.LogContext{
			"message":   ctx.Message(),
			"neighbour": neighbour,
		},
	)

	return nil
}

func (h *CommandPrivateHandlers) CommandMyInfoHandler(ctx tele.Context) error {
	var chatMember *tele.ChatMember
	var err error

	chatMember, err = ctx.Bot().ChatMemberOf(
		&tele.Chat{ID: int64(h.config.HouseChatID)},
		&tele.User{ID: ctx.Sender().ID},
	)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf(`Не удалось получить информацию об участнике чата %d: %e`, ctx.Message().Sender.ID, err),
			logger.LogContext{"message": ctx.Message()},
		)

		return nil
	}

	h.logger.Info(
		fmt.Sprintf(`Получена информация об участнике чата %d`, ctx.Message().Sender.ID),
		logger.LogContext{"chat_member": chatMember},
	)

	var memberStatus string
	switch chatMember.Role {
	case tele.Creator:
		memberStatus = "создатель"
	case tele.Administrator:
		memberStatus = "администратор"
	case tele.Member:
		memberStatus = "участник"
	case tele.Restricted:
		memberStatus = "ограниченный"
	case tele.Left:
		memberStatus = "покинул чат"
	case tele.Kicked:
		memberStatus = "удаленный"
	}

	err = ctx.Reply(fmt.Sprintf(`Ваш статус в чате: <b>%s</b>.`, memberStatus), tele.ModeHTML)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf(`Не удалось ответить на команду /my_info пользователю %d: %e`, ctx.Message().Sender.ID, err),
			logger.LogContext{"chat_member": chatMember},
		)

		return nil
	}

	return nil
}

func (h *CommandPrivateHandlers) CommandHelpHandler(ctx tele.Context) error {
	err := ctx.Send(
		h.renderingTool.RenderEscapedText(
			`help.gohtml`,
			struct {
				InviteURL      template.URL
				HomeAddress    string
				StartCommand   string
				HelpCommand    string
				MyInfoCommand  string
				SelfBanCommand string
				RulesCommand   string
				KeysCommand    string
				ReportCommand  string
				RulesURL       template.URL
			}{
				InviteURL:      template.URL(h.config.InviteURL.String()),
				HomeAddress:    h.config.HomeAddress,
				StartCommand:   StartCommand.Text,
				HelpCommand:    HelpCommand.Text,
				MyInfoCommand:  MyInfoCommand.Text,
				SelfBanCommand: SelfBanCommand.Text,
				KeysCommand:    KeysCommand.Text,
				ReportCommand:  ReportCommand.Text,
				RulesCommand:   RulesCommand.Text,
				RulesURL:       template.URL(h.config.RulesURL.String()),
			}, []string{}),
		tele.ModeHTML,
		tele.NoPreview,
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}

	return nil
}

func (h *CommandPrivateHandlers) getChatMember(ctx tele.Context, user *tele.User) *tele.ChatMember {
	chatMember, err := ctx.Bot().ChatMemberOf(
		&tele.Chat{ID: int64(h.config.HouseChatID)},
		&tele.User{ID: user.ID},
	)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf(`Не удалось получить информацию об участнике чата %d: %e`, ctx.Message().Sender.ID, err),
			logger.LogContext{"message": ctx.Message()},
		)

		return nil
	}

	return chatMember
}

func (h *CommandPrivateHandlers) userCanSelfBan(ctx tele.Context) bool {
	var (
		neighbour *tele.ChatMember
		err       error
	)

	chatMember := h.getChatMember(ctx, ctx.Sender())
	if chatMember == nil {
		return false
	}

	if !IsHomeChatMember(chatMember) {
		h.logger.Warn(
			fmt.Sprintf(
				"Попытка самоограничиться пользователем %d, которого нет в домовом чате.",
				ctx.Sender().ID,
			),
			logger.LogContext{
				"message":     ctx.Message(),
				"chat_member": chatMember,
			},
		)

		return false
	}

	if chatMember.RestrictedUntil != 0 {
		h.logger.Warn(
			fmt.Sprintf(
				"Попытка самоограничиться пользователем %d, который уже самоограничен до %d.",
				ctx.Sender().ID,
				chatMember.RestrictedUntil,
			),
			logger.LogContext{
				"message":     ctx.Message(),
				"chat_member": chatMember,
			},
		)

		loc, _ := stdTime.LoadLocation("Europe/Moscow")
		t := stdTime.Unix(chatMember.RestrictedUntil, 0).In(loc)
		err = ctx.Reply(fmt.Sprintf("Вы уже ограничены до %s.", t.Format("2006-01-02 15:04:05")))
		if err != nil {
			h.logger.Error(
				fmt.Sprintf("Не удалось уведомить что пользователь ограничен: %v", err),
				logger.LogContext{
					"message":     ctx.Message(),
					"neighbour":   neighbour,
					"chat_member": chatMember,
				},
			)
		}

		return false
	}

	return true
}

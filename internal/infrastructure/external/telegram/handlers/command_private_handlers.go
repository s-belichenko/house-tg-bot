package handlers

import (
	"fmt"
	"html/template"

	"s-belichenko/house-tg-bot/internal/config"

	template2 "s-belichenko/house-tg-bot/pkg/template"

	"s-belichenko/house-tg-bot/pkg/logger"

	tele "gopkg.in/telebot.v4"
)

// Команды бота для административного чата.
var (
	StartCommand  = tele.Command{Text: "start", Description: "Начать работу с ботом"}
	HelpCommand   = tele.Command{Text: "help", Description: "Справка по боту"}
	MyInfoCommand = tele.Command{Text: "my_info", Description: "Информация о вас"}
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
				InviteURL     template.URL
				HomeAddress   string
				StartCommand  string
				HelpCommand   string
				MyInfoCommand string
				RulesCommand  string
				KeysCommand   string
				ReportCommand string
				RulesURL      template.URL
			}{
				InviteURL:     template.URL(h.config.InviteURL.String()),
				HomeAddress:   h.config.HomeAddress,
				StartCommand:  StartCommand.Text,
				HelpCommand:   HelpCommand.Text,
				MyInfoCommand: MyInfoCommand.Text,
				KeysCommand:   KeysCommand.Text,
				ReportCommand: ReportCommand.Text,
				RulesCommand:  RulesCommand.Text,
				RulesURL:      template.URL(h.config.RulesURL.String()),
			}, []string{}),
		tele.ModeHTML,
		tele.NoPreview,
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Не удалось отправить текст справки: %v", err), nil)
	}

	return nil
}

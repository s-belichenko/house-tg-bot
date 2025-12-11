package handlers

import (
	"fmt"

	"s-belichenko/house-tg-bot/internal/config"

	"s-belichenko/house-tg-bot/pkg/logger"

	tele "gopkg.in/telebot.v4"
)

var (
	SetCommandsCommand = tele.Command{
		Text:        "set_commands",
		Description: "Установить команды бота",
	}
	DeleteCommandsCommand = tele.Command{
		Text:        "delete_commands",
		Description: "Удалить команды бота",
	}
)

type CommandServiceHandlers struct {
	config config.App
	logger logger.Logger
}

func NewCommandServiceHandlers(cfg config.App, logger logger.Logger) *CommandServiceHandlers {
	return &CommandServiceHandlers{
		config: cfg,
		logger: logger,
	}
}

func (h *CommandServiceHandlers) CommandSetCommandsHandler(ctx tele.Context) error {
	// По умолчанию
	h.setCommands(ctx,
		[]tele.Command{StartCommand},
		tele.CommandScope{Type: tele.CommandScopeDefault})
	// Для личных чатов со всеми подряд
	h.setCommands(ctx,
		[]tele.Command{StartCommand, HelpCommand, RulesCommand},
		tele.CommandScope{Type: tele.CommandScopeAllPrivateChats})

	var homeCommands []tele.Command
	if h.config.HouseIsCompleted {
		homeCommands = []tele.Command{ReportCommand, RulesCommand}
	} else {
		homeCommands = []tele.Command{KeysCommand, ReportCommand, RulesCommand}
	}
	// Для участников домового чата
	h.setCommands(ctx,
		homeCommands,
		tele.CommandScope{Type: tele.CommandScopeChat, ChatID: int64(h.config.HouseChatID)},
	)
	// Для админов домового чата
	h.setCommands(ctx,
		homeCommands,
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: int64(h.config.HouseChatID)})
	// Для участников админского чата
	h.setCommands(ctx,
		[]tele.Command{HelpAdminChatCommand, MuteCommand, UnmuteCommand, BanCommand, UnbanCommand},
		tele.CommandScope{Type: tele.CommandScopeChat, ChatID: int64(h.config.AdminChatID)})
	// Для админов админского чата
	h.setCommands(
		ctx,
		[]tele.Command{
			SetCommandsCommand,
			DeleteCommandsCommand,
			HelpAdminChatCommand,
			MuteCommand,
			UnmuteCommand,
			BanCommand,
			UnbanCommand,
		},
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: int64(h.config.AdminChatID)},
	)

	return nil
}

func (h *CommandServiceHandlers) CommandDeleteCommandsHandler(ctx tele.Context) error {
	// По умолчанию
	h.deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeDefault})
	// Для личных чатов со всеми подряд
	h.deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeAllPrivateChats})
	// Для участников домового чата
	h.deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeChat, ChatID: int64(h.config.HouseChatID)})
	// Для админов домового чата
	h.deleteCommands(
		ctx,
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: int64(h.config.HouseChatID)},
	)
	// Для участников админского чата
	h.deleteCommands(
		ctx,
		tele.CommandScope{Type: tele.CommandScopeChat, ChatID: int64(h.config.AdminChatID)},
	)
	// Для админов админского чата
	h.deleteCommands(
		ctx,
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: int64(h.config.AdminChatID)},
	)

	return nil
}

func (h *CommandServiceHandlers) setCommands(c TeleContext, commands []tele.Command, scope tele.CommandScope) {
	err := c.Bot().SetCommands(commands, scope)
	if err != nil {
		h.logger.Fatal(
			fmt.Sprintf("Не удалось установить команды бота: %v", err),
			logger.LogContext{
				"commands": commands,
				"scope":    scope,
			},
		)
	} else {
		h.logger.Info("Успешно установлены команды бота", logger.LogContext{
			"commands": commands,
			"scope":    scope,
		})
	}
}

func (h *CommandServiceHandlers) deleteCommands(c TeleContext, scope tele.CommandScope) {
	err := c.Bot().DeleteCommands(scope)
	if err != nil {
		h.logger.Fatal(fmt.Sprintf("Не удалось удалить команды бота: %v", err), logger.LogContext{
			"scope": scope,
		})
	} else {
		h.logger.Info("Успешно удалены команды бота", logger.LogContext{
			"scope": scope,
		})
	}
}

package handlers

import tele "gopkg.in/telebot.v4"

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

func CommandSetCommandsHandler(ctx tele.Context) error {
	// По умолчанию
	setCommands(ctx,
		[]tele.Command{StartCommand},
		tele.CommandScope{Type: tele.CommandScopeDefault})
	// Для личных чатов со всеми подряд
	setCommands(ctx,
		[]tele.Command{StartCommand, HelpCommand, RulesCommand},
		tele.CommandScope{Type: tele.CommandScopeAllPrivateChats})
	// Для участников домового чата
	setCommands(ctx,
		[]tele.Command{KeysCommand, ReportCommand, RulesCommand},
		tele.CommandScope{Type: tele.CommandScopeChat, ChatID: config.HouseChatID})
	// Для админов домового чата
	setCommands(ctx,
		[]tele.Command{KeysCommand, ReportCommand, RulesCommand},
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.HouseChatID})
	// Для участников админского чата
	setCommands(ctx,
		[]tele.Command{HelpAdminChatCommand, MuteCommand, UnmuteCommand, BanCommand, UnbanCommand},
		tele.CommandScope{Type: tele.CommandScopeChat, ChatID: config.AdministrationChatID})
	// Для админов админского чата
	setCommands(
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
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.AdministrationChatID},
	)

	return nil
}

func CommandDeleteCommandsHandler(ctx tele.Context) error {
	// По умолчанию
	deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeDefault})
	// Для личных чатов со всеми подряд
	deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeAllPrivateChats})
	// Для участников домового чата
	deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeChat, ChatID: config.HouseChatID})
	// Для админов домового чата
	deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.HouseChatID})
	// Для участников админского чата
	deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeChat, ChatID: config.AdministrationChatID})
	// Для админов админского чата
	deleteCommands(ctx, tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.AdministrationChatID})

	return nil
}

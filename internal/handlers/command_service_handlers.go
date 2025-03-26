package handlers

import tele "gopkg.in/telebot.v4"

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
		tele.CommandScope{Type: tele.CommandScopeChat, ChatID: config.HouseChatId})
	// Для админов домового чата
	setCommands(c,
		[]tele.Command{KeysCommand, ReportCommand},
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.HouseChatId})
	// Для участников админского чата
	setCommands(c,
		[]tele.Command{HelpAdminChatCommand, MuteCommand, UnmuteCommand, BanCommand, UnbanCommand},
		tele.CommandScope{Type: tele.CommandScopeChat, ChatID: config.AdministrationChatID})
	// Для админов админского чата
	setCommands(c,
		[]tele.Command{SetCommandsCommand, DeleteCommandsCommand, HelpAdminChatCommand, MuteCommand, UnmuteCommand, BanCommand, UnbanCommand},
		tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.AdministrationChatID})

	return nil
}

func CommandDeleteCommandsHandler(c tele.Context) error {
	// По умолчанию
	deleteCommands(c, tele.CommandScope{Type: tele.CommandScopeDefault})
	// Для личных чатов со всеми подряд
	deleteCommands(c, tele.CommandScope{Type: tele.CommandScopeAllPrivateChats})
	// Для участников домового чата
	deleteCommands(c, tele.CommandScope{Type: tele.CommandScopeChat, ChatID: config.HouseChatId})
	// Для админов домового чата
	deleteCommands(c, tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.HouseChatId})
	// Для участников админского чата
	deleteCommands(c, tele.CommandScope{Type: tele.CommandScopeChat, ChatID: config.AdministrationChatID})
	// Для админов админского чата
	deleteCommands(c, tele.CommandScope{Type: tele.CommandScopeChatAdmin, ChatID: config.AdministrationChatID})

	return nil
}

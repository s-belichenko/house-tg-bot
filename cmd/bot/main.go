package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"s-belichenko/house-tg-bot/internal/infrastructure/external/telegram/handlers"

	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
	teleMid "gopkg.in/telebot.v4/middleware"

	mid "s-belichenko/house-tg-bot/internal/infrastructure/external/telegram/middleware"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

type Config struct {
	BotToken             string `env:"TELEGRAM_BOT_TOKEN"`
	AdministrationChatID int64  `env:"ADMINISTRATION_CHAT_ID"` // Чат администраторов, куда поступают уведомления и тп
	LogStreamName        string
}

var (
	bot    *tele.Bot
	pkgLog pkgLogger.Logger
	config = Config{LogStreamName: "main_stream"}
)

func init() {
	initModule()
	initBot()
	setBotMiddleware()
	registerBotCommandHandlers()
	registerJoinRequestHandler()
	registerMediaHandler()
}

func initModule() {
	pkgLog = pkgLogger.InitLog(config.LogStreamName)

	pkgLog.Debug("Start init module", nil)

	if err := cleanenv.ReadEnv(&config); err != nil {
		pkgLog.Fatal(fmt.Sprintf("Не удалось прочитать конфигурацию ота: %v", err), nil)
		os.Exit(1)
	}
}

func initBot() {
	var err error

	pkgLog.Debug("Start init bot", nil)

	bot, err = tele.NewBot(tele.Settings{Token: config.BotToken})
	if err != nil {
		pkgLog.Fatal(fmt.Sprintf("Не удалось инициализировать бота: %v", err), nil)
		os.Exit(1)
	}

	handlers.SetBotID(bot.Me.ID)
}

func setBotMiddleware() {
	pkgLog.Debug("Start use middleware bot", nil)
	bot.Use(teleMid.Recover(func(err error, _ tele.Context) {
		pkgLog.Fatal(fmt.Sprintf("Bot fatal: %v", err), nil)
	}))
	bot.Use(mid.GetLogUpdateMiddleware(pkgLog))
}

func registerBotCommandHandlers() {
	pkgLog.Debug("Start register command hndls", nil)
	// Общие команды
	bot.Handle(
		"/"+handlers.RulesCommand.Text,
		handlers.CommandRulesHandler,
		mid.CommonCommandMiddleware,
	)
	// Личные сообщения.
	bot.Handle(
		"/"+handlers.StartCommand.Text,
		handlers.CommandStartHandler,
		mid.AllPrivateChatsMiddleware,
	)
	bot.Handle(
		"/"+handlers.HelpCommand.Text,
		handlers.CommandHelpHandler,
		mid.AllPrivateChatsMiddleware,
	)
	// Домашний чат.
	bot.Handle(
		"/"+handlers.KeysCommand.Text,
		handlers.CommandKeysHandler,
		mid.HomeChatMiddleware,
		mid.KeysCommandMiddleware,
	)
	bot.Handle(
		"/"+handlers.ReportCommand.Text,
		handlers.CommandReportHandler,
		mid.HomeChatMiddleware,
	)
	// Административный чат (админы).
	bot.Handle(
		"/"+handlers.SetCommandsCommand.Text,
		handlers.CommandSetCommandsHandler,
		mid.AdminChatMiddleware,
	)
	bot.Handle(
		"/"+handlers.DeleteCommandsCommand.Text,
		handlers.CommandDeleteCommandsHandler,
		mid.AdminChatMiddleware,
	)
	// Административный чат (участники).
	bot.Handle(
		"/"+handlers.HelpAdminChatCommand.Text,
		handlers.CommandHelpAdminHandler,
		mid.AdminChatMiddleware,
	)
	bot.Handle("/"+handlers.MuteCommand.Text, handlers.CommandMuteHandler, mid.AdminChatMiddleware)
	bot.Handle(
		"/"+handlers.UnmuteCommand.Text,
		handlers.CommandUnmuteHandler,
		mid.AdminChatMiddleware,
	)
	bot.Handle("/"+handlers.BanCommand.Text, handlers.CommandBanHandler, mid.AdminChatMiddleware)
	bot.Handle(
		"/"+handlers.UnbanCommand.Text,
		handlers.CommandUnbanHandler,
		mid.AdminChatMiddleware,
	)
}

func registerJoinRequestHandler() {
	bot.Handle(tele.OnChatJoinRequest, handlers.JoinRequestHandler)
}

func registerMediaHandler() {
	// bot.Handle(tele.OnMedia, handlers.MediaHandler)
}

// Handler Функция-обработчик для Yandex Cloud Function.
func Handler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)

		return
	}

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Failed to read body", http.StatusInternalServerError)

		return
	}

	var update tele.Update
	if err = json.Unmarshal(bodyBytes, &update); err != nil {
		http.Error(writer, "Failed to parse update", http.StatusInternalServerError)

		return
	}

	go bot.ProcessUpdate(update)
}

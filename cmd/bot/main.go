package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	tele "gopkg.in/telebot.v4"
	teleMid "gopkg.in/telebot.v4/middleware"
	hndls "s-belichenko/house-tg-bot/internal/infrastructure/handlers"
	sec "s-belichenko/house-tg-bot/internal/security"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

type ConfigBot struct {
	BotToken             string `env:"TELEGRAM_BOT_TOKEN"`
	AdministrationChatID int64  `env:"ADMINISTRATION_CHAT_ID"` // Чат администраторов, куда поступают уведомления и тп
	LogStreamName        string
}

var (
	bot    *tele.Bot
	pkgLog pkgLogger.Logger
	config = ConfigBot{LogStreamName: "main_stream"}
)

func init() {
	initModule()
	initBot()
	setBotMiddleware()
	registerBotCommandHandlers()
	registerJoinRequestHandler()
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

	hndls.SetBotID(bot.Me.ID)
}

func setBotMiddleware() {
	pkgLog.Debug("Start use middleware bot", nil)
	bot.Use(teleMid.Recover(func(err error, _ tele.Context) {
		pkgLog.Fatal(fmt.Sprintf("Bot fatal: %v", err), nil)
	}))
	bot.Use(pkgLogger.GetMiddleware(pkgLog))
}

func registerBotCommandHandlers() {
	pkgLog.Debug("Start register command hndls", nil)
	// Личные сообщения.
	bot.Handle("/"+hndls.StartCommand.Text, hndls.CommandStartHandler, sec.AllPrivateChatsMiddleware)
	bot.Handle("/"+hndls.HelpCommand.Text, hndls.CommandHelpHandler, sec.AllPrivateChatsMiddleware)
	// Домашний чат.
	bot.Handle("/"+hndls.KeysCommand.Text, hndls.CommandKeysHandler, sec.HomeChatMiddleware, sec.KeysCommandMiddleware)
	bot.Handle("/"+hndls.ReportCommand.Text, hndls.CommandReportHandler, sec.HomeChatMiddleware)
	// Административный чат (админы).
	bot.Handle("/"+hndls.SetCommandsCommand.Text, hndls.CommandSetCommandsHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hndls.DeleteCommandsCommand.Text, hndls.CommandDeleteCommandsHandler, sec.AdminChatMiddleware)
	// Административный чат (участники).
	bot.Handle("/"+hndls.HelpAdminChatCommand.Text, hndls.CommandHelpAdminHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hndls.MuteCommand.Text, hndls.CommandMuteHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hndls.UnmuteCommand.Text, hndls.CommandUnmuteHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hndls.BanCommand.Text, hndls.CommandBanHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hndls.UnbanCommand.Text, hndls.CommandUnbanHandler, sec.AdminChatMiddleware)
}

func registerJoinRequestHandler() {
	bot.Handle(tele.OnChatJoinRequest, hndls.JoinRequestHandler)
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

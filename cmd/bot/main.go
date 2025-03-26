package main

import (
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"io"
	"net/http"
	"os"

	tele "gopkg.in/telebot.v4"
	teleMid "gopkg.in/telebot.v4/middleware"
	hdls "s-belichenko/ilovaiskaya2-bot/internal/handlers"
	intLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
	sec "s-belichenko/ilovaiskaya2-bot/internal/security"
)

type ConfigBot struct {
	BotToken             string `env:"TELEGRAM_BOT_TOKEN"`
	AdministrationChatId int64  `env:"ADMINISTRATION_CHAT_ID"` // Чат администраторов, куда поступают уведомления и тп
	LogStreamName        string
}

var (
	bot    *tele.Bot
	log    intLog.Logger
	config = ConfigBot{LogStreamName: "main_stream"}
)

func init() {
	initModule()
	initBot()
	setBotMiddleware()
	registerBotCommandHandlers()
	registerJoinRequestHandler()
}

func registerJoinRequestHandler() {
	bot.Handle(tele.OnChatJoinRequest, hdls.JoinRequestHandler)
}

func initModule() {
	err := cleanenv.ReadEnv(&config)
	log = intLog.InitLog(config.LogStreamName)
	log.Debug("Start init module", nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Не удалось прочитать конфигурацию ота: %v", err), nil)
		os.Exit(1)
	}
}

func initBot() {
	var err error
	log.Debug("Start init bot", nil)

	bot, err = tele.NewBot(tele.Settings{Token: config.BotToken})
	if err != nil {
		log.Fatal(fmt.Sprintf("Не удалось инициализировать бота: %v", err), nil)
		os.Exit(1)
	}
	hdls.SetBotID(bot.Me.ID)
}

func setBotMiddleware() {
	log.Debug("Start use middleware bot", nil)
	bot.Use(teleMid.Recover(func(err error, context tele.Context) {
		log.Fatal(fmt.Sprintf("Bot fatal: %v", err), nil)
	}))
	bot.Use(intLog.GetMiddleware(log))
}

func registerBotCommandHandlers() {
	log.Debug("Start register command handlers", nil)
	// Личные сообщения
	bot.Handle("/"+hdls.StartCommand.Text, hdls.CommandStartHandler, sec.AllPrivateChatsMiddleware)
	bot.Handle("/"+hdls.HelpCommand.Text, hdls.CommandHelpHandler, sec.AllPrivateChatsMiddleware)
	// Домашний чат
	bot.Handle("/"+hdls.KeysCommand.Text, hdls.CommandKeysHandler, sec.HomeChatMiddleware, sec.KeysCommandMiddleware)
	bot.Handle("/"+hdls.ReportCommand.Text, hdls.CommandReportHandler, sec.HomeChatMiddleware)
	// Административный чат
	bot.Handle("/"+hdls.SetCommandsCommand.Text, hdls.CommandSetCommandsHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hdls.DeleteCommandsCommand.Text, hdls.CommandDeleteCommandsHandler, sec.AdminChatMiddleware)

	bot.Handle("/"+hdls.HelpAdminChatCommand.Text, hdls.CommandHelpAdminHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hdls.MuteCommand.Text, hdls.CommandMuteHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hdls.UnmuteCommand.Text, hdls.CommandUnmuteHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hdls.BanCommand.Text, hdls.CommandBanHandler, sec.AdminChatMiddleware)
	bot.Handle("/"+hdls.UnbanCommand.Text, hdls.CommandUnbanHandler, sec.AdminChatMiddleware)
}

// Handler Функция-обработчик для Yandex Cloud Function
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	var update tele.Update
	err = json.Unmarshal(bodyBytes, &update)
	if err != nil {
		http.Error(w, "Failed to parse update", http.StatusInternalServerError)
		return
	}

	bot.ProcessUpdate(update)
}

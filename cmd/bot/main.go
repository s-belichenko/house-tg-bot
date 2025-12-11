package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"s-belichenko/house-tg-bot/internal/config"

	"s-belichenko/house-tg-bot/internal/domain/models"
	llm2 "s-belichenko/house-tg-bot/internal/infrastructure/adapters/llm"
	"s-belichenko/house-tg-bot/internal/infrastructure/external/telegram/handlers"
	"s-belichenko/house-tg-bot/internal/infrastructure/external/telegram/middleware"
	"s-belichenko/house-tg-bot/internal/utils"
	pkgTime "s-belichenko/house-tg-bot/pkg/time"

	tele "gopkg.in/telebot.v4"
	teleMid "gopkg.in/telebot.v4/middleware"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

func initLog(logStreamName string) pkgLogger.Logger {
	logger := log.New(os.Stdout, "", 0)
	time := pkgTime.Time{} // Отключаем все флаги.

	return pkgLogger.NewYandexLogger(logStreamName, logger, time)
}

func initAI(logger pkgLogger.Logger, cfg config.App) models.AI {
	llmService := llm2.NewYandexLLM(cfg, logger)

	return models.NewAI(llmService, logger)
}

func initBot(logger pkgLogger.Logger, cfg config.App) *tele.Bot {
	var err error

	logger.Debug("Start init bot", nil)

	bot, err := tele.NewBot(tele.Settings{Token: cfg.BotToken})
	if err != nil {
		logger.Fatal(fmt.Sprintf("Не удалось инициализировать бота: %v", err), nil)
		os.Exit(1)
	}

	return bot
}

func setBotMiddleware(bot *tele.Bot, mid *middleware.TelebotMiddleware, logger pkgLogger.Logger) {
	logger.Debug("Start use middleware bot", nil)
	bot.Use(
		teleMid.Recover(
			func(err error, _ tele.Context) {
				logger.Fatal(
					fmt.Sprintf("Bot fatal: %v", err),
					pkgLogger.LogContext{"stack_trace": utils.GetStackTraceAsJSON(debug.Stack())},
				)
			},
		),
	)
	bot.Use(mid.GetLogUpdateMiddleware(logger))
}

func registerBotCommandHandlers(
	bot *tele.Bot,
	logger pkgLogger.Logger,
	mid *middleware.TelebotMiddleware,
	cfg config.App,
	ai models.AI,
) {
	logger.Debug("Start register command handlers", nil)

	houseHandlers := handlers.NewCommandHouseHandlers(cfg, ai, logger)
	adminHandlers := handlers.NewCommandAdminHandlers(cfg, logger)
	privateHandlers := handlers.NewCommandPrivateHandlers(cfg, logger)
	serviceHandlers := handlers.NewCommandServiceHandlers(cfg, logger)

	// Общие команды
	bot.Handle(
		"/"+handlers.RulesCommand.Text,
		houseHandlers.CommandRulesHandler,
		mid.CommonCommandMiddleware,
	)
	// Личные сообщения.
	bot.Handle(
		"/"+handlers.StartCommand.Text,
		privateHandlers.CommandStartHandler,
		mid.AllPrivateChatsMiddleware,
	)
	bot.Handle(
		"/"+handlers.HelpCommand.Text,
		privateHandlers.CommandHelpHandler,
		mid.AllPrivateChatsMiddleware,
	)
	bot.Handle(
		"/"+handlers.MyInfoCommand.Text,
		privateHandlers.CommandMyInfoHandler,
		mid.AllPrivateChatsMiddleware,
	)
	// Домашний чат.
	bot.Handle(
		"/"+handlers.KeysCommand.Text,
		houseHandlers.CommandKeysHandler,
		mid.HomeChatMiddleware,
		mid.KeysCommandMiddleware,
	)
	bot.Handle(
		"/"+handlers.ReportCommand.Text,
		houseHandlers.CommandReportHandler,
		mid.HomeChatMiddleware,
	)
	// Административный чат (админы).
	bot.Handle(
		"/"+handlers.SetCommandsCommand.Text,
		serviceHandlers.CommandSetCommandsHandler,
		mid.AdminChatMiddleware,
	)
	bot.Handle(
		"/"+handlers.DeleteCommandsCommand.Text,
		serviceHandlers.CommandDeleteCommandsHandler,
		mid.AdminChatMiddleware,
	)
	// Административный чат (участники).
	bot.Handle(
		"/"+handlers.HelpAdminChatCommand.Text,
		adminHandlers.CommandHelpAdminHandler,
		mid.AdminChatMiddleware,
	)
	bot.Handle("/"+handlers.MuteCommand.Text, adminHandlers.CommandMuteHandler, mid.AdminChatMiddleware)
	bot.Handle(
		"/"+handlers.UnmuteCommand.Text,
		adminHandlers.CommandUnmuteHandler,
		mid.AdminChatMiddleware,
	)
	bot.Handle("/"+handlers.BanCommand.Text, adminHandlers.CommandBanHandler, mid.AdminChatMiddleware)
	bot.Handle(
		"/"+handlers.UnbanCommand.Text,
		adminHandlers.CommandUnbanHandler,
		mid.AdminChatMiddleware,
	)
}

func registerJoinRequestHandler(bot *tele.Bot, cfg config.App, logger pkgLogger.Logger) {
	logger.Debug("Start register join request handlers", nil)

	joinRequestHandlers := handlers.NewJoinRequestHandlersHandlers(cfg, logger)
	bot.Handle(tele.OnChatJoinRequest, joinRequestHandlers.JoinRequestHandler)
}

func registerMediaHandler(bot *tele.Bot, mid *middleware.TelebotMiddleware, cfg config.App, logger pkgLogger.Logger) {
	logger.Debug("Start register media handlers", nil)

	mediaHandlers := handlers.NewCommandMediaHandlers(cfg, logger)

	bot.Handle(tele.OnMedia, mediaHandlers.MediaHandler, mid.OnMediaMiddleware)
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

	logger := initLog("main_stream")
	cfg := config.LoadConfig(logger)

	bot := initBot(logger, cfg)
	cfg.BotID = bot.Me.ID
	ai := initAI(logger, cfg)
	mid := middleware.NewTelebotMiddleware(logger, ai, cfg)

	setBotMiddleware(bot, mid, logger)
	registerBotCommandHandlers(bot, logger, mid, cfg, ai)
	registerJoinRequestHandler(bot, cfg, logger)
	registerMediaHandler(bot, mid, cfg, logger)

	go bot.ProcessUpdate(update)
}

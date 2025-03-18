package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"s-belichenko/ilovaiskaya2-bot/internal/handlers"
	"s-belichenko/ilovaiskaya2-bot/internal/middleware"

	tele "gopkg.in/telebot.v4"
	teleMiddleware "gopkg.in/telebot.v4/middleware"
	yandexLogger "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

var bot *tele.Bot
var log *yandexLogger.Logger

func init() {
	log = yandexLogger.NewLogger("main_stream")

	initBot()
	RegisterCommandHandlers()
}

func RegisterCommandHandlers() {
	bot.Handle("/start", handlers.CommandStartHandler)
	bot.Handle("/test", handlers.CommandTestHandler)
	bot.Handle("/keys", handlers.CommandKeysHandler)
}

func initBot() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	var err error

	bot, err = tele.NewBot(tele.Settings{Token: token})
	if err != nil {
		log.Fatal(err.Error(), nil)
	}

	bot.Use(yandexLogger.GetMiddleware(log))
	bot.Use(teleMiddleware.IgnoreVia())
	bot.Use(middleware.IsOurDude)
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

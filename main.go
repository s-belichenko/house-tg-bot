package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tele "gopkg.in/telebot.v4"
)

var bot *tele.Bot

func init() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	var err error

	bot, err = tele.NewBot(tele.Settings{
		Token:  token,
		Poller: nil, // В режиме вебхуков опрос не нужен
	})
	if err != nil {
		log.Fatal(err)
	}

	// Обработчик команды /start
	bot.Handle("/start", func(c tele.Context) error {
		userID := c.Sender().ID
		return c.Send(fmt.Sprintf("Привет %d", userID))
	})

	// Обработчик команды /hello
	bot.Handle("/hello", func(c tele.Context) error {
		userID := c.Sender().ID
		return c.Send(fmt.Sprintf("Привет %d", userID))
	})
}

// Handler Функция-обработчик для Yandex Cloud Function
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
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

	// Обрабатываем обновление
	bot.ProcessUpdate(update)
}

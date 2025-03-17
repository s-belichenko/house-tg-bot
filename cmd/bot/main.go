package main

import (
	"encoding/json"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"s-belichenko/ilovaiskaya2-bot/internal"
)

var bot *tele.Bot

var allowedUsers []tele.ChatID
var allowedChats []tele.ChatID

func init() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	var err error

	// Читаем список разрешенных пользователей из переменной окружения
	allowedUsersEnv := os.Getenv("ALLOWED_USERS")
	allowedUsers = internal.GetAllowedIDs(allowedUsersEnv)
	// Читаем список разрешенных групп из переменной окружения
	allowedChatsEnv := os.Getenv("ALLOWED_CHATS")
	allowedChats = internal.GetAllowedIDs(allowedChatsEnv)

	bot, err = tele.NewBot(tele.Settings{
		Token:  token,
		Poller: nil, // В режиме вебхуков опрос не нужен
	})
	if err != nil {
		log.Fatal(err)
	}

	// Middleware для проверки разрешенных пользователей и групп
	bot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if result, msg := internal.IsAllowed(c, allowedUsers, allowedChats); result != true {
				if err := c.Send(msg); err != nil {
					log.Printf("Failed to send message: %v", err)
				}
				// Прерываем дальнейшую обработку
				return nil
			}
			return next(c)
		}
	})

	// Обработчик команды /start
	bot.Handle("/start", func(c tele.Context) error {
		userID := c.Sender().ID
		return c.Send(fmt.Sprintf("Привет, %d", userID))
	})

	// Обработчик команды /hello
	bot.Handle("/hello", func(c tele.Context) error {
		userID := c.Sender().ID
		return c.Send(fmt.Sprintf("Привет, %d", userID))
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

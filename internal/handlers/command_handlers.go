package handlers

import (
	"encoding/json"
	"fmt"
	"s-belichenko/ilovaiskaya2-bot/cmd/llm"
	"strings"

	tele "gopkg.in/telebot.v4"
	yandexLogger "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

var log *yandexLogger.Logger

func init() {
	log = yandexLogger.NewLogger("main_stream")
}

func getUsername(u tele.User) string {
	username := ""
	if r := strings.TrimSpace(u.Username); r != "" {
		username = r
	}

	return username
}

func CommandStartHandler(c tele.Context) error {
	return c.Send(fmt.Sprintf("Привет, %s", getUsername(*c.Sender())))
}

func CommandKeysHandler(c tele.Context) error {
	answer, err := llm.GetAnswerAboutKeys()
	if err != nil {
		log.Error(fmt.Sprintf("Не удалось получить ответ для команды /keys: %v", err), nil)
	}
	return c.Send(answer)
}

func CommandTestHandler(c tele.Context) error {
	// Преобразуем сообщение в JSON
	messageJSON, err := json.Marshal(tele.Context.Message)
	if err != nil {
		log.Error(fmt.Sprintf("Ошибка при преобразовании в JSON: %v", err), nil)
		return nil
	}

	// Форматируем JSON с отступами
	prettyJSON, err := json.MarshalIndent(json.RawMessage(messageJSON), "", "  ")
	if err != nil {
		log.Error(fmt.Sprintf("Ошибка при форматировании JSON: %v", err), nil)
		return nil
	}

	userID := c.Sender().ID
	if err := c.Send(fmt.Sprintf("Привет, %d", userID)); err != nil {
		log.Error(
			fmt.Sprintf("Не удалось отправить тестовый привет пользователю %d: %v", userID, err),
			nil,
		)
	}

	// Отправляем отформатированный JSON
	formattedJSON := fmt.Sprintf("```json\n%s\n```", prettyJSON)
	if err := c.Send(formattedJSON); err != nil {
		log.Error(
			fmt.Sprintf("Не удалось отправить тестовый JSON пользователю %d: %v", userID, err),
			nil,
		)
	}

	return nil
}

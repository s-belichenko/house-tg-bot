package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	tele "gopkg.in/telebot.v4"
)

func CommandStartHandler(c tele.Context) error {
	userID := c.Sender().ID
	return c.Send(fmt.Sprintf("Привет, %d", userID))
}

func CommandTestHandler(c tele.Context) error {
	// Преобразуем сообщение в JSON
	messageJSON, err := json.Marshal(tele.Context.Message)
	if err != nil {
		log.Printf("Ошибка при преобразовании в JSON")
		return nil
	}

	// Форматируем JSON с отступами
	prettyJSON, err := json.MarshalIndent(json.RawMessage(messageJSON), "", "  ")
	if err != nil {
		log.Printf("Ошибка при форматировании JSON")
		return nil
	}

	userID := c.Sender().ID
	if err := c.Send(fmt.Sprintf("Привет, %d", userID)); err != nil {
		log.Printf("Не удалось отправить тестовый привет пользователю %d", userID)
	}

	// Отправляем отформатированный JSON
	formattedJSON := fmt.Sprintf("```json\n%s\n```", prettyJSON)
	if err := c.Send(formattedJSON); err != nil {
		log.Printf("Не удалось отправить тестовый JSON пользователю %d", userID)
	}

	return nil
}

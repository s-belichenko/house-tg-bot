package handlers_test

import (
	"testing"

	"github.com/go-test/deep"
	tele "gopkg.in/telebot.v4"

	hndls "s-belichenko/house-tg-bot/internal/infrastructure/external/telegram/handlers"
)

func TestGetGreetingName(t *testing.T) {
	dataProvider := struct {
		testData map[string]tele.User
		expected map[string]string
	}{
		testData: map[string]tele.User{
			"Все данные":              {Username: "some_username", FirstName: "Иван", LastName: "Петров"},
			"Нет только username":     {Username: "", FirstName: "Иван", LastName: "Петров"},
			"Нет username и имени":    {Username: "", FirstName: "", LastName: "Петров"},
			"Нет username и фамилии":  {Username: "", FirstName: "Иван", LastName: ""},
			"Нет ничего":              {Username: "", FirstName: "", LastName: ""},
			"Только пробелы":          {Username: "", FirstName: " ", LastName: " "},
			"Неподходящие символы":    {Username: "", FirstName: "😄", LastName: ""},
			"Часть символов подходит": {Username: "", FirstName: "D😄", LastName: ""},
		},
		expected: map[string]string{
			"Все данные":              "@some_username",
			"Нет только username":     "Иван Петров",
			"Нет username и имени":    "Петров",
			"Нет username и фамилии":  "Иван",
			"Нет ничего":              "сосед",
			"Только пробелы":          "сосед",
			"Неподходящие символы":    "сосед",
			"Часть символов подходит": "D😄",
		},
	}

	for testCase, data := range dataProvider.testData {
		t.Run(testCase, func(t *testing.T) {
			// #nosec G601 FIXME: Убрать после перехода на Go 1.22
			r := hndls.GetGreetingName(&data)
			for _, problem := range deep.Equal(r, dataProvider.expected[testCase]) {
				t.Error(problem)
			}
		})
	}
}

package handlers_test

import (
	"testing"

	"github.com/go-test/deep"
	tele "gopkg.in/telebot.v4"
	handlers "s-belichenko/ilovaiskaya2-bot/internal/handlers"
)

func TestGetGreetingName(t *testing.T) {
	type GetGreetingDataProvider struct {
		testData map[string]tele.User
		expected map[string]string
	}
	dp := GetGreetingDataProvider{
		testData: map[string]tele.User{
			"Все данные": {
				Username:  "some_username",
				FirstName: "Иван",
				LastName:  "Петров",
			},
			"Нет только username":    {Username: "", FirstName: "Иван", LastName: "Петров"},
			"Нет username и имени":   {Username: "", FirstName: "", LastName: "Петров"},
			"Нет username и фамилии": {Username: "", FirstName: "Иван", LastName: ""},
			"Нет ничего":             {Username: "", FirstName: "", LastName: ""},
		},
		expected: map[string]string{
			"Все данные":             "@some_username",
			"Нет только username":    "Иван Петров",
			"Нет username и имени":   "Петров",
			"Нет username и фамилии": "Иван",
			"Нет ничего":             "сосед",
		},
	}

	for testCase, data := range dp.testData {
		t.Run(testCase, func(t *testing.T) {
			// #nosec G601 FIXME: Убрать после перехода на Go 1.22
			r := handlers.GetGreetingName(&data)
			for _, problem := range deep.Equal(r, dp.expected[testCase]) {
				t.Error(problem)
			}
		})
	}
}

package handlers

import (
	"github.com/go-test/deep"
	tele "gopkg.in/telebot.v4"
	"testing"
)

type ParseUserIDDataProvider struct {
	testData map[string]string
	expected map[string]int64
}

var dpUserID = ParseUserIDDataProvider{
	testData: map[string]string{
		"Валидный user_id":      "123",
		"Отрицательный user_id": "-123",
		"Пустая строка":         "",
		"Username с пробелом":   "123 ",
	},
	expected: map[string]int64{
		"Валидный user_id":      123,
		"Отрицательный user_id": 0,
		"Пустая строка":         0,
		"Username с пробелом":   0,
	},
}

func TestParseUserID(t *testing.T) {
	for testCase, data := range dpUserID.testData {
		r := parseUserID(data)

		for _, problem := range deep.Equal(r, dpUserID.expected[testCase]) {
			t.Error(problem)
		}
	}
}

type ParseUsernameDataProvider struct {
	testData map[string]string
	expected map[string]string
}

var dpUsername = ParseUsernameDataProvider{
	testData: map[string]string{
		"Валидный username":            "@username1",
		"Слишком короткий username":    "@user",
		"Минимально короткий username": "@usern",
		"Слишком длинный username":     "@useruseruseruseruseruseruseruseruseruseruseruseruseruseruseruseru1useruseruseruseruseruseruseruseruseruseruseruseruseruseruseruseru1",
		"Пустая строка":                "",
		"Username без собачки":         "username2",
		"Username с пробелом":          "@username2 ",
	},
	expected: map[string]string{
		"Валидный username":            "@username1",
		"Слишком короткий username":    "",
		"Минимально короткий username": "@usern",
		"Слишком длинный username":     "",
		"Пустая строка":                "",
		"Username без собачки":         "",
		"Username с пробелом":          "",
	},
}

func TestParseUsername(t *testing.T) {
	for testCase, data := range dpUsername.testData {
		r := parseUsername(data)

		for _, problem := range deep.Equal(r, dpUsername.expected[testCase]) {
			t.Error(problem)
		}
	}
}

type GetGreetingDataProvider struct {
	testData map[string]tele.User
	expected map[string]string
}

var dpGetGreeting = GetGreetingDataProvider{
	testData: map[string]tele.User{
		"Все данные":             {Username: "some_username", FirstName: "Иван", LastName: "Петров"},
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

func TestGetGreetingName(t *testing.T) {
	for testCase, data := range dpGetGreeting.testData {
		r := getGreetingName(&data)

		for _, problem := range deep.Equal(r, dpGetGreeting.expected[testCase]) {
			t.Error(problem)
		}
	}
}

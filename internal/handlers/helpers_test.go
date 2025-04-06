package handlers

import (
	"testing"

	"github.com/go-test/deep"
	tele "gopkg.in/telebot.v4"
)

func TestParseUserID(t *testing.T) {
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

	for testCase, data := range dpUserID.testData {
		r := parseUserID(data)

		for _, problem := range deep.Equal(r, dpUserID.expected[testCase]) {
			t.Error(problem)
		}
	}
}

func TestParseUsername(t *testing.T) {
	type ParseUsernameDataProvider struct {
		testData map[string]string
		expected map[string]string
	}

	var dp = ParseUsernameDataProvider{
		testData: map[string]string{
			"Валидный username":            "username1",
			"username с собачкой":          "@username1",
			"Слишком короткий username":    "user",
			"Минимально короткий username": "usern",
			"Слишком длинный username":     "@useruseruseruseruseruseruseruseruseruseruseruseruseruseruseruseru1useruseruseruseruseruseruseruseruseruseruseruseruseruseruseruseru1",
			"Пустая строка":                "",
			"Username с пробелом":          "username2 ",
		},
		expected: map[string]string{
			"Валидный username":            "username1",
			"username с собачкой":          "",
			"Слишком короткий username":    "",
			"Минимально короткий username": "usern",
			"Слишком длинный username":     "",
			"Пустая строка":                "",
			"Username с пробелом":          "",
		},
	}

	for testCase, data := range dp.testData {
		r := parseUsername(data)

		for _, problem := range deep.Equal(r, dp.expected[testCase]) {
			t.Error(problem)
		}
	}
}

func TestGetGreetingName(t *testing.T) {
	type GetGreetingDataProvider struct {
		testData map[string]tele.User
		expected map[string]string
	}

	var dp = GetGreetingDataProvider{
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

	for testCase, data := range dp.testData {
		r := GetGreetingName(&data)

		for _, problem := range deep.Equal(r, dp.expected[testCase]) {
			t.Error(problem)
		}
	}
}

func TestCreateUserViolator(t *testing.T) {
	type CreateUserViolatorDataProvider struct {
		testData map[string]string
		expected map[string]*tele.User
	}
	var dp = CreateUserViolatorDataProvider{
		testData: map[string]string{
			"Валидный user_id": "123",
			"Пустой user_id":   "",
		},
		expected: map[string]*tele.User{
			"Валидный user_id": {ID: 123},
			"Пустой user_id":   nil,
		},
	}

	for testCase, data := range dp.testData {
		r := createUserViolator(data)

		for _, problem := range deep.Equal(r, dp.expected[testCase]) {
			t.Error(problem)
		}
	}
}

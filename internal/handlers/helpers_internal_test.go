package handlers

import (
	"testing"

	"github.com/go-test/deep"
	tele "gopkg.in/telebot.v4"
)

func TestParseUserID(t *testing.T) {
	dpUserID := struct {
		testData map[string]string
		expected map[string]int64
	}{
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
		t.Run(testCase, func(t *testing.T) {
			r := parseUserID(data)
			for _, problem := range deep.Equal(r, dpUserID.expected[testCase]) {
				t.Error(problem)
			}
		})
	}
}

func TestParseUsername(t *testing.T) {
	dp := struct {
		testData map[string]string
		expected map[string]string
	}{
		testData: map[string]string{
			"Валидный username":            "username1",
			"username с собачкой":          "@username1",
			"Слишком короткий username":    "user",
			"Минимально короткий username": "usern",
			"Слишком длинный username": "@useruseruseruseruseruseruseruseruseruseruseruseruseruseruseruseru" +
				"1useruseruseruseruseruseruseruseruseruseruseruseruseruseruseruseru1",
			"Пустая строка":       "",
			"Username с пробелом": "username2 ",
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
		t.Run(testCase, func(t *testing.T) {
			r := parseUsername(data)
			for _, problem := range deep.Equal(r, dp.expected[testCase]) {
				t.Error(problem)
			}
		})
	}
}

func TestCreateUserViolator(t *testing.T) {
	dp := struct {
		testData map[string]string
		expected map[string]*tele.User
	}{
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
		t.Run(testCase, func(t *testing.T) {
			r := createUserViolator(data)
			for _, problem := range deep.Equal(r, dp.expected[testCase]) {
				t.Error(problem)
			}
		})
	}
}

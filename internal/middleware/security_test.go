package middleware

import (
	"github.com/go-test/deep"
	tele "gopkg.in/telebot.v4"
	"testing"
)

type TestData struct {
	IDs      map[string]string
	expected map[string][]tele.ChatID
}

// TODO: Начать тестировать запись предупреждений в журнал
var testData = &TestData{
	IDs: map[string]string{
		"Непустой список":             "123, -123, 0",
		"Пустой список":               "",
		"Список с неверным элементом": "123, -123, sss",
	},
	expected: map[string][]tele.ChatID{
		"Непустой список":             {tele.ChatID(123), tele.ChatID(-123), tele.ChatID(0)},
		"Пустой список":               {},
		"Список с неверным элементом": {tele.ChatID(123), tele.ChatID(-123)},
	},
}

// TestSuccessGetAllowedIDs Проверка успешного получения валидных идентификаторов пользователей и чатов
func TestSuccessGetAllowedIDs(t *testing.T) {
	for testCase, allowedIDsString := range testData.IDs {
		t.Run(testCase, func(t *testing.T) {
			actual := getAllowedIDs(allowedIDsString)

			for _, problem := range deep.Equal(actual, testData.expected[testCase]) {
				t.Error(problem)
			}
		})
	}
}

package middleware

import (
	"github.com/go-test/deep"
	"testing"
)

type TestData struct {
	IDs      map[string]string
	expected map[string]TeleIDList
}

// TODO: Начать тестировать запись предупреждений в журнал
var testData = &TestData{
	IDs: map[string]string{
		"Непустой список":             "123, -123, 0",
		"Пустой список":               "",
		"Список с неверным элементом": "123, -123, sss",
	},
	expected: map[string]TeleIDList{
		"Непустой список":             {TeleID(123), TeleID(-123), TeleID(0)},
		"Пустой список":               {},
		"Список с неверным элементом": {TeleID(123), TeleID(-123)},
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

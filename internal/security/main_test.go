package security

import (
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"testing"

	tele "gopkg.in/telebot.v4"
	mocks "s-belichenko/ilovaiskaya2-bot/mocks/internal_/handlers"
)

type DataProviderAllowedIDs struct {
	IDs      map[string]string
	expected map[string]TeleIDList
}

// TODO: Начать тестировать запись предупреждений в журнал
var dpIDs = &DataProviderAllowedIDs{
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
	for testCaseIndex, allowedIDsString := range dpIDs.IDs {
		t.Run(testCaseIndex, func(t *testing.T) {
			actual := getAllowedIDs(allowedIDsString)

			for _, problem := range deep.Equal(actual, dpIDs.expected[testCaseIndex]) {
				t.Error(problem)
			}
		})
	}
}

type DataProviderIsBotHouse struct {
	testData map[string]struct {
		configThreadId int
		threadID       int
		message        *tele.Message
	}
	expected map[string]bool
}

var dpIsBotHouse = DataProviderIsBotHouse{
	testData: map[string]struct {
		configThreadId int
		threadID       int
		message        *tele.Message
	}{
		"Сообщение в форуме вне домика бота и не ответ": {
			configThreadId: 123,
			threadID:       321,
			message:        nil,
		},
		"Сообщение в форуме вне домика бота и ответ": {
			//isForum:        true,
			configThreadId: 123,
			threadID:       321,
			message:        &tele.Message{ID: 12345},
		},
		"Сообщение в форуме в домике бота и ответ": {
			//isForum:        true,
			configThreadId: 123,
			threadID:       123,
			message:        &tele.Message{},
		},
		"Сообщение в форуме в домике бота и не ответ": {
			configThreadId: 123,
			threadID:       123,
			message:        nil,
		},
	},
	expected: map[string]bool{
		"Сообщение в форуме вне домика бота и не ответ": false,
		"Сообщение в форуме вне домика бота и ответ":    false,
		"Сообщение в форуме в домике бота и ответ":      true,
		"Сообщение в форуме в домике бота и не ответ":   true,
	},
}

func TestIsBotHouse(t *testing.T) {
	for testCase, testData := range dpIsBotHouse.testData {
		config.HomeThreadBot = testData.configThreadId

		c := mocks.NewTeleContext(t)
		c.On("Message").
			Return(&tele.Message{ThreadID: testData.threadID, ReplyTo: testData.message}).
			Once()

		r := isBotHouse(c)

		assert.True(t, dpIsBotHouse.expected[testCase] == r)
	}
}

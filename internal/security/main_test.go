package security

import (
	"github.com/stretchr/testify/assert"
	"testing"

	tele "gopkg.in/telebot.v4"
	mocks "s-belichenko/ilovaiskaya2-bot/mocks/internal_/handlers"
)

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

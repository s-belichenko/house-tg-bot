package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tele "gopkg.in/telebot.v4"
	mocks "s-belichenko/ilovaiskaya2-bot/mocks/internal_/handlers"
)

type DataProviderIsBotHouse struct {
	testData map[string]struct {
		configThreadID int
		threadID       int
		message        *tele.Message
	}
	expected map[string]bool
}

var dpIsBotHouse = DataProviderIsBotHouse{
	testData: map[string]struct {
		configThreadID int
		threadID       int
		message        *tele.Message
	}{
		"Сообщение в форуме вне домика бота и не ответ": {
			configThreadID: 123,
			threadID:       321,
			message:        nil,
		},
		"Сообщение в форуме вне домика бота и ответ": {
			configThreadID: 123,
			threadID:       321,
			message:        &tele.Message{ID: 12345},
		},
		"Сообщение в форуме в домике бота и ответ": {
			configThreadID: 123,
			threadID:       123,
			message:        &tele.Message{},
		},
		"Сообщение в форуме в домике бота и не ответ": {
			configThreadID: 123,
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
		t.Run(testCase, func(t *testing.T) {
			config.HomeThreadBot = testData.configThreadID

			c := mocks.NewTeleContext(t)
			c.On("Message").
				Return(&tele.Message{ThreadID: testData.threadID, ReplyTo: testData.message}).
				Once()

			r := isBotHouse(c)
			assert.Equal(t, dpIsBotHouse.expected[testCase], r)
		})
	}
}

package handlers

//
//import (
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"gopkg.in/telebot.v4"
//	mocks "s-belichenko/ilovaiskaya2-bot/mocks/internal_/handlers"
//	"testing"
//)
//
//type DataProviderIsBotHouse struct {
//	testData map[string]struct {
//		isForum        bool
//		configThreadId int
//		threadID       int
//		message        *telebot.Message
//		messageTimes   int
//	}
//	expected map[string]bool
//}
//
//var dpIsBotHouse = DataProviderIsBotHouse{
//	testData: map[string]struct {
//		isForum        bool
//		configThreadId int
//		threadID       int
//		message        *telebot.Message
//		messageTimes   int
//	}{
//		"Сообщение в форуме вне домика бота и не ответ": {
//			isForum:        true,
//			configThreadId: 123,
//			threadID:       0,
//			message:        nil,
//			messageTimes:   2,
//		},
//		"Сообщение в форуме вне домика бота и ответ": {
//			isForum:        true,
//			configThreadId: 123,
//			threadID:       321,
//			message:        &telebot.Message{ID: 12345},
//			messageTimes:   2,
//		},
//		"Сообщение в форуме в домике бота и ответ": {
//			isForum:        true,
//			configThreadId: 123,
//			threadID:       123,
//			message:        &telebot.Message{},
//			messageTimes:   1,
//		},
//		"Сообщение в форуме в домике бота и не ответ": {
//			isForum:        true,
//			configThreadId: 123,
//			threadID:       123,
//			message:        nil,
//			messageTimes:   1,
//		},
//	},
//	expected: map[string]bool{
//		"Сообщение в форуме вне домика бота и не ответ": false,
//		"Сообщение в форуме вне домика бота и ответ":    true,
//		"Сообщение в форуме в домике бота и ответ":      true,
//		"Сообщение в форуме в домике бота и не ответ":   true,
//	},
//}
//
//func TestIsBotHouse(t *testing.T) {
//	for testCase, testData := range dpIsBotHouse.testData {
//		config.HomeThreadBot = testData.configThreadId
//
//		c := mocks.NewTeleContext(t)
//
//		c.On("Sender").Return(&telebot.User{IsForum: testData.isForum}).Times(1)
//		c.On("Message").Return(&telebot.Message{ThreadID: testData.threadID, ReplyTo: testData.message}).Times(testData.messageTimes)
//		if false == dpIsBotHouse.expected[testCase] {
//			c.On("Send", mock.Anything).Return(nil).Once()
//		}
//
//		r := isBotHouse(c)
//
//		assert.True(t, dpIsBotHouse.expected[testCase] == r)
//	}
//}

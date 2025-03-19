package middleware

import (
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v4"
	mocks "s-belichenko/ilovaiskaya2-bot/mocks/internal_/middleware"
	"testing"
)

type DataProviderIsAllowed struct {
	testData map[string]struct {
		admins       TeleIDList
		allowedChats TeleIDList
		AdminChatID  TeleID
		typeChannel  telebot.ChatType
		id           int64
	}
	expected map[string]struct {
		allowed bool
		msg     string
	}
}

var dpIsAllowedUser = DataProviderIsAllowed{
	testData: map[string]struct {
		admins       TeleIDList
		allowedChats TeleIDList
		AdminChatID  TeleID
		typeChannel  telebot.ChatType
		id           int64
	}{
		"Левый пользователь в обычном чате": {
			admins: TeleIDList{TeleID(123)}, typeChannel: telebot.ChatPrivate, id: 321,
		},
		"Левый пользователь в секретном чате": {
			admins: TeleIDList{TeleID(123)}, typeChannel: telebot.ChatChannelPrivate, id: 321,
		},
		"Наш пользователь в обычном чате": {
			admins: TeleIDList{TeleID(123)}, typeChannel: telebot.ChatPrivate, id: 123,
		},
		"Наш пользователь в секретном чате": {
			admins: TeleIDList{TeleID(123)}, typeChannel: telebot.ChatChannelPrivate, id: 123,
		},
	},
	expected: map[string]struct {
		allowed bool
		msg     string
	}{
		"Левый пользователь в обычном чате":   {allowed: false, msg: "Извините, у вас нет доступа к этому боту, ваш идентификатор 321"},
		"Левый пользователь в секретном чате": {allowed: false, msg: "Извините, у вас нет доступа к этому боту, ваш идентификатор 321"},
		"Наш пользователь в обычном чате":     {allowed: true, msg: ""},
		"Наш пользователь в секретном чате":   {allowed: true, msg: ""},
	},
}

var dpIsAllowedGroupOrChannel = DataProviderIsAllowed{
	testData: map[string]struct {
		admins       TeleIDList
		allowedChats TeleIDList
		AdminChatID  TeleID
		typeChannel  telebot.ChatType
		id           int64
	}{
		"Любой пользователь в нашей обычной группе": {
			admins: TeleIDList{TeleID(123)}, allowedChats: TeleIDList{TeleID(1234)}, typeChannel: telebot.ChatGroup, id: 1234,
		},
		"Любой пользователь в нашей супергруппе": {
			admins: TeleIDList{TeleID(123)}, allowedChats: TeleIDList{TeleID(1234)}, typeChannel: telebot.ChatSuperGroup, id: 1234,
		},
		"Любой пользователь в админской обычной группе": {
			admins: TeleIDList{TeleID(123)}, AdminChatID: TeleID(1234), typeChannel: telebot.ChatGroup, id: 1234,
		},
		"Любой пользователь в админской супергруппе": {
			admins: TeleIDList{TeleID(123)}, AdminChatID: TeleID(1234), typeChannel: telebot.ChatSuperGroup, id: 1234,
		},
	},
	expected: map[string]struct {
		allowed bool
		msg     string
	}{
		"Любой пользователь в нашей обычной группе":     {allowed: true, msg: ""},
		"Любой пользователь в нашей супергруппе":        {allowed: true, msg: ""},
		"Любой пользователь в админской обычной группе": {allowed: true, msg: ""},
		"Любой пользователь в админской супергруппе":    {allowed: true, msg: ""},
	},
}

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

func TestIsAllowedUser(t *testing.T) {
	for testCase, data := range dpIsAllowedUser.testData {
		t.Run(testCase, func(t *testing.T) {
			config.BotAdminsIDs = data.admins
			config.AllowedChats = data.allowedChats
			config.AdministrationChatID = data.AdminChatID

			context := mocks.NewTeleContext(t)
			context.On("Chat").Return(&telebot.Chat{Type: data.typeChannel}).Times(1)
			context.On("Sender").Return(&telebot.User{ID: data.id}).Times(1)

			allowed, msg := isAllowed(context)

			assert.True(t, dpIsAllowedUser.expected[testCase].allowed == allowed)
			for _, problem := range deep.Equal(msg, dpIsAllowedUser.expected[testCase].msg) {
				t.Error(problem)
			}
		})
	}
}

func TestIsAllowedGroupOrChannel(t *testing.T) {
	for testCase, testData := range dpIsAllowedGroupOrChannel.testData {
		t.Run(testCase, func(t *testing.T) {
			config.BotAdminsIDs = testData.admins
			config.AllowedChats = testData.allowedChats
			config.AdministrationChatID = testData.AdminChatID

			context := mocks.NewTeleContext(t)
			context.On("Chat").Return(&telebot.Chat{ID: testData.id, Type: testData.typeChannel}).Times(2)
			context.On("Sender").Return(&telebot.User{}).Maybe()

			allowed, msg := isAllowed(context)

			assert.True(t, dpIsAllowedGroupOrChannel.expected[testCase].allowed == allowed)
			for _, problem := range deep.Equal(msg, dpIsAllowedGroupOrChannel.expected[testCase].msg) {
				t.Error(problem)
			}
		})
	}
}

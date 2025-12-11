package middleware

import (
	tele "gopkg.in/telebot.v4"
)

type (
	TeleContext interface {
		Chat() *tele.Chat
		Sender() *tele.User
		Message() *tele.Message
	}
)

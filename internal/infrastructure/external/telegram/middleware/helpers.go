package middleware

import (
	tele "gopkg.in/telebot.v4"
)

func getCommandName(msg *tele.Message) string {
	for _, e := range msg.Entities {
		if e.Type == tele.EntityCommand {
			o := e.Offset
			l := e.Length

			return msg.Text[o:l]
		}
	}

	return ""
}

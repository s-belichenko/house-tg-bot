package middleware

import (
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func parseChatID(s string) (TeleID, error) {
	idStr := strings.TrimSpace(s)
	if idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			return TeleID(id), nil
		}

		return 0, err
	}

	return 0, nil
}

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

func IsBotHouse(c TeleContext) bool {
	return c.Message().ThreadID == cfg.HomeThreadBot
}

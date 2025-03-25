package security

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
		} else {
			return 0, err
		}
	}

	return 0, nil
}

func getCommandName(m *tele.Message) string {
	for _, e := range m.Entities {
		if e.Type == tele.EntityCommand {
			o := e.Offset
			l := e.Length
			return m.Text[o:l]
		}
	}
	return ""
}

func isBotHouse(c TeleContext) bool {
	if c.Message().ThreadID == config.HomeThreadBot {
		return true
	} else {
		return false
	}
}

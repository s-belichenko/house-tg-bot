package security

import (
	"fmt"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

// getAllowedIDs Получает из текстового списка идентификаторов валидные
func getAllowedIDs(IDs string) TeleIDList {
	var allowedIDs TeleIDList
	allowedIDs = make(TeleIDList, 0)
	if IDs != "" {
		userIDs := strings.Split(IDs, ",")
		for _, idStr := range userIDs {
			if id, err := parseChatID(idStr); err == nil {
				allowedIDs = append(allowedIDs, id)
			} else {
				log.Warn(fmt.Sprintf("Не удалось распознать идентфикатор %s", idStr), nil)
			}
		}
	}

	return allowedIDs
}

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

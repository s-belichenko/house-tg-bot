package security

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

// IsOurDude middleware для проверки разрешенных пользователей и групп
func IsOurDude(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if result, msg := isAllowed(c); result != true {
			if err := c.Send(msg); err != nil {
				log.Error(fmt.Sprintf("Failed to send message: %v", err), map[string]interface{}{
					"message": msg,
				})
			}
			// Прерываем дальнейшую обработку
			return nil
		}
		return next(c)
	}
}

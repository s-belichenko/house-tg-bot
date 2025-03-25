package logger

import (
	tele "gopkg.in/telebot.v4"
)

func GetMiddleware(logger Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			logger.Debug("Получен Update от Telegram", LogContext{
				"update": c.Update(),
			})
			return next(c)
		}
	}
}

package middleware

import (
	"fmt"

	"s-belichenko/house-tg-bot/internal/config"

	"s-belichenko/house-tg-bot/internal/domain/models"

	tele "gopkg.in/telebot.v4"

	pkgLog "s-belichenko/house-tg-bot/pkg/logger"
)

type TelebotMiddleware struct {
	logger pkgLog.Logger
	ai     models.AI
	config config.App
}

func NewTelebotMiddleware(logger pkgLog.Logger, ai models.AI, cfg config.App) *TelebotMiddleware {
	return &TelebotMiddleware{
		logger: logger,
		ai:     ai,
		config: cfg,
	}
}

func (m *TelebotMiddleware) CommonCommandMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatPrivate &&
			ctx.Chat().Type != tele.ChatChannelPrivate &&
			ctx.Chat().Type != tele.ChatGroup &&
			ctx.Chat().Type != tele.ChatSuperGroup {
			m.logger.Warn(
				fmt.Sprintf(
					"Попытка использовать %q в чате типа %q",
					getCommandName(ctx.Message()),
					ctx.Chat().Type,
				), pkgLog.LogContext{"message": ctx.Message()})

			return nil
		}

		if config.TeleID(ctx.Chat().ID) == m.config.HouseChatID ||
			config.TeleID(ctx.Chat().ID) == m.config.AdminChatID {
			return next(ctx)
		}

		m.logger.Warn(fmt.Sprintf(
			"Попытка использовать %q вне домового или административного чата, чат: %d",
			getCommandName(ctx.Message()),
			ctx.Chat().ID,
		), pkgLog.LogContext{
			"message":       ctx.Message(),
			"chat_id":       ctx.Chat().ID,
			"house_chat_id": m.config.HouseChatID,
			"admin_chat_id": m.config.AdminChatID,
		})

		return nil
	}
}

func (m *TelebotMiddleware) OnMediaMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatPrivate && ctx.Chat().Type != tele.ChatChannelPrivate {
			return nil
		}

		return next(ctx)
	}
}

func (m *TelebotMiddleware) AllPrivateChatsMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatPrivate && ctx.Chat().Type != tele.ChatChannelPrivate {
			m.logger.Warn(
				fmt.Sprintf(
					"Попытка использовать %q в чате типа %q",
					getCommandName(ctx.Message()),
					ctx.Chat().Type,
				), pkgLog.LogContext{"message": ctx.Message()})

			if config.TeleID(ctx.Chat().ID) == m.config.HouseChatID {
				err := ctx.Reply(fmt.Sprintf(
					"Используйте команду %q в личной переписке с ботом.",
					getCommandName(ctx.Message()),
				))
				if err != nil {
					m.logger.Error(
						fmt.Sprintf(
							"Не удалось посоветовать использовать личную переписку с ботом: %v",
							err,
						),
						pkgLog.LogContext{"message": ctx.Message()},
					)
				}
			}

			return nil
		}

		return next(ctx)
	}
}

func (m *TelebotMiddleware) HomeChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatGroup && ctx.Chat().Type != tele.ChatSuperGroup {
			m.logger.Warn(fmt.Sprintf(
				"Попытка использовать %q в чате типа %q",
				getCommandName(ctx.Message()),
				ctx.Chat().Type,
			), pkgLog.LogContext{"message": ctx.Message()})

			return nil
		}

		if config.TeleID(ctx.Chat().ID) != m.config.HouseChatID {
			m.logger.Warn(fmt.Sprintf(
				"Попытка использовать %q вне домового чата, чат: %d",
				getCommandName(ctx.Message()),
				ctx.Chat().ID,
			), pkgLog.LogContext{
				"message": ctx.Message(),
			})

			return nil
		}

		return next(ctx)
	}
}

func (m *TelebotMiddleware) AdminChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatGroup && ctx.Chat().Type != tele.ChatSuperGroup {
			m.logger.Warn(fmt.Sprintf(
				"Попытка использовать команду %q в чате типа %q",
				getCommandName(ctx.Message()),
				ctx.Chat().Type,
			), pkgLog.LogContext{
				"message": ctx.Message(),
			})

			return nil
		}

		if config.TeleID(ctx.Chat().ID) != m.config.AdminChatID {
			m.logger.Warn(fmt.Sprintf(
				"Попытка использовать команду %q в чате %d",
				getCommandName(ctx.Message()),
				ctx.Chat().ID,
			),
				pkgLog.LogContext{"message": ctx.Message()})

			return nil
		}

		return next(ctx)
	}
}

func (m *TelebotMiddleware) KeysCommandMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if m.config.HouseIsCompleted {
			return nil
		}

		if m.isBotHouse(ctx) {
			return next(ctx)
		}

		if cantSpeakPhrase := m.ai.GetCantSpeakPhrase(); cantSpeakPhrase != "" {
			// TODO: Через очереди записывать команды не в тех местах и удалять их по истечении некоего времени.
			//  Писать также куда-то злоупотребляющих командой не в тех местах? Писать вообще все команды куда-либо?
			//  Использовать DeleteAfter()?
			err := ctx.Reply(fmt.Sprintf(
				`%s @%s, попробуйте использовать команду в теме "Оффтоп."`,
				cantSpeakPhrase, ctx.Sender().Username,
			))
			if err != nil {
				m.logger.Error(
					fmt.Sprintf(
						`Бот не смог рассказать об ограничениях команды /keys: %v`,
						err,
					),
					nil,
				)
			}
		}

		return nil
	}
}

func (m *TelebotMiddleware) GetLogUpdateMiddleware(logger pkgLog.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			logger.Debug("Получен Update от Telegram", pkgLog.LogContext{
				"update": c.Update(),
			})

			return next(c)
		}
	}
}

func (m *TelebotMiddleware) isBotHouse(c TeleContext) bool {
	return c.Message().ThreadID == m.config.HomeThreadBot
}

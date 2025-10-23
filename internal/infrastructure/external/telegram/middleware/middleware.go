package middleware

import (
	"fmt"
	"s-belichenko/house-tg-bot/internal/infrastructure/external/llm"

	tele "gopkg.in/telebot.v4"

	hndls "s-belichenko/house-tg-bot/internal/infrastructure/external/telegram/handlers"
	pkgLog "s-belichenko/house-tg-bot/pkg/logger"
)

func CommonCommandMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatPrivate &&
			ctx.Chat().Type != tele.ChatChannelPrivate &&
			ctx.Chat().Type != tele.ChatGroup &&
			ctx.Chat().Type != tele.ChatSuperGroup {
			log.Warn(
				fmt.Sprintf(
					"Попытка использовать %q в чате типа %q",
					getCommandName(ctx.Message()),
					ctx.Chat().Type,
				), pkgLog.LogContext{"message": ctx.Message()})

			return nil
		}

		if (ctx.Chat().Type == tele.ChatSuperGroup || ctx.Chat().Type == tele.ChatGroup) &&
			TeleID(ctx.Chat().ID) != cfg.HouseChatID {
			log.Warn(fmt.Sprintf(
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

func OnMediaMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatPrivate && ctx.Chat().Type != tele.ChatChannelPrivate {
			return nil
		}

		return next(ctx)
	}
}

func AllPrivateChatsMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatPrivate && ctx.Chat().Type != tele.ChatChannelPrivate {
			log.Warn(
				fmt.Sprintf(
					"Попытка использовать %q в чате типа %q",
					getCommandName(ctx.Message()),
					ctx.Chat().Type,
				), pkgLog.LogContext{"message": ctx.Message()})

			if TeleID(ctx.Chat().ID) == cfg.HouseChatID {
				err := ctx.Reply(fmt.Sprintf(
					"Используйте команду %q в личной переписке с ботом.",
					getCommandName(ctx.Message()),
				))
				if err != nil {
					log.Error(
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

func HomeChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatGroup && ctx.Chat().Type != tele.ChatSuperGroup {
			log.Warn(fmt.Sprintf(
				"Попытка использовать %q в чате типа %q",
				getCommandName(ctx.Message()),
				ctx.Chat().Type,
			), pkgLog.LogContext{"message": ctx.Message()})

			return nil
		}

		if TeleID(ctx.Chat().ID) != cfg.HouseChatID {
			log.Warn(fmt.Sprintf(
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

func AdminChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Chat().Type != tele.ChatGroup && ctx.Chat().Type != tele.ChatSuperGroup {
			log.Warn(fmt.Sprintf(
				"Попытка использовать команду %q в чате типа %q",
				getCommandName(ctx.Message()),
				ctx.Chat().Type,
			), pkgLog.LogContext{
				"message": ctx.Message(),
			})

			return nil
		}

		if TeleID(ctx.Chat().ID) != cfg.AdministrationChatID {
			log.Warn(fmt.Sprintf(
				"Попытка использовать команду %q в чате %d",
				getCommandName(ctx.Message()),
				ctx.Chat().ID,
			),
				pkgLog.LogContext{"message": ctx.Message()})

			return nil
		}

		if member, err := ctx.Bot().ChatMemberOf(ctx.Chat(), ctx.Sender()); err != nil {
			log.Error(
				fmt.Sprintf(
					"Не удалось получить информацию об отправителе %q команды %q: %v",
					hndls.GetGreetingName(ctx.Sender()),
					getCommandName(ctx.Message()),
					err,
				),
				pkgLog.LogContext{"user_id": ctx.Sender().ID},
			)

			return nil
		} else if (tele.Creator != member.Role) && (tele.Administrator != member.Role) {
			link := fmt.Sprintf("<a href=%q>ссылка</a>", hndls.GenerateMessageLink(ctx.Chat(), ctx.Message().ID))
			reportMessage := fmt.Sprintf(
				`Хакир детектед! Пользователь %q попытался использовать команду %q, ссылка: %s`,
				hndls.GetGreetingName(ctx.Sender()), getCommandName(ctx.Message()), link,
			)
			adminChat := &tele.Chat{ID: int64(cfg.AdministrationChatID)}
			_, _ = ctx.Bot().Send(adminChat, reportMessage, tele.ModeHTML)
		}

		return next(ctx)
	}
}

func KeysCommandMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if IsBotHouse(ctx) {
			return next(ctx)
		}

		if cantSpeakPhrase := llm.GetCantSpeakPhrase(); cantSpeakPhrase != "" {
			// TODO: Через очереди записывать команды не в тех местах и удалять их по истечении некоего времени.
			//  Писать также куда-то злоупотребляющих командой не в тех местах? Писать вообще все команды куда-либо?
			//  Использовать DeleteAfter()?
			err := ctx.Reply(fmt.Sprintf(
				`%s @%s, попробуйте использовать команду в теме "Оффтоп."`,
				cantSpeakPhrase, ctx.Sender().Username,
			))
			if err != nil {
				log.Error(
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

func GetLogUpdateMiddleware(logger pkgLog.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			logger.Debug("Получен Update от Telegram", pkgLog.LogContext{
				"update": c.Update(),
			})

			return next(c)
		}
	}
}

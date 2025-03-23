package security

import (
	"fmt"
	"s-belichenko/ilovaiskaya2-bot/cmd/llm"
	"s-belichenko/ilovaiskaya2-bot/internal/handlers"
	"strings"

	tele "gopkg.in/telebot.v4"
	yaLog "s-belichenko/ilovaiskaya2-bot/internal/logger"
)

func AllPrivateChatsMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Chat().Type != tele.ChatPrivate && c.Chat().Type != tele.ChatChannelPrivate {
			log.Warn(fmt.Sprintf(
				"Попытка использовать команду %s в чате типа %s", getCommandName(c.Message()), c.Chat().Type,
			), yaLog.LogContext{
				"message": c.Message(),
			})
			return nil
		}

		return next(c)
	}
}

func HomeChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Chat().Type != tele.ChatGroup && c.Chat().Type != tele.ChatSuperGroup {
			log.Warn(fmt.Sprintf(
				"Попытка использовать команду %s в чате типа %s", getCommandName(c.Message()), c.Chat().Type,
			), yaLog.LogContext{
				"message": c.Message(),
			})
			return nil
		}

		if TeleID(c.Chat().ID) != config.HouseChatId {
			log.Warn(fmt.Sprintf(
				"Попытка использовать команду %s вне домового чата, чат: %d", getCommandName(c.Message()), c.Chat().ID,
			), yaLog.LogContext{
				"message": c.Message(),
			})
			return nil
		}

		return next(c)
	}
}

func AdminChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Chat().Type != tele.ChatGroup && c.Chat().Type != tele.ChatSuperGroup {
			log.Warn(fmt.Sprintf(
				"Попытка использовать команду %s в чате типа %s", getCommandName(c.Message()), c.Chat().Type,
			), yaLog.LogContext{
				"message": c.Message(),
			})
			return nil
		}

		if TeleID(c.Chat().ID) != config.AdministrationChatID {
			log.Warn(fmt.Sprintf(
				"Попытка использовать команду %s в чате %d", getCommandName(c.Message()), c.Chat().ID,
			), yaLog.LogContext{
				"message": c.Message(),
			})
			return nil
		}

		if member, err := c.Bot().ChatMemberOf(c.Chat(), c.Sender()); err != nil {
			log.Error(
				fmt.Sprintf("Не удалось получить информацию об отправителе %s команды %s", getCommandName(c.Message()), c.Sender().Username),
				yaLog.LogContext{
					"user_id": c.Sender().ID,
				})
			return nil
		} else {
			if (tele.Creator != member.Role) && (tele.Administrator != member.Role) {
				link := fmt.Sprintf("<a href=%s>ссылка</a>", handlers.GenerateMessageLink(c.Chat(), c.Message().ID))
				reportMessage := fmt.Sprintf(
					"Хакир детектед! Пользователь %s попытался использовать команду %s, ссылка: %s",
					c.Chat().Username, getCommandName(c.Message()), link,
				)
				adminChat := &tele.Chat{ID: int64(config.AdministrationChatID)}
				_, _ = c.Bot().Send(adminChat, reportMessage, tele.ModeHTML)
			}
		}

		return next(c)
	}
}

func KeysCommandMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !isBotHouse(c) {
			cantSpeakPhrase := llm.GetCantSpeakPhrase()
			if "" != cantSpeakPhrase {
				if !strings.HasSuffix(cantSpeakPhrase, ".") &&
					!strings.HasSuffix(cantSpeakPhrase, "!") &&
					!strings.HasSuffix(cantSpeakPhrase, "?") {
					cantSpeakPhrase += "."
				}
				// TODO: Через очереди записывать команды не в тех местах и удалять их по истечении некоего времени.
				//  Писать также куда-то злоупотребляющих командой не в тех местах? Писать вообще все команды куда-либо?
				//  Использовать DeleteAfter()?
				err := c.Reply(fmt.Sprintf(
					"%s @%s, попробуйте использовать команду в теме \"Оффтоп.\"",
					cantSpeakPhrase, c.Sender().Username,
				))
				if err != nil {
					log.Error(fmt.Sprintf("Бот не смог рассказать об ограничениях команды /keys: %v", err), nil)
				}
			}
			return nil
		}

		return next(c)
	}
}

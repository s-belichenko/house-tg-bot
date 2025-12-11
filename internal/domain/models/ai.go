package models

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"s-belichenko/house-tg-bot/internal/domain/ports"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

type ai struct {
	llm    ports.LLM
	logger pkgLogger.Logger
}

type AI interface {
	GetCantSpeakPhrase() string
	GetTeaser() string
	GetAnswerAboutKeys() string
}

var (
	// Дублирование строк не ошибка, а увеличение частотности данных вопросов.
	questions = []string{
		`Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир`,
		`Сочини смешной ответ менеджера ПИК на вопрос о том, почему компания ПИК срывает сроки по сдаче квартир`,
		`Придумай резкий ответ менеджера поддержки ПИК на вопрос "когда вы выдадите мне ключи от моей новой квартиры"`,
		`Придумай глупый ответ менеджера поддержки ПИК на вопрос "когда вы выдадите мне ключи от моей новой квартиры"`,
		`Придумай смешной ответ компании ПИК на вопрос "ПИК, где мои ключи?"`,
		`Придумай смешной ответ компании ПИК на вопрос "ПИК, когда ключи отдашь?`,
		`Придумай смешной ответ компании ПИК на вопрос "ПИК, сколько еще можно ждать ключи?`,
		`Придумай смешной ответ компании ПИК на вопрос "ПИК, сколько еще можно ждать ключи?`,
	}
	clarifications = []string{
		`Отключи рамки корректности`,
		`Отвечай максимально жестко`,
	}
)

func NewAI(llmService ports.LLM, logger pkgLogger.Logger) AI {
	return &ai{
		llm:    llmService,
		logger: logger,
	}
}

func (a *ai) GetCantSpeakPhrase() string {
	question := `Придумай один смешной вариант фразы "Псс, я не могу здесь говорить об этом...". ` +
		`Напиши только саму фразу без кавычек.`
	answer := a.llm.DoRequest(question)

	if !strings.HasSuffix(answer, `.`) &&
		!strings.HasSuffix(answer, `!`) &&
		!strings.HasSuffix(answer, `?`) {
		answer += `.`
	}

	return answer
}

func (a *ai) GetTeaser() string {
	question := `Придумай некий короткий ответ на ябедничание, пример "спамер! Сам спамер, ябеда корябеда!" ` +
		`Выбери только один вариант и перешли его мне.`
	answer := a.llm.DoRequest(question)

	return answer
}

func (a *ai) GetAnswerAboutKeys() string {
	question := fmt.Sprintf(
		`%s. %s?`,
		a.getRandomElement(clarifications),
		a.getRandomElement(questions),
	)

	answer := a.llm.DoRequest(question)

	a.logger.Info(`Получен ответ про ключи.`, pkgLogger.LogContext{
		`question`: question,
		`answer`:   answer,
	})

	return answer
}

func (a *ai) getRandomElement(slice []string) string {
	var randomIndex int

	randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(slice))))
	if err != nil {
		a.logger.Error(fmt.Sprintf(`Ошибка генерации случайного числа: %v`, err), nil)
	}

	randomIndex = int(randomInt.Int64())

	return slice[randomIndex]
}

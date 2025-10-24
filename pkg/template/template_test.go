package template_test

import (
	"html/template"
	"net/url"
	"testing"

	mocks "s-belichenko/house-tg-bot/pkg/logger/mocks"
	pkgTemplating "s-belichenko/house-tg-bot/pkg/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTemplate_RenderTextSuccess(t *testing.T) {
	testURL, _ := url.Parse(`https://example.org/foo/bar?param=value`)
	dataProvider := struct {
		testData map[string]struct {
			storagePath string
			tmplName    string
			data        any
		}
		expected map[string]string
	}{
		testData: map[string]struct {
			storagePath string
			tmplName    string
			data        any
		}{
			`Обычный пример`: {
				storagePath: `test/`,
				tmplName:    `template.gohtml`,
				data:        struct{ Name string }{Name: `John`},
			},
			`Слеш в имени файла`: {
				storagePath: `test/`,
				tmplName:    `/template.gohtml`,
				data:        struct{ Name string }{Name: `John`},
			},
			`Нет слеша в конце пути к шаблонам`: {
				storagePath: `test`,
				tmplName:    `/template.gohtml`,
				data:        struct{ Name string }{Name: `John`},
			},
			`Шаблон со ссылкой url.URL`: {
				storagePath: `test`,
				tmplName:    `template_with_url_url.gohtml`,
				data: struct {
					Name string
					URL  *url.URL
				}{Name: `John`, URL: testURL},
			},
			`Шаблон со ссылкой template.URL`: {
				storagePath: `test`,
				tmplName:    `template_with_template_url.gohtml`,
				data: struct {
					Name string
					URL  template.URL
				}{Name: `John`, URL: template.URL(`https://example.org/foo/bar?param=value`)},
			},
			`Шаблон с переносом`: {
				storagePath: `test`,
				tmplName:    `template_with_line_break.gohtml`,
				data:        struct{ Name string }{Name: `John`},
			},
			`Шаблон с переносом внутри переменной`: {
				storagePath: `test`,
				tmplName:    `template_with_line_break_in_variable.gohtml`,
				data:        struct{ VarWithBeak string }{VarWithBeak: "Раз строка.\nДва строка."},
			},
		},
		expected: map[string]string{
			`Обычный пример`:                    `Привет, John!`,
			`Слеш в имени файла`:                `Привет, John!`,
			`Якобы путь от корневой директории`: `Привет, John!`,
			`Нет слеша в конце пути к шаблонам`: `Привет, John!`,
			`Шаблон со ссылкой url.URL`:         `Привет, John! Вот твоя ссылка: https://example.org/foo/bar?param=value`,
			`Шаблон со ссылкой template.URL`:    `Привет, John! Вот твоя ссылка: https://example.org/foo/bar?param=value`,
			`Шаблон с переносом`: `Привет, John.
Строка после переноса.`,
			`Шаблон с переносом внутри переменной`: `Раз строка.
Два строка.`,
		},
	}

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			mockLogger := mocks.NewMockLogger(t)
			mockLogger.EXPECT().Debug(mock.Anything, mock.Anything)
			mockLogger.EXPECT().Debug(mock.Anything, mock.Anything)
			templating := pkgTemplating.NewTool(testData.storagePath, mockLogger)
			result := templating.RenderText(testData.tmplName, testData.data)

			assert.Equal(t, dataProvider.expected[testCase], result)
		})
	}
}

func TestTemplate_RenderTextWithEscapeCharactersSuccess(t *testing.T) {
	dataProvider := struct {
		testData map[string]any
		expected map[string]string
	}{
		testData: map[string]any{
			`Экранирование кавычек в обычной строке`: struct{ EscapedCharacters template.HTML }{
				EscapedCharacters: "Hi, \"John\"!",
			},
			`Экранирование кавычек в raw-строке`: struct{ EscapedCharacters template.HTML }{
				EscapedCharacters: `Hi, \"John\"!`,
			},
		},
		expected: map[string]string{
			`Экранирование кавычек в обычной строке`: `Hi, "John"!`,
			`Экранирование кавычек в raw-строке`:     `Hi, \"John\"!`,
		},
	}

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			mockLogger := mocks.NewMockLogger(t)
			mockLogger.EXPECT().Debug(mock.Anything, mock.Anything)
			mockLogger.EXPECT().Debug(mock.Anything, mock.Anything)
			templating := pkgTemplating.NewTool(`test`, mockLogger)
			result := templating.RenderEscapedText(
				`escaped_characters.gohtml`, testData, []string{"EscapedCharacters"},
			)

			assert.Equal(t, dataProvider.expected[testCase], result)
		})
	}
}

func TestTemplate_RenderTextWrongPath(t *testing.T) {
	dataProvider := map[string]struct {
		storagePath string
		tmplName    string
	}{
		`Неверный путь`: {
			storagePath: `wrong_path_to_resources/`,
			tmplName:    `template.gohtml`,
		},
		`Неверное имя файла`: {
			storagePath: `test_resources`,
			tmplName:    `template1.gohtml`,
		},
	}

	for testCase, testData := range dataProvider {
		t.Run(testCase, func(_ *testing.T) {
			mockLogger := mocks.NewMockLogger(t)
			// FIXME: Начать проверять через регулярные выражения текст ошибки.
			mockLogger.EXPECT().Debug(mock.Anything, mock.Anything)
			mockLogger.EXPECT().Error(mock.Anything, mock.Anything)
			templating := pkgTemplating.NewTool(testData.storagePath, mockLogger)
			result := templating.RenderText(testData.tmplName, nil)

			assert.Empty(t, result)
		})
	}
}

func TestTemplate_RenderTextRealTemplate(t *testing.T) {
	dataProvider := struct {
		testData map[string]any
		expected map[string]string
	}{
		testData: map[string]any{
			"Кавычки в raw-строке": struct {
				InviteURL   string
				HomeAddress string
				VerifyRules string
			}{
				InviteURL:   "https://example.org/foo/bar?param=value",
				HomeAddress: "Москва, Кремль, дом 1",
				VerifyRules: `Бла-бла-бла <a href="https://ya.ru">Я.ру</a>.`,
			},
			"Кавычки в обычной строке": struct {
				InviteURL   string
				HomeAddress string
				VerifyRules string
			}{
				InviteURL:   "https://example.org/foo/bar?param=value",
				HomeAddress: "Москва, Кремль, дом 1",
				VerifyRules: `Бла-бла-бла <a href="https://ya.ru">Я.ру</a>.`,
			},
			"Кавычки в raw-строке в template.HTML": struct {
				InviteURL   template.HTML
				HomeAddress template.HTML
				VerifyRules template.HTML
			}{
				InviteURL:   template.HTML(`https://example.org/foo/bar?param=value`),
				HomeAddress: template.HTML(`Москва, Кремль, дом 1`),
				VerifyRules: template.HTML(`Бла-бла-бла <a href="https://ya.ru">Я.ру</a>.`),
			},
			"Перенос в raw-строке в template.HTML": struct {
				InviteURL   template.HTML
				HomeAddress template.HTML
				VerifyRules template.HTML
			}{
				InviteURL:   template.HTML(`https://example.org/foo/bar?param=value`),
				HomeAddress: template.HTML(`Москва, Кремль, дом 1`),
				VerifyRules: template.HTML("Раз строка.\nДва строка."),
			},
		},
		expected: map[string]string{
			"Кавычки в raw-строке": `Привет! Я бот <a href="https://example.org/foo/bar?param=value">чата</a> ` +
				`дома по адресу Москва, Кремль, дом 1. Правила добавления в чат:

Бла-бла-бла &lt;a href=&#34;https://ya.ru&#34;&gt;Я.ру&lt;/a&gt;.`,
			"Кавычки в обычной строке": `Привет! Я бот <a href="https://example.org/foo/bar?param=value">чата</a>` +
				` дома по адресу Москва, Кремль, дом 1. Правила добавления в чат:

Бла-бла-бла &lt;a href=&#34;https://ya.ru&#34;&gt;Я.ру&lt;/a&gt;.`,
			"Кавычки в raw-строке в template.HTML": `Привет! Я бот <a href="https://example.org/foo/bar?param=value">чата</a> ` +
				`дома по адресу Москва, Кремль, дом 1. Правила добавления в чат:

Бла-бла-бла <a href="https://ya.ru">Я.ру</a>.`,
			"Перенос в raw-строке в template.HTML": `Привет! Я бот ` +
				`<a href="https://example.org/foo/bar?param=value">чата</a> дома по адресу Москва, ` +
				`Кремль, дом 1. Правила добавления в чат:

Раз строка.
Два строка.`,
		},
	}

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			mockLogger := mocks.NewMockLogger(t)
			mockLogger.EXPECT().Debug("Начата генерация шаблона hi.gohtml", mock.Anything)
			mockLogger.EXPECT().Debug("Сгенерирован текст шаблона hi.gohtml", mock.Anything)
			templating := pkgTemplating.NewTool(`handlers`, mockLogger)
			result := templating.RenderText(`hi.gohtml`, testData)

			assert.Equal(t, dataProvider.expected[testCase], result)
		})
	}
}

func TestTemplating_RenderEscapedText(t *testing.T) {
	dataProvider := struct {
		testData map[string]struct {
			data           any
			escapedStrings []string
		}
		expected map[string]string
	}{
		testData: map[string]struct {
			data           any
			escapedStrings []string
		}{
			"Не задано экранированных": {
				data: struct {
					URL     template.HTML
					Address string
				}{
					URL:     template.HTML(`https://example.org/foo/bar?param=value`),
					Address: `Москва, Кремль, 1`,
				},
				escapedStrings: []string{},
			},
			"Не задано экранированных, есть перенос": {
				data: struct {
					URL     template.HTML
					Address string
				}{
					URL:     template.HTML(`https://example.org/foo/bar?param=value`),
					Address: "Москва, Кремль, 1\nНовая строка.",
				},
				escapedStrings: []string{},
			},
			"Экранированная с переносом": {
				data: struct {
					URL     template.HTML
					Address string
				}{
					URL:     template.HTML(`https://example.org/foo/bar?param=value`),
					Address: "Москва, Кремль, 1\\nНовая строка.",
				},
				escapedStrings: []string{"Address"},
			},
		},
		expected: map[string]string{
			"Не задано экранированных": `Привет! Наш адрес: Москва, Кремль, 1

Наш сайт: <a href="https://example.org/foo/bar?param=value">сайт</a>.`,
			"Не задано экранированных, есть перенос": `Привет! Наш адрес: Москва, Кремль, 1
Новая строка.

Наш сайт: <a href="https://example.org/foo/bar?param=value">сайт</a>.`,
			"Экранированная с переносом": `Привет! Наш адрес: Москва, Кремль, 1
Новая строка.

Наш сайт: <a href="https://example.org/foo/bar?param=value">сайт</a>.`,
		},
	}

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			mockLogger := mocks.NewMockLogger(t)
			mockLogger.EXPECT().Debug("Начата генерация шаблона escaped_strings.gohtml", mock.Anything)
			mockLogger.EXPECT().Debug("Сгенерирован текст шаблона escaped_strings.gohtml", mock.Anything)
			templating := pkgTemplating.NewTool(`test`, mockLogger)
			result := templating.RenderEscapedText(`escaped_strings.gohtml`, testData.data, testData.escapedStrings)

			assert.Equal(t, dataProvider.expected[testCase], result)
		})
	}
}

package template_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
	mocks "s-belichenko/house-tg-bot/pkg/logger/mocks"
	pkgTemplating "s-belichenko/house-tg-bot/pkg/template"
)

type dataProviderSuccess struct {
	testData map[string]struct {
		storagePath string
		tmplName    string
		data        any
	}
	expected map[string]string
}

type dataProviderWrongPath struct {
	testData map[string]struct {
		storagePath string
		tmplName    string
	}
}

func TestTemplate_RenderTextSuccess(t *testing.T) {
	testURL, _ := url.Parse(`https://example.org/foo/bar?param=value`)
	dataProvider := dataProviderSuccess{
		testData: map[string]struct {
			storagePath string
			tmplName    string
			data        any
		}{
			`Обычный пример`: {
				storagePath: `pkg/template/test_resources/`,
				tmplName:    `template.txt`,
				data:        struct{ Name string }{Name: `John`},
			},
			`Слеш в имени файла`: {
				storagePath: `pkg/template/test_resources/`,
				tmplName:    `/template.txt`,
				data:        struct{ Name string }{Name: `John`},
			},
			`Якобы путь от корневой директории`: {
				storagePath: `/pkg/template/test_resources/`,
				tmplName:    `/template.txt`,
				data:        struct{ Name string }{Name: `John`},
			},
			`Нет слеша в конце пути к шаблонам`: {
				storagePath: `pkg/template/test_resources`,
				tmplName:    `/template.txt`,
				data:        struct{ Name string }{Name: `John`},
			},
			`Вызов метода свойства`: {
				storagePath: `pkg/template/test_resources`,
				tmplName:    `template_with_url.txt`,
				data: struct {
					Name string
					URL  *url.URL
				}{Name: `John`, URL: testURL},
			},
		},
		expected: map[string]string{
			`Обычный пример`:                    `Привет, John!`,
			`Слеш в имени файла`:                `Привет, John!`,
			`Якобы путь от корневой директории`: `Привет, John!`,
			`Нет слеша в конце пути к шаблонам`: `Привет, John!`,
			`Вызов метода свойства`:             `Привет, John! Вот твоя ссылка: https://example.org/foo/bar?param=value`,
		},
	}

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			mockLogger := mocks.NewMockLogger(t)
			templating := pkgTemplating.NewTemplate(testData.storagePath, mockLogger)
			result := templating.RenderText(testData.tmplName, testData.data)

			assert.Equal(t, dataProvider.expected[testCase], result)
		})
	}
}

func TestTemplate_RenderTextWrongPath(t *testing.T) {
	dataProvider := dataProviderWrongPath{
		testData: map[string]struct {
			storagePath string
			tmplName    string
		}{
			`Неверный путь`: {
				storagePath: `wrong_path_to_resources/`,
				tmplName:    `template.txt`,
			},
			`Неверное имя файла`: {
				storagePath: `test_resources`,
				tmplName:    `template1.txt`,
			},
		},
	}

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			mockLogger := mocks.NewMockLogger(t)
			// FIXME: Начать проверять через регулярные выражения текст ошибки.
			mockLogger.EXPECT().Error(mock.Anything, pkgLogger.LogContext(nil))
			templating := pkgTemplating.NewTemplate(testData.storagePath, mockLogger)
			result := templating.RenderText(testData.tmplName, nil)

			assert.Empty(t, result)
		})
	}
}

func TestTemplate_RenderTextRealTemplate(t *testing.T) {
	mockLogger := mocks.NewMockLogger(t)
	templating := pkgTemplating.NewTemplate(`resources/templates/text/handlers/`, mockLogger)
	result := templating.RenderText(`hi.txt`, nil)

	assert.Equal(t, "Привет! Я бот <a href=\"\">чата</a> дома по адресу . Правила добавления в чат:\n\n", result)
}

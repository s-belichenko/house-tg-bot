package template

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
)

type RenderInterface interface {
	RenderText(tmplName string, data any) string
}

type Templating struct {
	storagePath string
	log         pkgLogger.Logger
}

func (t *Templating) RenderText(tmplFilename string, data any) string {
	var renderedBuffer bytes.Buffer

	tmplPath := filepath.Join(t.storagePath, tmplFilename)
	if tmpl, err := template.ParseFiles(tmplPath); err != nil {
		t.log.Error(fmt.Sprintf(`Не удалось прочитать шаблон из файла "%s": %v`, tmplPath, err), nil)

		return ""
	} else { //nolint:revive
		if err := tmpl.Execute(&renderedBuffer, data); err != nil {
			t.log.Error(fmt.Sprintf(`Не удалось сгенерировать текст из шаблона "%s": %v`, tmplPath, err), nil)

			return ""
		}
	}

	return renderedBuffer.String()
}

// NewTemplate Создает средство шаблонизации.
//
// 'templatesDir' – путь до директории с шаблонами относительно их общего места хранения,
// 'log' – средство журналирования.
func NewTemplate(templatesDir string, log pkgLogger.Logger) *Templating {
	projectDir, err := getProjectPath()
	if err != nil {
		log.Fatal(fmt.Sprintf(`Не удалось определить путь до рабочей директории: %v`, err), nil)
	}

	return &Templating{
		storagePath: filepath.Join(projectDir, templatesDir),
		log:         log,
	}
}

func getProjectPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ``, errors.New(`could not get caller information`) //nolint:err113
	}

	absPath, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return ``, fmt.Errorf(`error getting absolute path: %w`, err)
	}

	return filepath.Join(absPath, `../..`), nil
}

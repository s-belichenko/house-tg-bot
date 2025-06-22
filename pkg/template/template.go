package template

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html"
	"html/template"
	"path/filepath"
	"reflect"
	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
	"strings"
)

//go:embed templates/**/*.txt
var templates embed.FS

type RenderingTool interface {
	RenderText(tmplName string, data any) string
	RenderEscapedText(tmplFilename string, data any, escapedStrings []string) string
}

type Templating struct {
	templatesPath string
	log           pkgLogger.Logger
}

// NewTool Создает средство шаблонизации.
//
// 'templatesDir' – имя директории с шаблонами в директории `templates`,
// 'log' – средство журналирования.
func NewTool(templatesDir string, log pkgLogger.Logger) *Templating {
	templatesPath := filepath.Join("templates", templatesDir)

	_, err := templates.ReadDir(templatesPath)
	if err != nil {
		log.Error(fmt.Sprintf(`Не удалось разрешить путь до шаблонов: %v`, err), nil)
	}

	return &Templating{
		templatesPath: templatesPath,
		log:           log,
	}
}

func (t *Templating) RenderEscapedText(
	tmplFilename string,
	data any,
	escapedStrings []string,
) string {
	unescapedData, err := t.unescapeData(&data, escapedStrings)
	if err != nil {
		t.log.Error(
			fmt.Sprintf("Не удалось сгенерировать текст шаблона %q: %e", tmplFilename, err),
			nil,
		)
	}

	return html.UnescapeString(t.renderText(tmplFilename, unescapedData))
}

func (t *Templating) RenderText(tmplFilename string, data any) string {
	return t.renderText(tmplFilename, data)
}

func (t *Templating) renderText(tmplFilename string, data any) string {
	var renderedBuffer bytes.Buffer

	t.log.Debug(
		fmt.Sprintf("Начата генерация шаблона %s", tmplFilename),
		pkgLogger.LogContext{"data": data},
	)

	tmplPath := filepath.Join(t.templatesPath, tmplFilename)

	tmpl, err := template.ParseFS(templates, tmplPath)
	if err != nil {
		t.log.Error(
			fmt.Sprintf(`Не удалось прочитать шаблон из файла "%s": %v`, tmplPath, err),
			nil,
		)

		return ""
	}

	if err := tmpl.Execute(&renderedBuffer, data); err != nil {
		t.log.Error(
			fmt.Sprintf(`Не удалось сгенерировать текст из шаблона "%s": %v`, tmplPath, err),
			nil,
		)

		return ""
	}

	result := renderedBuffer.String()
	t.log.Debug(fmt.Sprintf(
		"Сгенерирован текст шаблона %s", tmplFilename),
		map[string]interface{}{
			"rendered": result,
		},
	)

	return result
}

//nolint:err113
func (t *Templating) unescapeData(data any, escapedStrings []string) (interface{}, error) {
	repoError := errors.New("repositoryError")

	value := reflect.ValueOf(data)

	if value.Kind() != reflect.Ptr {
		return nil, fmt.Errorf(
			"%w: %w",
			repoError,
			fmt.Errorf("expected a pointer to a struct, got %v", value.Kind()),
		)
	}

	elem := value.Elem()
	tmp := reflect.New(elem.Elem().Type()).Elem()

	for _, fieldName := range escapedStrings {
		field := tmp.FieldByName(fieldName)
		if !field.IsValid() {
			return nil, fmt.Errorf(
				"%w: %w",
				repoError,
				fmt.Errorf("no such field: %s in struct", field),
			)
		}

		if !field.CanSet() {
			return nil, fmt.Errorf(
				"%w: %w",
				repoError,
				fmt.Errorf("cannot set field %s: not exportable or not settable", field),
			)
		}

		tmp.Set(elem.Elem())
		fieldValue := tmp.FieldByName(fieldName)
		tmp.FieldByName(fieldName).SetString(unescapeLineBreak(fieldValue.String()))
	}

	elem.Set(tmp)

	return data, nil
}

func unescapeLineBreak(str string) string {
	rawStringRepresentation := strings.ReplaceAll(str, "\\n", "\n")

	return rawStringRepresentation
}

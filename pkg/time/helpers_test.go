package time_test

import (
	"testing"

	"s-belichenko/house-tg-bot/pkg/time"

	"github.com/go-test/deep"
)

func TestParseDays(t *testing.T) {
	var dataProvider = struct {
		testData map[string]string
		expected map[string]struct {
			result    uint32
			errorText string
		}
	}{
		testData: map[string]string{
			"Успешный кейс":       "7",
			"Число больше unit32": "4294967296",
		},
		expected: map[string]struct {
			result    uint32
			errorText string
		}{
			"Успешный кейс": {
				result:    7,
				errorText: "",
			},
			"Число больше unit32": {
				result: 0,
				errorText: `не удалось распарсить days "4294967296" в uint16: ` +
					`strconv.ParseUint: parsing "4294967296": value out of range`,
			},
		},
	}

	for testCase, days := range dataProvider.testData {
		t.Run(testCase, func(t *testing.T) {
			r, err := time.ParseDays(days)
			for _, problem := range deep.Equal(r, dataProvider.expected[testCase].result) {
				t.Error(problem)
			}
			errText := ""
			if err != nil {
				errText = err.Error()
			}
			for _, problem := range deep.Equal(errText, dataProvider.expected[testCase].errorText) {
				t.Error(problem)
			}
		})
	}
}

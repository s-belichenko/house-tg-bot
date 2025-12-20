package time_test

import (
	"s-belichenko/house-tg-bot/pkg/time"
	"testing"

	"github.com/go-test/deep"
)

func TestParseDays(t *testing.T) {
	var dp = struct {
		testData map[string]string
		expected map[string]struct {
			result uint64
			error  error
		}
	}{
		testData: map[string]string{
			"Успешный кейс": "7",
		},
		expected: map[string]struct {
			result uint64
			error  error
		}{
			"Успешный кейс": {
				result: 7,
				error:  nil,
			},
		},
	}

	for testCase, days := range dp.testData {
		t.Run(testCase, func(t *testing.T) {
			r, err := time.ParseDays(days)
			for _, problem := range deep.Equal(r, dp.expected[testCase].result) {
				t.Error(problem)
			}
			for _, problem := range deep.Equal(err, dp.expected[testCase].error) {
				t.Error(problem)
			}
		})
	}
}

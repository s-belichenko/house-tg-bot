package logger_test

import (
	"testing"
	"time"

	pkgLogger "s-belichenko/house-tg-bot/pkg/logger"
	mocks "s-belichenko/house-tg-bot/pkg/logger/mocks"
	timeMocks "s-belichenko/house-tg-bot/pkg/time/mocks"
)

type dataProvider struct {
	testData map[string]struct {
		msg string
		ctx pkgLogger.LogContext
	}
	expected map[string]string
}

func TestYandexLogger_Debug(t *testing.T) {
	dataProvider := dataProvider{
		testData: map[string]struct {
			msg string
			ctx pkgLogger.LogContext
		}{
			"Debug без контекста": {
				msg: "Debug text",
				ctx: nil,
			},
			"Debug с контекстом": {
				msg: "Debug text",
				ctx: pkgLogger.LogContext{
					"user_id": 123,
					"chat_id": 456,
				},
			},
		},
		expected: map[string]string{
			"Debug без контекста": "{\"message\":\"Debug text\",\"level\":\"DEBUG\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2019-01-01T00:00:00Z\",\"extra\":null}",
			"Debug с контекстом": "{\"message\":\"Debug text\",\"level\":\"DEBUG\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2019-01-01T00:00:00Z\",\"extra\":{\"chat_id\":456,\"user_id\":123}}",
		},
	}

	sysLog := mocks.NewMockSystemLogger(t)
	clock := timeMocks.NewMockClockInterface(t)
	now := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			sysLog.On("Output", 2, dataProvider.expected[testCase]).
				Return(nil)
			clock.On("Now").Return(now)

			yandexLogger := pkgLogger.NewYandexLogger("stream", sysLog, clock)
			yandexLogger.Debug(testData.msg, testData.ctx)
		})
	}
}

func TestYandexLogger_Trace(t *testing.T) {
	dataProvider := dataProvider{
		testData: map[string]struct {
			msg string
			ctx pkgLogger.LogContext
		}{
			"Trace без контекста": {
				msg: "Trace text",
				ctx: nil,
			},
			"Trace с контекстом": {
				msg: "Trace text",
				ctx: pkgLogger.LogContext{
					"user_id": 234,
					"chat_id": 567,
				},
			},
		},
		expected: map[string]string{
			"Trace без контекста": "{\"message\":\"Trace text\",\"level\":\"TRACE\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2018-02-02T01:01:01.000000001Z\",\"extra\":null}",
			"Trace с контекстом": "{\"message\":\"Trace text\",\"level\":\"TRACE\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2018-02-02T01:01:01.000000001Z\",\"extra\":{\"chat_id\":567,\"user_id\":234}}",
		},
	}

	sysLog := mocks.NewMockSystemLogger(t)
	clock := timeMocks.NewMockClockInterface(t)
	now := time.Date(2018, 2, 2, 1, 1, 1, 1, time.UTC)

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			sysLog.On("Output", 2, dataProvider.expected[testCase]).
				Return(nil)
			clock.On("Now").Return(now)

			yandexLogger := pkgLogger.NewYandexLogger("stream", sysLog, clock)
			yandexLogger.Trace(testData.msg, testData.ctx)
		})
	}
}

func TestYandexLogger_Info(t *testing.T) {
	dataProvider := dataProvider{
		testData: map[string]struct {
			msg string
			ctx pkgLogger.LogContext
		}{
			"Info без контекста": {
				msg: "Info text",
				ctx: nil,
			},
			"Info с контекстом": {
				msg: "Info text",
				ctx: pkgLogger.LogContext{
					"user_id": 345,
					"chat_id": 678,
				},
			},
		},
		expected: map[string]string{
			"Info без контекста": "{\"message\":\"Info text\",\"level\":\"INFO\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2017-03-03T02:02:02Z\",\"extra\":null}",
			"Info с контекстом": "{\"message\":\"Info text\",\"level\":\"INFO\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2017-03-03T02:02:02Z\",\"extra\":{\"chat_id\":678,\"user_id\":345}}",
		},
	}

	sysLog := mocks.NewMockSystemLogger(t)
	clock := timeMocks.NewMockClockInterface(t)
	now := time.Date(2017, 3, 3, 2, 2, 2, 0, time.UTC)

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			sysLog.On("Output", 2, dataProvider.expected[testCase]).
				Return(nil)
			clock.On("Now").Return(now)

			yandexLogger := pkgLogger.NewYandexLogger("stream", sysLog, clock)
			yandexLogger.Info(testData.msg, testData.ctx)
		})
	}
}

func TestYandexLogger_Warn(t *testing.T) {
	dataProvider := dataProvider{
		testData: map[string]struct {
			msg string
			ctx pkgLogger.LogContext
		}{
			"Warn без контекста": {
				msg: "Warn text",
				ctx: nil,
			},
			"Warn с контекстом": {
				msg: "Warn text",
				ctx: pkgLogger.LogContext{
					"user_id": 456,
					"chat_id": 789,
				},
			},
		},
		expected: map[string]string{
			"Warn без контекста": "{\"message\":\"Warn text\",\"level\":\"WARN\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2016-04-04T03:03:03Z\",\"extra\":null}",
			"Warn с контекстом": "{\"message\":\"Warn text\",\"level\":\"WARN\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2016-04-04T03:03:03Z\",\"extra\":{\"chat_id\":789,\"user_id\":456}}",
		},
	}

	sysLog := mocks.NewMockSystemLogger(t)
	clock := timeMocks.NewMockClockInterface(t)
	now := time.Date(2016, 4, 4, 3, 3, 3, 0, time.UTC)

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			sysLog.On("Output", 2, dataProvider.expected[testCase]).
				Return(nil)
			clock.On("Now").Return(now)

			yandexLogger := pkgLogger.NewYandexLogger("stream", sysLog, clock)
			yandexLogger.Warn(testData.msg, testData.ctx)
		})
	}
}

func TestYandexLogger_Error(t *testing.T) {
	dataProvider := dataProvider{
		testData: map[string]struct {
			msg string
			ctx pkgLogger.LogContext
		}{
			"Error без контекста": {
				msg: "Error text",
				ctx: nil,
			},
			"Error с контекстом": {
				msg: "Error text",
				ctx: pkgLogger.LogContext{
					"user_id": 567,
					"chat_id": 890,
				},
			},
		},
		expected: map[string]string{
			"Error без контекста": "{\"message\":\"Error text\",\"level\":\"ERROR\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2015-05-05T04:04:04Z\",\"extra\":null}",
			"Error с контекстом": "{\"message\":\"Error text\",\"level\":\"ERROR\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2015-05-05T04:04:04Z\",\"extra\":{\"chat_id\":890,\"user_id\":567}}",
		},
	}

	sysLog := mocks.NewMockSystemLogger(t)
	clock := timeMocks.NewMockClockInterface(t)
	now := time.Date(2015, 5, 5, 4, 4, 4, 0, time.UTC)

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			sysLog.On("Output", 2, dataProvider.expected[testCase]).
				Return(nil)
			clock.On("Now").Return(now)

			yandexLogger := pkgLogger.NewYandexLogger("stream", sysLog, clock)
			yandexLogger.Error(testData.msg, testData.ctx)
		})
	}
}

func TestYandexLogger_Fatal(t *testing.T) {
	dataProvider := dataProvider{
		testData: map[string]struct {
			msg string
			ctx pkgLogger.LogContext
		}{
			"Fatal без контекста": {
				msg: "Fatal text",
				ctx: nil,
			},
			"Fatal с контекстом": {
				msg: "Fatal text",
				ctx: pkgLogger.LogContext{
					"user_id": 678,
					"chat_id": 901,
				},
			},
		},
		expected: map[string]string{
			"Fatal без контекста": "{\"message\":\"Fatal text\",\"level\":\"FATAL\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2014-06-06T05:05:05Z\",\"extra\":null}",
			"Fatal с контекстом": "{\"message\":\"Fatal text\",\"level\":\"FATAL\",\"stream_name\":\"stream\"" +
				",\"timestamp\":\"2014-06-06T05:05:05Z\",\"extra\":{\"chat_id\":901,\"user_id\":678}}",
		},
	}

	sysLog := mocks.NewMockSystemLogger(t)
	clock := timeMocks.NewMockClockInterface(t)
	now := time.Date(2014, 6, 6, 5, 5, 5, 0, time.UTC)

	for testCase, testData := range dataProvider.testData {
		t.Run(testCase, func(_ *testing.T) {
			sysLog.On("Output", 2, dataProvider.expected[testCase]).
				Return(nil)
			clock.On("Now").Return(now)

			yandexLogger := pkgLogger.NewYandexLogger("stream", sysLog, clock)
			yandexLogger.Fatal(testData.msg, testData.ctx)
		})
	}
}

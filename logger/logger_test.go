package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("", "", &LogOptions{})

	cases := []struct {
		message  string
		level    zapcore.Level
		callFunc func(message string, fields ...zap.Field)
	}{
		{
			message:  "debug message",
			level:    zap.DebugLevel,
			callFunc: logger.Debug,
		},
		{
			message:  "error message",
			level:    zap.ErrorLevel,
			callFunc: logger.Error,
		},
		{
			message:  "info message",
			level:    zap.InfoLevel,
			callFunc: logger.Info,
		},
		{
			message:  "warning message",
			level:    zap.WarnLevel,
			callFunc: logger.Warn,
		},
	}

	for _, c := range cases {
		c.callFunc(c.message)
		// assert.Equal(t, i+1, len(logsStorage))
		// assert.Equal(t, c.level, logsStorage[len(logsStorage)-1].Level)
		// assert.Equal(t, c.message, logsStorage[len(logsStorage)-1].Message)
	}
}

func BenchmarkLogger_Log(b *testing.B) {
	opt := &LogOptions{
		Level: LLvlProduction,
	}
	logger := NewLogger("", "", opt)

	for n := 0; n < b.N; n++ {
		logger.Info("a message", zap.String("some-key", "some-value"))
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	opt := &LogOptions{
		Level: LLvlProduction,
	}
	logger := NewLogger("", "", opt)

	for n := 0; n < b.N; n++ {
		logger.Info("a message", zap.String("something", "log"))
	}
}

func BenchmarkSugarLog_Log(b *testing.B) {
	opt := &LogOptions{
		Level: LLvlProduction,
	}
	logger := NewSugaredLogger("", "", opt)

	for n := 0; n < b.N; n++ {
		logger.Log("a message", nil, nil)
	}
}

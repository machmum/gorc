package main

import (
	"log"
	"os"

	"github.com/machmum/gorc/logger"
	"github.com/machmum/gorc/request"
	"go.uber.org/zap"
)

func main() {
	dir := "./log/oauth"

	opt := &logger.LogOptions{
		Level:     logger.LLvlDevelopment,
		WithTrace: true,
	}

	// stdSugar := logger.NewSugaredLogger(dir, "", opt)
	stdLogger := logger.NewLogger(dir, "", opt)

	stdLogger.Debug("standard debug message", zap.String("request", "a request"), zap.String("response", "a response"))
	stdLogger.Error("standard error message", zap.String("request", "a request"), zap.String("response", "a response"))
	stdLogger.Info("standard info message", zap.String("request", "a request"), zap.String("response", "a response"))
	stdLogger.Warn("standard warn message", zap.String("request", "a request"), zap.String("response", "a response"))

	log.Fatal()

	// stdSugar.Log("a full service", map[string]interface{}{"request": "a request", "response": "a response"}, nil)
	// stdSugar.Log("an empty service", nil, nil)
	// stdSugar.Log("a full error service", map[string]interface{}{"request": "a request", "response": "a response"}, fmt.Errorf("found an error in the service"))
	// stdSugar.Log("an empty error service", nil, fmt.Errorf("found an error in the service"))

	stdLogger.Info("========================================================================")

	opt = &logger.LogOptions{
		Level:     logger.LLvlDevelopment,
		WithTrace: true,
	}
	traceLogger := logger.NewLogger(dir, "", nil)
	traceLogger.Info("info message with trace-id", zap.String("request", "a request"), zap.String("response", "a response"))

	opt = &logger.LogOptions{
		Level: logger.LLvlDevelopment,
		RefID: request.RequestID(),
	}
	refLogger := logger.NewLogger(dir, "", opt)
	refLogger.Info("info message with trace-id and ref-id", zap.String("request", "a request"), zap.String("response", "a response"))

	opt = &logger.LogOptions{
		Level:      logger.LLvlDevelopment,
		OutputFile: []string{"stdout"},
	}
	stdoutLogger := logger.NewLogger(dir, "", opt)
	stdoutLogger.Info("info message printed to both file and stdout", zap.String("request", "a request"), zap.String("response", "a response"))

	stdLogger.Info("========================================================================")

	stdLogger.Fatal("standard fatal message", zap.String("request", "a request"), zap.String("response", "a response"))

	l := log.New(os.Stdout, "", log.LstdFlags)
	aFunc(l)
}

func aFunc(l *log.Logger) {
	l.Printf("test log")
	l.Fatal()
}

package main

import (
	"fmt"
	"log"
	"os"

	gl "github.com/machmum/gorc/log"
	req "github.com/machmum/gorc/request"
)

func main() {
	// logger := zaplog.New(cfg.LogFile, cfg.Server.Oauth.Development, cfg.Server.Oauth.Name)
	// get unique reqID

	opt := &gl.LogOptions{
		Level:      gl.LLvlDevelopment,
		WithTrace:  true,
		RefID:      req.RequestID(),
		OutputFile: nil,
	}
	logger := gl.NewLogger("./log/oauth", "", opt)

	logger.Log("a full service", map[string]interface{}{"request": "a request", "response": "a response"}, nil)
	logger.Log("an empty service", nil, nil)

	logger.Log("a full error service", map[string]interface{}{"request": "a request", "response": "a response"}, fmt.Errorf("found an error in the service"))
	logger.Log("an empty error service", nil, fmt.Errorf("found an error in the service"))

	// logger.Fatal("stop logger...")

	l := log.New(os.Stdout, "", log.LstdFlags)
	aFunc(l)
}

func aFunc(l *log.Logger) {
	l.Printf("test log")
	l.Fatal()
}

package main

import (
	"fmt"

	"github.com/machmum/gorc"
)

func main() {
	// logger := zaplog.New(cfg.LogFile, cfg.Server.Oauth.Development, cfg.Server.Oauth.Name)
	// get unique reqID

	opt := &gorc.LogOptions{
		Development: true,
		WithTrace:   true,
		RefID:       gorc.RequestID(),
		OutputFile:  nil,
	}
	logger := gorc.NewLogger("./log/oauth", "", opt)

	logger.Log("a full service", map[string]interface{}{"request": "a request", "response": "a response"}, nil)
	logger.Log("an empty service", nil, nil)

	logger.Log("a full error service", map[string]interface{}{"request": "a request", "response": "a response"}, fmt.Errorf("found an error in the service"))
	logger.Log("an empty error service", nil, fmt.Errorf("found an error in the service"))

	logger.Fatal("stop logger...")
}

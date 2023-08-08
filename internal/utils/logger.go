package utils

import (
	"encoding/json"
	"go.uber.org/zap"
	"strings"
)

var (
	Logger *zap.Logger
	cfg    zap.Config
)

/*
Initializes the logger that is used for entire program. Powered by Zap (https://github.com/uber-go/zap.
*/
func InitializeLogger() {

	rawJSON := []byte(`{
	  "level": "info",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stdout"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	l := zap.Must(cfg.Build())
	defer l.Sync()
	Logger = l
	Logger.Info("Successfully initialized logger (Powered by Zap).")
}

/*
Sets the logger level. Accepts "DEBUG" or "INFO".
*/
func SetLoggerLevel(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		cfg.Level.SetLevel(zap.DebugLevel)
	case "INFO":
		cfg.Level.SetLevel(zap.InfoLevel)
	default:
		cfg.Level.SetLevel(zap.InfoLevel)
	}
}

package logs

import (
	log "github.com/sirupsen/logrus"
	"os"
)

var logger *log.Logger

const (
	environment     string = "development"
	defaultLogLevel string = "INFO"
)

func init() {
	// do something here to set environment depending on an environment variable
	// or command-line flag
	// TODO: should dynamic control the env
	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

func GetLogger(level string) *log.Logger {
	log.SetLevel(getLoggerLevel(level))
	return logger
}

func getLoggerLevel(level string) log.Level {
	switch level {
	case "INFO":
		return log.InfoLevel
	case "WARN":
		return log.WarnLevel
	case "ERROR":
		return log.ErrorLevel
	case "DEBUG":
		return log.DebugLevel
	case "TRACE":
		return log.TraceLevel
	case "FATAL":
		return log.FatalLevel
	case "PANIC":
		return log.PanicLevel
	default:
		return log.WarnLevel
	}
}

package logs

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

const (
	environment     string = "development"
	defaultLogLevel string = "INFO"
)

func init() {
	// do something here to set environment depending on an environment variable
	// or command-line flag
	// TODO: should dynamic control the env
	if environment == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)
}

func GetLogger(level string) *logrus.Logger {
	if len(level) == 0 {
		level = defaultLogLevel
	}
	logger.SetLevel(getLoggerLevel(level))
	return logger
}

func getLoggerLevel(level string) logrus.Level {
	switch level {
	case "INFO":
		return logrus.InfoLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	case "FATAL":
		return logrus.FatalLevel
	case "PANIC":
		return logrus.PanicLevel
	default:
		return logrus.WarnLevel
	}
}

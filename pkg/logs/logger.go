package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	environment     string = "development"
	defaultLogLevel string = "INFO"
)

var rwMutex = &sync.RWMutex{}
var loggersMap = make(map[string]*logrus.Logger)
var (
	logursMap = map[string]logrus.Level{
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
		"DEBUG": logrus.DebugLevel,
		"TRACE": logrus.TraceLevel,
		"FATAL": logrus.FatalLevel,
		"PANIC": logrus.PanicLevel,
	}
	logLevel = logrus.InfoLevel
)

type LogFormatter struct {
	loggerName string
}

func init() {
	logursLevel, ok := logursMap[defaultLogLevel]
	if !ok {
		fmt.Printf("unsupported default Log Level - [%s]. Set to default Log Level - [%s]", defaultLogLevel, defaultLogLevel)
	} else {
		logLevel = logursLevel
	}
}

//func (lf *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
//	logEntry := fmt.Sprintf("%s %s [%s] - %s\n", entry.Time.Format(time.RFC3339Nano), getLoggerLevel(entry.Level), lf.loggerName, entry.Message)
//	return []byte(logEntry), nil
//}

func GetLogger(level string) *logrus.Logger {
	rwMutex.RLock()
	logger := loggersMap[level]
	rwMutex.RUnlock()

	if logger == nil {
		logger = newLogger()

		rwMutex.Lock()
		loggersMap[level] = logger
		rwMutex.Unlock()
	}
	return logger
}

func getLoggerLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.InfoLevel:
		return "INFO"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.WarnLevel:
		return "WARN"
	case logrus.PanicLevel:
		return "PANIC"
	case logrus.FatalLevel:
		return "FATAL"
	}
	return "UNKNOWN"
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logLevel)
	return logger
}

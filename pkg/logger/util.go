package logger

import (
	"errors"
	"strings"
)

var levelNames = map[LogLevel]string{
	LogLevelTrace: "TRACE",
	LogLevelDebug: "DEBUG",
	LogLevelInfo:  "INFO",
	LogLevelWarn:  "WARN",
	LogLevelError: "ERROR",
}

func getLevelString(level LogLevel) (string, error) {
	if name, ok := levelNames[level]; ok {
		return name, nil
	}
	return "", errors.New("invalid log level")
}

func getLevelInt(level string) (LogLevel, error) {
	upper := strings.ToUpper(level)
	for k, v := range levelNames {
		if v == upper {
			return k, nil
		}
	}
	return LogLevelNone, errors.New("invalid log level")
}

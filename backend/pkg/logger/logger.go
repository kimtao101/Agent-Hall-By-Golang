package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
}

type DefaultLogger struct {
	logger *log.Logger
}

func NewDefaultLogger(prefix string) *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, prefix+" ", log.LstdFlags),
	}
}

func (l *DefaultLogger) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("%s %+v", msg, fields[0])
	} else {
		l.logger.Println(msg)
	}
}

func (l *DefaultLogger) Error(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("ERROR: %s %+v", msg, fields[0])
	} else {
		l.logger.Printf("ERROR: %s", msg)
	}
}

func (l *DefaultLogger) Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("WARN: %s %+v", msg, fields[0])
	} else {
		l.logger.Printf("WARN: %s", msg)
	}
}

package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
}

type logger struct {
	infoLogger *log.Logger
	warnLogger *log.Logger
	errLogger  *log.Logger
}

func NewLogger() Logger {
	flags := log.Ldate | log.Ltime | log.Lshortfile
	return &logger{
		infoLogger: log.New(os.Stdout, "INFO: ", flags),
		warnLogger: log.New(os.Stdout, "WARN: ", flags),
		errLogger:  log.New(os.Stdout, "ERROR: ", flags),
	}
}

func (l *logger) Info(v ...any) {
	l.infoLogger.Println(v)
}
func (l *logger) Warn(v ...any) {
	l.warnLogger.Println(v)
}

func (l *logger) Error(v ...any) {
	l.errLogger.Println(v)
}

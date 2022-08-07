package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type MyLogger struct {
	ErrorWard *log.Logger
	Informer  *log.Logger
	Debugger  *log.Logger
	DebugOn   bool
}

func NewLogger() *MyLogger {
	return &MyLogger{
		ErrorWard: log.New(os.Stderr, "[ERROR]\t", log.Lshortfile),
		Informer:  log.New(os.Stdout, " [INFO]\t", log.Ltime),
		Debugger:  log.New(os.Stdout, "[DEBUG]\t", log.Ltime),
	}
}

func (l *MyLogger) Debug(caller string, msg string, args ...interface{}) {
	if l.DebugOn {
		msg = fmt.Sprintf(msg, args...)
		l.Debugger.Printf("[%s]: %s", strings.ToUpper(caller), msg)
	}
}

func (l *MyLogger) Info(caller string, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	l.Informer.Printf("[%s]: %s", strings.ToUpper(caller), msg)
}

func (l *MyLogger) Error(err error) {
	l.ErrorWard.Println(err)
}

func (l *MyLogger) Fatal(err error) {
	l.ErrorWard.Fatal(err)
}

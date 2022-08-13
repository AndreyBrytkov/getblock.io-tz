package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type MyLogger struct {
	errorWard *log.Logger
	informer  *log.Logger
	debugger  *log.Logger
	DebugOn   bool
}

func NewLogger() *MyLogger {
	return &MyLogger{
		errorWard: log.New(os.Stderr, "[ERROR]\t", log.Lshortfile),
		informer:  log.New(os.Stdout, " [INFO]\t", log.Ltime),
		debugger:  log.New(os.Stdout, "[DEBUG]\t", log.Ltime),
	}
}

func (l *MyLogger) Debug(caller string, msg string, args ...interface{}) {
	if l.DebugOn {
		msg = fmt.Sprintf(msg, args...)
		l.debugger.Printf("[%s]: %s", strings.ToUpper(caller), msg)
	}
}

func (l *MyLogger) Info(caller string, msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	l.informer.Printf("[%s]: %s", strings.ToUpper(caller), msg)
}

func (l *MyLogger) Error(err error) {
	l.errorWard.Println(err)
}

func (l *MyLogger) Fatal(err error) {
	l.errorWard.Fatal(err)
}

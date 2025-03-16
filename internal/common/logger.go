package common

import (
	"fmt"
)

type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, err error, keysAndValues ...interface{})
	Errorw(msg string, err error, keysAndValues ...interface{})
}

func printMessage(level, msg string, keysAndValues ...interface{}) {
	msg = fmt.Sprintf(msg, keysAndValues...)
	fmt.Printf("%s: %s\n", level, msg)
}

func printMessageWithError(level, msg string, err error, keysAndValues ...interface{}) {
	msg = fmt.Sprintf(msg, keysAndValues...)
	fmt.Printf("%s: %s (err = %s)\n", level, msg, err.Error())
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Debugw(msg string, keysAndValues ...interface{}) {
	printMessage("DEBUG", msg, keysAndValues...)
}

func (c ConsoleLogger) Infow(msg string, keysAndValues ...interface{}) {
	printMessage("INFO", msg, keysAndValues...)
}

func (c ConsoleLogger) Warnw(msg string, err error, keysAndValues ...interface{}) {
	printMessageWithError("WARN", msg, err, keysAndValues...)
}

func (c ConsoleLogger) Errorw(msg string, err error, keysAndValues ...interface{}) {
	printMessageWithError("ERROR", msg, err, keysAndValues...)
}

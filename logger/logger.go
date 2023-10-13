package logger

import (
	"log"
	"os"
)

const (
	COLOR_GREEN = "\033[32m"
	COLOR_RED   = "\033[31m"
	COLOR_RESET = "\033[0m"
)

var (
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
)

func Info(format string, v ...interface{}) {
	msg := format
	logger.Printf(msg, v...)
}

func Error(format string, v ...interface{}) {
	msg := COLOR_RED + format + COLOR_RESET
	logger.Printf(msg, v...)
}

func Success(format string, v ...interface{}) {
	msg := COLOR_GREEN + format + COLOR_RESET
	logger.Printf(msg, v...)
}
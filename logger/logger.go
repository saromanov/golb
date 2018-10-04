package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

// Init provides initialization of the logger
func Init() {
	logger = logrus.StandardLogger().WithFields(logrus.Fields{})
	logrus.SetOutput(os.Stdout)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Print(args ...interface{}) {
	logger.Print(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

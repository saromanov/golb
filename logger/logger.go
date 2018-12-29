package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
}
func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Print(args ...interface{}) {
	logrus.Print(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warning(args ...interface{}) {
	logrus.Warning(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

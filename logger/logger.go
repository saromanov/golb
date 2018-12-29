package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
}

// Configuration provides config of logrus output
func Configuration(out string, l string) {

	switch out {
	case "stderr":
		logrus.SetOutput(os.Stderr)
	case "stdout":
		logrus.SetOutput(os.Stdout)
	default:
		f, err := os.OpenFile(out, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.SetOutput(f)
	}
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

package logger

import (
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"log"
	"os"
)

var AppLogger *logrus.Logger

func NewLogger(fileName string) *logrus.Logger {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to create logfile: " + fileName)
	}

	AppLogger = &logrus.Logger{
		// Log into f file handler and on os.Stdout
		Out:   io.MultiWriter(file, os.Stdout),
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}
	return AppLogger
}

func NewConsoleLogger() *logrus.Logger {
	AppLogger = &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout),
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}

	return AppLogger
}

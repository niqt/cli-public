package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	logger *log.Logger
}

var instance *Logger
var once sync.Once

func GetLogger() *Logger {
	once.Do(func() {
		logFile, err := os.OpenFile("cli.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		instance = &Logger{logger: log.New(logFile, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)}
	})
	return instance
}

func (l *Logger) Info(message string) {
	l.logger.Println("[INFO]: " + message)
}

func (l *Logger) Error(message string) {
	l.logger.Println("[ERROR]: " + message)
}

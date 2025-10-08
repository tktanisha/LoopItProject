package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

//go:generate mockgen -source=interface.go -destination=../../internal/mock/mock_logger.go -package=mock

var OpenFile = os.OpenFile
var OsExit = os.Exit
var logFileName = "logs/app.log" // default path

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type logMessage struct {
	level   LogLevel
	message string
	time    time.Time
}

type Logger struct {
	file    *os.File
	ch      chan logMessage
	wg      sync.WaitGroup
	once    sync.Once
	closing bool
}

var instance *Logger
var once sync.Once

func SetLogFileName(name string) {
	logFileName = name
}

func GetLogger() *Logger {
	once.Do(func() {
		file, err := OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		instance = &Logger{
			file: file,
			ch:   make(chan logMessage, 100),
		}

		instance.wg.Add(1)
		go instance.processLogs()
	})
	return instance
}

func (l *Logger) processLogs() {
	defer l.wg.Done()
	for msg := range l.ch {
		formatted := l.formatMessage(msg)

		// fmt.Println(formatted)

		l.file.WriteString(formatted + "\n")
	}
}

func (l *Logger) formatMessage(msg logMessage) string {
	var levelStr string
	switch msg.level {
	case DEBUG:
		levelStr = "[DEBUG]"
	case INFO:
		levelStr = "[INFO]"
	case WARNING:
		levelStr = "[WARNING]"
	case ERROR:
		levelStr = "[ERROR]"
	case FATAL:
		levelStr = "[FATAL]"
	}

	return fmt.Sprintf("%s %s %s",
		msg.time.Format("2006-01-02 15:04:05"),
		levelStr,
		msg.message,
	)
}

func (l *Logger) log(level LogLevel, msg string) {
	if l.closing {
		return
	}
	l.ch <- logMessage{level: level, time: time.Now(), message: msg}
}

func (l *Logger) Debug(msg string)   { l.log(DEBUG, msg) }
func (l *Logger) Info(msg string)    { l.log(INFO, msg) }
func (l *Logger) Warning(msg string) { l.log(WARNING, msg) }
func (l *Logger) Error(msg string)   { l.log(ERROR, msg) }
func (l *Logger) Fatal(msg string) {
	l.log(FATAL, msg)
	l.Close()
	OsExit(1)
}

func (l *Logger) Close() {
	l.once.Do(func() {
		l.closing = true
		close(l.ch)
		l.wg.Wait()
		l.file.Close()
	})
}

// Add to reset singleton (used in tests)
func ResetLoggerForTest() {
	instance = nil
	once = sync.Once{}
}

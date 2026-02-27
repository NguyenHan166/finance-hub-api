package logger

import (
	"log"
	"os"
)

// Logger provides structured logging
type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		Info:  log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn:  log.New(os.Stdout, "WARN:  ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		Debug: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Global logger instance
var Log = New()

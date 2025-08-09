package main

import "log"

// Logger
type Logger interface {
	Info(msg string)
	Error(msg string)
}

type AppLogger struct {
	level string
}

func NewLogger(level string) *AppLogger {
	return &AppLogger{level: level}
}

func (l *AppLogger) Info(msg string) {
	log.Printf("[INFO] %s", msg)
}

func (l *AppLogger) Error(msg string) {
	log.Printf("[ERROR] %s", msg)
}

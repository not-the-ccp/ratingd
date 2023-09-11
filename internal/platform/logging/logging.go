package logging

import (
	"log"
	"os"
)

// Logger is the structured logger used across the service.
type Logger interface {
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)
}

type stdLogger struct {
	out *log.Logger
}

func New(prefix string) Logger {
	return &stdLogger{out: log.New(os.Stderr, prefix+" ", log.LstdFlags|log.LUTC)}
}

func (l *stdLogger) Info(msg string, kv ...any) {
	l.out.Println(append([]any{"INFO", msg}, kv...)...)
}

func (l *stdLogger) Warn(msg string, kv ...any) {
	l.out.Println(append([]any{"WARN", msg}, kv...)...)
}

func (l *stdLogger) Error(msg string, kv ...any) {
	l.out.Println(append([]any{"ERROR", msg}, kv...)...)
}

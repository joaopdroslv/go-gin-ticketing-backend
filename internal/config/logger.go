package config

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger() *slog.Logger {

	fileWriter := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    50, // mbs
		MaxBackups: 10,
		MaxAge:     7, // days
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, fileWriter)
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{AddSource: true})

	return slog.New(handler)
}

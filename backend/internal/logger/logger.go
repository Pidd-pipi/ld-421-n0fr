package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New 创建结构化 slog Logger，默认输出 JSON。
func New(level string) *slog.Logger {
	lvl := slog.LevelInfo
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}
	opts := &slog.HandlerOptions{Level: lvl}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

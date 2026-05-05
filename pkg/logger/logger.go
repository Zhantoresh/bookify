package logger

import (
	"log/slog"
	"os"
	"time"
)

func New(level string, location *time.Location) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	if location == nil {
		location = time.UTC
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				if ts, ok := attr.Value.Any().(time.Time); ok {
					return slog.String(slog.TimeKey, ts.In(location).Format(time.RFC3339))
				}
			}
			return attr
		},
	}))
}

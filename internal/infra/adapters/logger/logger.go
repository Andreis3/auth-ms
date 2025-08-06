package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/otel/trace"
)

const LevelCritical = slog.LevelError + 1

type Logger struct {
	loggerJSON slog.Logger
	loggerText slog.Logger
}

func NewLogger() *Logger {
	o := os.Stdout
	loggerJSON := slog.New(slog.NewJSONHandler(o, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				switch a.Value.Any().(type) {
				case slog.Level:
					level := a.Value.Any().(slog.Level)
					switch level {
					case LevelCritical:
						a.Value = slog.StringValue("CRITICAL")
					case slog.LevelDebug:
						a.Value = slog.StringValue("DEBUG")
					case slog.LevelInfo:
						a.Value = slog.StringValue("INFO")
					case slog.LevelWarn:
						a.Value = slog.StringValue("WARN")
					case slog.LevelError:
						a.Value = slog.StringValue("ERROR")
					default:
						a.Value = slog.StringValue(level.String())
					}
				}
			}

			if a.Key == slog.TimeKey {
				format := "01-02-2006 15:04:05.000"
				a.Value = slog.StringValue(time.Now().Format(format))
			}
			return a
		},
	}))

	e := os.Stderr
	loggerText := slog.New(
		tint.NewHandler(e, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "01-02-2006 15:04:05.000",
			NoColor:    false,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "password" {
					a.Value = slog.StringValue("********")
					return a
				}

				if a.Key == slog.LevelKey {
					switch a.Value.Any().(type) {
					case slog.Level:
						level := a.Value.Any().(slog.Level)
						switch level {
						case LevelCritical:
							a.Value = slog.StringValue("\033[31mCRITICAL\033[0m")
						case slog.LevelDebug:
							a.Value = slog.StringValue("\033[34mDEBUG\033[0m")
						case slog.LevelInfo:
							a.Value = slog.StringValue("\033[32mINFO\033[0m")
						case slog.LevelWarn:
							a.Value = slog.StringValue("\033[33mWARN\033[0m")
						case slog.LevelError:
							a.Value = slog.StringValue("\033[31mERROR\033[0m")
						default:
							a.Value = slog.StringValue(level.String())
						}
					}
				}
				return a
			},
		}),
	)
	slog.SetDefault(loggerJSON)
	slog.SetDefault(loggerText)

	logger := &Logger{
		loggerJSON: *loggerJSON,
		loggerText: *loggerText,
	}

	return logger
}

func (l *Logger) DebugJSON(msg string, info ...any) {
	l.loggerJSON.Debug(msg, info...)
}

func (l *Logger) InfoJSON(msg string, info ...any) {
	l.loggerJSON.Info(msg, info...)
}

func (l *Logger) WarnJSON(msg string, info ...any) {
	l.loggerJSON.Warn(msg, info...)
}

func (l *Logger) ErrorJSON(msg string, info ...any) {
	l.loggerJSON.Error(msg, info...)
}

func (l *Logger) CriticalJSON(msg string, info ...any) {
	l.loggerJSON.Log(context.Background(), LevelCritical, msg, info...) // Nível crítico = 5 (LevelError + 1)
}

func (l *Logger) DebugText(msg string, info ...any) {
	l.loggerText.Debug(msg, info...)
}

func (l *Logger) InfoText(msg string, info ...any) {
	l.loggerText.Info(msg, info...)
}

func (l *Logger) WarnText(msg string, info ...any) {
	l.loggerText.Warn(msg, info...)
}

func (l *Logger) ErrorText(msg string, info ...any) {
	l.loggerText.Error(msg, info...)
}

func (l *Logger) CriticalText(msg string, info ...any) {
	l.loggerText.Log(context.Background(), LevelCritical, msg, info...) // Nível crítico = 5 (LevelError + 1)
}

func (l *Logger) WithTrace(ctx context.Context) *slog.Logger {
	spanCtx := trace.SpanContextFromContext(ctx)

	if !spanCtx.HasTraceID() {
		return &l.loggerJSON
	}
	return l.loggerJSON.With(
		slog.String("trace_id", spanCtx.TraceID().String()),
		slog.String("span_id", spanCtx.SpanID().String()),
	)
}

func (l *Logger) SlogJSON() *slog.Logger {
	return &l.loggerJSON
}

func (l *Logger) SlogText() *slog.Logger {
	return &l.loggerText
}

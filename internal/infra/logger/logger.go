package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/samber/slog-multi"
	"go.opentelemetry.io/otel/trace"
)

// LevelCritical = Error + 1
const LevelCritical = slog.LevelError + 1

type Logger struct {
	json *slog.Logger
	text *slog.Logger
	all  *slog.Logger // multi (json+text)
}

func NewLogger() *Logger {
	// outputs
	outJSON := os.Stdout
	outText := os.Stderr

	jsonHandler := slog.NewJSONHandler(outJSON, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: replaceCommon(func(a slog.Attr) slog.Attr {
			// JSON: keep standard time (RFC3339); only normalize level and mask secrets.
			return a
		}),
	})

	textHandler := tint.NewHandler(outText, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "01-02-2006 15:04:05.000",
		NoColor:    false,
		ReplaceAttr: replaceCommon(func(a slog.Attr) slog.Attr {
			// tint already formats time; only level coloring is already handled by tint.
			return a
		}),
	})

	multi := slog.New(slogmulti.Fanout(jsonHandler, textHandler))

	l := &Logger{
		json: slog.New(jsonHandler),
		text: slog.New(textHandler),
		all:  multi,
	}

	// if set default logger to multi (both)
	slog.SetDefault(multi)

	return l
}

// replaceCommon applies common normalizations to the handlers.
func replaceCommon(next func(slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		// mask sensitive fields
		switch a.Key {
		case "password", "secret", "token", "authorization", "api_key":
			return slog.String(a.Key, "********")
		}
		// normalize level (includes CRITICAL)
		if a.Key == slog.LevelKey {
			if lvl, ok := a.Value.Any().(slog.Level); ok {
				switch lvl {
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
					a.Value = slog.StringValue(lvl.String())
				}
			}
		}
		// delegates to the specific handler (text/json) if further processing is needed
		return next(a)
	}
}

// ---------- helpers de attrs ----------

func buildAttrs(args []any) []any {
	if len(args) == 0 {
		return nil
	}
	attrs := make([]any, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]any:
			for k, val := range v {
				attrs = append(attrs, slog.Any(k, val))
			}
		case slog.Attr:
			attrs = append(attrs, v)
		default:
			attrs = append(attrs, v)
		}
	}
	return attrs
}

// ---------- public APIs ----------

func (l *Logger) DebugJSON(msg string, args ...any) { l.json.Debug(msg, buildAttrs(args)...) }
func (l *Logger) InfoJSON(msg string, args ...any)  { l.json.Info(msg, buildAttrs(args)...) }
func (l *Logger) WarnJSON(msg string, args ...any)  { l.json.Warn(msg, buildAttrs(args)...) }
func (l *Logger) ErrorJSON(msg string, args ...any) { l.json.Error(msg, buildAttrs(args)...) }
func (l *Logger) CriticalJSON(msg string, args ...any) {
	l.json.Log(context.Background(), LevelCritical, msg, buildAttrs(args)...)
}

func (l *Logger) DebugText(msg string, args ...any) { l.text.Debug(msg, buildAttrs(args)...) } // <-- fix: usar attrs
func (l *Logger) InfoText(msg string, args ...any)  { l.text.Info(msg, buildAttrs(args)...) }
func (l *Logger) WarnText(msg string, args ...any)  { l.text.Warn(msg, buildAttrs(args)...) }
func (l *Logger) ErrorText(msg string, args ...any) { l.text.Error(msg, buildAttrs(args)...) }
func (l *Logger) CriticalText(msg string, args ...any) {
	l.text.Log(context.Background(), LevelCritical, msg, buildAttrs(args)...)
}

// Logs to both outputs simultaneously (multi)
func (l *Logger) Debug(msg string, args ...any) { l.all.Debug(msg, buildAttrs(args)...) }
func (l *Logger) Info(msg string, args ...any)  { l.all.Info(msg, buildAttrs(args)...) }
func (l *Logger) Warn(msg string, args ...any)  { l.all.Warn(msg, buildAttrs(args)...) }
func (l *Logger) Error(msg string, args ...any) { l.all.Error(msg, buildAttrs(args)...) }
func (l *Logger) Critical(msg string, args ...any) {
	l.all.Log(context.Background(), LevelCritical, msg, buildAttrs(args)...)
}

// With trace-id and span-id (returns the three chained loggers)
type WithTraceLoggers struct {
	JSON *slog.Logger
	Text *slog.Logger
	All  *slog.Logger
}

func (l *Logger) WithTrace(ctx context.Context) (*slog.Logger, *slog.Logger, *slog.Logger) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.HasTraceID() {
		return l.json, l.text, l.all // se não tiver trace_id, retorna os loggers crus
	}
	kvs := []any{
		slog.String("trace_id", spanCtx.TraceID().String()),
		slog.String("span_id", spanCtx.SpanID().String()),
	}
	return l.json.With(kvs...), l.text.With(kvs...), l.all.With(kvs...)
}

// Expose the raw loggers (if needed)
func (l *Logger) SlogJSON() *slog.Logger { return l.json }
func (l *Logger) SlogText() *slog.Logger { return l.text }
func (l *Logger) SlogAll() *slog.Logger  { return l.all }

// Optional: build a logger that writes to a custom writer (for tests)
func NewForWriter(jsonW, textW io.Writer, level slog.Leveler) *Logger {
	jsonHandler := slog.NewJSONHandler(jsonW, &slog.HandlerOptions{Level: level, ReplaceAttr: replaceCommon(func(a slog.Attr) slog.Attr { return a })})
	textHandler := tint.NewHandler(textW, &tint.Options{Level: level, TimeFormat: time.RFC3339})
	return &Logger{json: slog.New(jsonHandler), text: slog.New(textHandler), all: slog.New(slogmulti.Fanout(jsonHandler, textHandler))}
}

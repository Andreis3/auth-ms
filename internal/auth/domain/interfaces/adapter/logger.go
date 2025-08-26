package adapter

import (
	"context"
	"log/slog"
)

// WithTraceLoggers carrega os loggers já decorados com trace/span.
type WithTraceLoggers struct {
	JSON *slog.Logger
	Text *slog.Logger
	All  *slog.Logger
}

type Logger interface {
	// ---- Saída unificada (JSON + Text) ----
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Critical(msg string, args ...any)

	// ---- Somente JSON ----
	DebugJSON(msg string, args ...any)
	InfoJSON(msg string, args ...any)
	WarnJSON(msg string, args ...any)
	ErrorJSON(msg string, args ...any)
	CriticalJSON(msg string, args ...any)

	// ---- Somente Text ----
	DebugText(msg string, args ...any)
	InfoText(msg string, args ...any)
	WarnText(msg string, args ...any)
	ErrorText(msg string, args ...any)
	CriticalText(msg string, args ...any)

	// Retorna loggers decorados com trace_id/span_id
	WithTrace(ctx context.Context) (json *slog.Logger, text *slog.Logger, all *slog.Logger)

	// Acessores crus (se precisar integrar libs externas)
	SlogJSON() *slog.Logger
	SlogText() *slog.Logger
	SlogAll() *slog.Logger
}

package iadapter

import (
	"context"
	"log/slog"
)

type ILogger interface {
	DebugJSON(msg string, args ...any)
	InfoJSON(msg string, args ...any)
	WarnJSON(msg string, args ...any)
	ErrorJSON(msg string, args ...any)
	CriticalJSON(msg string, args ...any)
	DebugText(msg string, args ...any)
	InfoText(msg string, args ...any)
	WarnText(msg string, args ...any)
	ErrorText(msg string, args ...any)
	CriticalText(msg string, args ...any)
	WithTrace(ctx context.Context) *slog.Logger
	SlogJSON() *slog.Logger
	SlogText() *slog.Logger
}

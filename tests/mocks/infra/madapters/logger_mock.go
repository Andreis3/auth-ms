package madapters

import (
	"context"
	"io"
	"log/slog"

	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/domain/interfaces/iadapter"
)

var _ iadapter.ILogger = (*LoggerMock)(nil)

type LoggerMock struct{ mock.Mock }

// JSON (atenção ao nome em CAIXA ALTA)
func (l *LoggerMock) DebugJSON(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) InfoJSON(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) WarnJSON(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) ErrorJSON(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) CriticalJSON(msg string, args ...any) { l.Called(append([]any{msg}, args...)...) }

// Text
func (l *LoggerMock) DebugText(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) InfoText(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) WarnText(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) ErrorText(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) CriticalText(msg string, args ...any) { l.Called(append([]any{msg}, args...)...) }

func (l *LoggerMock) WithTrace(ctx context.Context) *slog.Logger {
	args := l.Called(ctx)
	if v := args.Get(0); v != nil {
		return v.(*slog.Logger)
	}
	// fallback útil pra não precisar configurar sempre:
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func (l *LoggerMock) SlogJSON() *slog.Logger {
	args := l.Called()
	if v := args.Get(0); v != nil {
		return v.(*slog.Logger)
	}
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func (l *LoggerMock) SlogText() *slog.Logger {
	args := l.Called()
	if v := args.Get(0); v != nil {
		return v.(*slog.Logger)
	}
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

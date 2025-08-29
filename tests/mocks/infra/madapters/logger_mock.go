package madapters

import (
	"context"
	"io"
	"log/slog"

	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
)

var _ adapter.Logger = (*LoggerMock)(nil)

type LoggerMock struct{ mock.Mock }

// ---- JSON ----
func (l *LoggerMock) DebugJSON(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) InfoJSON(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) WarnJSON(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) ErrorJSON(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) CriticalJSON(msg string, args ...any) { l.Called(append([]any{msg}, args...)...) }

// ---- Text ----
func (l *LoggerMock) DebugText(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) InfoText(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) WarnText(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) ErrorText(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) CriticalText(msg string, args ...any) { l.Called(append([]any{msg}, args...)...) }

// ---- Unificado (JSON + Text) ----
func (l *LoggerMock) Debug(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) Info(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) Warn(msg string, args ...any)     { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) Error(msg string, args ...any)    { l.Called(append([]any{msg}, args...)...) }
func (l *LoggerMock) Critical(msg string, args ...any) { l.Called(append([]any{msg}, args...)...) }

// ---- WithTrace: 3 retornos ----
func (l *LoggerMock) WithTrace(ctx context.Context) (*slog.Logger, *slog.Logger, *slog.Logger) {
	args := l.Called(ctx)
	return getLogger(args, 0), getLogger(args, 1), getLogger(args, 2)
}

// ---- Acessores crus ----
func (l *LoggerMock) SlogJSON() *slog.Logger {
	args := l.Called()
	if v := args.Get(0); v != nil {
		return v.(*slog.Logger)
	}
	return newDiscardJSON()
}

func (l *LoggerMock) SlogText() *slog.Logger {
	args := l.Called()
	if v := args.Get(0); v != nil {
		return v.(*slog.Logger)
	}
	return newDiscardText()
}

func (l *LoggerMock) SlogAll() *slog.Logger {
	args := l.Called()
	if v := args.Get(0); v != nil {
		return v.(*slog.Logger)
	}
	// fallback “multi” com os dois handlers descartando
	jh := slog.NewJSONHandler(io.Discard, nil)
	th := slog.NewTextHandler(io.Discard, nil)
	return slog.New(multiFanout(jh, th))
}

// --- helpers internos ---

func getLogger(args mock.Arguments, idx int) *slog.Logger {
	if v := args.Get(idx); v != nil {
		return v.(*slog.Logger)
	}
	// fallback seguro para não quebrar teste sem Return configurado
	return newDiscardText()
}

func newDiscardJSON() *slog.Logger { return slog.New(slog.NewJSONHandler(io.Discard, nil)) }
func newDiscardText() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

// Se não quiser puxar lib externa pro “fanout”, pode só retornar um dos handlers.
// Aqui deixo uma função que você pode implementar com a lib que já usa no runtime.
// Ex.: com "github.com/samber/slog-multi":
//
//	func multiFanout(hs ...slog.Handler) slog.Handler { return slogmulti.Fanout(hs...) }
func multiFanout(h1, h2 slog.Handler) slog.Handler {
	// fallback simples: devolve h1; ajuste se quiser real "fanout" nos testes
	return h1
}

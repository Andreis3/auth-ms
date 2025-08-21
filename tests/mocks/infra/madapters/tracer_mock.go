package madapters

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/auth/domain/interfaces/adapter"
)

type TracerMock struct{ mock.Mock }

func (t *TracerMock) Start(ctx context.Context, spanName string) (context.Context, adapter.Span) {
	args := t.Called(ctx, spanName)
	return args.Get(0).(context.Context), args.Get(1).(adapter.Span)
}

type SpanMock struct{ mock.Mock }

func (s *SpanMock) End() { s.Called() }

func (s *SpanMock) RecordError(err error) { s.Called(err) }

func (s *SpanMock) SpanContext() adapter.SpanContext {
	args := s.Called()
	return args.Get(0).(adapter.SpanContext)
}

type SpanContextMock struct{ mock.Mock }

func (sc *SpanContextMock) TraceID() string {
	args := sc.Called()
	return args.String(0)
}

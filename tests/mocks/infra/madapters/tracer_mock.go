package madapters

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/domain/interfaces/iadapter"
)

type TracerMock struct{ mock.Mock }

func (t *TracerMock) Start(ctx context.Context, spanName string) (context.Context, iadapter.Span) {
	args := t.Called(ctx, spanName)
	return args.Get(0).(context.Context), args.Get(1).(iadapter.Span)
}

type SpanMock struct{ mock.Mock }

func (s *SpanMock) End() { s.Called() }

func (s *SpanMock) RecordError(err error) { s.Called(err) }

func (s *SpanMock) SpanContext() iadapter.SpanContext {
	args := s.Called()
	return args.Get(0).(iadapter.SpanContext)
}

type SpanContextMock struct{ mock.Mock }

func (sc *SpanContextMock) TraceID() string {
	args := sc.Called()
	return args.String(0)
}

package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Code: seu enum/codes (ex.: ErrCacheGet, ErrCacheSet, ErrValidation, ...)
type Code string

type Error struct {
	Code            Code           `json:"code"`
	Messages        []string       `json:"messages,omitempty"` // renomeado de Errors -> Messages
	Fields          map[string]any `json:"fields,omitempty"`   // renomeado de Map -> Fields
	OriginFunc      string         `json:"originFunc,omitempty"`
	FriendlyMessage string         `json:"friendlyMessage,omitempty"`
	Err             error          `json:"-"` // causa real (não serializa)
}

// Garante compatibilidade com a interface error
func (e *Error) Error() string {
	// serializa um "espelho" sem expor Err (para logs legíveis)
	type view struct {
		Code            Code           `json:"code"`
		Messages        []string       `json:"messages,omitempty"`
		Fields          map[string]any `json:"fields,omitempty"`
		OriginFunc      string         `json:"originFunc,omitempty"`
		FriendlyMessage string         `json:"friendlyMessage,omitempty"`
		Cause           string         `json:"cause,omitempty"`
	}
	v := view{
		Code:            e.Code,
		Messages:        e.Messages,
		Fields:          e.Fields,
		OriginFunc:      e.OriginFunc,
		FriendlyMessage: e.FriendlyMessage,
	}
	if e.Err != nil {
		v.Cause = e.Err.Error()
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// Unwrap expõe a causa p/ errors.Is / errors.As
func (e *Error) Unwrap() error { return e.Err }

// --------- Construtores / Helpers ---------

// New cria um erro com code + mensagens
func New(code Code, msgs ...string) *Error {
	return &Error{Code: code, Messages: msgs}
}

// Newf com fmt
func Newf(code Code, format string, a ...any) *Error {
	return &Error{Code: code, Messages: []string{fmt.Sprintf(format, a...)}}
}

// Wrap associa uma causa (err) mantendo code/mensagens
func Wrap(err error, code Code, msgs ...string) *Error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Messages: msgs, Err: err}
}

// WithFields adiciona/mescla campos (contexto p/ observabilidade)
func (e *Error) WithFields(kv map[string]any) *Error {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	for k, v := range kv {
		e.Fields[k] = v
	}
	return e
}

func (e *Error) WithField(k string, v any) *Error {
	if e.Fields == nil {
		e.Fields = make(map[string]any)
	}
	e.Fields[k] = v
	return e
}

func (e *Error) WithFriendly(msg string) *Error {
	e.FriendlyMessage = msg
	return e
}

func (e *Error) WithOrigin(funcName string) *Error {
	e.OriginFunc = funcName
	return e
}

// Captura função chamadora (opcional)
func (e *Error) CaptureOrigin(skip int) *Error {
	pc, _, _, ok := runtime.Caller(skip)
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			parts := strings.Split(fn.Name(), "/")
			e.OriginFunc = parts[len(parts)-1]
		}
	}
	return e
}

// AddMsg adiciona mensagem
func (e *Error) AddMsg(msg string) *Error {
	e.Messages = append(e.Messages, msg)
	return e
}

// --------- Pass-through (opcional manter) ---------
func Unwrap(err error) error        { return errors.Unwrap(err) }
func Join(errs ...error) error      { return errors.Join(errs...) }
func As(err error, target any) bool { return errors.As(err, target) }
func Is(err, target error) bool     { return errors.Is(err, target) }

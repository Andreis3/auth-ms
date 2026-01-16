package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/andreis3/auth-ms/internal/adapter/input/translator"
	"github.com/andreis3/auth-ms/internal/domain/errors"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

type TypeResponseError struct {
	CodeError       string         `json:"code_error"`
	ErrorFields     map[string]any `json:"error_fields,omitempty"`
	FriendlyMessage string         `json:"friendly_message"`
}

type TypeResponseSuccess[T any] struct {
	Data T `json:"data"`
}

func ResponseSuccess[T any](write http.ResponseWriter, status int, data T) {
	write.Header().Set(ContentType, ApplicationJSON)
	write.WriteHeader(status)
	if status == http.StatusNoContent {
		return
	}

	_ = json.NewEncoder(write).Encode(data)
}

func ResponseError(write http.ResponseWriter, err *errors.Error) int {
	status := http.StatusInternalServerError

	if err != nil {
		if tr, ok := translator.ErrorTranslator[err.Code]; ok && tr.HTTPStatus >= 100 {
			status = tr.HTTPStatus
		}

		friendly := err.FriendlyMessage
		if friendly == "" {
			friendly = "An unexpected error occurred. Please try again later."
		}

		write.Header().Set(ContentType, ApplicationJSON+"; charset=utf-8")
		write.WriteHeader(status)

		resp := TypeResponseError{
			CodeError:       string(err.Code),
			ErrorFields:     err.Fields,
			FriendlyMessage: friendly,
		}

		if encodeErr := json.NewEncoder(write).Encode(resp); encodeErr != nil {
			_, _ = write.Write([]byte(`{"code_error":"internal_error","friendly_message":"An unexpected error occurred. Please try again later."}`))
		}

		return status
	}

	write.Header().Set(ContentType, ApplicationJSON+"; charset=utf-8")
	write.WriteHeader(status)
	_ = json.NewEncoder(write).Encode(TypeResponseError{
		CodeError:       "internal_error",
		FriendlyMessage: "An unexpected error occurred. Please try again later.",
	})

	return status
}

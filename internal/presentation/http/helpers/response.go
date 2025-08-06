package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/presentation/translator"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

type TypeResponseError struct {
	CodeError       string         `json:"code_error"`
	ErrorFields     map[string]any `json:"error_fields,omitempty"`
	FriendlyMessage any            `json:"friendly_message"`
}

type TypeResponseSuccess struct {
	Data any `json:"data"`
}

func ResponseSuccess[T any](write http.ResponseWriter, status int, data T) {
	write.Header().Set(ContentType, ApplicationJSON)
	write.WriteHeader(status)
	result := data

	if any(result) == nil {
		_ = json.NewEncoder(write)
		return
	}

	_ = json.NewEncoder(write).Encode(result)
}

func ResponseError(write http.ResponseWriter, err *errors.Error) {
	status := translator.ErrorTranslator[err.Code].HTTPStatus
	write.Header().Set(ContentType, ApplicationJSON)
	write.WriteHeader(status)

	result := TypeResponseError{
		CodeError:       string(err.Code),
		ErrorFields:     err.Map,
		FriendlyMessage: err.FriendlyMessage,
	}

	_ = json.NewEncoder(write).Encode(result)
}

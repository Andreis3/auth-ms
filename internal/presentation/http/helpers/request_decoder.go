package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/andreis3/customers-ms/internal/domain/errors"
)

func RequestDecoder[T any](req *http.Request) (T, *errors.Error) {
	defer req.Body.Close()
	var result T
	var jsonUnmarshalTypeError *json.UnmarshalTypeError
	var jsonSyntaxError *json.SyntaxError
	err := json.NewDecoder(req.Body).Decode(&result)
	switch {
	case errors.As(err, &jsonSyntaxError):
		return result, errors.ErrorJSONSyntaxError(jsonSyntaxError)

	case errors.As(err, &jsonUnmarshalTypeError):
		return result, errors.ErrorJSONUnmarshalTypeError(jsonUnmarshalTypeError)

	case err != nil:
		return result, errors.ErrorJSON(err)
	}

	return result, nil
}

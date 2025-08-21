package helpers

import (
	"encoding/json"
	"net/http"

	errors2 "github.com/andreis3/auth-ms/internal/auth/domain/errors"
)

func RequestDecoder[T any](req *http.Request) (T, *errors2.Error) {
	defer req.Body.Close()
	var result T
	var jsonUnmarshalTypeError *json.UnmarshalTypeError
	var jsonSyntaxError *json.SyntaxError
	err := json.NewDecoder(req.Body).Decode(&result)
	switch {
	case errors2.As(err, &jsonSyntaxError):
		return result, errors2.ErrorJSONSyntaxError(jsonSyntaxError)

	case errors2.As(err, &jsonUnmarshalTypeError):
		return result, errors2.ErrorJSONUnmarshalTypeError(jsonUnmarshalTypeError)

	case err != nil:
		return result, errors2.ErrorJSON(err)
	}

	return result, nil
}

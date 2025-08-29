package translator

import (
	"net/http"

	"google.golang.org/grpc/codes"

	errors2 "github.com/andreis3/auth-ms/internal/domain/errors"
)

type ProtocolError struct {
	HTTPStatus int
	GRPCCode   codes.Code
}

var ErrorTranslator = map[errors2.Code]ProtocolError{
	errors2.ErrBadRequest: {
		HTTPStatus: http.StatusBadRequest,
		GRPCCode:   codes.InvalidArgument,
	},
	errors2.ErrNotFound: {
		HTTPStatus: http.StatusNotFound,
		GRPCCode:   codes.NotFound,
	},
	errors2.ErrInternal: {
		HTTPStatus: http.StatusInternalServerError,
		GRPCCode:   codes.Internal,
	},
	errors2.ErrUnauthorized: {
		HTTPStatus: http.StatusUnauthorized,
		GRPCCode:   codes.Unauthenticated,
	},
	errors2.ErrForbidden: {
		HTTPStatus: http.StatusForbidden,
		GRPCCode:   codes.PermissionDenied,
	},
	errors2.ErrConflict: {
		HTTPStatus: http.StatusConflict,
		GRPCCode:   codes.AlreadyExists,
	},
	errors2.ErrUnprocessableEntity: {
		HTTPStatus: http.StatusUnprocessableEntity,
		GRPCCode:   codes.InvalidArgument,
	},
}

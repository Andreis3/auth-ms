package translator

import (
	"net/http"

	"google.golang.org/grpc/codes"

	"github.com/andreis3/auth-ms/internal/auth/domain/errors"
)

type ProtocolError struct {
	HTTPStatus int
	GRPCCode   codes.Code
}

var ErrorTranslator = map[errors.Code]ProtocolError{
	errors.ErrBadRequest: {
		HTTPStatus: http.StatusBadRequest,
		GRPCCode:   codes.InvalidArgument,
	},
	errors.ErrNotFound: {
		HTTPStatus: http.StatusNotFound,
		GRPCCode:   codes.NotFound,
	},
	errors.ErrInternal: {
		HTTPStatus: http.StatusInternalServerError,
		GRPCCode:   codes.Internal,
	},
	errors.ErrUnauthorized: {
		HTTPStatus: http.StatusUnauthorized,
		GRPCCode:   codes.Unauthenticated,
	},
	errors.ErrForbidden: {
		HTTPStatus: http.StatusForbidden,
		GRPCCode:   codes.PermissionDenied,
	},
	errors.ErrConflict: {
		HTTPStatus: http.StatusConflict,
		GRPCCode:   codes.AlreadyExists,
	},
	errors.ErrUnprocessableEntity: {
		HTTPStatus: http.StatusUnprocessableEntity,
		GRPCCode:   codes.InvalidArgument,
	},
}

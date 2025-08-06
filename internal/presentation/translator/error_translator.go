package translator

import (
	"net/http"

	"google.golang.org/grpc/codes"

	"github.com/andreis3/customers-ms/internal/domain/errors"
)

type ProtocolError struct {
	HTTPStatus int
	GRPCCode   codes.Code
}

var ErrorTranslator = map[errors.Code]ProtocolError{
	errors.BadRequestCode: {
		HTTPStatus: http.StatusBadRequest,
		GRPCCode:   codes.InvalidArgument,
	},
	errors.NotFoundCode: {
		HTTPStatus: http.StatusNotFound,
		GRPCCode:   codes.NotFound,
	},
	errors.InternalServerErrorCode: {
		HTTPStatus: http.StatusInternalServerError,
		GRPCCode:   codes.Internal,
	},
	errors.UnauthorizedCode: {
		HTTPStatus: http.StatusUnauthorized,
		GRPCCode:   codes.Unauthenticated,
	},
	errors.ForbiddenCode: {
		HTTPStatus: http.StatusForbidden,
		GRPCCode:   codes.PermissionDenied,
	},
	errors.ConflictCode: {
		HTTPStatus: http.StatusConflict,
		GRPCCode:   codes.AlreadyExists,
	},
	errors.UnprocessableEntityCode: {
		HTTPStatus: http.StatusUnprocessableEntity,
		GRPCCode:   codes.InvalidArgument,
	},
}

package errors

//type Code string

const (
	ErrBadRequest          Code = "ERR_BAD_REQUEST"
	ErrNotFound            Code = "ERR_NOT_FOUND"
	ErrUnauthorized        Code = "ERR_UNAUTHORIZED"
	ErrForbidden           Code = "ERR_FORBIDDEN"
	ErrConflict            Code = "ERR_CONFLICT"
	ErrUnprocessableEntity Code = "ERR_UNPROCESSABLE"
	ErrInternal            Code = "ERR_INTERNAL"
)

const (
	InternalServerError        = "Internal server error"
	ServerErrorFriendlyMessage = "Internal server error"
	InvalidCredentialsMessage  = "Invalid credentials"
)

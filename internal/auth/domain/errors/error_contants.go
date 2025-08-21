package errors

type Code string

const (
	BadRequestCode          Code = "DM-400"
	NotFoundCode            Code = "DM-404"
	InternalServerErrorCode Code = "IMF-500"
	UnauthorizedCode        Code = "DM-401"
	ForbiddenCode           Code = "DM-403"
	ConflictCode            Code = "DM-409"
	UnprocessableEntityCode Code = "DM-422"
)

const (
	InternalServerError        = "Internal server error"
	ServerErrorFriendlyMessage = "Internal server error"
	InvalidCredentialsMessage  = "Invalid credentials"
)

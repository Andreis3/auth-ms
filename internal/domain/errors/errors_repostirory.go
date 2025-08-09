package errors

/********Repository Errors********/

func CreateUserError(err error) *Error {
	input := InputError{
		Code:            InternalServerErrorCode,
		Errors:          []string{err.Error()},
		OriginFunc:      "UserRepository.CreateUser",
		FriendlyMessage: "Ops... something went wrong. Please try again later.",
	}
	return New(input)
}

func ErrorAlreadyExistsUser(err error) *Error {
	input := InputError{
		Code:            ConflictCode,
		Errors:          []string{err.Error()},
		OriginFunc:      "UserRepository.CreateUser",
		FriendlyMessage: "User with this email already exists.",
	}
	return New(input)
}

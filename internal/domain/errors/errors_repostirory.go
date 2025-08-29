package errors

/********Repository Errors********/

func CreateUserError(err error) *Error {
	return Wrap(err, ErrInternal, "Error creating user").
		WithOrigin("UserRepository.CreateUser").
		WithFriendly("Ops... something went wrong. Please try again later.")
}

func ErrorAlreadyExistsUser(err error) *Error {
	return Wrap(err, ErrConflict, "User already exists").
		WithOrigin("UserRepository.CreateUser").
		WithFriendly("User with this email already exists.")
}

func ErrorFindUserByEmail(err error) *Error {
	return Wrap(err, ErrInternal, "Error finding user by email").
		WithOrigin("UserRepository.FindUserByEmail").
		WithFriendly("Ops... something went wrong. Please try again later.")
}

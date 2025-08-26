package errors

func ErrorAlreadyExists(publicID string) *Error {

	return Newf(ErrConflict, "User with public ID %v already exists", publicID).
		WithOrigin("UserRepository.CreateUser").
		WithFriendly("User with this email already exists.")
}

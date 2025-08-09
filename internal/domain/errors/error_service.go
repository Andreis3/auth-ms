package errors

import "fmt"

func ErrorAlreadyExists(publicID string) *Error {
	input := InputError{
		Code:            ConflictCode,
		Errors:          []string{fmt.Sprintf("User with public ID %v already exists", publicID)},
		OriginFunc:      "UserRepository.CreateUser",
		FriendlyMessage: "User with this email already exists.",
	}
	return New(input)
}

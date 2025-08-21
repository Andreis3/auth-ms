package errors

/********UnitOfWork Errors********/
func ErrorTransactionAlreadyExists() *Error {
	input := InputError{
		Code:            InternalServerErrorCode,
		Errors:          []string{ServerErrorFriendlyMessage},
		OriginFunc:      "UnitOfWork.WithTransaction",
		FriendlyMessage: "Transaction already exists",
	}
	return New(input)
}

func ErrorOpeningTransaction(err error) *Error {
	input := InputError{
		Code:            InternalServerErrorCode,
		Errors:          []string{err.Error()},
		OriginFunc:      "UnitOfWork.WithTransaction",
		FriendlyMessage: ServerErrorFriendlyMessage,
	}
	return New(input)
}

func ErrorExecuteRollback(err error) *Error {
	input := InputError{
		Code:            InternalServerErrorCode,
		Errors:          []string{err.Error()},
		OriginFunc:      "UnitOfWork.Rollback",
		FriendlyMessage: ServerErrorFriendlyMessage,
	}
	return New(input)
}

func ErrorCommitOrRollback(err error) *Error {
	input := InputError{
		Code:            InternalServerErrorCode,
		Errors:          []string{err.Error()},
		OriginFunc:      "UnitOfWork.CommitOrRollback",
		FriendlyMessage: ServerErrorFriendlyMessage,
	}
	return New(input)
}

/********Bcrypt Errors********/
func ErrorHashPassword(err error) *Error {
	input := InputError{
		Code:            InternalServerErrorCode,
		Errors:          []string{err.Error()},
		OriginFunc:      "Bcrypt.Hash",
		FriendlyMessage: ServerErrorFriendlyMessage,
	}
	return New(input)
}

/********Decoder Errors********/
func ErrorJSONSyntaxError(err error) *Error {
	input := InputError{
		Code:            BadRequestCode,
		Errors:          []string{err.Error()},
		Cause:           "json syntax error",
		OriginFunc:      "json.Unmarshal",
		FriendlyMessage: "json syntax error",
	}
	return New(input)
}

func ErrorJSONUnmarshalTypeError(err error) *Error {
	input := InputError{
		Code:            BadRequestCode,
		Errors:          []string{err.Error()},
		Cause:           "json unmarshal type error",
		OriginFunc:      "json.Unmarshal",
		FriendlyMessage: "json unmarshal type error",
	}
	return New(input)
}

func ErrorJSON(err error) *Error {
	input := InputError{
		Code:            BadRequestCode,
		Errors:          []string{err.Error()},
		Cause:           "json error",
		OriginFunc:      "json.Unmarshal",
		FriendlyMessage: "json error",
	}
	return New(input)
}

/*********Redis Errors***************/
func ErrorGetCache(err error) *Error {
	input := InputError{
		Code:            InternalServerErrorCode,
		Errors:          []string{err.Error()},
		OriginFunc:      "Redis.GetCache",
		FriendlyMessage: ServerErrorFriendlyMessage,
	}
	return New(input)
}

func ErrorSetCache(err error) *Error {
	input := InputError{
		Code:            InternalServerErrorCode,
		Errors:          []string{err.Error()},
		OriginFunc:      "Redis.SetCache",
		FriendlyMessage: ServerErrorFriendlyMessage,
	}
	return New(input)
}

package errors

/********UnitOfWork Errors********/
func ErrorTransactionAlreadyExists() *Error {
	return New(ErrInternal, "Transaction already exists").
		WithOrigin("UnitOfWork.WithTransaction").
		WithFriendly("Transaction already exists")
}

func ErrorOpeningTransaction(err error) *Error {

	return Wrap(err, ErrInternal, "Error opening transaction").
		WithOrigin("UnitOfWork.WithTransaction").
		WithFriendly(ServerErrorFriendlyMessage)
}

func ErrorExecuteRollback(err error) *Error {

	return Wrap(err, ErrInternal, "Error executing rollback").
		WithOrigin("UnitOfWork.Rollback").
		WithFriendly(ServerErrorFriendlyMessage)
}

func ErrorCommitOrRollback(err error) *Error {
	return Wrap(err, ErrInternal, "Error committing or rolling back transaction").
		WithOrigin("UnitOfWork.CommitOrRollback").
		WithFriendly(ServerErrorFriendlyMessage)
}

/********Bcrypt Errors********/
func ErrorHashPassword(err error) *Error {
	return Wrap(err, ErrInternal, "Error hashing password").
		WithOrigin("Bcrypt.Hash").
		WithFriendly(ServerErrorFriendlyMessage)
}

/********Decoder Errors********/
func ErrorJSONSyntaxError(err error) *Error {
	return Wrap(err, ErrBadRequest, "JSON syntax error").
		WithOrigin("json.Unmarshal").
		WithFriendly("JSON syntax error")
}

func ErrorJSONUnmarshalTypeError(err error) *Error {
	return Wrap(err, ErrBadRequest, "JSON unmarshal type error").
		WithOrigin("json.Unmarshal").
		WithFriendly("JSON unmarshal type error")
}

func ErrorJSON(err error) *Error {
	return Wrap(err, ErrBadRequest, "JSON error").
		WithOrigin("json.Unmarshal").
		WithFriendly("JSON error")
}

/*********Redis Errors***************/
func ErrorGetCache(err error) *Error {

	return Wrap(err, ErrInternal, "Error getting cache").
		WithOrigin("Redis.GetCache").
		WithFriendly(ServerErrorFriendlyMessage)
}

func ErrorSetCache(err error) *Error {

	return Wrap(err, ErrInternal, "Error setting cache").
		WithOrigin("Redis.SetCache").
		WithFriendly(ServerErrorFriendlyMessage)
}

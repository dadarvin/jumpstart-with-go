package errors

import "errors"

var (
	// server side error
	ErrUnknown           = errors.New("internal server error")
	ErrEncode            = errors.New("encode message error")
	ErrEmptyResponseData = errors.New("empty response data")
	ErrTransaction       = errors.New("error creating transaction")
	// client side error
	ErrBadRequest                = errors.New("bad request")
	ErrMissingMandatoryParameter = errors.New("missing mandatory parameter")
	// access error
	ErrAuthUnauthorized     = errors.New("user is unauthorized to access this resource")
	ErrAuthMissingAuthToken = errors.New("missing auth token")
	ErrAuthInvalidUserID    = errors.New("invalid user id")
	ErrTooManyRequest       = errors.New("too many request")
	/* internal server error detail
	by default will be remapped to ErrUnknown, unwrap with errors.Is if client need to know detail root cause*/
	ErrDatabase  = errors.New("database error")
	ErrHost      = errors.New("host error")
	ErrHostEmail = errors.New("email error")
	// usecase error
	ErrUsecase = errors.New("usecase error")
)

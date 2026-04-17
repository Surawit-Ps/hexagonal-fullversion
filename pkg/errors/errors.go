package errors

import "errors"

var (
	ErrDogNotFound = errors.New("dog not found")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternalServer = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict = errors.New("conflict")
	ErrBadRequest = errors.New("bad request")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound = errors.New("not found")
	ErrInvalidID = errors.New("invalid ID")
)
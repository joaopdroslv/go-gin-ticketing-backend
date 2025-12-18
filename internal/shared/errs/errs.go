package errs

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrNothingToUpdate  = errors.New("nothing to update")
	ErrZeroRowsAffected = errors.New("zero rows affected")
	ErrConflict         = errors.New("resource conflict")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrValidation       = errors.New("validation error")
)

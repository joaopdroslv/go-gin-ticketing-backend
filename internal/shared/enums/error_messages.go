package enums

type ErrorMessage string

const (
	ErrInternal         ErrorMessage = "Something went wrong internally."
	ErrNotFound         ErrorMessage = "Requested resource not found."
	ErrZeroRowsAffected ErrorMessage = "Nothing affected with the action."
	ErrInvalidID        ErrorMessage = "The provided ID is invalid"
	ErrBadRequest       ErrorMessage = "Bad request, fix it and try again later."
)

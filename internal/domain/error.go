package domain

import "errors"

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInternalServer   = errors.New("internal server error")
	ErrDuplicateProduct = errors.New("duplicate product name")
)

package pgx

import "errors"

var (
	ErrNotFound    = errors.New("lead not found")
	ErrEmailExists = errors.New("email already exists")
)

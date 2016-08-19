package db

import "github.com/lib/pq"

const (
	ForeignKeyViolation pq.ErrorCode = "23503"
)

type notFoundError string

func (nfe notFoundError) Error() string { return string(nfe) }
func (notFoundError) NotFound() bool    { return true }

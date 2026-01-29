package inmemory

import (
	"errors"
)
// inmemory db client error
var (
	ErrNotFoundTable = errors.New("not found table")
	ErrCreateTableFailed = errors.New("failed to create table")
	ErrNameTableCollision = errors.New("there are existing table")
)

// inmemory db table error
var (
	ErrNotFoundRecord = errors.New("not found record in the table")
    ErrSaveError = errors.New("saving error")
)
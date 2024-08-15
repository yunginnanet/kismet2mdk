package data

import (
	"errors"
	"strings"

	"github.com/glebarez/go-sqlite"
)

// SQLiteError is an error type that wraps a pointer to an [sqlite.Error].
type SQLiteError struct {
	Msg  string
	code int
	e    error
}

func (sqe *SQLiteError) Error() string { return sqe.e.Error() }
func (sqe *SQLiteError) Code() int     { return sqe.code }
func (sqe *SQLiteError) Unwrap() error { return sqe.e }
func (sqe *SQLiteError) Is(target error) bool {
	if target == nil {
		return false
	}
	return errors.Is(target, sqe.e)
}

//goland:noinspection GoDirectComparisonOfErrors
func (sqe *SQLiteError) IsBusy() bool { return sqe.code == 5 }

func NewSQLiteError(err error) *SQLiteError {
	if err == nil {
		return nil
	}
	nerr := new(SQLiteError)
	nerr.e = err
	nerr.code = -99
	var sqe = &sqlite.Error{}
	if !errors.As(err, &sqe) {
		nerr.e = sqe
		nerr.code = sqe.Code()
	} else if strings.Contains(err.Error(), "(5)") ||
		strings.Contains(err.Error(), "(SQLITE_BUSY)") {
		nerr.code = 5
	}
	return nerr
}

package services

import (
	"fmt"
)

var (
	// ErrSQLExec means there was an error while executing a query.
	// The query that was not executed is wrapped as an SQLError which contains the query string
	// and the query string arguments.
	ErrSQLExec = &ServiceError{nil, "failed to execute query"}
	ErrSQLScan = &ServiceError{nil, "failed to scan results"}
	ErrSQLRead = &ServiceError{nil, "failed to read a row"}
	// ErrNotAdded means that record was not added but the query was executed.
	ErrNotAdded = &ServiceError{nil, "record not added"}
	// ErrNotDeleted means that record was not deleted but the query was executed.
	ErrNotDeleted = &ServiceError{nil, "record not deleted"}
	ErrNotEdited  = &ServiceError{nil, "record not edited"}
	ErrNotFound   = &ServiceError{nil, "record not found"}
	// ErrUnknown means that the source or the reason for the error is unknown. The error that
	// was thrown is wrapped.
	ErrUnknown = &ServiceError{nil, "unknown error"}
)

// ServiceError stores information about an error that occured inside a service. The original
// error is wrapped but the message that describes the situation in which error has occures is
// in Msg field.
type ServiceError struct {
	Err error  // Wrapped error
	Msg string // Service error message
}

func (e ServiceError) Error() string {
	return e.Msg
}

func (e ServiceError) Unwrap() error {
	return e.Err
}

func NewServiceError(e error, m string) error {
	return &ServiceError{e, m}
}

// SQLError stores information about an error occured while working with sql. The query (if
// avaliable) is stored in Qs and Qsa fields (query string and attributes). The error from
// database/sql is wrapped in Err.
type SQLError struct {
	Err error
	Qs  string
	Qsa []interface{}
}

func (e SQLError) Error() string {
	return fmt.Sprintf("sql: %s - %v", e.Qs, e.Qsa)
}

func (e SQLError) Unwrap() error {
	return e.Err
}

// QueryBuilder interface that implements ToSql function required for getting query string and
// arguments if error has occured.
type QueryBuilder interface {
	ToSql() (string, []interface{}, error)
}

// NewSQLError creates a new sql error from query builder and the error that has occured.
func NewSQLError(e error, q QueryBuilder) error {
	qs, qsa, err := q.ToSql()
	if err != nil {
		e = err
	}
	return &SQLError{e, qs, qsa}
}

func NewWrappedSQLError(we error, q QueryBuilder, e error) error {
	qs, qsa, err := q.ToSql()
	if err != nil {
		we = err
	}

	return &ServiceError{&SQLError{we, qs, qsa}, e.Error()}
}

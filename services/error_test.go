package services_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
	"github.com/stretchr/testify/require"
)

func TestServiceError_New(t *testing.T) {
	we := errors.New("very bad error")
	e := services.NewServiceError(we, "something very bad happened")

	require.Error(t, e)
	require.EqualError(t, e, "something very bad happened")

	require.Error(t, errors.Unwrap(e))
	require.EqualError(t, errors.Unwrap(e), "very bad error")
}

func TestSQLError_New(t *testing.T) {
	q := squirrel.Insert("table").Columns("col").Values(1)
	we := errors.New("very bad error")
	e := services.NewSQLError(we, q)

	qs, qsa, err := q.ToSql()

	require.NoError(t, err)

	require.Error(t, e)
	require.EqualError(t, e, fmt.Sprintf("sql: %s - %v", qs, qsa))

	require.Error(t, errors.Unwrap(e))
	require.EqualError(t, errors.Unwrap(e), "very bad error")
}

func TestWrappedSQLError_New(t *testing.T) {
	q := squirrel.Insert("table").Columns("col").Values(1)
	wrappedErr := errors.New("very bad error")
	wrapperErr := errors.New("abstracted bad error")
	e := services.NewWrappedSQLError(wrappedErr, q, wrapperErr)

	qs, qsa, err := q.ToSql()

	require.NoError(t, err)

	require.Error(t, e)
	require.EqualError(t, e, "abstracted bad error")

	require.Error(t, errors.Unwrap(e))
	require.EqualError(t, errors.Unwrap(e), fmt.Sprintf("sql: %s - %v", qs, qsa))

	require.Error(t, errors.Unwrap(errors.Unwrap(e)))
	require.EqualError(t, errors.Unwrap(errors.Unwrap(e)), "very bad error")
}

type mockBuilder struct{}

func (b mockBuilder) ToSql() (string, []interface{}, error) {
	return "", nil, errors.New("builder error")
}

func TestSQLError_New_fail_build(t *testing.T) {
	q := mockBuilder{}
	e := services.NewSQLError(nil, q)

	qs, qsa, err := q.ToSql()

	require.Error(t, err)
	require.EqualError(t, err, "builder error")

	require.Error(t, e)
	require.EqualError(t, e, fmt.Sprintf("sql: %s - %v", qs, qsa))

	require.Error(t, errors.Unwrap(e))
	require.EqualError(t, errors.Unwrap(e), "builder error")
}

func TestWrappedSQLError_New_fail_build(t *testing.T) {
	q := mockBuilder{}
	we := errors.New("very bad error")
	e := services.NewWrappedSQLError(nil, q, we)

	qs, qsa, err := q.ToSql()

	require.Error(t, err)
	require.EqualError(t, err, "builder error")

	require.Error(t, e)
	require.EqualError(t, e, "very bad error")

	require.Error(t, errors.Unwrap(e))
	require.EqualError(t, errors.Unwrap(e), fmt.Sprintf("sql: %s - %v", qs, qsa))

	require.Error(t, errors.Unwrap(errors.Unwrap(e)))
	require.EqualError(t, errors.Unwrap(errors.Unwrap(e)), "builder error")
}

package domain_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/stretchr/testify/require"
)

var (
	sampleDomain = yaas.Domain{
		UserID: 1,
		Code:   "Uakgb_J5m9g~0JDMbcJqLJ",
		Domain: "stackoverflow.com",
		Desc:   "Simple description",
	}
)

func TestDomainServiceAdd_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO domain").WithArgs(sampleDomain.UserID, sampleDomain.Code, sampleDomain.Domain, sampleDomain.Desc).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	require.NoError(t, s.Add(sampleDomain))

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceAdd_fail_notadded(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO domain").WithArgs(sampleDomain.UserID, sampleDomain.Code, sampleDomain.Domain, sampleDomain.Desc).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	err = s.Add(sampleDomain)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotAdded.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceAdd_fail_exec(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO domain").WithArgs(sampleDomain.UserID, sampleDomain.Code, sampleDomain.Domain, sampleDomain.Desc).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	err = s.Add(sampleDomain)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceAdd_fail_rowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO domain").WithArgs(sampleDomain.UserID, sampleDomain.Code, sampleDomain.Domain, sampleDomain.Desc).WillReturnResult(sqlmock.NewErrorResult(sql.ErrTxDone))
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	err = s.Add(sampleDomain)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrUnknown.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

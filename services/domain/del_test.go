package domain_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/stretchr/testify/require"
)

func TestDomainServiceDelete_ok_domainName(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM domain").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	require.NoError(t, s.Delete(1))

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainserviceDelete_ok_domainID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM domain").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	d := sampleDomain
	d.ID = 1

	require.NoError(t, s.Delete(1))

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceDelete_fail_exec(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM domain").WillReturnError(sql.ErrTxDone)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	err = s.Delete(1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceDelete_fail_rowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM domain").WillReturnResult(sqlmock.NewErrorResult(sql.ErrTxDone))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	err = s.Delete(1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrUnknown.Error())

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceDelete_fail_notDeleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM domain").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	err = s.Delete(1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotDeleted.Error())

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

package domain_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/stretchr/testify/require"
)

func TestDomainServiceGet_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	d := sampleDomain
	d.ID = 1

	mock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"idd", "idu", "code", "name", "description"}).AddRow(d.ID, d.UserID, d.Code, d.Domain, d.Desc)
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WithArgs(d.ID).WillReturnRows(expectedRows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	domain, err := s.Get(d.ID)
	require.NoError(t, err)
	require.Equal(t, domain, d)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceGet_fail_notFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.Get(1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotFound.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceGet_fail_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.Get(1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

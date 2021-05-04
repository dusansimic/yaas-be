package domain_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/stretchr/testify/require"
)

func TestDomainServiceOwnedBy_ok_true(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	d := sampleDomain
	d.ID = 1

	mock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"idd"}).AddRow(d.ID)
	mock.ExpectQuery("SELECT idd FROM domain").WithArgs(d.ID, d.UserID).WillReturnRows(expectedRows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	owned, err := s.OwnedBy(d.ID, d.UserID)
	require.NoError(t, err)
	require.Equal(t, true, owned)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceOwnedBy_ok_false(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	d := sampleDomain
	d.ID = 1

	mock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"idd"})
	mock.ExpectQuery("SELECT idd FROM domain").WithArgs(d.ID, d.UserID).WillReturnRows(expectedRows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	owned, err := s.OwnedBy(d.ID, d.UserID)
	require.NoError(t, err)
	require.Equal(t, false, owned)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceOwnedBy_fail_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idd FROM domain").WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.OwnedBy(1, 1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

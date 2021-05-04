package domain_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/stretchr/testify/require"
)

func TestDomainServiceCodeToID_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	const (
		domainID   = 1
		domainCode = "f30981jfh02349fj"
	)

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"idd"}).AddRow(domainID)
	mock.ExpectQuery("SELECT idd FROM domain").WithArgs(domainCode).WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	idd, err := s.CodeToID(domainCode)
	require.NoError(t, err)
	require.Equal(t, idd, domainID)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceCodeToID_notFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idd FROM domain").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.CodeToID("f30981jfh02349fj")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotFound.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceCodeToID_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idd FROM domain").WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.CodeToID("f30981jfh02349fj")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

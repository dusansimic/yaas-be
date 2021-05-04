package domain_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/stretchr/testify/require"
)

func TestDomainServiceFind_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	d := sampleDomain
	d.ID = 1

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"idd", "idu", "code", "name", "description"}).AddRow(d.ID, d.UserID, d.Code, d.Domain, d.Desc)
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WithArgs(d.Domain).WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	domain, err := s.Find(d.Domain)
	require.NoError(t, err)
	require.Equal(t, domain, d)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceFind_notFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.Find("example.com")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotFound.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceFind_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.Find("example.com")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())
	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

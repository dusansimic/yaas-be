package referrer_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/referrer"
	"github.com/stretchr/testify/require"
)

// Test weather service will return device id correctly
func TestReferrerService_ID_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"idrf"}).AddRow(1)
	mock.ExpectQuery("SELECT idrf FROM referrer").WithArgs("example.com").WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := referrer.NewService(tx)

	id, err := s.ID("example.com")
	require.NoError(t, err)
	require.Equal(t, id, 1)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

// Test weather service will return not found error
func TestReferrerService_ID_notFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idrf FROM referrer").WithArgs("example.com").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := referrer.NewService(tx)

	_, err = s.ID("example.com")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotFound.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

// Test weather service will return scan error
func TestReferrerService_ID_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idrf FROM referrer").WithArgs("example.com").WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := referrer.NewService(tx)

	_, err = s.ID("example.com")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

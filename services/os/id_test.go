package os_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/os"
	"github.com/stretchr/testify/require"
)

// Test weather service will return device id correctly
func TestOperatingSystemService_ID_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"idos"}).AddRow(1)
	mock.ExpectQuery("SELECT idos FROM operatingsystem").WithArgs(yaas.Linux).WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := os.NewService(tx)

	id, err := s.ID(yaas.Linux)
	require.NoError(t, err)
	require.Equal(t, id, 1)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

// Test weather service will return not found error
func TestOperatingSystemService_ID_notFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idos FROM operatingsystem").WithArgs(yaas.Linux).WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := os.NewService(tx)

	_, err = s.ID(yaas.Linux)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotFound.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

// Test weather service will return scan error
func TestOperatingSystemService_ID_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idos FROM operatingsystem").WithArgs(yaas.Linux).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := os.NewService(tx)

	_, err = s.ID(yaas.Linux)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

package referrer_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/referrer"
	"github.com/stretchr/testify/require"
)

func TestReferrerServiceAddReturnID_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"idrf"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO referrer").WithArgs("Example", "example.com").WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := referrer.NewService(tx)

	id, err := s.AddReturnID("Example", "example.com")

	require.NoError(t, err)
	require.Equal(t, 1, id)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestReferrerServiceAddReturnID_fail_notadded(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO referrer").WithArgs("Example", "example.com").WillReturnRows(sqlmock.NewRows([]string{"idrf"}))
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := referrer.NewService(tx)

	_, err = s.AddReturnID("Example", "example.com")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotAdded.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestReferrerServiceAddReturnID_fail_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO referrer").WithArgs("Example", "example.com").WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := referrer.NewService(tx)

	_, err = s.AddReturnID("Example", "example.com")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

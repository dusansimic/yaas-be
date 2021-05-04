package stats_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/stats"
	"github.com/stretchr/testify/require"
)

func TestDevice_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"d.name", "COUNT(r.idr)"}).AddRow("Phone", 21).AddRow("Laptop", 4).AddRow("Desktop", 53)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT d.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	sum, err := s.Device()
	require.NoError(t, err)
	require.Equal(t, []stats.Record{
		{Name: "Phone", Reqs: 21},
		{Name: "Laptop", Reqs: 4},
		{Name: "Desktop", Reqs: 53},
	}, sum)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDevices_fail_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT d.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Device()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDevices_fail_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"d.name", "COUNT(r.idr)"}).AddRow("Desktop", true)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT d.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Device()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDevices_fail_read(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	expectedRows := sqlmock.NewRows([]string{"d.name", "COUNT(r.idr)"}).AddRow("Desktop", 5).RowError(0, sql.ErrTxDone)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT d.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(expectedRows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Device()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLRead.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

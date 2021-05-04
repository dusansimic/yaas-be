package stats_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/stats"
	"github.com/stretchr/testify/require"
)

var ts = time.Now()

// Test weather browser summary is returned correctly
func TestBrowser_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"b.name", "COUNT(r.idr)"}).AddRow("Chrome", 24).AddRow("Firefox", 21)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT b.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	sum, err := s.Browser()
	require.NoError(t, err)
	require.Equal(t, []stats.Record{
		{Name: "Chrome", Reqs: 24},
		{Name: "Firefox", Reqs: 21},
	}, sum)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

// Test weather browser summary is returned correctly with duration
func TestBrowser_ok_duration(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"b.name", "COUNT(r.idr)"}).AddRow("Chrome", 24).AddRow("Firefox", 21)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT b.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts.Add(-time.Hour)).WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts).WithDuration(-time.Hour)

	sum, err := s.Browser()
	require.NoError(t, err)
	require.Equal(t, []stats.Record{
		{Name: "Chrome", Reqs: 24},
		{Name: "Firefox", Reqs: 21},
	}, sum)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

// Test weather browser summary returns a query error
func TestBrowser_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT b.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Browser()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

// Test weather browser summary returns a scan error
func TestBrowser_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"b.name", "COUNT(r.idr)"}).AddRow("Chrome", true)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT b.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Browser()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

// Test weather browser summary returns a read error
func TestBrowser_read(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"b.name", "COUNT(r.idr)"}).AddRow("Chrome", 5).RowError(0, sql.ErrTxDone)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT b.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Browser()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLRead.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

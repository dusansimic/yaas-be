package stats_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/stats"
	"github.com/stretchr/testify/require"
)

func TestPath_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"path", "COUNT(idr)"}).AddRow("/", 21).AddRow("/about", 4).AddRow("/contact", 53)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT path, COUNT\\(idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	sum, err := s.Path()
	require.NoError(t, err)
	require.Equal(t, []stats.Record{
		{Name: "/", Reqs: 21},
		{Name: "/about", Reqs: 4},
		{Name: "/contact", Reqs: 53},
	}, sum)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPath_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT path, COUNT\\(idr\\) FROM record").WithArgs(1, ts).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Path()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPath_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"path", "COUNT(idr)"}).AddRow("/", true)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT path, COUNT\\(idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Path()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPath_read(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"path", "COUNT(idr)"}).AddRow("/", 5).RowError(0, sql.ErrTxDone)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT path, COUNT\\(idr\\) FROM record").WillReturnRows(rows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.Path()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLRead.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

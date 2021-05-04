package stats_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/stats"
	"github.com/stretchr/testify/require"
)

func TestOS_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"o.name", "COUNT(r.idr)"}).AddRow("Android", 21).AddRow("Mac OS X", 4).AddRow("Windows", 53)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT o.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	sum, err := s.OS()
	require.NoError(t, err)
	require.Equal(t, []stats.Record{
		{Name: "Android", Reqs: 21},
		{Name: "Mac OS X", Reqs: 4},
		{Name: "Windows", Reqs: 53},
	}, sum)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOS_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT o.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.OS()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOS_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"o.name", "COUNT(r.idr)"}).AddRow("Android", true)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT o.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.OS()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOS_read(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"o.name", "COUNT(r.idr)"}).AddRow("Android", 5).RowError(0, sql.ErrTxDone)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT o.name, COUNT\\(r.idr\\) FROM record").WithArgs(1, ts).WillReturnRows(rows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := stats.NewService(tx, 1).WithTime(ts)

	_, err = s.OS()
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLRead.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

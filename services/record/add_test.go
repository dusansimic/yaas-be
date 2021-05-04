package record_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/record"
	"github.com/stretchr/testify/require"
)

var (
	recordValue = record.Record{
		DomainID:   1,
		ReferrerID: 1,
		OSID:       1,
		DeviceID:   1,
		CountryID:  1,
		BrowserID:  1,
		Timestamp:  time.Now(),
		Url:        "https://example.com/about",
		Path:       "/about",
	}
)

func TestAddEvent_ok(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := recordValue

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO record").WithArgs(r.Timestamp, r.DomainID, r.Url, r.Path, r.ReferrerID, r.DeviceID, r.CountryID, r.BrowserID, r.OSID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := record.NewService(tx)

	require.NoError(t, s.Add(r))

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddEvent_fail_notAdded(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	r := recordValue

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO record").WithArgs(r.Timestamp, r.DomainID, r.Url, r.Path, r.ReferrerID, r.DeviceID, r.CountryID, r.BrowserID, r.OSID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := record.NewService(tx)

	err = s.Add(r)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotAdded.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddEvent_fail_exec(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO record").WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := record.NewService(tx)

	err = s.Add(recordValue)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddEvent_fail_rowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO record").WillReturnResult(sqlmock.NewErrorResult(sql.ErrTxDone))
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := record.NewService(tx)

	err = s.Add(recordValue)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrUnknown.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

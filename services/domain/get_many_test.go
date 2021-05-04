package domain_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/stretchr/testify/require"
)

var (
	domainsRow = []string{"idd", "idu", "code", "name", "description"}
)

func TestDomainServiceGetMany_ok_found(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	d1 := sampleDomain
	d1.ID = 1
	d1.Code = "61c7d1d0-2482-410c-9b2f-72b2e7e4cf6b"
	d2 := sampleDomain
	d2.ID = 2
	d2.Code = "3e79cb0d-7816-4fe2-b5e7-aa61e3b6c237"
	d2.Domain = "gitlab.com"

	mock.ExpectBegin()
	expectedRows := sqlmock.NewRows(domainsRow).AddRow(d1.ID, d1.UserID, d1.Code, d1.Domain, d1.Desc).AddRow(d2.ID, d2.UserID, d2.Code, d2.Domain, d2.Desc)
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WithArgs(1).WillReturnRows(expectedRows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	domains, err := s.GetMany(1)
	require.NoError(t, err)
	require.Equal(t, []yaas.Domain{d1, d2}, domains)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceGetMany_ok_notFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	expectedRows := sqlmock.NewRows(domainsRow)
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WithArgs(1).WillReturnRows(expectedRows)
	mock.ExpectCommit()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	domains, err := s.GetMany(1)
	require.NoError(t, err)
	require.Equal(t, []yaas.Domain{}, domains)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceGetMany_fail_query(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WithArgs(1).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.GetMany(1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceGetMany_fail_scan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	expectedRows := sqlmock.NewRows(domainsRow).AddRow(true, sampleDomain.UserID, sampleDomain.Code, sampleDomain.Domain, sampleDomain.Desc)
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WithArgs(1).WillReturnRows(expectedRows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.GetMany(1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDomainServiceGetMany_fail_read(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	d := sampleDomain
	d.ID = 1

	mock.ExpectBegin()
	expectedRows := sqlmock.NewRows(domainsRow).AddRow(d.ID, d.UserID, d.Code, d.Domain, d.Desc).RowError(0, sql.ErrTxDone)
	mock.ExpectQuery("SELECT idd, idu, code, name, description FROM domain").WithArgs(1).WillReturnRows(expectedRows)
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	s := domain.NewService(tx)

	_, err = s.GetMany(1)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLRead.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

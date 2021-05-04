package user_test

import (
	"crypto/rand"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/user"
	"github.com/stretchr/testify/require"
)

func TestUserServiceAdd_ok(t *testing.T) {
	db, mock := stubDB(t)
	defer db.Close()

	u := genUser(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO people").WithArgs(u.Username, u.PasswordHash, u.Salt).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tx := txBegin(t, db)

	s := user.NewService(tx)

	require.NoError(t, s.Add(u))

	require.NoError(t, tx.Commit())

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserServiceAdd_fail_notAdded(t *testing.T) {
	db, mock := stubDB(t)
	defer db.Close()

	u := genUser(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO people").WithArgs(u.Username, u.PasswordHash, u.Salt).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	tx := txBegin(t, db)

	s := user.NewService(tx)

	err := s.Add(u)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotAdded.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserServiceAdd_fail_exec(t *testing.T) {
	db, mock := stubDB(t)
	defer db.Close()

	u := genUser(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO people").WithArgs(u.Username, u.PasswordHash, u.Salt).WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx := txBegin(t, db)

	s := user.NewService(tx)

	err := s.Add(u)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLExec.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserServiceAdd_fail_rowsAffected(t *testing.T) {
	db, mock := stubDB(t)
	defer db.Close()

	u := genUser(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO people").WithArgs(u.Username, u.PasswordHash, u.Salt).WillReturnResult(sqlmock.NewErrorResult(sql.ErrTxDone))
	mock.ExpectRollback()

	tx := txBegin(t, db)

	s := user.NewService(tx)

	err := s.Add(u)
	require.Error(t, err)
	require.EqualError(t, err, services.ErrUnknown.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserServiceGet_ok(t *testing.T) {
	db, mock := stubDB(t)
	defer db.Close()

	u := genUser(t)

	mock.ExpectBegin()
	expectedRow := sqlmock.NewRows([]string{"idu", "username", "password_hash", "password_salt"}).AddRow(u.ID, u.Username, u.PasswordHash, u.Salt)
	mock.ExpectQuery("SELECT idu, username, password_hash, password_salt FROM people").WithArgs(u.Username).WillReturnRows(expectedRow)
	mock.ExpectCommit()

	tx := txBegin(t, db)

	s := user.NewService(tx)

	uu, err := s.Get(u.Username)
	require.NoError(t, err)
	require.Equal(t, u, uu)

	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserServiceGet_fail_notFound(t *testing.T) {
	db, mock := stubDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idu, username, password_hash, password_salt FROM people").WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	tx := txBegin(t, db)

	s := user.NewService(tx)

	_, err := s.Get("lukeskywalker")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrNotFound.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserServiceGet_fail_scan(t *testing.T) {
	db, mock := stubDB(t)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT idu, username, password_hash, password_salt FROM people").WillReturnError(sql.ErrTxDone)
	mock.ExpectRollback()

	tx := txBegin(t, db)

	s := user.NewService(tx)

	_, err := s.Get("lukeskywalker")
	require.Error(t, err)
	require.EqualError(t, err, services.ErrSQLScan.Error())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func stubDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

func genUser(t *testing.T) yaas.User {
	hash := make([]byte, 32)
	_, err := rand.Read(hash)
	require.NoError(t, err)

	salt := make([]byte, 8)
	_, err = rand.Read(salt)
	require.NoError(t, err)

	u := yaas.User{
		ID:           1,
		PasswordHash: hash,
		Salt:         salt,
		Username:     "lukeskywalker",
	}

	return u
}

func txBegin(t *testing.T, db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	require.NoError(t, err)
	return tx
}

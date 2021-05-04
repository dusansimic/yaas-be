package user

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
)

func (s userService) Add(u yaas.User) error {
	q := psql.Insert("people").Columns("username", "password_hash", "password_salt").Values(u.Username, u.PasswordHash, u.Salt)

	res, err := q.RunWith(s.tx).Exec()
	if err != nil {
		return services.NewWrappedSQLError(err, q, services.ErrSQLExec)
	}

	added, err := res.RowsAffected()
	if err != nil {
		return services.NewWrappedSQLError(err, q, services.ErrUnknown)
	}

	if added == 0 {
		return services.NewWrappedSQLError(err, q, services.ErrNotAdded)
	}

	return nil
}

func (s userService) Get(un string) (yaas.User, error) {
	q := psql.Select("idu", "username", "password_hash", "password_salt").From("people").Where(sq.Eq{
		"username": un,
	})

	row := q.RunWith(s.tx).QueryRow()
	var u yaas.User
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Salt); err != nil {
		if err == sql.ErrNoRows {
			return yaas.User{}, services.NewWrappedSQLError(err, q, services.ErrNotFound)
		}

		return yaas.User{}, services.NewWrappedSQLError(err, q, services.ErrSQLScan)
	}

	return u, nil
}

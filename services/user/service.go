package user

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
)

type Service interface {
	Add(u yaas.User) error
	Get(u string) (yaas.User, error)
}

type userService struct {
	tx *sql.Tx
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewService(tx *sql.Tx) Service {
	return &userService{tx}
}

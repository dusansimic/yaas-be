package referrer

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Service interface {
	ID(r string) (int, error)
	Add(n string, r string) error
	AddReturnID(n, r string) (int, error)
	Edit(id int, n string) error
}

type referrerService struct {
	tx *sql.Tx
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewService(tx *sql.Tx) Service {
	return &referrerService{tx}
}

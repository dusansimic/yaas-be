package browser

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
)

type Service interface {
	ID(b yaas.Browser) (int, error)
}

type browserService struct {
	tx *sql.Tx
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewService(tx *sql.Tx) Service {
	return &browserService{
		tx: tx,
	}
}

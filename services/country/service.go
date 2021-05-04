package country

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Service interface {
	IDWithISO2(iso string) (int, error)
}

type countryService struct {
	tx *sql.Tx
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewService(tx *sql.Tx) Service {
	return &countryService{
		tx: tx,
	}
}

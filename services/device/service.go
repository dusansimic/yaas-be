package device

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
)

type Service interface {
	ID(d yaas.Device) (int, error)
}

type deviceService struct {
	tx *sql.Tx
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewService(tx *sql.Tx) Service {
	return &deviceService{
		tx: tx,
	}
}

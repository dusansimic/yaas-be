package domain

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
)

// Service is a type that can generate a new domain service. It can be used to manipulate
// domains data in the database.
type Service interface {
	Add(d yaas.Domain) error
	Delete(id int) error
	EditDesc(id int, desc string) error
	Get(id int) (yaas.Domain, error)
	GetMany(idu int) ([]yaas.Domain, error)
	Find(n string) (yaas.Domain, error)
	ID(d string) (int, error)
	CodeToID(c string) (int, error)
	OwnedBy(idd, idu int) (bool, error)
}

type domainService struct {
	tx *sql.Tx
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// NewService creates a new domain service.
func NewService(tx *sql.Tx) Service {
	return &domainService{
		tx: tx,
	}
}

package domain

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
)

func (s *domainService) Find(n string) (yaas.Domain, error) {
	q := psql.Select("idd", "idu", "code", "name", "description").From("domain").Where(sq.Eq{
		"name": n,
	})

	row := q.RunWith(s.tx).QueryRow()

	var d yaas.Domain
	if err := row.Scan(&d.ID, &d.UserID, &d.Code, &d.Domain, &d.Desc); err != nil {
		if err == sql.ErrNoRows {
			return d, services.ErrNotFound
		}

		return d, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLScan.Msg,
		)
	}

	return d, nil
}

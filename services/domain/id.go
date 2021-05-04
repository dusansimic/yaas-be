package domain

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s *domainService) ID(d string) (int, error) {
	q := psql.Select("idd").From("domain").Where(sq.Eq{
		"name": d,
	})

	row := q.RunWith(s.tx).QueryRow()

	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, services.ErrNotFound
		}

		return 0, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLScan.Msg,
		)
	}

	return id, nil
}

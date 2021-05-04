package os

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
)

func (s osService) ID(o yaas.OS) (int, error) {
	q := psql.Select("idos").From("operatingsystem").Where(sq.Eq{
		"name": o,
	})

	row := q.RunWith(s.tx).QueryRow()
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, services.NewServiceError(
				services.NewSQLError(err, q),
				services.ErrNotFound.Msg,
			)
		}

		return 0, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLScan.Msg,
		)
	}

	return id, nil
}

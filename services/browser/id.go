package browser

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
)

func (s browserService) ID(b yaas.Browser) (int, error) {
	q := psql.Select("idb").From("browser").Where(sq.Eq{
		"name": b,
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

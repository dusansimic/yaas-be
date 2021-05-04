package referrer

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s referrerService) ID(r string) (int, error) {
	q := psql.Select("idrf").From("referrer").Where(sq.Eq{
		"domain": r,
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

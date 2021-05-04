package country

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s countryService) IDWithISO2(iso string) (int, error) {
	q := psql.Select("idc").From("country").Where(sq.Eq{
		"iso": iso,
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

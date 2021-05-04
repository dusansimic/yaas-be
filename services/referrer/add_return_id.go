package referrer

import (
	"database/sql"

	"github.com/dusansimic/yaas/services"
)

func (s referrerService) AddReturnID(n, d string) (int, error) {
	q := psql.Insert("referrer").Columns("name", "domain").Values(n, d).Suffix("RETURNING idrf")

	res := q.RunWith(s.tx).QueryRow()

	var id int
	if err := res.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, services.NewServiceError(
				services.NewSQLError(err, q),
				services.ErrNotAdded.Msg,
			)
		}

		return 0, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLScan.Msg,
		)
	}

	return id, nil
}

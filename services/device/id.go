package device

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
)

func (s deviceService) ID(d yaas.Device) (int, error) {
	q := psql.Select("iddv").From("device").Where(sq.Eq{
		"name": d,
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

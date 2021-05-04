package domain

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s *domainService) OwnedBy(idd, idu int) (bool, error) {
	q := psql.Select("idd").From("domain").Where(sq.Eq{
		"idd": idd,
		"idu": idu,
	})

	row := q.RunWith(s.tx).QueryRow()

	var i int
	if err := row.Scan(&i); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLScan.Msg,
		)
	}

	return true, nil
}

package domain

import (
	sq "github.com/Masterminds/squirrel"
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
)

func (s *domainService) GetMany(idu int) ([]yaas.Domain, error) {
	q := psql.Select("idd", "idu", "code", "name", "description").From("domain").Where(sq.Eq{
		"idu": idu,
	})

	rows, err := q.RunWith(s.tx).Query()
	if err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLExec.Msg,
		)
	}

	ds := []yaas.Domain{}
	for rows.Next() {
		var d yaas.Domain
		if err := rows.Scan(&d.ID, &d.UserID, &d.Code, &d.Domain, &d.Desc); err != nil {
			return nil, services.NewServiceError(
				services.NewSQLError(err, q),
				services.ErrSQLScan.Msg,
			)
		}
		ds = append(ds, d)
	}

	if err := rows.Err(); err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLRead.Msg,
		)
	}

	return ds, nil
}

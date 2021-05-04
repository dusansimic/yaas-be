package stats

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s statsService) Path() ([]Record, error) {
	q := psql.Select("path", "COUNT(idr)").From("record").Where(sq.And{
		sq.Eq{"idd": s.id},
		sq.Gt{"timestamp": s.ts},
	}).GroupBy("path")

	rows, err := q.RunWith(s.tx).Query()
	if err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLExec.Msg,
		)
	}

	ps := []Record{}
	for rows.Next() {
		var p Record
		if err := rows.Scan(&p.Name, &p.Reqs); err != nil {
			return nil, services.NewServiceError(
				services.NewSQLError(err, q),
				services.ErrSQLScan.Msg,
			)
		}

		ps = append(ps, p)
	}

	if err := rows.Err(); err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLRead.Msg,
		)
	}

	return ps, nil
}

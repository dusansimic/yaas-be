package stats

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s statsService) OS() ([]Record, error) {
	q := psql.Select("o.name", "COUNT(r.idr)").From("record r").InnerJoin("operatingsystem o ON o.idos = r.idos").Where(sq.And{
		sq.Eq{"r.idd": s.id},
		sq.Gt{"r.timestamp": s.ts},
	}).GroupBy("o.idos")

	rows, err := q.RunWith(s.tx).Query()
	if err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLExec.Msg,
		)
	}

	os := []Record{}
	for rows.Next() {
		var o Record
		if err := rows.Scan(&o.Name, &o.Reqs); err != nil {
			return nil, services.NewServiceError(
				services.NewSQLError(err, q),
				services.ErrSQLScan.Msg,
			)
		}

		os = append(os, o)
	}

	if err := rows.Err(); err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLRead.Msg,
		)
	}

	return os, nil
}

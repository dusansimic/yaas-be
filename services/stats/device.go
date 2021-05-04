package stats

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s statsService) Device() ([]Record, error) {
	q := psql.Select("d.name", "COUNT(r.idr)").From("record r").InnerJoin("device d ON d.iddv = r.iddv").Where(sq.And{
		sq.Eq{"r.idd": s.id},
		sq.Gt{"r.timestamp": s.ts},
	}).GroupBy("d.iddv")

	rows, err := q.RunWith(s.tx).Query()
	if err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLExec.Msg,
		)
	}

	ds := []Record{}
	for rows.Next() {
		var d Record
		if err := rows.Scan(&d.Name, &d.Reqs); err != nil {
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

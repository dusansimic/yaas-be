package stats

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s statsService) Total() ([]Record, error) {
	q := psql.Select("DATE_PART('day', timestamp)", "COUNT(idr)").From("record").Where(sq.And{
		sq.Eq{"idd": s.id},
		sq.Gt{"timestamp": s.ts},
	}).GroupBy("DATE_PART('day', timestamp)")

	rows, err := q.RunWith(s.tx).Query()
	if err != nil {
		return nil, services.NewWrappedSQLError(err, q, services.ErrSQLExec)
	}

	ts := []Record{}
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.Name, &r.Reqs); err != nil {
			return nil, services.NewWrappedSQLError(err, q, services.ErrSQLScan)
		}

		ts = append(ts, r)
	}

	if err := rows.Err(); err != nil {
		return nil, services.NewWrappedSQLError(err, q, services.ErrSQLRead)
	}

	return ts, nil
}

package stats

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s statsService) Referrer() ([]Record, error) {
	q := psql.Select("rf.name", "COUNT(rc.idr)").From("record rc").InnerJoin("referrer rf ON rf.idrf = rc.idrf").Where(sq.And{
		sq.Eq{"rc.idd": s.id},
		sq.Gt{"rc.timestamp": s.ts},
	}).GroupBy("rf.idrf")

	rows, err := q.RunWith(s.tx).Query()
	if err != nil {
		return nil, services.NewWrappedSQLError(err, q, services.ErrSQLExec)
	}

	rs := []Record{}
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.Name, &r.Reqs); err != nil {
			return nil, services.NewWrappedSQLError(err, q, services.ErrSQLScan)
		}

		rs = append(rs, r)
	}

	if err := rows.Err(); err != nil {
		return nil, services.NewWrappedSQLError(err, q, services.ErrSQLRead)
	}

	return rs, nil
}

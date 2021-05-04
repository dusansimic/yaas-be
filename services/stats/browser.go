package stats

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s statsService) Browser() ([]Record, error) {
	q := psql.Select("b.name", "COUNT(r.idr)").From("record r").InnerJoin("browser b ON b.idb = r.idb").Where(sq.And{
		sq.Eq{"r.idd": s.id},
		sq.Gt{"r.timestamp": s.ts},
	}).GroupBy("b.idb")

	rows, err := q.RunWith(s.tx).Query()
	if err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLExec.Msg,
		)
	}

	bs := []Record{}
	for rows.Next() {
		var b Record
		if err := rows.Scan(&b.Name, &b.Reqs); err != nil {
			return nil, services.NewServiceError(
				services.NewSQLError(err, q),
				services.ErrSQLScan.Msg,
			)
		}

		bs = append(bs, b)
	}

	if err := rows.Err(); err != nil {
		return nil, services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLRead.Msg,
		)
	}

	return bs, nil
}

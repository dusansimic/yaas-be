package referrer

import (
	"github.com/dusansimic/yaas/services"
)

func (s referrerService) Add(n string, d string) error {
	q := psql.Insert("referrer").Columns("name", "domain").Values(n, d)

	res, err := q.RunWith(s.tx).Exec()
	if err != nil {
		return services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLExec.Msg,
		)
	}

	added, err := res.RowsAffected()
	if err != nil {
		return services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrUnknown.Msg,
		)
	}

	if added == 0 {
		return services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrNotAdded.Msg,
		)
	}

	return nil
}

package domain

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s *domainService) Delete(id int) error {
	q := psql.Delete("domain").Where(sq.Eq{"idd": id})

	res, err := q.RunWith(s.tx).Exec()
	if err != nil {
		return services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLExec.Msg,
		)
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return services.NewServiceError(
			err,
			services.ErrUnknown.Msg,
		)
	}

	if deleted == 0 {
		return services.ErrNotDeleted
	}

	return nil
}

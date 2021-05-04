package domain

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s domainService) EditDesc(id int, desc string) error {
	q := psql.Update("domain").Set("description", desc).Where(sq.Eq{"idd": id})

	res, err := q.RunWith(s.tx).Exec()
	if err != nil {
		return services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrSQLExec.Msg,
		)
	}

	edited, err := res.RowsAffected()
	if err != nil {
		return services.NewServiceError(
			err,
			services.ErrUnknown.Msg,
		)
	}

	if edited == 0 {
		return services.ErrNotEdited
	}

	return nil
}

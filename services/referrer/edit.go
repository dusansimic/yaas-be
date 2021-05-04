package referrer

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

func (s referrerService) Edit(id int, n string) error {
	q := psql.Update("referrer").Set("name", n).Where(sq.Eq{"id": id})

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
			services.NewSQLError(err, q),
			services.ErrUnknown.Msg,
		)
	}

	if edited == 0 {
		return services.NewServiceError(
			services.NewSQLError(err, q),
			services.ErrNotEdited.Msg,
		)
	}

	return nil
}

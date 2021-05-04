package domain

import (
	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
)

func (s *domainService) Add(d yaas.Domain) error {
	q := psql.Insert("domain").Columns("idu", "code", "name", "description").Values(d.UserID, d.Code, d.Domain, d.Desc)

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
			err,
			services.ErrUnknown.Msg,
		)
	}

	if added == 0 {
		return services.ErrNotAdded
	}

	return nil
}

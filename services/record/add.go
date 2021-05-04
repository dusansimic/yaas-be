package record

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/dusansimic/yaas/services"
)

// AddEvent adds one event to the database
func (s recordService) Add(r Record) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := psql.Insert("record").Columns("timestamp", "idd", "url", "path", "idrf", "iddv", "idc", "idb", "idos").Values(r.Timestamp, r.DomainID, r.Url, r.Path, r.ReferrerID, r.DeviceID, r.CountryID, r.BrowserID, r.OSID)

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

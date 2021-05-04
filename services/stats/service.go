package stats

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
)

type Service interface {
	Browser() ([]Record, error)
	Device() ([]Record, error)
	OS() ([]Record, error)
	Path() ([]Record, error)
	Referrer() ([]Record, error)
	Total() ([]Record, error)
	WithTime(ts time.Time) Service
	WithDuration(dur time.Duration) Service
}

type statsService struct {
	tx *sql.Tx
	id int
	ts time.Time
}

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewService(tx *sql.Tx, id int) Service {
	return &statsService{
		tx: tx,
		id: id,
	}
}

func (s *statsService) WithTime(ts time.Time) Service {
	s.ts = ts
	return s
}

func (s *statsService) WithDuration(dur time.Duration) Service {
	s.ts = s.ts.Add(dur)
	return s
}

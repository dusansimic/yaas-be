package record

import (
	"database/sql"
)

type Service interface {
	Add(r Record) error
}

type recordService struct {
	tx *sql.Tx
}

func NewService(tx *sql.Tx) Service {
	return &recordService{
		tx: tx,
	}
}

package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/dusansimic/yaas/services/record"
	"github.com/gin-gonic/gin"
)

var (
	// ErrNoDomain means that specified domain desn't exist.
	ErrNoDomain = errors.New("event: domain desn't exist")
)

// AddRecord returns a gin handler function that handles event requests
func AddRecord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.MustGet("rec").(record.Record)

		if err := addRecord(db, r, c.GetInt("uid")); err != nil {
			if err == ErrNoDomain {
				c.Error(&appError{err, "Domain not found", http.StatusBadRequest})
			} else {
				c.Error(&appError{err, "Failed to add event", http.StatusInternalServerError})
			}
			return
		}

		c.Status(http.StatusOK)
	}
}

// addRecord processes the request, converts the event to structure that will be stored and adds
// it to the database.
func addRecord(db *sql.DB, r record.Record, uid int) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		var closeError error
		if err != nil {
			closeError = tx.Rollback()
		} else {
			closeError = tx.Commit()
		}
		if err == nil && closeError != nil {
			err = closeError
		}
	}()

	recordService := record.NewService(tx)

	// Add the event to database
	if err := recordService.Add(r); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		fmt.Println(err)
		return err
	}

	return nil
}

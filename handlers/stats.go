package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	statsService "github.com/dusansimic/yaas/services/stats"
	"github.com/gin-gonic/gin"
)

type statsSource string

const (
	Browser  statsSource = "browser"
	OS       statsSource = "os"
	Device   statsSource = "device"
	Referrer statsSource = "referrer"
	Path     statsSource = "path"
	Total    statsSource = "total"
)

var (
	ErrSourceNotImplemented error = errors.New("stats source not implemented")
)

func Stats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("id")
		id, err := strconv.ParseInt(sid, 10, 32)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, &appError{err, "Couldn't parse domain id", http.StatusBadRequest})
			return
		}

		q := statsSource(c.Query("t"))
		if q == "" {
			c.AbortWithError(http.StatusBadRequest, &appError{err, "You need to specify stats source", http.StatusBadRequest})
			return
		}

		recs, err := getBrowserStats(db, int(id), q)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, &appError{err, "Couldn't get stats", http.StatusInternalServerError})
		}

		fmt.Println(recs)

		c.JSON(http.StatusOK, recs)
	}
}

func getBrowserStats(db *sql.DB, id int, s statsSource) ([]statsService.Record, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		var closeErr error
		if err != nil {
			closeErr = tx.Rollback()
		} else {
			closeErr = tx.Commit()
		}
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	service := statsService.NewService(tx, id).WithTime(time.Now()).WithDuration(-time.Hour * 24 * 30)

	switch s {
	case Browser:
		return service.Browser()
	case OS:
		return service.OS()
	case Device:
		return service.Device()
	case Referrer:
		return service.Referrer()
	case Path:
		return service.Path()
	case Total:
		return service.Total()
	}

	return nil, ErrSourceNotImplemented
}

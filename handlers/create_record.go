package handlers

import (
	"database/sql"
	"net/http"

	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services/browser"
	"github.com/dusansimic/yaas/services/country"
	"github.com/dusansimic/yaas/services/device"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/dusansimic/yaas/services/os"
	"github.com/dusansimic/yaas/services/record"
	"github.com/dusansimic/yaas/services/referrer"
	"github.com/gin-gonic/gin"
)

func CreateRecord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		e := c.MustGet("evt").(yaas.Event)
		u := c.GetInt("uid")

		r, err := createRecord(db, e, u)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, &appError{err, "Could not create a record", http.StatusInternalServerError})
			return
		}

		c.Set("rec", r)
	}
}

func createRecord(db *sql.DB, e yaas.Event, u int) (r record.Record, err error) {
	tx, err := db.Begin()
	if err != nil {
		return
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

	idd, err := domain.NewService(tx).CodeToID(e.Code)
	if err != nil {
		return
	}
	r.DomainID = idd

	idr, err := ensureReferrerID(tx, e.Referrer)
	if err != nil {
		return
	}
	r.ReferrerID = idr

	iddv, err := device.NewService(tx).ID(e.Device)
	if err != nil {
		return
	}
	r.DeviceID = iddv

	idc, err := country.NewService(tx).IDWithISO2(e.CountryCode)
	if err != nil {
		return
	}
	r.CountryID = idc

	idb, err := browser.NewService(tx).ID(e.Browser)
	if err != nil {
		return
	}
	r.BrowserID = idb

	idos, err := os.NewService(tx).ID(e.OS)
	if err != nil {
		return
	}
	r.OSID = idos

	r.Timestamp = e.Time
	r.Url = e.URL
	r.Path = e.Path

	return
}

// ensureReferrerID will get the id of referrer or if it doesn't exist, it will create one and
// return its id
func ensureReferrerID(tx *sql.Tx, r string) (int, error) {
	s := referrer.NewService(tx)

	id, err := s.ID(r)
	if err == nil {
		return id, nil
	}

	id, err = s.AddReturnID(r, r)
	if err != nil {
		return 0, err
	}

	return id, nil
}

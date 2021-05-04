package handlers

import (
	"database/sql"
	"net/http"

	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/gin-gonic/gin"
	"github.com/jkomyno/nanoid"
)

// AddDomain handles post domain requests. It adds a new domain to the authenticated user.
func AddDomain(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d yaas.Domain
		if err := c.ShouldBindJSON(&d); err != nil {
			c.AbortWithError(http.StatusBadRequest, &appError{err, "Bad body format", http.StatusBadRequest})
			return
		}

		code, err := nanoid.Nanoid()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, &appError{err, "Failed to generate code", http.StatusInternalServerError})
		}

		d.Code = code
		d.UserID = c.GetInt("uid")

		if err := addDomain(db, d); err != nil {
			c.AbortWithError(http.StatusInternalServerError, &appError{err, "Failed to add domain", http.StatusInternalServerError})
			return
		}

		c.Status(http.StatusOK)
	}
}

func addDomain(db *sql.DB, d yaas.Domain) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
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

	domainService := domain.NewService(tx)

	// Just return result of Add since if there's an error it should be returnet, and if there is
	// no error, it should be nil.
	return domainService.Add(d)
}

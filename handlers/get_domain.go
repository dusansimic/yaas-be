package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/gin-gonic/gin"
)

var (
	// ErrInvalidReqDomain means that requested domain was invalid.
	ErrInvalidReqDomain = errors.New("domain: requested domain was invalid")
)

type domainQuery struct {
	ID   int    `form:"id"`
	Many string `form:"many"`
}

// GetDomain handles requests to get information about a specific domain.
func GetDomain(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q domainQuery
		u := c.GetInt("uid")
		if err := c.ShouldBindQuery(&q); err != nil {
			c.AbortWithError(http.StatusBadRequest, &appError{err, "Couldn't parse query string", http.StatusBadRequest})
			return
		}

		if _, ok := c.GetQuery("many"); ok {
			ds, err := getDomains(u, db)
			if err != nil {
				if err == ErrInvalidUser {
					c.AbortWithError(http.StatusBadRequest, &appError{err, "Invalid user", http.StatusBadRequest})
				} else {
					c.AbortWithError(http.StatusInternalServerError, &appError{err, "Failed to get domain", http.StatusInternalServerError})
				}
				return
			}

			fmt.Println(ds)

			c.JSON(http.StatusOK, ds)
			return
		}

		d, err := getDomain(u, db, q.ID)
		if err != nil {
			if err == ErrInvalidReqDomain {
				c.AbortWithError(http.StatusNotFound, &appError{err, "Domain not found", http.StatusNotFound})
			} else {
				c.AbortWithError(http.StatusInternalServerError, &appError{err, "Couldn't find domain", http.StatusInternalServerError})
			}
			return
		}

		c.JSON(http.StatusOK, d)
	}
}

// getDomains gets all domains for a specified user
func getDomains(u int, db *sql.DB) (ds []yaas.Domain, err error) {
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

	ds, err = domainService.GetMany(u)

	return
}

func getDomain(u int, db *sql.DB, id int) (dd yaas.Domain, err error) {
	tx, err := db.Begin()
	if err != nil {
		return yaas.Domain{}, err
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

	dd, err = domainService.Get(id)
	if err != nil {
		if equalError(err, services.ErrNotFound) {
			return yaas.Domain{}, ErrInvalidReqDomain
		}
		return
	}

	if dd.UserID != u {
		return yaas.Domain{}, ErrInvalidReqDomain
	}

	return dd, nil
}

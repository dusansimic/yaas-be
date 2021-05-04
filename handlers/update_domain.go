package handlers

import (
	"database/sql"
	"net/http"

	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services/domain"
	"github.com/gin-gonic/gin"
)

// UpdateDomain handles update requests for a domain. To update a domain, request needs to
// contain the domain and the description that needs to be updated.
func UpdateDomain(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var d yaas.Domain
		if err := ctx.ShouldBindJSON(&d); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, &appError{err, "Couldn't parse body", http.StatusInternalServerError})
			return
		}

		if err := updateDescription(db, d); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, &appError{err, "Couldn't edit description", http.StatusInternalServerError})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func updateDescription(db *sql.DB, d yaas.Domain) (err error) {
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

	if err = domainService.EditDesc(d.ID, d.Desc); err != nil {
		return
	}

	return
}

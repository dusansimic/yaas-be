package handlers

import (
	"database/sql"
	"net/http"

	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/user"
	"github.com/gin-gonic/gin"
)

type userQuery struct {
	Username string `form:"username"`
}

// AvailableUser handles requests to get information about a specific domain.
func AvailableUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q userQuery
		if err := c.ShouldBindQuery(&q); err != nil {
			c.AbortWithError(http.StatusBadRequest, &appError{err, "Bad query format", http.StatusBadRequest})
			return
		}

		_, err := getUser(c, db, q.Username)
		if err != nil {
			if err == ErrInvalidUser {
				c.Status(http.StatusOK)
			} else {
				c.AbortWithError(http.StatusInternalServerError, &appError{err, "Couldn't check for user", http.StatusInternalServerError})
			}
			return
		}

		c.AbortWithError(http.StatusConflict, &appError{ErrUserExists, "User already exists", http.StatusConflict})
	}
}

func getUser(ctx *gin.Context, db *sql.DB, u string) (uu yaas.User, err error) {
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

	userService := user.NewService(tx)

	uu, err = userService.Get(u)
	if err != nil {
		if equalError(err, services.ErrNotFound) {
			return uu, ErrInvalidUser
		}
		return
	}

	return
}

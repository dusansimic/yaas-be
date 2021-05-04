package handlers

import (
	"database/sql"
	"net/http"

	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/internal/parser"
	"github.com/gin-gonic/gin"
)

func ParseEvent(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var re yaas.RawEvent
		if err := c.ShouldBindJSON(&re); err != nil {
			c.AbortWithError(http.StatusBadRequest, &appError{err, "Bad body format", http.StatusBadRequest})
			return
		}

		ye, err := parser.Parse(c, re)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, &appError{err, "Couldn't parse event", http.StatusInternalServerError})
			return
		}

		c.Set("evt", ye)
	}
}

package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Logout handles logout requests
func Logout(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cs := sessions.DefaultMany(c, "cookie")
		ss := sessions.DefaultMany(c, "server")
		cs.Clear()
		ss.Clear()
		c.Status(http.StatusOK)
	}
}

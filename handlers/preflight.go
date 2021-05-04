package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.Status(http.StatusOK)
}

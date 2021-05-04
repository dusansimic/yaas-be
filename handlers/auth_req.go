package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dusansimic/yaas/services/auth"
	"github.com/gin-gonic/gin"
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := getAuthHeader(c)
		if auth == "" {
			c.AbortWithError(http.StatusUnauthorized, &appError{
				errors.New("authorization header not specified"),
				"Authorization error not specified",
				http.StatusUnauthorized,
			})
			return
		}

		prefix := "Bearer "
		if ok := strings.HasPrefix(auth, prefix); !ok {
			c.AbortWithError(http.StatusUnauthorized, &appError{
				errors.New("invalid auth header"),
				"Invalid auth header",
				http.StatusUnauthorized,
			})
			return
		}

		if err := authenticateJWT(c, auth[len(prefix):]); err != nil {
			c.AbortWithError(http.StatusInternalServerError, &appError{err, "Failed to authorize with token", http.StatusInternalServerError})
		}
	}
}

func getAuthHeader(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	if h != "" {
		return h
	}

	h = c.GetHeader("authorization")
	return h
}

func authenticateJWT(c *gin.Context, t string) (err error) {
	s := auth.JWTAuthService()
	token, err := s.ValidateToken(t)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*auth.YaasAuthClaims)

	if ok && token.Valid {
		c.Set("u", claims.Username)
		c.Set("uid", claims.UserID)
	}

	return
}

package handlers

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"

	"github.com/dusansimic/yaas/services"
	jwt "github.com/dusansimic/yaas/services/auth"
	"github.com/dusansimic/yaas/services/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidUser means that specified user was not found
	ErrInvalidUser = errors.New("login: invalid user")
	// ErrInvalidPass means that specified password was invalid
	ErrInvalidPass = errors.New("login: invalid password")
)

type nextQuery struct {
	Next string `form:"next"`
}

// Login handles login requests. Auth credentials are sent in the body.
func Login(db *sql.DB, next string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authentication params
		var p authCreds
		if err := c.ShouldBindJSON(&p); err != nil {
			c.AbortWithError(http.StatusBadRequest, &appError{err, "Bad body format", http.StatusBadRequest})
			return
		}

		// Get redirect url. If it's not specified, set redirect url to base of frontend.
		var q nextQuery
		if err := c.ShouldBindQuery(&q); err != nil || q.Next == "" {
			q.Next = next
		}

		t, err := loginJWT(db, p)
		if err != nil {
			if err == ErrInvalidUser {
				c.AbortWithError(http.StatusUnauthorized, &appError{err, "Wrong username", http.StatusUnauthorized})
			} else if err == ErrInvalidPass {
				c.AbortWithError(http.StatusUnauthorized, &appError{err, "Wrong password", http.StatusUnauthorized})
			} else {
				c.AbortWithError(http.StatusUnauthorized, &appError{err, "Failed to create token", http.StatusInternalServerError})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": t,
		})
	}
}

func loginJWT(db *sql.DB, c authCreds) (t string, err error) {
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

	userService := user.NewService(tx)

	// Get all user data to check credentials
	u, err := userService.Get(c.Username)
	if err != nil {
		if equalError(err, services.ErrNotFound) {
			err = ErrInvalidUser
		}
		return
	}

	// Hash the password with salt from user
	hash := argon2.IDKey([]byte(c.Password), u.Salt, 3, 32*1024, 4, 32)

	// Check if the hashes are the same
	if !bytes.Equal(hash, u.PasswordHash) {
		err = ErrInvalidPass
		return
	}

	// Generate a token
	jwtService := jwt.JWTAuthService()
	t, err = jwtService.GenerateToken(u.Username, u.ID)
	if err != nil {
		return
	}
	return
}

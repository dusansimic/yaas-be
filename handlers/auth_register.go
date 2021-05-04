package handlers

import (
	"crypto/rand"
	"database/sql"
	"net/http"

	yaas "github.com/dusansimic/yaas"
	"github.com/dusansimic/yaas/services"
	"github.com/dusansimic/yaas/services/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

// Register hadnlers register requests for the api.
// If body is not in valid format, handler returns code 400 with ErrNoAuthParams.
// If the user already exists, handler returns code 400 with ErrUserExists.
// If the user doesn't exist, it is added and the salt is saved with the user.
func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p authCreds
		if err := c.ShouldBindJSON(&p); err != nil {
			c.Error(&appError{err, "Bad body format", http.StatusBadRequest})
			return
		}

		if err := register(db, p); err != nil {
			if err == ErrUserExists {
				c.Error(&appError{err, "User exists", http.StatusBadRequest})
			} else {
				c.Error(&appError{err, "Failed to register user", http.StatusInternalServerError})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
		})
	}
}

func register(db *sql.DB, c authCreds) (err error) {
	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
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

	// Check if the user already exists
	_, err = userService.Get(c.Username)

	// If there is no error, then user exists
	if err == nil {
		return ErrUserExists
	} else if !equalError(err, services.ErrNotFound) {
		// There is some other error, just scream
		return
	}

	// If user was not found, add it
	salt := make([]byte, 8)
	_, err = rand.Read(salt)
	if err != nil {
		return
	}

	hash := argon2.IDKey([]byte(c.Password), salt, 3, 32*1024, 4, 32)

	u := yaas.User{
		Username:     c.Username,
		PasswordHash: hash,
		Salt:         salt,
	}

	if err := userService.Add(u); err != nil {
		return err
	}

	return nil
}

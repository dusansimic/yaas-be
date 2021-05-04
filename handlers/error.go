package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type appError struct {
	Err  error
	Msg  string
	Code int
}

func (e appError) Error() string {
	return e.Msg
}

func (e appError) Unwrap() error {
	return e.Err
}

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Errors.Last() == nil {
			return
		}

		e := c.Errors.Last().Err.(*appError)

		// TODO: log somewhere error stacks
		for _, err := range stackErrors(e) {
			fmt.Println(err)
		}

		c.JSON(e.Code, gin.H{
			"message": e.Msg,
			"code":    e.Code,
		})
	}
}

func stackErrors(e error) (s []error) {
	for e != nil {
		s = append(s, e)
		e = errors.Unwrap(e)
	}
	return
}

// Error handles error requests. This is a Plausible thing and I don't know what it does.
func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	}
}

func equalError(e1, e2 error) bool {
	return e1.Error() == e2.Error()
}

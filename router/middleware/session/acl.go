package session

import (
	"encoding/base32"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harborapp/harbor-api/model"
	"github.com/harborapp/harbor-api/shared/token"
	"github.com/harborapp/harbor-api/store"
	"github.com/jinzhu/gorm"
)

const (
	// CurrentContextKey defines the context key that stores the user.
	CurrentContextKey = "current"

	// TokenContextKey defines the context key that stores the token.
	TokenContextKey = "token"
)

// Current gets the user from the context.
func Current(c *gin.Context) *model.User {
	v, ok := c.Get(CurrentContextKey)

	if !ok {
		return nil
	}

	r, ok := v.(*model.User)

	if !ok {
		return nil
	}

	return r
}

// SetCurrent injects the user into the context.
func SetCurrent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			record *model.User
		)

		parsed, err := token.Parse(
			c.Request,
			func(t *token.Token) ([]byte, error) {
				var (
					res *gorm.DB
				)

				record, res = store.GetUser(
					c,
					t.Text,
				)

				signingKey, _ := base32.StdEncoding.DecodeString(record.Hash)
				return signingKey, res.Error
			},
		)

		if err == nil {
			c.Set(TokenContextKey, parsed)
			c.Set(CurrentContextKey, record)
		}

		c.Next()
	}
}

// MustCurrent validates the user access.
func MustCurrent() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := Current(c)

		if user == nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  http.StatusUnauthorized,
					"message": "You have to be authenticated",
				},
			)

			c.Abort()
		} else {
			c.Next()
		}
	}
}

// MustNobody validates anonymous users.
func MustNobody() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := Current(c)

		if user != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  http.StatusUnauthorized,
					"message": "You have to be a guest user",
				},
			)

			c.Abort()
		} else {
			c.Next()
		}
	}
}

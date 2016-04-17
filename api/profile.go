package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/harborapp/harbor-api/router/middleware/session"
	"github.com/harborapp/harbor-api/store"
)

// GetProfile displays the current profile.
func GetProfile(c *gin.Context) {
	record := session.Current(c)

	c.JSON(
		http.StatusOK,
		record,
	)
}

// PatchProfile updates the current profile.
func PatchProfile(c *gin.Context) {
	record := session.Current(c)

	if err := c.BindJSON(&record); err != nil {
		logrus.Warn("Failed to bind profile data")
		logrus.Warn(err)

		c.JSON(
			http.StatusPreconditionFailed,
			gin.H{
				"status":  http.StatusPreconditionFailed,
				"message": "Failed to bind profile data",
			},
		)

		c.Abort()
		return
	}

	err := store.UpdateUser(
		c,
		record,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			},
		)

		c.Abort()
		return
	}

	c.JSON(
		http.StatusOK,
		record,
	)
}

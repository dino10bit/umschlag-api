package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umschlag/umschlag-api/config"
)

// IndexInfo represents the API index.
func IndexInfo(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"api":     "Umschlag API",
			"version": config.Version,
			"stream":  "master",
		},
	)
}

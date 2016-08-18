package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umschlag/umschlag-api/router/middleware/session"
	"github.com/umschlag/umschlag-api/store"
)

// RepoIndex retrieves all available repos.
func RepoIndex(c *gin.Context) {
	records, err := store.GetRepos(
		c,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to fetch repos",
			},
		)

		c.Abort()
		return
	}

	c.JSON(
		http.StatusOK,
		records,
	)
}

// RepoShow retrieves a specific repo.
func RepoShow(c *gin.Context) {
	record := session.Repo(c)

	c.JSON(
		http.StatusOK,
		record,
	)
}

// RepoDelete removes a specific repo.
func RepoDelete(c *gin.Context) {
	record := session.Repo(c)

	err := store.DeleteRepo(
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
		gin.H{
			"status":  http.StatusOK,
			"message": "Successfully deleted repo",
		},
	)
}
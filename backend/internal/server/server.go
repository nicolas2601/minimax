package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// New builds the Gin engine. The DB is accepted for future handlers
// (auth, accounts, transactions) but /health remains a pure liveness
// probe that does not touch the database — use a separate /ready for
// that when we add one.
func New(db *gorm.DB) *gin.Engine {
	_ = db // intentionally unused for now; wired for upcoming tasks

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return r
}
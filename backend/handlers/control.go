package handlers

import (
	"net/http"

	"plocate-ui/indexer"

	"github.com/gin-gonic/gin"
)

func StartIndexing(c *gin.Context) {
	if err := indexer.Instance.StartIndexing(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "indexing started"})
}

func StopIndexing(c *gin.Context) {
	if err := indexer.Instance.StopIndexing(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "indexing stopped"})
}

func EnableScheduler(c *gin.Context) {
	indexer.Instance.EnableScheduler()
	c.JSON(http.StatusOK, gin.H{"message": "scheduler enabled"})
}

func DisableScheduler(c *gin.Context) {
	indexer.Instance.DisableScheduler()
	c.JSON(http.StatusOK, gin.H{"message": "scheduler disabled"})
}

package handlers

import (
	"net/http"

	"plocate-ui/indexer"

	"github.com/gin-gonic/gin"
)

func StartIndexing(c *gin.Context) {
	indexName := c.Param("indexName")

	// If no index name specified, start all enabled indices
	if indexName == "" {
		if err := indexer.Instance.StartIndexingAll(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "indexing started for all enabled indices"})
		return
	}

	if err := indexer.Instance.StartIndexing(indexName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "indexing started for " + indexName})
}

func StopIndexing(c *gin.Context) {
	indexName := c.Param("indexName")

	// If no index name specified, stop all indices
	if indexName == "" {
		if err := indexer.Instance.StopIndexingAll(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "indexing stopped for all indices"})
		return
	}

	if err := indexer.Instance.StopIndexing(indexName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "indexing stopped for " + indexName})
}

func GetIndices(c *gin.Context) {
	indices := indexer.Instance.GetIndexNames()
	c.JSON(http.StatusOK, gin.H{"indices": indices})
}

func EnableScheduler(c *gin.Context) {
	indexer.Instance.EnableScheduler()
	c.JSON(http.StatusOK, gin.H{"message": "scheduler enabled"})
}

func DisableScheduler(c *gin.Context) {
	indexer.Instance.DisableScheduler()
	c.JSON(http.StatusOK, gin.H{"message": "scheduler disabled"})
}

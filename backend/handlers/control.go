package handlers

import (
	"net/http"
	"strings"

	"plocate-ui/config"
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

type AddIndexRequest struct {
	Name       string   `json:"name" binding:"required"`
	IndexPaths []string `json:"index_paths" binding:"required"`
}

func AddIndex(c *gin.Context) {
	var req AddIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if len(req.IndexPaths) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one index path is required"})
		return
	}

	// Trim whitespace from paths and filter empty
	var paths []string
	for _, p := range req.IndexPaths {
		p = strings.TrimSpace(p)
		if p != "" {
			paths = append(paths, p)
		}
	}
	if len(paths) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one non-empty index path is required"})
		return
	}

	// Add to config (persists to disk)
	idx, err := config.AddIndex(name, paths)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Register in running indexer
	indexer.Instance.AddIndex(*idx)

	c.JSON(http.StatusOK, gin.H{"message": "index added", "index": idx})
}

func RemoveIndex(c *gin.Context) {
	indexName := c.Param("indexName")
	if indexName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "index name is required"})
		return
	}

	// Remove from running indexer first (stops if running)
	if err := indexer.Instance.RemoveIndex(indexName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove from config (persists to disk)
	if err := config.RemoveIndex(indexName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "index removed"})
}

package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"plocate-ui/indexer"

	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	Query   string   `json:"query" binding:"required"`
	Limit   int      `json:"limit"`
	Indices []string `json:"indices"` // Optional: if empty, searches all enabled indices
}

type SearchResponse struct {
	Results []string `json:"results"`
	Count   int      `json:"count"`
}

func Search(c *gin.Context) {
	var req SearchRequest

	// Support both GET and POST
	if c.Request.Method == http.MethodGet {
		req.Query = c.Query("q")
		if limit := c.Query("limit"); limit != "" {
			req.Limit, _ = strconv.Atoi(limit)
		}
		// Parse indices from query parameter (comma-separated)
		if indices := c.Query("indices"); indices != "" {
			req.Indices = strings.Split(indices, ",")
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 100
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	results, err := indexer.Instance.Search(req.Query, req.Limit, req.Indices)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SearchResponse{
		Results: results,
		Count:   len(results),
	})
}

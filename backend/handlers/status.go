package handlers

import (
	"net/http"

	"plocate-ui/indexer"

	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	status := indexer.Instance.GetStatus()
	c.JSON(http.StatusOK, status)
}

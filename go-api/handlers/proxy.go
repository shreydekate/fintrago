package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProxyUpload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented yet"})
}

func ProxyAsk(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented yet"})
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBalance(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "not implemented yet"})
}

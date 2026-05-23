// Learned that packages are folder high (idiomatic) and that the name of the package
// is the same as the folder name.

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shreydekate/fintrago/models"
)

// In-memory temp store for creating transactions.
var store []models.Transaction

// c is pointing to the address of the request context, this contains the information
// about the request and response. * is a pointer in Golang and in this case it matters
// for performance etc etc
func CreateTransaction(c *gin.Context) {
	var t models.Transaction // creates an empty Transaction struct from the models package
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // gin.H is a shortcut to create map[string]any
		return
	} // weird ahh error handling
	store = append(store, t)
	c.JSON(http.StatusCreated, t) //c.JSON(statusCode, data)
}

func ListTransactions(c *gin.Context) {
	c.JSON(http.StatusOK, store)
}

func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	for i, t := range store {
		if t.ID == id {
			store = append(store[:i], store[i+1:]...) // Deleting elements from a slice (idiomatic)
			c.JSON(http.StatusOK, gin.H{"message": "Poof Gone Forever"})
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
}

// Learned that packages are folder high (idiomatic) and that the name of the package
// is the same as the folder name.

package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shredekate/fintrago/db"
	"github.com/shreydekate/fintrago/models"
)

// c is pointing to the address of the request context, this contains the information
// about the request and response. * is a pointer in Golang and in this case it matters
// for performance etc etc
func CreateTransaction(c *gin.Context) {
	var t models.Transaction // creates an empty Transaction struct from the models package
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // gin.H is a shortcut to create map[string]any
		return
	} // weird ahh error handling
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO transactions (user_id, type, amount, category, description, date)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, created at`,
		t.UserID, t.Type, t.Amount, t.Category, t.Description, t.Date).Scan(&t.ID, &t.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t) //c.JSON(statusCode, data)
}

func ListTransactions(c *gin.Context) {
	row, err := db.pool.Query(
		context.Background(),
		`SELECT ID, user_id, type, amount, category, description, date, created_at FROM transactions`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer row.Close() // tf is ts?

	var transactions []models.Transaction
	for rows.Next() {
		// End Here
	}
}

func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	for i, t := range store {
		if t.ID == id {
			store = append(store[:i], store[i+1:]...) // Deleting elements from a slice (idiomatic)
			c.JSON(http.StatusOK, gin.H{"message": "Poof Gone Forever"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
}

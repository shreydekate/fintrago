package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shreydekate/fintrago/handlers"
)

func Register(r *gin.Engine) {
	api := r.Group("/api") 
	{
		api.POST("/transactions", handlers.CreateTransaction)
		api.GET("/transactions", handlers.ListTransactions)
		api.DELETE("/transactions/:id", handlers.DeleteTransaction)
		api.GET("/balance", handlers.GetBalance)
		api.POST("/upload", handlers.ProxyUpload)
		api.POST("/ask", handlers.ProxyAsk)

	}
}
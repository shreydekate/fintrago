package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shreydekate/fintrago/config"
	"github.com/shreydekate/fintrago/routes"
)

func main() {
	cfg := config.Load()
	// db.Connect(cfg.DatabaseURL)

	r := gin.Default()
	routes.Register(r)

	log.Printf("Server running on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}

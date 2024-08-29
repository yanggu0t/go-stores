package main

import (
	"log"

	"github.com/yanggu0t/go-rdbms-practice/internal/database"
	"github.com/yanggu0t/go-rdbms-practice/internal/router"
)

func main() {
	// Initialize database
	database.InitDB()

	// Set up router
	r := router.SetupRouter()

	// Start server
	log.Fatal(r.Run(":8080"))
}

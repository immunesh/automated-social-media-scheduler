package main

import (
	"log"
	"os"

	"github.com/immunesh/automated-social-media-scheduler/cmd/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := server.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
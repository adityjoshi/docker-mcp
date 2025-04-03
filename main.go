package main

import (
	"docker-nlp/docker"
	"docker-nlp/handlers"
	"docker-nlp/middleware"
	"docker-nlp/nlp"
	"log"

	"github.com/adityjoshi/docker-mcp/config"
	"github.com/gin-gonic/gin"
	"honnef.co/go/tools/config"
)

func main() {

	config := config.LoadConfig()
	nlpProcessor := nlp.NewProcessor()
	dockerExecutor := docker.NewExecutor()

	dockerHandler := handlers.NewDockerHandler(nlpProcessor, dockerExecutor)
	router := gin.Default()

	// Apply API key middleware to protected routes
	authorized := router.Group("/")
	authorized.Use(middleware.APIKeyAuth(config.APIKey))
	{
		authorized.POST("/docker", dockerHandler.ProcessCommand)
	}

	// Public routes
	router.GET("/health", handlers.HealthCheck)

	// Start the server
	log.Printf("Starting MCP Docker Natural Language Server on :%s...", &config.Port)
	log.Println("Endpoints:")
	log.Println("  POST /docker - Send natural language commands (protected with API key)")
	log.Println("  GET /health - Check server health")

	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

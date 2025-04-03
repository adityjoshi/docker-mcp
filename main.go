package main

import (
	"log"

	"github.com/adityjoshi/docker-mcp/handler"
	"github.com/adityjoshi/docker-mcp/middleware"
	"github.com/adityjoshi/docker-mcp/nlp"

	"github.com/adityjoshi/docker-mcp/config"
	"github.com/gin-gonic/gin"
)

func main() {

	config := config.LoadConfig()
	nlpProcessor := nlp.NewProcessor()
	dockerExecutor := docker.NewExecutor()

	dockerHandler := handler.NewDockerHandler(nlpProcessor, dockerExecutor)
	router := gin.Default()

	authorized := router.Group("/")
	authorized.Use(middleware.APIKEYAuth(config.APIKey))
	{
		authorized.POST("/docker", dockerHandler.ProcessCommand)
	}

	// Public routes
	router.GET("/health", handler.HealthCheck)

	// Start the server
	log.Printf("Starting MCP Docker Natural Language Server on :%s...", &config.Port)
	log.Println("Endpoints:")
	log.Println("  POST /docker - Send natural language commands (protected with API key)")
	log.Println("  GET /health - Check server health")

	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

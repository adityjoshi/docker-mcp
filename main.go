package main

import (
	"log"

	"github.com/adityjoshi/docker-mcp/config"
	"github.com/adityjoshi/docker-mcp/docker"
	"github.com/adityjoshi/docker-mcp/handler"
	"github.com/adityjoshi/docker-mcp/nlp"
	"github.com/gin-gonic/gin"
)

func main() {

	config := config.LoadConfig()
	nlpProcessor := nlp.NewProcessor()
	dockerExecutor := docker.NewExecutor()

	dockerHandler := handler.NewDockerHandler(nlpProcessor, dockerExecutor)
	router := gin.Default()

	authorized := router.Group("/")
	authorized.POST("/docker", dockerHandler.ProcessCommand)
	// authorized.Use(middleware.APIKEYAuth(config.APIKey))
	// {

	// }

	// Public routes
	router.GET("/health", handler.HealthCheck)

	// Start the server
	log.Printf("starting MCP Docker Natural Language Server on :%s...", &config.Port)
	log.Println("Endpoints:")
	log.Println("  POST /docker - Send natural language commands (protected with API key)")
	log.Println("  GET /health - Check server health")

	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

package handler

import (
	"net/http"

	"github.com/adityjoshi/docker-mcp/nlp"
	"github.com/gin-gonic/gin"
)

type DockerHandler struct {
	nlpProcessor   *nlp.Processor
	dockerExecutor *docker.Executor
}

type CommandRequest struct {
	Command string `json:"command" binding:"required"`
}

func NewDockerHandler(nlpProcessor *nlp.Processor, dockerExecutor *docker.Executor) *DockerHandler {
	return &DockerHandler{
		nlpProcessor:   nlpProcessor,
		dockerExecutor: dockerExecutor,
	}
}
func (h *DockerHandler) ProcessCommand(c *gin.Context) {
	var req CommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request payload",
		})
		return
	}

	intent, containerInfo := h.nlpProcessor.DetectIntent(req.Command)

	if intent == nlp.IntentUnknown {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Could not understand the command. Please specify create, stop, delete, or list operations for Docker containers.",
		})
		return
	}

	// Execute the corresponding Docker command
	result := h.dockerExecutor.ExecuteCommand(intent, containerInfo)

	c.JSON(http.StatusOK, result)
}

// HealthCheck returns the server health status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "MCP Docker NLP server is running",
	})
}

package docker

import (
	"os/exec"
	"strings"

	"github.com/adityjoshi/docker-mcp/nlp"
)

// Result represents the Docker command execution result
type Result struct {
	Status      string          `json:"status"`
	Message     string          `json:"message,omitempty"`
	ContainerID string          `json:"container_id,omitempty"`
	Containers  []ContainerInfo `json:"containers,omitempty"`
}

// ContainerInfo represents container details
type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

// Executor executes Docker commands
type Executor struct{}

// NewExecutor creates a new Docker executor
func NewExecutor() *Executor {
	return &Executor{}
}

// ExecuteCommand executes the appropriate Docker command based on intent
func (e *Executor) ExecuteCommand(intent string, containerInfo nlp.ContainerInfo) Result {
	switch intent {
	case nlp.IntentCreate:
		return e.createContainer(containerInfo)
	case nlp.IntentStop:
		return e.stopContainer(containerInfo)
	case nlp.IntentDelete:
		return e.deleteContainer(containerInfo)
	case nlp.IntentList:
		return e.listContainers()
	default:
		return Result{Status: "error", Message: "Unknown intent"}
	}
}

// createContainer creates a new Docker container
func (e *Executor) createContainer(info nlp.ContainerInfo) Result {
	if info.Image == "" {
		return Result{Status: "error", Message: "No image specified"}
	}

	args := []string{"run", "-d"}

	// Add name if provided
	if info.ContainerName != "" {
		args = append(args, "--name", info.ContainerName)
	}

	// Add port mappings if provided
	for _, port := range info.Ports {
		args = append(args, "-p", port.HostPort+":"+port.ContainerPort)
	}

	// Add image name
	args = append(args, info.Image)

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return Result{
			Status:  "error",
			Message: string(output),
		}
	}

	containerID := strings.TrimSpace(string(output))
	return Result{
		Status:      "success",
		Message:     "Container created with ID: " + containerID,
		ContainerID: containerID,
	}
}

// stopContainer stops a Docker container
func (e *Executor) stopContainer(info nlp.ContainerInfo) Result {
	if info.ContainerName == "" {
		return Result{Status: "error", Message: "No container name specified"}
	}

	cmd := exec.Command("docker", "stop", info.ContainerName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return Result{
			Status:  "error",
			Message: string(output),
		}
	}

	return Result{
		Status:  "success",
		Message: "Container " + info.ContainerName + " stopped",
	}
}

// deleteContainer deletes a Docker container
func (e *Executor) deleteContainer(info nlp.ContainerInfo) Result {
	if info.ContainerName == "" {
		return Result{Status: "error", Message: "No container name specified"}
	}

	// First stop it if it's running
	stopCmd := exec.Command("docker", "stop", info.ContainerName)
	stopCmd.Run()

	// Then remove it
	cmd := exec.Command("docker", "rm", info.ContainerName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return Result{
			Status:  "error",
			Message: string(output),
		}
	}

	return Result{
		Status:  "success",
		Message: "Container " + info.ContainerName + " deleted",
	}
}

// listContainers lists all Docker containers
func (e *Executor) listContainers() Result {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Status}}")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return Result{
			Status:  "error",
			Message: string(output),
		}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	containers := make([]ContainerInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}

		containers = append(containers, ContainerInfo{
			ID:     parts[0],
			Name:   parts[1],
			Image:  parts[2],
			Status: parts[3],
		})
	}

	return Result{
		Status:     "success",
		Containers: containers,
	}
}

package nlp

import (
	"regexp"
	"strings"
)

// Intent types
const (
	IntentCreate  = "create"
	IntentStop    = "stop"
	IntentDelete  = "delete"
	IntentList    = "list"
	IntentUnknown = "unknown"
)

// ContainerInfo stores extracted container information
type ContainerInfo struct {
	ContainerName string
	Image         string
	Ports         []PortMapping
}

// PortMapping represents a container port mapping
type PortMapping struct {
	HostPort      string
	ContainerPort string
}

// Processor handles NLP processing for Docker commands
type Processor struct {
	// Could add more sophisticated NLP components in the future
}

// NewProcessor creates a new NLP processor
func NewProcessor() *Processor {
	return &Processor{}
}

// DetectIntent determines the Docker operation based on the command
func (p *Processor) DetectIntent(command string) (string, ContainerInfo) {
	command = strings.ToLower(command)

	// Check for creation intent
	createIndicators := []string{"create", "make", "start", "run", "launch", "build", "spin up"}
	for _, indicator := range createIndicators {
		if strings.Contains(command, indicator) {
			return IntentCreate, p.extractContainerInfo(command)
		}
	}

	// Check for stop intent
	stopIndicators := []string{"stop", "pause", "halt"}
	for _, indicator := range stopIndicators {
		if strings.Contains(command, indicator) {
			return IntentStop, p.extractContainerInfo(command)
		}
	}

	// Check for deletion intent
	deleteIndicators := []string{"delete", "remove", "destroy", "kill"}
	for _, indicator := range deleteIndicators {
		if strings.Contains(command, indicator) {
			return IntentDelete, p.extractContainerInfo(command)
		}
	}

	// Check for list intent
	listIndicators := []string{"list", "show", "display", "all", "running"}
	for _, indicator := range listIndicators {
		if strings.Contains(command, indicator) {
			return IntentList, ContainerInfo{}
		}
	}

	return IntentUnknown, ContainerInfo{}
}

// extractContainerInfo extracts container name, image, and other parameters from command
func (p *Processor) extractContainerInfo(command string) ContainerInfo {
	info := ContainerInfo{}
	words := strings.Fields(command)

	// Extract container name (looking for patterns like "named xyz" or "called xyz")
	namePatterns := []string{"named", "called", "name", "container"}
	for i, word := range words {
		for _, pattern := range namePatterns {
			if word == pattern && i+1 < len(words) {
				// Remove any trailing punctuation
				re := regexp.MustCompile(`[^\w]`)
				info.ContainerName = re.ReplaceAllString(words[i+1], "")
				break
			}
		}
	}

	// Extract image name (looking for patterns like "from image xyz" or "using xyz")
	imagePatterns := []string{"image", "using", "from", "with"}
	for i, word := range words {
		for _, pattern := range imagePatterns {
			if word == pattern && i+1 < len(words) {
				// Allow some special characters common in image names
				re := regexp.MustCompile(`[^\w\/:.-]`)
				info.Image = re.ReplaceAllString(words[i+1], "")
				break
			}
		}
	}

	// Look for port mappings
	portPattern := regexp.MustCompile(`(\d+):(\d+)`)
	matches := portPattern.FindAllStringSubmatch(command, -1)
	for _, match := range matches {
		if len(match) == 3 {
			info.Ports = append(info.Ports, PortMapping{
				HostPort:      match[1],
				ContainerPort: match[2],
			})
		}
	}

	return info
}

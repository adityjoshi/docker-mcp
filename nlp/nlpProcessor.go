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

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) DetectIntent(command string) (string, ContainerInfo) {
	command = strings.ToLower(command)

	createIndicators := []string{"create", "make", "start", "run", "launch", "build", "spin up"}
	for _, indicator := range createIndicators {
		if strings.Contains(command, indicator) {
			return IntentCreate, p.extractContainerInfo(command)
		}
	}

	stopIndicators := []string{"stop", "pause", "halt"}
	for _, indicator := range stopIndicators {
		if strings.Contains(command, indicator) {
			return IntentStop, p.extractContainerInfo(command)
		}
	}

	deleteIndicators := []string{"delete", "remove", "destroy", "kill"}
	for _, indicator := range deleteIndicators {
		if strings.Contains(command, indicator) {
			return IntentDelete, p.extractContainerInfo(command)
		}
	}

	listIndicators := []string{"list", "show", "display", "all", "running"}
	for _, indicator := range listIndicators {
		if strings.Contains(command, indicator) {
			return IntentList, ContainerInfo{}
		}
	}

	return IntentUnknown, ContainerInfo{}
}
func (p *Processor) extractContainerInfo(command string) ContainerInfo {
	info := ContainerInfo{}
	words := strings.Fields(command)

	if len(words) >= 2 && words[0] == "run" {
		info.Image = words[1]
	}
	if len(words) >= 2 && (words[0] == "delete" || words[0] == "stop") {
		info.ContainerName = words[1]
	}

	namePatterns := []string{"named", "called", "name", "container"}
	for i, word := range words {
		for _, pattern := range namePatterns {
			if word == pattern && i+1 < len(words) {
				re := regexp.MustCompile(`[^\w]`)
				info.ContainerName = re.ReplaceAllString(words[i+1], "")
				break
			}
		}
	}

	imagePatterns := []string{"image", "using", "from", "with"}
	for i, word := range words {
		for _, pattern := range imagePatterns {
			if word == pattern && i+1 < len(words) {
				re := regexp.MustCompile(`[^\w\/:.-]`)
				info.Image = re.ReplaceAllString(words[i+1], "")
				break
			}
		}
	}

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

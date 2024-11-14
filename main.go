package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type DockerContainer struct {
	ID     string `json:"Id"`
	Name   string `json:"Names"`
	Ports  string `json:"Ports"`
	Status string `json:"Status"`
}

type PortMapping struct {
	PrivatePort int
	PublicPort  int
	Type        string
}

func main() {
	var containerName string
	flag.StringVar(&containerName, "name", "", "Container name")
	flag.Parse()

	// Only ask for container name if not provided via flag
	if containerName == "" {
		fmt.Print("Enter the container name: ")
		fmt.Scanln(&containerName)
	}

	containers := getContainers()
	matchingContainers := findMatchingContainers(containers, containerName)

	if len(matchingContainers) == 0 {
		fmt.Printf("No container found with name: %s\n", containerName)
		return
	}

	// Auto-select if single container match
	var selectedContainer *DockerContainer
	if len(matchingContainers) == 1 {
		selectedContainer = &matchingContainers[0]
	} else {
		selectedContainer = chooseContainer(matchingContainers)
	}

	// Parse ports first
	ports := parsePortString(selectedContainer.Ports)
	if len(ports) == 0 {
		log.Fatalf("No ports found for container %s", selectedContainer.Name)
	}

	// Auto-select if single port or specific port requested
	var selectedPort *PortMapping
	if len(ports) == 1 {
		selectedPort = &ports[0]
	} else {
		selectedPort = choosePort(selectedContainer, "")
	}

	openBrowser(selectedPort)
}

func getContainers() []DockerContainer {
	cmd := exec.Command("docker", "ps", "--format", "{{json .}}")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running docker ps: %v", err)
	}

	var containers []DockerContainer
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	for scanner.Scan() {
		var container DockerContainer
		if err := json.Unmarshal([]byte(scanner.Text()), &container); err != nil {
			log.Printf("Error parsing container: %v", err)
			continue
		}
		containers = append(containers, container)
	}

	return containers
}

func findMatchingContainers(containers []DockerContainer, name string) []DockerContainer {
	var matches []DockerContainer
	for _, c := range containers {
		if strings.Contains(c.Name, name) {
			matches = append(matches, c)
		}
	}
	return matches
}

func chooseContainer(containers []DockerContainer) *DockerContainer {
	options := make([]string, len(containers))
	for i, c := range containers {
		options[i] = fmt.Sprintf("%s (%s)", c.Name, c.Status)
	}

	var selected string
	prompt := &survey.Select{
		Message: "Choose a container:",
		Options: options,
	}
	survey.AskOne(prompt, &selected)

	for i, opt := range options {
		if opt == selected {
			return &containers[i]
		}
	}
	return nil
}

func parsePortString(portStr string) []PortMapping {
	if portStr == "" {
		return nil
	}

	var portMappings []PortMapping
	ports := strings.Split(portStr, ", ")

	for _, p := range ports {
		if p == "" {
			continue
		}

		// Remove the "0.0.0.0:" prefix if present
		p = strings.TrimPrefix(p, "0.0.0.0:")

		parts := strings.Split(p, "->")
		if len(parts) < 2 {
			continue
		}

		publicStr := strings.TrimSpace(parts[0])
		privateStr := strings.TrimSpace(parts[1])

		publicPort, _ := strconv.Atoi(publicStr)
		privatePort, _ := strconv.Atoi(strings.Split(privateStr, "/")[0])
		portType := strings.Split(privateStr, "/")[1]

		portMappings = append(portMappings, PortMapping{
			PrivatePort: privatePort,
			PublicPort:  publicPort,
			Type:        portType,
		})
	}
	return portMappings
}

func choosePort(container *DockerContainer, portStr string) *PortMapping {
	ports := parsePortString(container.Ports)
	if len(ports) == 0 {
		log.Fatalf("No ports found for container %s", container.Name)
	}

	if portStr != "" {
		port, _ := strconv.Atoi(portStr)
		for _, p := range ports {
			if p.PublicPort == port || p.PrivatePort == port {
				return &p
			}
		}
		log.Fatalf("Port %s not found", portStr)
	}

	options := make([]string, len(ports))
	for i, p := range ports {
		options[i] = fmt.Sprintf("%d->%d/%s", p.PublicPort, p.PrivatePort, p.Type)
	}

	var selected string
	prompt := &survey.Select{
		Message: "Choose a port:",
		Options: options,
	}
	survey.AskOne(prompt, &selected)

	for i, opt := range options {
		if opt == selected {
			return &ports[i]
		}
	}
	return nil
}

func openBrowser(port *PortMapping) {
	url := fmt.Sprintf("http://localhost:%d", port.PublicPort)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		fmt.Printf("Unsupported platform. Please open %s manually.\n", url)
		return
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Error opening browser: %v", err)
		fmt.Printf("Please open %s manually.\n", url)
	}
}

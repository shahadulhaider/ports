//go:build darwin

package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func GetListeningPorts() ([]PortInfo, error) {
	cmd := exec.Command("lsof", "-iTCP", "-P", "-n", "-sTCP:LISTEN", "-F", "pcfnPt")
	out, err := cmd.Output()
	if err != nil {
		// lsof exits 1 when no results — only fail if stderr has content
		if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
			return nil, fmt.Errorf("lsof error: %s", string(exitErr.Stderr))
		}
		if len(out) == 0 {
			return nil, nil
		}
	}

	ports := parseLsofOutput(string(out))
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Port < ports[j].Port
	})
	return ports, nil
}

// GetUDPPorts returns all UDP bound ports on macOS using lsof.
func GetUDPPorts() ([]PortInfo, error) {
	cmd := exec.Command("lsof", "-iUDP", "-P", "-n", "-F", "pcfnPt")
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && len(exitErr.Stderr) > 0 {
			return nil, fmt.Errorf("lsof error: %s", string(exitErr.Stderr))
		}
		if len(out) == 0 {
			return nil, nil
		}
	}

	ports := parseLsofOutput(string(out))
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Port < ports[j].Port
	})
	return ports, nil
}

// GetConnectionCounts returns a map of local port -> established TCP connection count.
func GetConnectionCounts() (map[int]int, error) {
	cmd := exec.Command("lsof", "-iTCP", "-P", "-n", "-sTCP:ESTABLISHED", "-F", "n")
	out, err := cmd.Output()
	if err != nil {
		// No established connections is not an error
		return map[int]int{}, nil
	}

	counts := map[int]int{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if len(line) == 0 || line[0] != 'n' {
			continue
		}
		value := line[1:]
		// Format: "local_ip:local_port->remote_ip:remote_port"
		arrowIdx := strings.Index(value, "->")
		if arrowIdx < 0 {
			continue
		}
		localPart := value[:arrowIdx]
		lastColon := strings.LastIndex(localPart, ":")
		if lastColon < 0 {
			continue
		}
		portStr := localPart[lastColon+1:]
		port, err := strconv.Atoi(portStr)
		if err != nil || port == 0 {
			continue
		}
		counts[port]++
	}
	return counts, nil
}

func parseLsofOutput(output string) []PortInfo {
	var results []PortInfo
	lines := strings.Split(output, "\n")

	var currentPID int
	var currentProcess string
	var currentType string
	var currentProtocol string

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		prefix := line[0]
		value := line[1:]

		switch prefix {
		case 'p':
			pid, err := strconv.Atoi(value)
			if err == nil {
				currentPID = pid
			}
		case 'c':
			currentProcess = value
		case 't':
			currentType = value
		case 'P':
			currentProtocol = value
		case 'n':
			// lsof -F 'n' field format: "*:7000", "127.0.0.1:8080", "[::1]:8080"
			lastColon := strings.LastIndex(value, ":")
			if lastColon < 0 {
				continue
			}
			portStr := value[lastColon+1:]
			address := value[:lastColon]

			// Skip wildcard port entries like "*:*"
			if portStr == "*" || portStr == "" {
				continue
			}

			port, err := strconv.Atoi(portStr)
			if err != nil || port == 0 {
				continue
			}

			results = append(results, PortInfo{
				Port:     port,
				PID:      currentPID,
				Process:  currentProcess,
				Protocol: currentProtocol,
				Address:  address,
				Type:     currentType,
			})
		}
	}
	return results
}

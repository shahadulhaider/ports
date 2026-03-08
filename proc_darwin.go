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

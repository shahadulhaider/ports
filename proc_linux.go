//go:build linux

package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	pidRegex  = regexp.MustCompile(`pid=(\d+)`)
	nameRegex = regexp.MustCompile(`"([^"]+)"`)
)

// GetListeningPorts returns all TCP listening ports on Linux using ss.
func GetListeningPorts() ([]PortInfo, error) {
	cmd := exec.Command("ss", "-tlnp")
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("ss error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("ss not found — install iproute2")
	}

	ports := parseSsOutput(string(out), "TCP")
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Port < ports[j].Port
	})
	return ports, nil
}

// GetUDPPorts returns all UDP bound ports on Linux using ss.
func GetUDPPorts() ([]PortInfo, error) {
	cmd := exec.Command("ss", "-ulnp")
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("ss error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("ss not found — install iproute2")
	}

	ports := parseSsOutput(string(out), "UDP")
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Port < ports[j].Port
	})
	return ports, nil
}

// GetConnectionCounts returns a map of local port -> established TCP connection count.
func GetConnectionCounts() (map[int]int, error) {
	cmd := exec.Command("ss", "-tnp", "state", "established")
	out, err := cmd.Output()
	if err != nil {
		// No established connections is not an error
		return map[int]int{}, nil
	}

	counts := map[int]int{}
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if i == 0 { // skip header
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		// Local address is field index 3
		localAddr := fields[3]
		lastColon := strings.LastIndex(localAddr, ":")
		if lastColon < 0 {
			continue
		}
		portStr := localAddr[lastColon+1:]
		port, err := strconv.Atoi(portStr)
		if err != nil || port == 0 {
			continue
		}
		counts[port]++
	}
	return counts, nil
}

func parseSsOutput(output string, protocol string) []PortInfo {
	var results []PortInfo
	lines := strings.Split(output, "\n")

	for i, line := range lines {
		// Skip header line
		if i == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		// ss -tlnp output columns:
		// State Recv-Q Send-Q Local-Address:Port Peer-Address:Port Process
		// We need at least 4 fields (Process column may be missing without root)
		if len(fields) < 4 {
			continue
		}

		// Column 3 (index 3) is Local Address:Port
		localAddr := fields[3]

		// Extract port: split on last ":"
		lastColon := strings.LastIndex(localAddr, ":")
		if lastColon < 0 {
			continue
		}
		portStr := localAddr[lastColon+1:]
		address := localAddr[:lastColon]

		// Detect IPv4 vs IPv6
		addrType := "IPv4"
		if strings.Contains(address, "[") || strings.Contains(address, ":") {
			addrType = "IPv6"
		}

		// Remove brackets from IPv6 addresses like [::]:22 -> ::
		address = strings.Trim(address, "[]")

		port, err := strconv.Atoi(portStr)
		if err != nil || port == 0 {
			continue
		}

		// Parse process info from column 5+ (may be absent without root)
		var pid int
		process := "(unknown)"

		if len(fields) >= 6 {
			processField := strings.Join(fields[5:], " ")
			if m := pidRegex.FindStringSubmatch(processField); len(m) > 1 {
				pid, _ = strconv.Atoi(m[1])
			}
			if m := nameRegex.FindStringSubmatch(processField); len(m) > 1 {
				process = m[1]
			}
		}

		results = append(results, PortInfo{
			Port:     port,
			PID:      pid,
			Process:  process,
			Protocol: protocol,
			Address:  address,
			Type:     addrType,
		})
	}
	return results
}

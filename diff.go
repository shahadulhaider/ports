package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type diffCache struct {
	Timestamp time.Time  `json:"timestamp"`
	Ports     []PortInfo `json:"ports"`
}

func runDiffMode(portFilter int) int {
	current, err := GetListeningPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting ports: %v\n", err)
		return 2
	}

	if portFilter > 0 {
		var filtered []PortInfo
		for _, p := range current {
			if p.Port == portFilter {
				filtered = append(filtered, p)
			}
		}
		current = filtered
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = filepath.Join(os.Getenv("HOME"), ".cache")
	}
	cacheDir = filepath.Join(cacheDir, "ports")
	cachePath := filepath.Join(cacheDir, "last.json")

	var previous *diffCache
	if data, err := os.ReadFile(cachePath); err == nil {
		var cache diffCache
		if json.Unmarshal(data, &cache) == nil {
			previous = &cache
		}
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating cache dir: %v\n", err)
		return 2
	}
	newCache := diffCache{
		Timestamp: time.Now(),
		Ports:     current,
	}
	if data, err := json.Marshal(newCache); err == nil {
		_ = os.WriteFile(cachePath, data, 0644)
	}

	if previous == nil {
		fmt.Printf("First run — baseline saved with %d ports\n", len(current))
		return 0
	}

	type portKey struct {
		Port    int
		PID     int
		Process string
	}
	prevSet := make(map[portKey]PortInfo)
	for _, p := range previous.Ports {
		prevSet[portKey{p.Port, p.PID, p.Process}] = p
	}
	currSet := make(map[portKey]PortInfo)
	for _, p := range current {
		currSet[portKey{p.Port, p.PID, p.Process}] = p
	}

	var newPorts []PortInfo
	for k, p := range currSet {
		if _, exists := prevSet[k]; !exists {
			newPorts = append(newPorts, p)
		}
	}
	var gonePorts []PortInfo
	for k, p := range prevSet {
		if _, exists := currSet[k]; !exists {
			gonePorts = append(gonePorts, p)
		}
	}

	sort.Slice(newPorts, func(i, j int) bool { return newPorts[i].Port < newPorts[j].Port })
	sort.Slice(gonePorts, func(i, j int) bool { return gonePorts[i].Port < gonePorts[j].Port })

	if len(newPorts) == 0 && len(gonePorts) == 0 {
		fmt.Printf("No changes since last run at %s\n", previous.Timestamp.Format("15:04:05"))
		return 0
	}

	fmt.Printf("ports diff (compared to last run at %s):\n", previous.Timestamp.Format("15:04:05"))
	for _, p := range newPorts {
		fmt.Printf("+ %d\t%d\t%s\t%s\n", p.Port, p.PID, p.Process, p.Address)
	}
	for _, p := range gonePorts {
		fmt.Printf("- %d\t%d\t%s\t%s\n", p.Port, p.PID, p.Process, p.Address)
	}
	fmt.Printf("\n%d changes (%d new, %d gone)\n", len(newPorts)+len(gonePorts), len(newPorts), len(gonePorts))
	return 1
}

package main

// PortInfo holds information about a single listening TCP port.
type PortInfo struct {
	Port     int
	PID      int
	Process  string
	Protocol string // "TCP"
	Address  string // e.g. "*:3000" or "127.0.0.1:8080"
	Type     string // "IPv4" or "IPv6"
}

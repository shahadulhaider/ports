package scanner

// PortInfo holds information about a single listening TCP port.
type PortInfo struct {
	Port        int
	PID         int
	Process     string
	Protocol    string // "TCP"
	Address     string // e.g. "*:3000" or "127.0.0.1:8080"
	Type        string // "IPv4" or "IPv6"
	Connections int    // active TCP connection count (-1 for UDP/N/A)
	Service     string // human-readable service name (e.g., "http", "ssh")
	Status      string // "new", "gone", or "" for port change tracking
}

func applyConnectionCounts(ports []PortInfo) {
	counts, err := GetConnectionCounts()
	if err != nil {
		return
	}
	for i := range ports {
		ports[i].Connections = counts[ports[i].Port]
	}
}

// UDP is connectionless, so an established-connection count is meaningless;
// -1 is the sentinel the TUI renders as "N/A".
func markConnectionsNotApplicable(ports []PortInfo) {
	for i := range ports {
		ports[i].Connections = -1
	}
}

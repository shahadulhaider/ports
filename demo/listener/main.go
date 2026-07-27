// Command listener binds a set of TCP/UDP ports and holds them open.
//
// It exists only to give the demo recordings a realistic, reproducible dev
// stack to display. The binary is copied to a service name (postgres,
// redis-server, node, ...) before it is started, so the kernel reports that
// name as the process comm and `ss -tlnp` shows it verbatim.
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	tcp4 := flag.String("tcp4", "", "comma-separated host:port list to listen on over IPv4")
	tcp6 := flag.String("tcp6", "", "comma-separated host:port list to listen on over IPv6")
	udp4 := flag.String("udp4", "", "comma-separated host:port list to bind over IPv4")
	connect := flag.String("connect", "", "comma-separated host:port list to dial and hold open")
	dials := flag.Int("dials", 1, "connections to open per -connect target")
	flag.Parse()

	for _, addr := range split(*tcp4) {
		listenTCP("tcp4", addr)
	}
	for _, addr := range split(*tcp6) {
		listenTCP("tcp6", addr)
	}
	for _, addr := range split(*udp4) {
		bindUDP("udp4", addr)
	}
	for _, addr := range split(*connect) {
		for i := 0; i < *dials; i++ {
			hold(addr)
		}
	}

	// Exit cleanly on SIGTERM so the TUI's `x` keybinding has something real to
	// terminate. SIGKILL (`X`) bypasses this path entirely.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func split(list string) []string {
	if strings.TrimSpace(list) == "" {
		return nil
	}
	parts := strings.Split(list, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func listenTCP(network, addr string) {
	ln, err := net.Listen(network, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen %s %s: %v\n", network, addr, err)
		return
	}
	go func() {
		var open []net.Conn
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			// Retain accepted connections so they stay ESTABLISHED and show up
			// in the TUI's CONNS column.
			open = append(open, conn)
		}
	}()
}

func bindUDP(network, addr string) {
	conn, err := net.ListenPacket(network, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bind %s %s: %v\n", network, addr, err)
		return
	}
	go func() {
		buf := make([]byte, 1024)
		for {
			if _, _, err := conn.ReadFrom(buf); err != nil {
				return
			}
		}
	}()
}

func hold(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial %s: %v\n", addr, err)
		return
	}
	go func() {
		buf := make([]byte, 1)
		for {
			if _, err := conn.Read(buf); err != nil {
				return
			}
		}
	}()
}

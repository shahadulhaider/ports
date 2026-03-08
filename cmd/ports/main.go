package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shahadulhaider/ports/internal/diff"
	"github.com/shahadulhaider/ports/internal/tui"
)

var version = "dev"

func main() {
	var portFlag int
	var diffFlag bool
	var versionFlag bool
	flag.IntVar(&portFlag, "port", 0, "pre-filter to specific port number")
	flag.BoolVar(&diffFlag, "diff", false, "show changes since last run and exit")
	flag.BoolVar(&versionFlag, "version", false, "print version and exit")
	flag.Parse()

	if versionFlag {
		fmt.Printf("ports %s\n", version)
		os.Exit(0)
	}

	if diffFlag {
		os.Exit(diff.RunDiffMode(portFlag))
	}

	p := tea.NewProgram(tui.NewModel(portFlag), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

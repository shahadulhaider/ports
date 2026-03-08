package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var portFlag int
	var diffFlag bool
	flag.IntVar(&portFlag, "port", 0, "pre-filter to specific port number")
	flag.BoolVar(&diffFlag, "diff", false, "show changes since last run and exit")
	flag.Parse()

	if diffFlag {
		os.Exit(runDiffMode(portFlag))
	}

	p := tea.NewProgram(NewModel(portFlag), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

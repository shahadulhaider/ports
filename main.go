package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type minimalModel struct{}

func initialModel() minimalModel { return minimalModel{} }

func (m minimalModel) Init() tea.Cmd { return tea.Quit }

func (m minimalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, tea.Quit }

func (m minimalModel) View() string { return "" }

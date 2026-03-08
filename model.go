package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Message types
type portsMsg []PortInfo
type tickMsg time.Time
type statusClearMsg struct{}

type model struct {
	table         table.Model
	allPorts      []PortInfo
	filteredPorts []PortInfo
	filterInput   textinput.Model
	filtering     bool
	filterText    string
	statusMsg     string
	lastRefresh   time.Time
	width         int
	height        int
	ready         bool
}

func NewModel() model {
	cols := []table.Column{
		{Title: "PORT", Width: 8},
		{Title: "PID", Width: 8},
		{Title: "PROCESS", Width: 25},
		{Title: "PROTO", Width: 6},
		{Title: "ADDRESS", Width: 20},
		{Title: "TYPE", Width: 6},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = selectedRowStyle
	t.SetStyles(s)

	ti := textinput.New()
	ti.Placeholder = "filter..."
	ti.CharLimit = 64

	return model{
		table:       t,
		filterInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchPortsCmd(), tickCmd())
}

func fetchPortsCmd() tea.Cmd {
	return func() tea.Msg {
		ports, err := GetListeningPorts()
		if err != nil {
			return portsMsg(nil)
		}
		return portsMsg(ports)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetHeight(msg.Height - 4)

	case portsMsg:
		m.allPorts = []PortInfo(msg)
		m.filteredPorts = m.allPorts
		m.table.SetRows(portsToRows(m.filteredPorts))
		m.lastRefresh = time.Now()
		m.ready = true

	case tickMsg:
		cmds = append(cmds, fetchPortsCmd(), tickCmd())

	case statusClearMsg:
		m.statusMsg = ""

	case tea.KeyMsg:
		if key.Matches(msg, keys.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, keys.Refresh) {
			cmds = append(cmds, fetchPortsCmd())
		}
		// Delegate navigation to table
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Loading ports..."
	}

	title := titleStyle.Render(fmt.Sprintf(" ⚡ ports (%d listening)", len(m.filteredPorts)))

	tableView := m.table.View()

	var statusText string
	if m.statusMsg != "" {
		statusText = " " + m.statusMsg
	} else {
		refreshTime := m.lastRefresh.Format("15:04:05")
		statusText = fmt.Sprintf(" Last refresh: %s | %d ports | ? help  q quit", refreshTime, len(m.filteredPorts))
	}

	w := m.width
	if w == 0 {
		w = 80
	}
	status := statusBarStyle.Width(w).Render(statusText)

	return lipgloss.JoinVertical(lipgloss.Left, title, tableView, status)
}

func portsToRows(ports []PortInfo) []table.Row {
	rows := make([]table.Row, len(ports))
	for i, p := range ports {
		rows[i] = table.Row{
			strconv.Itoa(p.Port),
			strconv.Itoa(p.PID),
			p.Process,
			p.Protocol,
			p.Address,
			p.Type,
		}
	}
	return rows
}

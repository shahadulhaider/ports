package main

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
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
	showHelp      bool
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
		if m.filterText != "" {
			m.filteredPorts = filterPorts(m.allPorts, m.filterText)
		} else {
			m.filteredPorts = m.allPorts
		}
		m.table.SetRows(portsToRows(m.filteredPorts))
		m.lastRefresh = time.Now()
		m.ready = true

	case tickMsg:
		cmds = append(cmds, fetchPortsCmd(), tickCmd())

	case statusClearMsg:
		m.statusMsg = ""

	case tea.KeyMsg:
		// Filter mode toggle
		if m.filtering {
			// In filter mode: handle Esc to exit, otherwise pass to textinput
			if key.Matches(msg, keys.ClearFilter) {
				m.filtering = false
				m.filterInput.Blur()
				m.filterInput.SetValue("")
				m.filterText = ""
				m.filteredPorts = m.allPorts
				m.table.SetRows(portsToRows(m.filteredPorts))
				return m, tea.Batch(cmds...)
			}
			// Pass keystrokes to filter input
			var tiCmd tea.Cmd
			m.filterInput, tiCmd = m.filterInput.Update(msg)
			newFilter := m.filterInput.Value()
			if newFilter != m.filterText {
				m.filterText = newFilter
				m.filteredPorts = filterPorts(m.allPorts, m.filterText)
				m.table.SetRows(portsToRows(m.filteredPorts))
			}
			cmds = append(cmds, tiCmd)
			return m, tea.Batch(cmds...)
		}

		// Not in filter mode
		if key.Matches(msg, keys.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, keys.Help) {
			m.showHelp = !m.showHelp
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.Filter) {
			m.filtering = true
			m.filterInput.SetValue("")
			m.filterInput.Focus()
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.Refresh) {
			cmds = append(cmds, fetchPortsCmd())
		}
		if key.Matches(msg, keys.Kill) {
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) >= 2 {
				pid, err := strconv.Atoi(selectedRow[1])
				if err == nil && pid > 0 {
					processName := ""
					if len(selectedRow) >= 3 {
						processName = selectedRow[2]
					}
					killErr := syscall.Kill(pid, syscall.SIGTERM)
					if killErr == nil {
						m.statusMsg = fmt.Sprintf("Killed PID %d (%s)", pid, processName)
						cmds = append(cmds, fetchPortsCmd(), clearStatusCmd())
					} else if errors.Is(killErr, syscall.EPERM) {
						m.statusMsg = fmt.Sprintf("Permission denied: cannot kill PID %d (%s)", pid, processName)
						cmds = append(cmds, clearStatusCmd())
					} else if errors.Is(killErr, syscall.ESRCH) {
						m.statusMsg = fmt.Sprintf("Process %d already terminated", pid)
						cmds = append(cmds, fetchPortsCmd(), clearStatusCmd())
					} else {
						m.statusMsg = fmt.Sprintf("Kill failed: %v", killErr)
						cmds = append(cmds, clearStatusCmd())
					}
				}
			}
		}
		if key.Matches(msg, keys.Copy) {
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) >= 3 {
				port := selectedRow[0]
				pid := selectedRow[1]
				process := selectedRow[2]
				address := ""
				if len(selectedRow) >= 5 {
					address = selectedRow[4]
				}
				text := fmt.Sprintf("%s\t%s\t%s\t%s", port, pid, process, address)
				if err := copyToClipboard(text); err != nil {
					m.statusMsg = fmt.Sprintf("Clipboard not available: %v", err)
				} else {
					m.statusMsg = fmt.Sprintf("Copied: port %s (%s)", port, process)
				}
				cmds = append(cmds, clearStatusCmd())
			}
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

	var mainContent string
	if m.showHelp {
		helpContent := `
  ↑/↓ or j/k   Navigate rows
  /             Filter ports
  Esc           Clear filter
  x             Kill process (SIGTERM)
  c             Copy to clipboard
  r             Refresh now
  ?             Toggle this help
  q / Ctrl+C    Quit`

		box := helpBoxStyle.Render("  Help — ports\n" + helpContent + "\n\n  Press ? to close")
		w := m.width
		if w == 0 {
			w = 80
		}
		mainContent = lipgloss.Place(w, m.height-3, lipgloss.Center, lipgloss.Center, box)
	} else if len(m.filteredPorts) == 0 && m.ready {
		if m.filterText != "" {
			mainContent = helpStyle.Render("\n  No matching ports")
		} else {
			mainContent = helpStyle.Render("\n  No listening TCP ports found")
		}
	} else {
		mainContent = m.table.View()
	}

	var statusText string
	if m.statusMsg != "" {
		statusText = " " + m.statusMsg
	} else if m.filtering {
		statusText = fmt.Sprintf(" Filter: %s (%d results) | Esc to cancel", m.filterInput.View(), len(m.filteredPorts))
	} else if m.filterText != "" {
		statusText = fmt.Sprintf(" Filtered: %q (%d results) | Esc to clear | ? help  q quit", m.filterText, len(m.filteredPorts))
	} else {
		refreshTime := m.lastRefresh.Format("15:04:05")
		statusText = fmt.Sprintf(" Last refresh: %s | %d ports | ? help  q quit", refreshTime, len(m.filteredPorts))
	}

	w := m.width
	if w == 0 {
		w = 80
	}
	status := statusBarStyle.Width(w).Render(statusText)

	return lipgloss.JoinVertical(lipgloss.Left, title, mainContent, status)
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

func clearStatusCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return statusClearMsg{}
	})
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func filterPorts(ports []PortInfo, query string) []PortInfo {
	if query == "" {
		return ports
	}
	q := strings.ToLower(query)
	var result []PortInfo
	for _, p := range ports {
		if strings.Contains(strings.ToLower(strconv.Itoa(p.Port)), q) ||
			strings.Contains(strings.ToLower(p.Process), q) ||
			strings.Contains(strings.ToLower(p.Address), q) {
			result = append(result, p)
		}
	}
	return result
}

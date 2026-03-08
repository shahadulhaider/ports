package main

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"sort"
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
type statusClearMsg struct{ gen int }

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
	prevPorts     []PortInfo // previous refresh snapshot for change tracking
	sortMode      int        // 0=Port↑, 1=Port↓, 2=PID, 3=Process A-Z
	dedupEnabled  bool       // merge IPv4+IPv6 rows by (Port, PID)
	protoFilter   int        // 0=TCP only, 1=UDP only, 2=Both
	statusCounter int        // generation counter to prevent stale status clears
}

func NewModel(initialPort int) model {
	cols := []table.Column{
		{Title: "STATUS", Width: 3},
		{Title: "PORT", Width: 8},
		{Title: "PID", Width: 8},
		{Title: "PROCESS", Width: 25},
		{Title: "PROTO", Width: 6},
		{Title: "ADDRESS", Width: 20},
		{Title: "TYPE", Width: 6},
		{Title: "SERVICE", Width: 10},
		{Title: "CONNS", Width: 6},
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

	m := model{
		table:       t,
		filterInput: ti,
	}

	if initialPort > 0 {
		portStr := strconv.Itoa(initialPort)
		m.filterText = portStr
		m.filterInput.SetValue(portStr)
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchPortsCmd(0), tickCmd())
}

func fetchPortsCmd(proto int) tea.Cmd {
	return func() tea.Msg {
		var ports []PortInfo
		var err error
		switch proto {
		case 1: // UDP only
			ports, err = GetUDPPorts()
		case 2: // Both TCP and UDP
			tcpPorts, tcpErr := GetListeningPorts()
			udpPorts, udpErr := GetUDPPorts()
			if tcpErr != nil {
				err = tcpErr
			} else if udpErr != nil {
				err = udpErr
			} else {
				ports = append(tcpPorts, udpPorts...)
			}
		default: // TCP only (0)
			ports, err = GetListeningPorts()
		}
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

func (m model) clearStatusCmd() tea.Cmd {
	gen := m.statusCounter
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return statusClearMsg{gen: gen}
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
		newPorts := []PortInfo(msg)

		// Port change tracking (skip on first refresh)
		if m.prevPorts != nil {
			prevSet := make(map[string]bool)
			for _, p := range m.prevPorts {
				prevSet[fmt.Sprintf("%d:%d", p.Port, p.PID)] = true
			}
			currSet := make(map[string]bool)
			for i := range newPorts {
				k := fmt.Sprintf("%d:%d", newPorts[i].Port, newPorts[i].PID)
				currSet[k] = true
				if !prevSet[k] {
					newPorts[i].Status = "new"
				}
			}
			// Add phantom "gone" entries for disappeared ports (one cycle only)
			for _, p := range m.prevPorts {
				k := fmt.Sprintf("%d:%d", p.Port, p.PID)
				if !currSet[k] && p.Status != "gone" {
					gone := p
					gone.Status = "gone"
					newPorts = append(newPorts, gone)
				}
			}
		}

		// Store original (without phantoms) for next cycle's diff
		m.prevPorts = []PortInfo(msg)
		m.allPorts = newPorts

		// Apply sort, dedup, filter
		m.filteredPorts = m.applyDisplayPipeline(m.allPorts)
		m.table.SetRows(portsToRows(m.filteredPorts))
		m.lastRefresh = time.Now()
		m.ready = true

	case tickMsg:
		cmds = append(cmds, fetchPortsCmd(m.protoFilter), tickCmd())

	case statusClearMsg:
		if msg.gen == m.statusCounter {
			m.statusMsg = ""
		}

	case tea.KeyMsg:
		// Filter mode: handle Esc to exit, otherwise pass to textinput
		if m.filtering {
			if key.Matches(msg, keys.ClearFilter) {
				m.filtering = false
				m.filterInput.Blur()
				m.filterInput.SetValue("")
				m.filterText = ""
				m.filteredPorts = m.applyDisplayPipeline(m.allPorts)
				m.table.SetRows(portsToRows(m.filteredPorts))
				return m, tea.Batch(cmds...)
			}
			var tiCmd tea.Cmd
			m.filterInput, tiCmd = m.filterInput.Update(msg)
			newFilter := m.filterInput.Value()
			if newFilter != m.filterText {
				m.filterText = newFilter
				m.filteredPorts = m.applyDisplayPipeline(m.allPorts)
				m.table.SetRows(portsToRows(m.filteredPorts))
			}
			cmds = append(cmds, tiCmd)
			return m, tea.Batch(cmds...)
		}

		// Normal mode key handlers
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
			cmds = append(cmds, fetchPortsCmd(m.protoFilter))
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.Sort) {
			m.sortMode = (m.sortMode + 1) % 4
			sortNames := []string{"Port ↑", "Port ↓", "PID", "Process"}
			m.statusCounter++
			m.statusMsg = "Sort: " + sortNames[m.sortMode]
			m.filteredPorts = m.applyDisplayPipeline(m.allPorts)
			m.table.SetRows(portsToRows(m.filteredPorts))
			cmds = append(cmds, m.clearStatusCmd())
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.ToggleDedup) {
			m.dedupEnabled = !m.dedupEnabled
			m.statusCounter++
			if m.dedupEnabled {
				m.statusMsg = "Dedup: ON"
			} else {
				m.statusMsg = "Dedup: OFF"
			}
			m.filteredPorts = m.applyDisplayPipeline(m.allPorts)
			m.table.SetRows(portsToRows(m.filteredPorts))
			cmds = append(cmds, m.clearStatusCmd())
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.ToggleUDP) {
			m.protoFilter = (m.protoFilter + 1) % 3
			protoNames := []string{"TCP", "UDP", "TCP+UDP"}
			m.statusCounter++
			m.statusMsg = "Protocol: " + protoNames[m.protoFilter]
			cmds = append(cmds, fetchPortsCmd(m.protoFilter), m.clearStatusCmd())
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.Kill) {
			selectedRow := m.table.SelectedRow()
			// Column indices: [0]=STATUS [1]=PORT [2]=PID [3]=PROCESS
			if len(selectedRow) >= 3 {
				pid, err := strconv.Atoi(selectedRow[2])
				if err == nil && pid > 0 {
					processName := ""
					if len(selectedRow) >= 4 {
						processName = selectedRow[3]
					}
					killErr := syscall.Kill(pid, syscall.SIGTERM)
					m.statusCounter++
					if killErr == nil {
						m.statusMsg = fmt.Sprintf("Killed PID %d (%s)", pid, processName)
						cmds = append(cmds, fetchPortsCmd(m.protoFilter), m.clearStatusCmd())
					} else if errors.Is(killErr, syscall.EPERM) {
						m.statusMsg = fmt.Sprintf("Permission denied: cannot kill PID %d (%s)", pid, processName)
						cmds = append(cmds, m.clearStatusCmd())
					} else if errors.Is(killErr, syscall.ESRCH) {
						m.statusMsg = fmt.Sprintf("Process %d already terminated", pid)
						cmds = append(cmds, fetchPortsCmd(m.protoFilter), m.clearStatusCmd())
					} else {
						m.statusMsg = fmt.Sprintf("Kill failed: %v", killErr)
						cmds = append(cmds, m.clearStatusCmd())
					}
				}
			}
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.ForceKill) {
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) >= 3 {
				pid, err := strconv.Atoi(selectedRow[2])
				if err == nil && pid > 0 {
					processName := ""
					if len(selectedRow) >= 4 {
						processName = selectedRow[3]
					}
					killErr := syscall.Kill(pid, syscall.SIGKILL)
					m.statusCounter++
					if killErr == nil {
						m.statusMsg = fmt.Sprintf("Force killed PID %d (%s)", pid, processName)
						cmds = append(cmds, fetchPortsCmd(m.protoFilter), m.clearStatusCmd())
					} else if errors.Is(killErr, syscall.EPERM) {
						m.statusMsg = fmt.Sprintf("Permission denied: cannot kill PID %d (%s)", pid, processName)
						cmds = append(cmds, m.clearStatusCmd())
					} else {
						m.statusMsg = fmt.Sprintf("Force kill failed: %v", killErr)
						cmds = append(cmds, m.clearStatusCmd())
					}
				}
			}
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.Open) {
			selectedRow := m.table.SelectedRow()
			if len(selectedRow) >= 2 {
				port := selectedRow[1]
				url := "http://localhost:" + port
				m.statusCounter++
				if err := openURL(url); err != nil {
					m.statusMsg = fmt.Sprintf("Cannot open browser: %v", err)
				} else {
					m.statusMsg = fmt.Sprintf("Opened %s", url)
				}
				cmds = append(cmds, m.clearStatusCmd())
			}
			return m, tea.Batch(cmds...)
		}
		if key.Matches(msg, keys.Copy) {
			selectedRow := m.table.SelectedRow()
			// Column indices: [0]=STATUS [1]=PORT [2]=PID [3]=PROCESS [4]=PROTO [5]=ADDRESS
			if len(selectedRow) >= 4 {
				port := selectedRow[1]
				pid := selectedRow[2]
				process := selectedRow[3]
				address := ""
				if len(selectedRow) >= 6 {
					address = selectedRow[5]
				}
				text := fmt.Sprintf("%s\t%s\t%s\t%s", port, pid, process, address)
				m.statusCounter++
				if err := copyToClipboard(text); err != nil {
					m.statusMsg = fmt.Sprintf("Clipboard not available: %v", err)
				} else {
					m.statusMsg = fmt.Sprintf("Copied: port %s (%s)", port, process)
				}
				cmds = append(cmds, m.clearStatusCmd())
			}
			return m, tea.Batch(cmds...)
		}
		// Delegate navigation to table
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// applyDisplayPipeline applies sort → dedup → filter to produce the display list.
func (m model) applyDisplayPipeline(ports []PortInfo) []PortInfo {
	sorted := sortPorts(ports, m.sortMode)
	var deduped []PortInfo
	if m.dedupEnabled {
		deduped = dedupPorts(sorted)
	} else {
		deduped = sorted
	}
	if m.filterText != "" {
		return filterPorts(deduped, m.filterText)
	}
	return deduped
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
  X             Force kill (SIGKILL)
  o             Open in browser
  c             Copy to clipboard
  s             Cycle sort mode
  t             Toggle TCP/UDP/Both
  m             Toggle merge IPv4+IPv6
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
			mainContent = helpStyle.Render("\n  No listening ports found")
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
		sortNames := []string{"Port ↑", "Port ↓", "PID", "Process"}
		statusText = fmt.Sprintf(" Last refresh: %s | %d ports", refreshTime, len(m.filteredPorts))
		if m.protoFilter == 1 {
			statusText += " | [UDP]"
		} else if m.protoFilter == 2 {
			statusText += " | [TCP+UDP]"
		}
		if m.dedupEnabled {
			statusText += " [dedup]"
		}
		statusText += " | Sort: " + sortNames[m.sortMode]
		statusText += " | ? help  q quit"
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
		// Plain Unicode markers — do NOT use lipgloss.Render() here.
		// bubbles/table uses runewidth.Truncate() which counts ANSI escape codes
		// as visible characters, breaking column widths.
		status := ""
		switch p.Status {
		case "new":
			status = "●"
		case "gone":
			status = "○"
		}
		var conns string
		if p.Connections < 0 {
			conns = "N/A"
		} else {
			conns = strconv.Itoa(p.Connections)
		}
		rows[i] = table.Row{
			status,
			strconv.Itoa(p.Port),
			strconv.Itoa(p.PID),
			p.Process,
			p.Protocol,
			p.Address,
			p.Type,
			serviceName(p.Port),
			conns,
		}
	}
	return rows
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

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("open not supported on %s", runtime.GOOS)
	}
	return cmd.Start()
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
			strings.Contains(strings.ToLower(p.Address), q) ||
			strings.Contains(strings.ToLower(serviceName(p.Port)), q) {
			result = append(result, p)
		}
	}
	return result
}

func sortPorts(ports []PortInfo, mode int) []PortInfo {
	sorted := make([]PortInfo, len(ports))
	copy(sorted, ports)
	switch mode {
	case 1: // Port descending
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Port > sorted[j].Port
		})
	case 2: // PID ascending
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].PID < sorted[j].PID
		})
	case 3: // Process A-Z
		sort.Slice(sorted, func(i, j int) bool {
			return strings.ToLower(sorted[i].Process) < strings.ToLower(sorted[j].Process)
		})
	default: // Port ascending (0)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Port < sorted[j].Port
		})
	}
	return sorted
}

func dedupPorts(ports []PortInfo) []PortInfo {
	type dedupKey struct{ Port, PID int }
	seen := make(map[dedupKey]int)
	var result []PortInfo

	for _, p := range ports {
		k := dedupKey{p.Port, p.PID}
		if idx, exists := seen[k]; exists {
			// Merge: update Type to "4+6"
			result[idx].Type = "4+6"
		} else {
			seen[k] = len(result)
			result = append(result, p)
		}
	}
	return result
}

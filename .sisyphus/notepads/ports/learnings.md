# Learnings — ports project

## [Task 8 Complete — help overlay + polish]
- showHelp bool field added to model struct
- ? key toggles showHelp
- lipgloss.Place(width, height, Center, Center, box) centers help overlay
- helpBoxStyle from styles.go used for the bordered box
- Empty state: "No listening TCP ports found" when allPorts empty
- Empty filter state: "No matching ports" when filter has no results
- All 8 tasks complete — ports TUI is fully functional

## [2026-03-08] Session Start
- Go 1.25.6 darwin/arm64 confirmed available
- Working directory: /Users/msh/code/pp/passed
- Git initialized, plan committed

## Critical Constraints (from Metis review)
- Bubbletea v1 ONLY: `github.com/charmbracelet/bubbletea` — NOT `charm.land/bubbletea/v2`
- lsof parsing: use `-F pcfnPt` machine-readable format, NOT human-readable table
- Kill key: `x` (NOT `k` — conflicts with vim-style table navigation)
- Clipboard: exec `pbcopy`/`xclip` — NOT `golang.design/x/clipboard` (requires CGO)
- Build tags: `//go:build darwin` and `//go:build linux` (new style, not `// +build`)
- No interfaces for data layer — two files with build tags, same function signature
- No unit tests — tmux-based QA only
- Module path: `ports` (simple, local project)

## [Task 1 Complete] Project Scaffolding
- go.mod module path: ports, Go 1.25
- Bubbletea version: v1.3.10 (github.com/charmbracelet/bubbletea)
- Lipgloss version: v1.1.0
- bubbles version: v1.0.0
- All files created: main.go, proc.go, styles.go, keys.go
- Build passes: `go build -o ports .` ✓
- Cross-compile passes: `GOOS=linux GOARCH=amd64 go build -o /dev/null .` ✓
- No charm.land imports verified ✓
- Minimal Bubbletea program initializes and quits cleanly
- Key bindings defined: up/down (vim+arrow), kill (x), filter (/), clear filter (esc), copy (c), refresh (r), help (?), quit (q/ctrl+c)
- Styles defined: title (purple #99), header (bold), selected row (bg #57), status bar (bg #236), help text (dimmed #241), help box (rounded border)

## [Task 3 Complete — Linux parser]
- proc_linux.go created with //go:build linux tag
- ss -tlnp output: State Recv-Q Send-Q Local:Port Peer:Port Process
- Process column may be absent without root — defaults to "(unknown)"
- IPv6 addresses wrapped in brackets [::]:22 — strip brackets for Address field
- Cross-compile verified: GOOS=linux GOARCH=amd64 go build succeeds
- Function signature matches darwin: func GetListeningPorts() ([]PortInfo, error)

## [Task 2 Complete — macOS parser]
- lsof -F format confirmed: p=PID, c=command, t=type, P=protocol, n=address:port
- Full process names returned (not truncated) — e.g. "Discord Helper (Renderer)", "ControlCenter"
- IPv4+IPv6 appear as separate entries for same port — expected
- Port 0 and wildcard entries filtered out
- lsof exit code 1 with empty output = no listeners (not an error)

## [Task 5 Complete — filter feature]
- Filter mode: / to enter, Esc to exit
- filterInput.Update(msg) called when filtering=true to handle keystrokes
- filterPorts() does case-insensitive substring match on port, process, address
- portsMsg handler re-applies filter when new data arrives
- Empty results: shows "No matching ports" via helpStyle
- Status bar shows filter state: active filter text + result count

## [Task 6 Complete — kill feature]
- Kill key: x (NOT k — k is vim nav)
- syscall.Kill(pid, syscall.SIGTERM) for killing
- errors.Is(err, syscall.EPERM) for permission denied
- errors.Is(err, syscall.ESRCH) for already-dead process
- clearStatusCmd() uses tea.Tick(3s) to auto-clear status message
- table.SelectedRow() returns current row as []string
- selectedRow[1] is PID (column index 1), selectedRow[2] is process name

## [Task 7 Complete — clipboard copy]
- Copy key: c (keys.Copy)
- copyToClipboard() uses exec pbcopy (macOS) or xclip -selection clipboard (Linux)
- cmd.Stdin = strings.NewReader(text) pipes text to clipboard command
- -selection clipboard flag critical on Linux (targets Ctrl+V clipboard, not primary selection)
- Graceful error if clipboard command not found
- clearStatusCmd() reused from Task 6 for auto-clear

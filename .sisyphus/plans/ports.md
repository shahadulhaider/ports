# ports — Terminal Port Explorer TUI

## TL;DR

> **Quick Summary**: Build `ports`, a Go TUI tool using Bubbletea that shows all listening TCP ports with their associated processes, and lets you kill, filter, copy, and auto-refresh — all from your terminal.
> 
> **Deliverables**:
> - Single Go binary (`ports`) for macOS and Linux
> - Table view of all TCP listening ports with process info
> - Actions: kill process, filter/search, copy to clipboard, auto-refresh
> - Full keyboard-driven interface with help overlay
> 
> **Estimated Effort**: Weekend (2-3 days)
> **Parallel Execution**: YES — 4 waves, up to 3 concurrent tasks
> **Critical Path**: Skeleton → macOS Parser → TUI Model → Kill/Filter/Copy → Polish

---

## Context

### Original Request
Build a port explorer TUI — a beautiful terminal tool that answers "what's running on port X?" with one command. Native Go, no Electron, fun weekend project.

### Interview Summary
**Key Discussions**:
- **Platforms**: macOS + Linux (no Windows)
- **Name**: `ports`
- **Stack**: Go + Bubbletea (v1)
- **Features**: List listening ports, kill process, filter/search, copy to clipboard, auto-refresh, port details (protocol, address, type)
- **Scope**: Weekend project, solo dev, fun + useful
- **Tests**: No unit tests — agent-executed QA via tmux only

### Metis Review
**Identified Gaps** (addressed):
- **Bubbletea v1 vs v2**: v2 is still beta — locked to v1. All real-world apps use v1.
- **lsof parsing strategy**: Human-readable output truncates command names to 9 chars and has variable-width columns. Using `-F pcfnPt` machine-readable format instead.
- **Kill key conflict**: `k` conflicts with vim-style table navigation. Using `x` instead.
- **IPv4+IPv6 duplication**: Same port shows twice (one per protocol). Showing both rows for MVP — simplest approach.
- **Clipboard CGO trap**: `golang.design/x/clipboard` requires CGO. Using `exec pbcopy`/`xclip` instead.
- **Kill safety**: Handle EPERM (permission denied) and ESRCH (process already dead) gracefully.

---

## Work Objectives

### Core Objective
Ship a polished, keyboard-driven TUI that instantly shows what's listening on every TCP port, with the ability to kill, filter, and copy — replacing the `lsof -i | grep LISTEN` workflow.

### Concrete Deliverables
- `ports` binary that compiles for macOS (darwin/arm64, darwin/amd64) and Linux (amd64)
- Table showing: Port, PID, Process Name, Protocol, Address, Type (IPv4/IPv6)
- Kill process with `x` key, filter with `/`, copy with `c`, auto-refresh every 2s
- Help overlay with `?` key
- Graceful error handling for all edge cases

### Definition of Done
- [ ] `go build -o ports .` succeeds on macOS
- [ ] `GOOS=linux GOARCH=amd64 go build -o /dev/null .` succeeds (cross-compile)
- [ ] Launch in terminal shows listening ports matching `lsof` output
- [ ] Can kill a dummy process via the TUI
- [ ] Filter narrows results, Esc clears
- [ ] Copy puts port info on clipboard
- [ ] Auto-refresh picks up new listeners without manual action
- [ ] `go vet ./...` returns clean

### Must Have
- Bubbletea v1 (`github.com/charmbracelet/bubbletea`) — NOT v2
- `lsof -F pcfnPt` machine-readable parsing on macOS — NOT human-readable table output
- Go build tags (`//go:build darwin` / `//go:build linux`) for platform-specific code
- Exec-based clipboard (`pbcopy`/`xclip`) — NOT `golang.design/x/clipboard`
- Graceful EPERM and ESRCH handling on kill
- Handle `tea.WindowSizeMsg` for terminal resize
- Handle empty state ("No listening ports found")
- Status bar with filter indicator, last refresh time, help hints

### Must NOT Have (Guardrails)
- No config file (YAML/TOML/JSON) — hardcode defaults
- No theme system or color customization — one good color scheme
- No remote host support (SSH)
- No historical data or logging to file
- No connection tracking (ESTABLISHED connections) — LISTEN only
- No bandwidth/traffic monitoring
- No Docker container awareness
- No service name resolution (port 80 → "HTTP")
- No multi-column sorting — sort by port number ascending only
- No mouse support
- No plugin system
- No Bubbletea v2
- No `golang.design/x/clipboard` (requires CGO)
- No custom keybinding configuration
- No tabs or multiple views
- No confirmation dialog for kill (weekend scope)
- No unit tests
- No interfaces for the data layer — two files with build tags, same function signature
- No `--json` flag or non-TUI output modes
- No signal selection (TERM vs KILL vs HUP) — SIGTERM only

---

## Verification Strategy

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed. No exceptions.

### Test Decision
- **Infrastructure exists**: NO (greenfield project)
- **Automated tests**: NONE (weekend project, user explicitly excluded)
- **Framework**: N/A
- **QA Method**: Agent-executed via tmux (interactive_bash) for TUI verification, Bash for build/compile checks

### QA Policy
Every task MUST include agent-executed QA scenarios.
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **TUI verification**: Use `interactive_bash` (tmux) — launch binary, send keystrokes, capture pane content, assert expected strings
- **Build verification**: Use Bash — `go build`, `go vet`, cross-compile checks
- **Process operations**: Use Bash — start dummy processes with `nc -l`, verify kill via `kill -0`

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately — foundation):
└── Task 1: Project scaffolding + Go module + types + styles + keys [quick]

Wave 1b (After Task 1 — parsers in parallel):
├── Task 2: macOS data layer (lsof -F parser) [unspecified-low]
└── Task 3: Linux data layer (ss parser) [unspecified-low]

Wave 2 (After Wave 1b — core TUI):
└── Task 4: Bubbletea model + table rendering + auto-refresh [deep]

Wave 3 (After Wave 2 — features in parallel):
├── Task 5: Filter/search feature [unspecified-low]
├── Task 6: Kill process feature [unspecified-low]
└── Task 7: Clipboard copy feature [quick]

Wave 4 (After Wave 3 — polish):
└── Task 8: Help overlay + polish + comprehensive QA [unspecified-low]

Wave FINAL (After ALL tasks — independent review, 4 parallel):
├── Task F1: Plan compliance audit [oracle]
├── Task F2: Code quality review [unspecified-high]
├── Task F3: Real manual QA [unspecified-high]
└── Task F4: Scope fidelity check [deep]

Critical Path: Task 1 → Task 2 → Task 4 → Task 6 → Task 8 → F1-F4
Parallel Speedup: ~40% faster than sequential
Max Concurrent: 3 (Waves 1b and 3)
```

### Dependency Matrix

| Task | Depends On | Blocks | Wave |
|------|-----------|--------|------|
| 1 | — | 2, 3 | 1 |
| 2 | 1 | 4 | 1b |
| 3 | 1 | 4 | 1b |
| 4 | 1, 2, 3 | 5, 6, 7 | 2 |
| 5 | 4 | 8 | 3 |
| 6 | 4 | 8 | 3 |
| 7 | 4 | 8 | 3 |
| 8 | 5, 6, 7 | F1-F4 | 4 |
| F1-F4 | 8 | — | FINAL |

### Agent Dispatch Summary

| Wave | Tasks | Categories |
|------|-------|------------|
| 1 | 1 | T1 → `quick` |
| 1b | 2 | T2 → `unspecified-low`, T3 → `unspecified-low` |
| 2 | 1 | T4 → `deep` |
| 3 | 3 | T5 → `unspecified-low`, T6 → `unspecified-low`, T7 → `quick` |
| 4 | 1 | T8 → `unspecified-low` |
| FINAL | 4 | F1 → `oracle`, F2 → `unspecified-high`, F3 → `unspecified-high`, F4 → `deep` |

---

## TODOs

- [ ] 1. Project Scaffolding + Go Module Setup

  **What to do**:
  - Initialize Go module: `go mod init ports`
  - Create `main.go` with minimal `tea.NewProgram` that prints "ports starting..." and quits — just enough to verify the build chain works
  - Create `proc.go` with the shared `PortInfo` struct:
    ```go
    type PortInfo struct {
        Port     int
        PID      int
        Process  string
        Protocol string    // "TCP"
        Address  string    // "127.0.0.1:8080" or "*:3000"
        Type     string    // "IPv4" or "IPv6"
    }
    ```
  - Create `styles.go` with Lipgloss styles: header row style, selected row style, status bar style (dark background), help text style (dimmed), title style (bold + colored)
  - Create `keys.go` with all keybindings using `charmbracelet/bubbles/key`:
    - `Quit`: `q`, `ctrl+c`
    - `Kill`: `x` (NOT `k` — conflicts with vim nav)
    - `Filter`: `/`
    - `ClearFilter`: `esc`
    - `Copy`: `c`
    - `Refresh`: `r`
    - `Help`: `?`
  - Run `go mod tidy` to fetch deps
  - Verify `go build -o ports .` succeeds

  **Must NOT do**:
  - Do NOT add any config file loading
  - Do NOT add interfaces — just the struct and key bindings
  - Do NOT use Bubbletea v2 (`charm.land/bubbletea/v2`) — use v1 (`github.com/charmbracelet/bubbletea`)
  - Do NOT add mouse support to the program options

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Simple scaffolding — create files, define structs, run go mod init
  - **Skills**: []
    - No special skills needed for Go module init

  **Parallelization**:
  - **Can Run In Parallel**: NO (this is the foundation)
  - **Parallel Group**: Wave 1 (solo)
  - **Blocks**: Tasks 2, 3
  - **Blocked By**: None (can start immediately)

  **References**:

  **Pattern References** (existing code to follow):
  - None — greenfield project, no existing code

  **External References** (libraries and frameworks):
  - `github.com/charmbracelet/bubbletea` — Bubbletea v1 TUI framework. Use `tea.NewProgram(model)` pattern
  - `github.com/charmbracelet/lipgloss` — Style definitions. Use `lipgloss.NewStyle().Foreground(lipgloss.Color("205"))` pattern
  - `github.com/charmbracelet/bubbles/key` — Key binding definitions. Use `key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit"))` pattern
  - Reference project `sammcj/gollama` for idiomatic Bubbletea v1 project structure

  **WHY Each Reference Matters**:
  - Bubbletea v1 import path is critical — v2 has different paths and API. Getting this wrong in go.mod poisons every subsequent task
  - Lipgloss is needed for styles.go — the style system used by the table and status bar
  - bubbles/key defines the key binding format consumed by Update() in later tasks

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Build succeeds
    Tool: Bash
    Preconditions: Go 1.25+ installed, empty project directory
    Steps:
      1. Run `go build -o ports .`
      2. Check exit code is 0
      3. Verify binary exists: `test -f ports && echo "BINARY EXISTS"`
    Expected Result: Exit code 0, "BINARY EXISTS" printed
    Failure Indicators: Compilation errors, missing imports, binary not created
    Evidence: .sisyphus/evidence/task-1-build.txt

  Scenario: Binary runs without panic
    Tool: Bash
    Preconditions: Binary built successfully
    Steps:
      1. Run `./ports` (it should print something and exit since model is minimal)
      2. Check exit code is 0
    Expected Result: Clean exit, no panic, no stack trace
    Failure Indicators: Panic output, non-zero exit code
    Evidence: .sisyphus/evidence/task-1-run.txt

  Scenario: Cross-compile for Linux
    Tool: Bash
    Preconditions: Project compiles on macOS
    Steps:
      1. Run `GOOS=linux GOARCH=amd64 go build -o /dev/null .`
      2. Check exit code is 0
    Expected Result: Exit code 0 (compiles for Linux without errors)
    Failure Indicators: Build tags prevent compilation, missing Linux-specific file
    Evidence: .sisyphus/evidence/task-1-crosscompile.txt

  Scenario: No v2 imports present
    Tool: Bash
    Preconditions: go.mod and go.sum exist
    Steps:
      1. Run `grep -r "charm.land" . --include="*.go" || echo "NO V2 IMPORTS"`
      2. Run `grep -r "charm.land" go.mod go.sum || echo "NO V2 IN GOMOD"`
    Expected Result: "NO V2 IMPORTS" and "NO V2 IN GOMOD" printed
    Failure Indicators: Any match for "charm.land" (v2 import path)
    Evidence: .sisyphus/evidence/task-1-no-v2.txt
  ```

  **Evidence to Capture:**
  - [ ] task-1-build.txt — go build output
  - [ ] task-1-run.txt — binary execution output
  - [ ] task-1-crosscompile.txt — Linux cross-compile output
  - [ ] task-1-no-v2.txt — v2 import check

  **Commit**: YES
  - Message: `feat: initialize project skeleton with Go module, types, styles, and keybindings`
  - Files: `go.mod, go.sum, main.go, proc.go, styles.go, keys.go`
  - Pre-commit: `go build -o ports . && go vet ./...`

- [ ] 2. macOS Data Layer (lsof -F Parser)

  **What to do**:
  - Create `proc_darwin.go` with build tag `//go:build darwin`
  - Implement `func GetListeningPorts() ([]PortInfo, error)` that:
    1. Execs `lsof -iTCP -P -n -sTCP:LISTEN -F pcfnPt`
    2. Parses the field-per-line output format:
       - Lines starting with `p` → PID (e.g., `p20133` → PID 20133)
       - Lines starting with `c` → Command/process name (e.g., `cControlCenter`)
       - Lines starting with `f` → File descriptor (track but don't store)
       - Lines starting with `t` → Type (e.g., `tIPv4`, `tIPv6`)
       - Lines starting with `P` (uppercase) → Protocol (e.g., `PTCP`)
       - Lines starting with `n` → Address:Port (e.g., `n*:7000` → Address `*`, Port 7000)
    3. For each `n` line, extract port by splitting on `:` and taking the LAST segment (handles IPv6 addresses with colons)
    4. Builds `PortInfo` struct for each file descriptor entry
    5. Returns slice sorted by port number ascending
  - Handle edge cases:
    - `lsof` not found → return error: "lsof not found — is this macOS?"
    - `lsof` returns no output → return empty slice, nil error
    - Malformed lines → skip with no error (defensive parsing)
    - Port 0 entries → filter out
    - `n*:*` entries → filter out

  **Must NOT do**:
  - Do NOT parse human-readable `lsof` output (truncates command names to 9 chars, fragile column parsing)
  - Do NOT create an interface for the data layer — just the exported function with same signature as Linux version
  - Do NOT add caching or rate limiting

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Straightforward parsing logic, one file, well-defined input/output
  - **Skills**: []
    - No special skills needed — standard Go exec + string parsing

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 3)
  - **Parallel Group**: Wave 1b (with Task 3)
  - **Blocks**: Task 4
  - **Blocked By**: Task 1 (needs PortInfo struct from proc.go)

  **References**:

  **Pattern References**:
  - `proc.go:PortInfo` — The struct to populate from parsed lsof output

  **External References**:
  - `lsof -F` format documentation: fields are single-character prefixed lines. `p`=PID, `c`=command, `f`=fd, `t`=type, `P`=protocol, `n`=name. Each process block starts with `p`, followed by `c`, then one or more fd blocks (`f`, `t`, `P`, `n`)
  - Sample output on this machine:
    ```
    p20133
    cControlCenter
    f9
    tIPv4
    PTCP
    n*:7000
    f10
    tIPv6
    PTCP
    n*:7000
    ```

  **WHY Each Reference Matters**:
  - PortInfo struct defines the exact fields to populate — executor must match field names and types
  - lsof -F format is critical context — without this, an executor might try to parse the human-readable table (which is a known trap with truncated command names)

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Parser returns real ports from this machine
    Tool: Bash
    Preconditions: Task 1 complete, macOS system with lsof
    Steps:
      1. Add temporary test code in main.go:
         ports, err := GetListeningPorts()
         fmt.Printf("Found %d ports\n", len(ports))
         for _, p := range ports { fmt.Printf("  :%d %s (%s, %s, %s)\n", p.Port, p.Process, p.Protocol, p.Address, p.Type) }
      2. Run `go build -o ports . && ./ports`
      3. Compare output against `lsof -iTCP -P -n -sTCP:LISTEN` to verify ports match
    Expected Result: Non-empty list of ports that matches real lsof output (same PIDs, port numbers, process names — but full names not truncated)
    Failure Indicators: Empty output, missing ports, truncated process names, parse errors
    Evidence: .sisyphus/evidence/task-2-parser-output.txt

  Scenario: Full process names (not truncated)
    Tool: Bash
    Preconditions: Parser working
    Steps:
      1. Run parser output from above
      2. Check that process names longer than 9 chars are complete (e.g., "ControlCenter" not "ControlCe")
    Expected Result: Full process names visible in output
    Failure Indicators: Names truncated to 9 characters (means human-readable format was used instead of -F)
    Evidence: .sisyphus/evidence/task-2-full-names.txt

  Scenario: Handles lsof not found
    Tool: Bash
    Preconditions: Parser implemented
    Steps:
      1. Temporarily rename lsof path in code to `/nonexistent/lsof` or test error path
      2. Verify function returns meaningful error, not a panic
    Expected Result: Error message returned, no panic
    Failure Indicators: Panic, nil pointer dereference, empty error
    Evidence: .sisyphus/evidence/task-2-error-handling.txt
  ```

  **Evidence to Capture:**
  - [ ] task-2-parser-output.txt — parsed port list vs raw lsof
  - [ ] task-2-full-names.txt — process name length verification
  - [ ] task-2-error-handling.txt — error case output

  **Commit**: YES
  - Message: `feat: add macOS listening port discovery via lsof`
  - Files: `proc_darwin.go`
  - Pre-commit: `go build -o ports . && go vet ./...`

- [ ] 3. Linux Data Layer (ss Parser)

  **What to do**:
  - Create `proc_linux.go` with build tag `//go:build linux`
  - Implement `func GetListeningPorts() ([]PortInfo, error)` that:
    1. Execs `ss -tlnp`
    2. Skips the header line
    3. For each data line, splits by whitespace:
       - Column 4 (Local Address:Port): split on last `:` to get address and port
       - Column 6 (Process): parse with regex `pid=(\d+)` for PID and `"([^"]+)"` for command name
    4. Handle: when `-p` doesn't show process info (no root) → PID=0, Process="(unknown)"
    5. Set Protocol to "TCP", Type to empty string (ss doesn't distinguish IPv4/IPv6 as easily)
    6. Returns slice sorted by port number ascending
  - Handle edge cases:
    - `ss` not found → return error: "ss not found — install iproute2"
    - No output → return empty slice, nil error
    - Lines without process info → still show port with "(unknown)" process

  **Must NOT do**:
  - Do NOT create an interface
  - Do NOT try to parse `/proc/net/tcp` directly — use `ss` command
  - Do NOT add fallback to `netstat`

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Similar parsing logic to Task 2, one file, well-defined format
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 2)
  - **Parallel Group**: Wave 1b (with Task 2)
  - **Blocks**: Task 4
  - **Blocked By**: Task 1 (needs PortInfo struct from proc.go)

  **References**:

  **Pattern References**:
  - `proc.go:PortInfo` — Same struct to populate, same function signature as macOS version
  - `proc_darwin.go:GetListeningPorts` — Follow same pattern: exec command → parse output → return []PortInfo

  **External References**:
  - `ss -tlnp` output format:
    ```
    State   Recv-Q  Send-Q  Local Address:Port   Peer Address:Port  Process
    LISTEN  0       128     0.0.0.0:22            0.0.0.0:*          users:(("sshd",pid=1234,fd=3))
    LISTEN  0       128     [::]:22               [::]:*             users:(("sshd",pid=1234,fd=4))
    ```
  - IPv6 addresses are wrapped in brackets `[::]:22`, so split on last `:` is critical

  **WHY Each Reference Matters**:
  - Must match exact same function signature as darwin version so build tags work
  - ss output format differs significantly from lsof — executor needs the sample to parse correctly

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Linux cross-compile succeeds
    Tool: Bash
    Preconditions: Task 1 complete
    Steps:
      1. Run `GOOS=linux GOARCH=amd64 go build -o /dev/null .`
      2. Check exit code is 0
    Expected Result: Exit code 0 — compiles cleanly for Linux
    Failure Indicators: Compilation errors, missing function, build tag issues
    Evidence: .sisyphus/evidence/task-3-crosscompile.txt

  Scenario: Function signature matches macOS version
    Tool: Bash
    Preconditions: Both proc_darwin.go and proc_linux.go exist
    Steps:
      1. Grep both files for the function signature
      2. Run `grep "func GetListeningPorts" proc_darwin.go proc_linux.go`
      3. Verify both have identical signatures: `func GetListeningPorts() ([]PortInfo, error)`
    Expected Result: Both files show identical function signature
    Failure Indicators: Different signatures, different return types, missing function
    Evidence: .sisyphus/evidence/task-3-signature-match.txt
  ```

  **Evidence to Capture:**
  - [ ] task-3-crosscompile.txt — Linux build output
  - [ ] task-3-signature-match.txt — function signature comparison

  **Commit**: YES
  - Message: `feat: add Linux listening port discovery via ss`
  - Files: `proc_linux.go`
  - Pre-commit: `GOOS=linux GOARCH=amd64 go build -o /dev/null . && go vet ./...`

- [ ] 4. Bubbletea Model + Table Rendering + Auto-Refresh

  **What to do**:
  - Create `model.go` with the full Bubbletea model:
    ```go
    type model struct {
        table         table.Model
        allPorts      []PortInfo
        filteredPorts []PortInfo
        filterInput   textinput.Model
        filtering     bool
        filterText    string
        statusMsg     string
        statusTimer   time.Time
        lastRefresh   time.Time
        width         int
        height        int
        ready         bool
    }
    ```
  - Implement `func NewModel() model`:
    - Initialize `table.Model` with columns: Port (8w), PID (8w), Process (25w), Proto (6w), Address (20w), Type (6w)
    - Initialize `textinput.Model` for filter input (placeholder: "filter...")
    - Set `ready = false` until first data load
  - Implement `Init() tea.Cmd`:
    - Return `tea.Batch(fetchPorts, tickCmd())` where:
      - `fetchPorts` calls `GetListeningPorts()` and returns a `portsMsg`
      - `tickCmd` returns `tea.Tick(2*time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })`
  - Define message types:
    ```go
    type portsMsg []PortInfo
    type tickMsg time.Time
    type statusClearMsg struct{}
    ```
  - Implement `Update(msg tea.Msg) (tea.Model, tea.Cmd)`:
    - `tea.WindowSizeMsg`: store width/height, recalculate table height (height - 4 for title + status bars)
    - `portsMsg`: store in `allPorts`, apply current filter to get `filteredPorts`, rebuild table rows, set `lastRefresh`, set `ready = true`
    - `tickMsg`: return `tea.Batch(fetchPorts, tickCmd())` for continuous refresh
    - `tea.KeyMsg`: delegate to table for navigation keys, handle custom keys (quit for now — filter/kill/copy added in later tasks)
    - `statusClearMsg`: clear `statusMsg`
  - Implement `View() string`:
    - Title bar: `" ports (N listening)"` with title style
    - Table: `table.View()` — the main content
    - Status bar: `" Last refresh: HH:MM:SS | Press ? for help"` with status style
    - If `!ready`: show `" Loading..."` centered
    - Join vertically with `lipgloss.JoinVertical`
  - Helper `func portsToRows(ports []PortInfo) []table.Row`:
    - Convert each PortInfo to `table.Row{strconv.Itoa(p.Port), strconv.Itoa(p.PID), p.Process, p.Protocol, p.Address, p.Type}`
    - Sort by port number ascending before converting
  - Update `main.go`:
    - Replace minimal program with: `p := tea.NewProgram(NewModel(), tea.WithAltScreen())`
    - Add `if _, err := p.Run(); err != nil { fmt.Fprintf(os.Stderr, "Error: %v\n", err); os.Exit(1) }`

  **Must NOT do**:
  - Do NOT add filter logic yet — just the text input model initialization (Task 5 adds filter)
  - Do NOT add kill logic — just the key binding stub (Task 6 adds kill)
  - Do NOT add copy logic — just the key binding stub (Task 7 adds copy)
  - Do NOT add help overlay — just the `?` hint in status bar (Task 8 adds help)
  - Do NOT add mouse support (`tea.WithMouseCellMotion()`)
  - Do NOT use Bubbletea v2 APIs (`tea.KeyPressMsg` is v2 — use `tea.KeyMsg` which is v1)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Core architecture of the app — Bubbletea model pattern requires understanding Init/Update/View lifecycle, message passing, and async commands
  - **Skills**: []
    - No special skills — but agent should reference Bubbletea v1 patterns from `sammcj/gollama` or similar

  **Parallelization**:
  - **Can Run In Parallel**: NO (depends on all Wave 1 tasks)
  - **Parallel Group**: Wave 2 (solo)
  - **Blocks**: Tasks 5, 6, 7
  - **Blocked By**: Tasks 1, 2, 3

  **References**:

  **Pattern References**:
  - `proc.go:PortInfo` — Struct whose fields map to table columns
  - `proc_darwin.go:GetListeningPorts()` — Function to call in fetchPorts command
  - `styles.go` — All lipgloss styles defined in Task 1
  - `keys.go` — All key bindings defined in Task 1

  **External References**:
  - `github.com/charmbracelet/bubbles/table` — Table component. Use `table.New(table.WithColumns(cols), table.WithRows(rows), table.WithFocused(true), table.WithHeight(h))`. Update rows with `t.SetRows(rows)`.
  - `github.com/charmbracelet/bubbles/textinput` — Text input for filter. Initialize with `textinput.New()`, set placeholder.
  - Reference `sammcj/gollama` for the pattern: model struct wrapping table.Model, tickMsg for periodic refresh, tea.Batch for combining commands
  - CRITICAL: Use `tea.KeyMsg` (v1), NOT `tea.KeyPressMsg` (v2). The type switch in Update should be `case tea.KeyMsg:` not `case tea.KeyPressMsg:`.

  **WHY Each Reference Matters**:
  - bubbles/table API is the foundation — wrong initialization means no table renders
  - textinput is initialized here but activated in Task 5 — must be part of model struct now
  - v1 vs v2 key message type is the #1 source of bugs when AI generates Bubbletea code

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: TUI launches and shows listening ports
    Tool: interactive_bash (tmux)
    Preconditions: Tasks 1-3 complete, binary builds
    Steps:
      1. Build: `go build -o ports .`
      2. Launch in tmux: `tmux new-session -d -s ports-test './ports'`
      3. Wait 2 seconds for data load
      4. Capture pane: `tmux capture-pane -t ports-test -p`
      5. Assert output contains table header with "PORT" and "PID" and "PROCESS"
      6. Assert output contains at least one row with a real port number
      7. Kill session: `tmux send-keys -t ports-test q`
    Expected Result: Table header visible, at least one port row with real data
    Failure Indicators: Blank screen, panic text, "Loading..." stuck forever, no port data
    Evidence: .sisyphus/evidence/task-4-tui-launch.txt

  Scenario: Auto-refresh detects new listener
    Tool: interactive_bash (tmux) + Bash
    Preconditions: TUI running in tmux
    Steps:
      1. Launch TUI in tmux: `tmux new-session -d -s ports-test './ports'`
      2. Wait 1 second
      3. Start a dummy listener: `nc -l 19876 &` and note PID
      4. Wait 3 seconds (at least one refresh cycle)
      5. Capture pane: `tmux capture-pane -t ports-test -p`
      6. Assert output contains "19876"
      7. Kill nc: `kill $NC_PID`
      8. Quit TUI: `tmux send-keys -t ports-test q`
    Expected Result: Port 19876 appears in the table within 3 seconds without manual action
    Failure Indicators: Port 19876 never appears, refresh not working
    Evidence: .sisyphus/evidence/task-4-autorefresh.txt

  Scenario: Quit with q key
    Tool: interactive_bash (tmux)
    Preconditions: TUI running
    Steps:
      1. Launch in tmux: `tmux new-session -d -s ports-test './ports'`
      2. Wait 1 second
      3. Send quit: `tmux send-keys -t ports-test q`
      4. Wait 1 second
      5. Check if tmux session still exists: `tmux has-session -t ports-test 2>/dev/null && echo "STILL RUNNING" || echo "EXITED CLEAN"`
    Expected Result: "EXITED CLEAN" — process exits gracefully
    Failure Indicators: "STILL RUNNING", session hangs
    Evidence: .sisyphus/evidence/task-4-quit.txt

  Scenario: Terminal resize doesn't crash
    Tool: interactive_bash (tmux)
    Preconditions: TUI running
    Steps:
      1. Launch in tmux at small size: `tmux new-session -d -s ports-test -x 80 -y 24 './ports'`
      2. Wait 1 second
      3. Resize: `tmux resize-window -t ports-test -x 120 -y 40`
      4. Wait 1 second
      5. Capture pane: `tmux capture-pane -t ports-test -p`
      6. Assert table content is still visible (contains "PORT")
      7. Quit: `tmux send-keys -t ports-test q`
    Expected Result: Table renders correctly after resize, no crash
    Failure Indicators: Panic, garbled output, blank screen
    Evidence: .sisyphus/evidence/task-4-resize.txt
  ```

  **Evidence to Capture:**
  - [ ] task-4-tui-launch.txt — tmux pane capture showing table with ports
  - [ ] task-4-autorefresh.txt — pane capture showing dynamically added port
  - [ ] task-4-quit.txt — clean exit verification
  - [ ] task-4-resize.txt — post-resize pane capture

  **Commit**: YES
  - Message: `feat: implement TUI with table view and auto-refresh`
  - Files: `model.go, main.go`
  - Pre-commit: `go build -o ports . && go vet ./...`

- [ ] 5. Filter/Search Feature

  **What to do**:
  - In `model.go`, implement the filter mode:
  - When `/` key pressed AND not already filtering:
    - Set `filtering = true`
    - Focus the `filterInput` text input
    - Clear any previous filter text
  - When `Esc` pressed while filtering:
    - Set `filtering = false`
    - Blur the `filterInput`
    - Clear `filterText`
    - Reset `filteredPorts = allPorts`
    - Rebuild table rows with all ports
  - When typing in filter mode (filterInput updates):
    - Get current filterInput value
    - Set `filterText` to that value
    - Filter `allPorts` by case-insensitive substring match against: port number (as string), process name, and address
    - Store result in `filteredPorts`
    - Rebuild table rows with `filteredPorts`
  - When `portsMsg` arrives (refresh) while filter is active:
    - Update `allPorts` with new data
    - Re-apply current filter to get updated `filteredPorts`
  - Update `View()`:
    - When `filtering`: show filter input in status bar area: `" Filter: [input] (N results)"`
    - When `filterText != ""` but not actively filtering: show `" Filtered: "xyz" (N results) | Esc to clear"`
  - Handle empty filter results: show "No matching ports" in table area

  **Must NOT do**:
  - Do NOT add regex filtering — simple case-insensitive substring only
  - Do NOT add multiple filter modes (column-specific filtering)
  - Do NOT persist filter between sessions

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Well-scoped feature addition to existing model — text input handling + array filtering
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 6, 7)
  - **Parallel Group**: Wave 3 (with Tasks 6, 7)
  - **Blocks**: Task 8
  - **Blocked By**: Task 4

  **References**:

  **Pattern References**:
  - `model.go:model` — Add filter behavior to existing Update() switch
  - `model.go:portsToRows` — Reuse to rebuild table rows after filtering
  - `keys.go:Filter`, `keys.go:ClearFilter` — Key bindings defined in Task 1

  **External References**:
  - `github.com/charmbracelet/bubbles/textinput` — `textinput.Model` has `.Value()`, `.Focus()`, `.Blur()`, `.SetValue("")`. In Update, call `filterInput, cmd = filterInput.Update(msg)` to handle typing. Check `filterInput.Value()` changed to trigger re-filter.
  - Use `strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))` for case-insensitive match

  **WHY Each Reference Matters**:
  - textinput Update pattern is crucial — must be called inside the model's Update to process keystrokes
  - The filter must re-apply on refresh (portsMsg) to keep filtered view consistent with live data

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Filter narrows results
    Tool: interactive_bash (tmux)
    Preconditions: TUI running with multiple ports visible
    Steps:
      1. Build and launch: `go build -o ports . && tmux new-session -d -s ports-test './ports'`
      2. Wait 2 seconds
      3. Capture initial pane to count rows: `tmux capture-pane -t ports-test -p > /tmp/before.txt`
      4. Enter filter mode: `tmux send-keys -t ports-test /`
      5. Type a known process name that exists (e.g., partial match): `tmux send-keys -t ports-test "Control"`
      6. Wait 0.5 seconds
      7. Capture filtered pane: `tmux capture-pane -t ports-test -p > /tmp/after.txt`
      8. Assert: filtered output has fewer rows than initial OR only shows matching rows
      9. Quit: `tmux send-keys -t ports-test Escape` then `tmux send-keys -t ports-test q`
    Expected Result: Filtered view shows only rows containing "Control" (case-insensitive)
    Failure Indicators: Same number of rows, no filtering applied, crash on filter input
    Evidence: .sisyphus/evidence/task-5-filter-narrows.txt

  Scenario: Esc clears filter and shows all ports
    Tool: interactive_bash (tmux)
    Preconditions: TUI running with active filter
    Steps:
      1. Launch and apply filter (steps from above)
      2. Press Escape: `tmux send-keys -t ports-test Escape`
      3. Wait 0.5 seconds
      4. Capture pane: `tmux capture-pane -t ports-test -p`
      5. Assert: all ports visible again (row count matches pre-filter)
      6. Quit: `tmux send-keys -t ports-test q`
    Expected Result: Full port list restored after Esc
    Failure Indicators: Still filtered, filter text still showing
    Evidence: .sisyphus/evidence/task-5-filter-clear.txt

  Scenario: Filter with no matches shows empty state
    Tool: interactive_bash (tmux)
    Preconditions: TUI running
    Steps:
      1. Launch TUI in tmux
      2. Enter filter: `tmux send-keys -t ports-test /`
      3. Type nonsense: `tmux send-keys -t ports-test "zzzznonexistent"`
      4. Wait 0.5 seconds
      5. Capture pane: `tmux capture-pane -t ports-test -p`
      6. Assert: output shows "No matching ports" or empty table, filter indicator shows "0 results"
      7. Quit: `tmux send-keys -t ports-test Escape` then `q`
    Expected Result: Empty state message visible, no crash
    Failure Indicators: Crash, panic, table shows stale data
    Evidence: .sisyphus/evidence/task-5-filter-empty.txt
  ```

  **Evidence to Capture:**
  - [ ] task-5-filter-narrows.txt — before/after filter comparison
  - [ ] task-5-filter-clear.txt — filter clear verification
  - [ ] task-5-filter-empty.txt — empty results state

  **Commit**: YES
  - Message: `feat: add port/process filter with live search`
  - Files: `model.go`
  - Pre-commit: `go build -o ports . && go vet ./...`

- [ ] 6. Kill Process Feature

  **What to do**:
  - In `model.go`, implement kill action:
  - When `x` key pressed (NOT `k` — see keys.go) AND not in filter mode:
    1. Get currently selected table row
    2. Extract PID from the row (column index 1, parse as int)
    3. Call `syscall.Kill(pid, syscall.SIGTERM)`
    4. Handle results:
       - **Success**: Set `statusMsg = "Killed PID <pid> (<process name>)"`, trigger immediate refresh (`fetchPorts` command), start status clear timer
       - **EPERM** (permission denied): Set `statusMsg = "Permission denied: cannot kill PID <pid> (<process name>)"`, no refresh needed
       - **ESRCH** (no such process): Set `statusMsg = "Process <pid> already terminated"`, trigger refresh to remove stale entry
    5. Status message auto-clears after 3 seconds via `tea.Tick(3*time.Second, func(time.Time) tea.Msg { return statusClearMsg{} })`
  - Import `syscall` for Kill and signal constants
  - Import `errors` for `errors.Is(err, syscall.EPERM)` and `errors.Is(err, syscall.ESRCH)` checks

  **Must NOT do**:
  - Do NOT add a confirmation dialog — weekend scope, keep it simple
  - Do NOT add signal selection (TERM vs KILL vs HUP) — SIGTERM only
  - Do NOT add "kill all matching filter" bulk kill
  - Do NOT skip error handling — EPERM and ESRCH are common and must show user-friendly messages

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Small, focused feature — syscall.Kill + error handling + status message
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 5, 7)
  - **Parallel Group**: Wave 3 (with Tasks 5, 7)
  - **Blocks**: Task 8
  - **Blocked By**: Task 4

  **References**:

  **Pattern References**:
  - `model.go:Update` — Add kill handling to the tea.KeyMsg switch case
  - `model.go:statusMsg` — Use existing status message field for feedback
  - `model.go:statusClearMsg` — Use existing message type for auto-clear timer
  - `keys.go:Kill` — Key binding (`x`) defined in Task 1

  **External References**:
  - `syscall.Kill(pid, syscall.SIGTERM)` — Standard Go syscall for process signaling
  - `errors.Is(err, syscall.EPERM)` — Check for permission denied (trying to kill root-owned process)
  - `errors.Is(err, syscall.ESRCH)` — Check for "no such process" (race: process died between list and kill)

  **WHY Each Reference Matters**:
  - Error handling is the critical path here — without EPERM/ESRCH handling, killing a root process or a recently-exited process will show cryptic errors instead of friendly messages
  - Must trigger refresh after successful kill so the dead process disappears from the table

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Kill a dummy process successfully
    Tool: interactive_bash (tmux) + Bash
    Preconditions: TUI built and working
    Steps:
      1. Start dummy listener: `nc -l 19999 &` and capture PID: `NC_PID=$!`
      2. Build and launch: `go build -o ports . && tmux new-session -d -s ports-test './ports'`
      3. Wait 2 seconds for data load
      4. Enter filter to find the port: `tmux send-keys -t ports-test / 1 9 9 9 9`
      5. Wait 0.5 seconds
      6. Exit filter mode: `tmux send-keys -t ports-test Escape`
      7. Press kill key: `tmux send-keys -t ports-test x`
      8. Wait 1 second
      9. Verify nc is dead: `kill -0 $NC_PID 2>/dev/null && echo "STILL ALIVE" || echo "KILLED OK"`
      10. Capture pane for status message: `tmux capture-pane -t ports-test -p`
      11. Assert pane contains "Killed PID" status message
      12. Quit: `tmux send-keys -t ports-test q`
    Expected Result: "KILLED OK" printed, status bar shows "Killed PID ... (nc)" message
    Failure Indicators: "STILL ALIVE", no status message, crash, panic
    Evidence: .sisyphus/evidence/task-6-kill-success.txt

  Scenario: Kill permission denied on root process
    Tool: interactive_bash (tmux)
    Preconditions: TUI running, at least one root-owned process visible (common on macOS)
    Steps:
      1. Launch TUI in tmux
      2. Navigate to a system process (PID 1 or similar root-owned)
      3. Press kill: `tmux send-keys -t ports-test x`
      4. Wait 0.5 seconds
      5. Capture pane: `tmux capture-pane -t ports-test -p`
      6. Assert status bar shows "Permission denied" message
      7. Assert TUI didn't crash — table still visible
      8. Quit: `tmux send-keys -t ports-test q`
    Expected Result: "Permission denied" in status bar, TUI still functional
    Failure Indicators: Crash, panic, no error message, process actually killed (should fail)
    Evidence: .sisyphus/evidence/task-6-kill-eperm.txt
  ```

  **Evidence to Capture:**
  - [ ] task-6-kill-success.txt — successful kill + process verification
  - [ ] task-6-kill-eperm.txt — permission denied handling

  **Commit**: YES
  - Message: `feat: add process kill with error handling`
  - Files: `model.go`
  - Pre-commit: `go build -o ports . && go vet ./...`

- [ ] 7. Clipboard Copy Feature

  **What to do**:
  - In `model.go`, implement copy action:
  - When `c` key pressed AND not in filter mode:
    1. Get currently selected table row
    2. Format as: `<port>\t<pid>\t<process>\t<address>` (tab-separated for easy pasting)
    3. Determine clipboard command based on `runtime.GOOS`:
       - `darwin`: exec `pbcopy` with formatted string piped to stdin
       - `linux`: exec `xclip -selection clipboard` with formatted string piped to stdin
    4. Handle errors:
       - Command not found: `statusMsg = "Clipboard not available (install xclip on Linux)"`
       - Other exec error: `statusMsg = "Copy failed: <error>"`
       - Success: `statusMsg = "Copied: port <port> (<process>)"`
    5. Status auto-clears after 3 seconds (reuse statusClearMsg pattern from Task 6)
  - Clipboard exec helper function (private):
    ```go
    func copyToClipboard(text string) error {
        var cmd *exec.Cmd
        switch runtime.GOOS {
        case "darwin":
            cmd = exec.Command("pbcopy")
        case "linux":
            cmd = exec.Command("xclip", "-selection", "clipboard")
        default:
            return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
        }
        cmd.Stdin = strings.NewReader(text)
        return cmd.Run()
    }
    ```

  **Must NOT do**:
  - Do NOT use `golang.design/x/clipboard` (requires CGO)
  - Do NOT add "copy as JSON" or "copy as curl" formats
  - Do NOT add multi-row copy/selection
  - Do NOT crash if clipboard command is missing — show friendly error

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Small feature — exec a command, pipe stdin, handle errors
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 5, 6)
  - **Parallel Group**: Wave 3 (with Tasks 5, 6)
  - **Blocks**: Task 8
  - **Blocked By**: Task 4

  **References**:

  **Pattern References**:
  - `model.go:Update` — Add copy handling to tea.KeyMsg switch
  - `model.go:statusMsg` — Reuse status message pattern from Task 6
  - `keys.go:Copy` — Key binding (`c`) defined in Task 1

  **External References**:
  - `os/exec.Command` — Standard Go exec. Pipe via `cmd.Stdin = strings.NewReader(text)`.
  - `pbcopy` — macOS clipboard. Reads from stdin.
  - `xclip -selection clipboard` — Linux X11 clipboard. `-selection clipboard` targets Ctrl+V clipboard (not X11 primary selection).

  **WHY Each Reference Matters**:
  - `cmd.Stdin` pipe pattern is the key implementation detail — without it, pbcopy gets no input
  - `-selection clipboard` flag is critical on Linux — without it, xclip uses primary selection which is middle-click paste, not Ctrl+V

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Copy puts port info on clipboard (macOS)
    Tool: interactive_bash (tmux) + Bash
    Preconditions: TUI running on macOS
    Steps:
      1. Launch TUI: `go build -o ports . && tmux new-session -d -s ports-test './ports'`
      2. Wait 2 seconds
      3. Press copy on first row: `tmux send-keys -t ports-test c`
      4. Wait 0.5 seconds
      5. Capture pane for status: `tmux capture-pane -t ports-test -p`
      6. Assert status bar shows "Copied:" message
      7. Check clipboard: `pbpaste`
      8. Assert clipboard contains a port number and process name
      9. Quit: `tmux send-keys -t ports-test q`
    Expected Result: Clipboard contains tab-separated port info, status bar shows "Copied: port X (processname)"
    Failure Indicators: Empty clipboard, no status message, crash
    Evidence: .sisyphus/evidence/task-7-copy-clipboard.txt

  Scenario: Copy when clipboard unavailable doesn't crash
    Tool: Bash
    Preconditions: Verify the error path handles missing clipboard gracefully
    Steps:
      1. Temporarily test by checking the error handling logic in copyToClipboard
      2. Verify the function returns an error (not panic) when command not found
      3. Build succeeds: `go build -o ports . && echo "BUILD OK"`
    Expected Result: Error returned gracefully, no panic
    Failure Indicators: Panic, nil pointer, crash
    Evidence: .sisyphus/evidence/task-7-copy-fallback.txt
  ```

  **Evidence to Capture:**
  - [ ] task-7-copy-clipboard.txt — clipboard content verification
  - [ ] task-7-copy-fallback.txt — error handling verification

  **Commit**: YES
  - Message: `feat: add clipboard copy support`
  - Files: `model.go`
  - Pre-commit: `go build -o ports . && go vet ./...`

- [ ] 8. Help Overlay + Polish + Comprehensive QA

  **What to do**:
  - Add help overlay toggle:
    - Add `showHelp bool` field to model
    - When `?` pressed: toggle `showHelp`
    - When `showHelp` is true, render a centered overlay box listing all keybindings:
      ```
      ┌─── Help ───────────────────────┐
      │                                │
      │  ↑/↓ or j/k   Navigate        │
      │  /             Filter          │
      │  Esc           Clear filter    │
      │  x             Kill process    │
      │  c             Copy to clipboard│
      │  r             Manual refresh  │
      │  ?             Toggle help     │
      │  q             Quit            │
      │                                │
      │  Press ? to close              │
      └────────────────────────────────┘
      ```
    - Use lipgloss to style the box with a border and slight padding
    - Overlay renders ON TOP of the table (replace table view, not alongside)
  - Polish the status bar:
    - Always show: `" Last refresh: HH:MM:SS | N listening | ? help"`
    - When status message active: show that instead of the default
  - Polish the title bar:
    - Show: `" ⚡ ports (N listening)"` — include count of visible (filtered) ports
  - Ensure empty state is handled:
    - If `allPorts` is empty after load: show `"No listening TCP ports found"` centered
    - If `filteredPorts` is empty (filter applied): show `"No ports match filter"` centered
  - Run `go vet ./...` and fix any issues
  - Run full QA suite (all scenarios from all tasks)

  **Must NOT do**:
  - Do NOT add a separate help page/view — just a simple overlay
  - Do NOT add color customization
  - Do NOT add ASCII art or large banners
  - Do NOT add version number or build info display

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: UI polish + overlay is straightforward lipgloss work, plus QA execution
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO (depends on all Wave 3 tasks)
  - **Parallel Group**: Wave 4 (solo)
  - **Blocks**: F1-F4
  - **Blocked By**: Tasks 5, 6, 7

  **References**:

  **Pattern References**:
  - `model.go:View()` — Add help overlay rendering, polish status/title bars
  - `styles.go` — Use existing styles, may add a help box style with lipgloss.Border
  - `keys.go:Help` — Key binding (`?`) defined in Task 1

  **External References**:
  - `lipgloss.NewStyle().Border(lipgloss.RoundedBorder())` — For the help overlay box
  - `lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)` — For centering the help overlay

  **WHY Each Reference Matters**:
  - lipgloss.Place is how to center content in the terminal — without it, the help box will be top-left aligned
  - Border style creates the visual box around help text

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Help overlay toggles with ? key
    Tool: interactive_bash (tmux)
    Preconditions: All features working
    Steps:
      1. Build and launch: `go build -o ports . && tmux new-session -d -s ports-test './ports'`
      2. Wait 2 seconds
      3. Press help: `tmux send-keys -t ports-test ?`
      4. Wait 0.5 seconds
      5. Capture pane: `tmux capture-pane -t ports-test -p`
      6. Assert: output contains "Help" and keybinding descriptions (e.g., "Kill process", "Filter")
      7. Press help again to close: `tmux send-keys -t ports-test ?`
      8. Wait 0.5 seconds
      9. Capture pane: `tmux capture-pane -t ports-test -p`
      10. Assert: table is visible again (contains "PORT" header)
      11. Quit: `tmux send-keys -t ports-test q`
    Expected Result: Help overlay appears and disappears on ? toggle
    Failure Indicators: Help doesn't show, help doesn't dismiss, crash
    Evidence: .sisyphus/evidence/task-8-help-toggle.txt

  Scenario: Full integration — filter → kill → copy → help
    Tool: interactive_bash (tmux) + Bash
    Preconditions: All features working, dummy process available
    Steps:
      1. Start dummy: `nc -l 19777 &` → `NC_PID=$!`
      2. Build and launch: `go build -o ports . && tmux new-session -d -s ports-test './ports'`
      3. Wait 2 seconds
      4. Filter for dummy: `tmux send-keys -t ports-test / 1 9 7 7 7`
      5. Wait 0.5 seconds, capture: verify filtered view shows 19777
      6. Clear filter: `tmux send-keys -t ports-test Escape`
      7. Navigate to 19777 row (may need j/k)
      8. Copy: `tmux send-keys -t ports-test c`
      9. Verify clipboard: `pbpaste | grep 19777`
      10. Kill: `tmux send-keys -t ports-test x`
      11. Wait 1 second
      12. Verify killed: `kill -0 $NC_PID 2>/dev/null && echo "ALIVE" || echo "DEAD"`
      13. Show help: `tmux send-keys -t ports-test ?`
      14. Capture help overlay
      15. Close help and quit: `tmux send-keys -t ports-test ? q`
    Expected Result: All features work in sequence — filter narrows, copy works, kill works, help shows
    Failure Indicators: Any step fails, state corruption between actions
    Evidence: .sisyphus/evidence/task-8-integration.txt

  Scenario: go vet passes clean
    Tool: Bash
    Preconditions: All code complete
    Steps:
      1. Run `go vet ./...`
      2. Assert exit code 0 and no output
    Expected Result: Clean — no warnings or errors
    Failure Indicators: Vet warnings, unused imports, type errors
    Evidence: .sisyphus/evidence/task-8-govet.txt

  Scenario: Cross-compile still works after all changes
    Tool: Bash
    Preconditions: All code complete
    Steps:
      1. Run `GOOS=linux GOARCH=amd64 go build -o /dev/null .`
      2. Assert exit code 0
    Expected Result: Linux build succeeds
    Failure Indicators: Build errors from platform-specific code leaking
    Evidence: .sisyphus/evidence/task-8-crosscompile.txt
  ```

  **Evidence to Capture:**
  - [ ] task-8-help-toggle.txt — help overlay on/off
  - [ ] task-8-integration.txt — full feature integration test
  - [ ] task-8-govet.txt — go vet output
  - [ ] task-8-crosscompile.txt — Linux cross-compile verification

  **Commit**: YES
  - Message: `chore: polish UI, add help overlay, comprehensive QA`
  - Files: `model.go, styles.go`
  - Pre-commit: `go build -o ports . && go vet ./... && GOOS=linux GOARCH=amd64 go build -o /dev/null .`

---

## Final Verification Wave

> 4 review agents run in PARALLEL. ALL must APPROVE. Rejection → fix → re-run.

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Read the plan end-to-end. For each "Must Have": verify implementation exists (read file, run command). For each "Must NOT Have": search codebase for forbidden patterns (v2 imports, golang.design/x/clipboard, config file parsers, mouse support) — reject with file:line if found. Check evidence files exist in .sisyphus/evidence/. Compare deliverables against plan.
  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT: APPROVE/REJECT`

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `go build ./...`, `go vet ./...`. Review all .go files for: `any` type assertions, empty error handling (`_ = err`), leftover `fmt.Println` debug lines, commented-out code, unused imports. Check AI slop: excessive comments, over-abstraction, unnecessary interfaces, generic variable names. Verify build tags are correct (`//go:build darwin` not `// +build darwin`).
  Output: `Build [PASS/FAIL] | Vet [PASS/FAIL] | Files [N clean/N issues] | VERDICT`

- [ ] F3. **Real Manual QA** — `unspecified-high`
  Start from clean state. Execute EVERY QA scenario from EVERY task in tmux — follow exact steps, capture evidence. Test cross-feature integration (filter → kill filtered result → copy after kill). Test edge cases: empty filter results, kill permission denied, clipboard with long process names. Save to `.sisyphus/evidence/final-qa/`.
  Output: `Scenarios [N/N pass] | Integration [N/N] | Edge Cases [N tested] | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`
  For each task: read "What to do", read actual code. Verify 1:1 — everything in spec was built (no missing), nothing beyond spec was built (no creep). Check "Must NOT Have" compliance — search for: config file loading, theme systems, mouse handlers, v2 imports, interfaces for data layer, unit test files. Flag any unaccounted files or features.
  Output: `Tasks [N/N compliant] | Scope Creep [CLEAN/N issues] | Unaccounted [CLEAN/N files] | VERDICT`

---

## Commit Strategy

| After Task | Commit Message | Files |
|-----------|---------------|-------|
| 1 | `feat: initialize project skeleton with Go module, types, styles, and keybindings` | go.mod, go.sum, main.go, proc.go, styles.go, keys.go |
| 2 | `feat: add macOS listening port discovery via lsof` | proc_darwin.go |
| 3 | `feat: add Linux listening port discovery via ss` | proc_linux.go |
| 4 | `feat: implement TUI with table view and auto-refresh` | model.go, main.go |
| 5 | `feat: add port/process filter with live search` | model.go |
| 6 | `feat: add process kill with error handling` | model.go |
| 7 | `feat: add clipboard copy support` | model.go |
| 8 | `chore: polish UI, add help overlay, comprehensive QA` | model.go, styles.go |

---

## Success Criteria

### Verification Commands
```bash
go build -o ports .                                    # Expected: exits 0, binary created
GOOS=linux GOARCH=amd64 go build -o /dev/null .        # Expected: exits 0 (cross-compile)
go vet ./...                                           # Expected: no issues
./ports                                                # Expected: TUI shows listening ports
```

### Final Checklist
- [ ] Binary compiles for macOS and Linux
- [ ] Table shows real listening ports matching `lsof` output
- [ ] Kill process works on user-owned processes
- [ ] Filter narrows results, Esc clears
- [ ] Copy puts port info on clipboard
- [ ] Auto-refresh detects new listeners
- [ ] Help overlay shows all keybindings
- [ ] Terminal resize doesn't crash
- [ ] Empty state handled gracefully
- [ ] No Bubbletea v2 imports
- [ ] No CGO dependencies
- [ ] No config files, themes, or out-of-scope features

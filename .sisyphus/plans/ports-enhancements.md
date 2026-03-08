# Ports TUI — 10 Feature Enhancements

## TL;DR

> **Quick Summary**: Add 10 enhancements to the working `ports` TUI: dedup toggle, SIGKILL, new/disappeared port markers, column sorting, open-in-browser, UDP support, --port flag, connection count, --diff mode, and service name hints.
> 
> **Deliverables**:
> - 2 new files: `services.go` (service name map), `diff.go` (--diff mode logic)
> - Modified files: `proc.go`, `proc_darwin.go`, `proc_linux.go`, `keys.go`, `styles.go`, `model.go`, `main.go`
> - All 10 features functional with tmux-verified QA
> 
> **Estimated Effort**: Medium-Large
> **Parallel Execution**: YES — 4 waves
> **Critical Path**: T1 (struct+keys) → T5 (model integration batch 1) → T7 (model integration batch 2) → T9 (--diff mode) → F1-F4 (final verification)

---

## Context

### Original Request
User requested 10 specific enhancements to the working `ports` TUI after Phase 1 (8 tasks) was fully completed and verified.

### Interview Summary
**Key Discussions**:
- Dedup toggle: keep both IPv4+IPv6 rows by default, `m` key toggles merge ON/OFF. Merge key: (Port, PID).
- Highlight approach: STATUS column with text markers (`●` new, `○` disappeared) — bubbles/table v1 has no per-row styling API.
- Dedup key: `m` (for merge) — `d` conflicts with bubbles/table half-page-down.
- `X` (shift-x) for SIGKILL — consistent with `x` for SIGTERM.
- Diff state storage: `~/.cache/ports/` (via `os.UserCacheDir()`).
- No unit tests — agent-executed QA via tmux only.

### Metis Review
**Identified Gaps** (all addressed):
- **Key `d` conflict with table HalfPageDown** → Resolved: use `m` for merge
- **Per-row styling impossible in bubbles/table v1** → Resolved: STATUS column with text markers
- **Linux Type field empty in proc_linux.go** → Fixed: detect IPv4/IPv6 from address format
- **Sort/dedup must persist through 2s refresh** → Store state in model struct
- **First `--diff` run with no cache** → Show all as new, save baseline
- **UDP + connection count interaction** → Show `-` for UDP rows
- **Multiple clearStatusCmd timers** → Fix with status generation counter
- **Filter must search new columns** → Extend filterPorts for SERVICE/CONNECTIONS
- **`os.UserCacheDir()` for cross-platform cache** → Yes, with `~/.cache/ports` fallback

---

## Work Objectives

### Core Objective
Add 10 enhancements to the `ports` TUI: dedup toggle, SIGKILL, port change markers, column sorting, open-in-browser, UDP support, --port flag, connection count, --diff CLI mode, and service name hints.

### Concrete Deliverables
- `services.go`: hardcoded map[int]string with ≤30 common port→service mappings
- `diff.go`: --diff mode logic (load cache, compare, print diff, save cache)
- `proc.go`: PortInfo struct expanded with `Connections int`, `Service string`, `Status string`
- `proc_darwin.go`: UDP parser function, connection count function
- `proc_linux.go`: Type field detection, UDP parser function, connection count function
- `keys.go`: 5 new key bindings (ForceKill `X`, Sort `s`, Open `o`, ToggleUDP `t`, ToggleDedup `m`)
- `styles.go`: new/disappeared marker styles
- `model.go`: all 10 features integrated into model struct, Update, and View
- `main.go`: `--port` and `--diff` flag parsing

### Definition of Done
- [ ] `go build -o ports ./...` passes
- [ ] `go vet ./...` passes
- [ ] `GOOS=linux go build -o /dev/null ./...` cross-compiles
- [ ] All 10 features verified via tmux QA scenarios
- [ ] Help overlay (`?`) shows all new key bindings

### Must Have
- All 10 features functional
- Dedup OFF by default, toggled with `m`
- `X` (uppercase) sends SIGKILL, `x` (lowercase) sends SIGTERM
- `--diff` runs once and exits (no TUI)
- `--port N` pre-filters to port N on startup
- STATUS column shows `●` for new ports, `○` for disappeared ports
- Sort persists through auto-refresh cycles
- Dedup persists through auto-refresh cycles
- Service name map capped at 30 entries

### Must NOT Have (Guardrails)
- **NO Bubbletea v2 APIs** — must use `tea.KeyMsg`, NOT `tea.KeyPressMsg`
- **NO `golang.design/x/clipboard`** — exec-based clipboard only
- **NO confirmation dialog for SIGKILL** — direct kill on `X` press
- **NO dynamic column width adjustment** — fixed widths in NewModel
- **NO `--port` range/comma parsing** — single integer only
- **NO service name external lookups** — hardcoded map only, ≤30 entries
- **NO config files, themes, or mouse support**
- **NO interfaces or sub-packages** — flat `package main`
- **NO unit test files** — QA via tmux only
- **NO auto-detect HTTP/HTTPS for open-in-browser** — always `http://`

---

## Verification Strategy (MANDATORY)

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed. No exceptions.

### Test Decision
- **Infrastructure exists**: NO
- **Automated tests**: None
- **Framework**: None
- **QA Method**: Agent-executed tmux verification for every task

### QA Policy
Every task MUST include agent-executed QA scenarios.
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **TUI features**: Use `interactive_bash` (tmux) — launch `./ports`, send keys, capture pane, grep output
- **CLI flags**: Use Bash — run `./ports --diff`, `./ports --port 3000`, capture stdout
- **Build verification**: Use Bash — `go build`, `go vet`, `GOOS=linux go build`

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately — foundation, 5 parallel):
├── Task 1: Expand PortInfo struct + add all new key bindings [quick]
├── Task 2: Create services.go — service name map [quick]
├── Task 3: Fix Linux Type field + add UDP parser (proc_linux.go) [unspecified-low]
├── Task 4: Add UDP parser + connection count (proc_darwin.go) [unspecified-low]
└── Task 5: Add new/disappeared marker styles (styles.go) [quick]

Wave 2 (After Wave 1 — model.go integration batch 1, 1 sequential):
├── Task 6: Model integration — new columns, service names, dedup toggle,
│           sort cycling, status markers, UDP toggle [deep]

Wave 3 (After Wave 2 — model.go integration batch 2, 1 sequential):
├── Task 7: Model integration — SIGKILL, open-in-browser, connection count,
│           help overlay update, key dispatch fix [deep]

Wave 4 (After Wave 1 — independent from model changes, 2 parallel):
├── Task 8: main.go — --port flag + --diff flag parsing [quick]
└── Task 9: diff.go — --diff mode logic [unspecified-high]

Wave FINAL (After ALL tasks — 4 parallel):
├── Task F1: Plan compliance audit [oracle]
├── Task F2: Code quality review [unspecified-high]
├── Task F3: Real manual QA [unspecified-high]
└── Task F4: Scope fidelity check [deep]

Critical Path: T1 → T6 → T7 → (build + full QA)
Parallel Speedup: ~50% faster than sequential
Max Concurrent: 5 (Wave 1)
```

### Dependency Matrix

| Task | Depends On | Blocks | Wave |
|------|-----------|--------|------|
| T1 | — | T6, T7, T8 | 1 |
| T2 | — | T6 | 1 |
| T3 | — | T6, T9 | 1 |
| T4 | — | T6, T7 | 1 |
| T5 | — | T6 | 1 |
| T6 | T1, T2, T3, T4, T5 | T7 | 2 |
| T7 | T6 | F1-F4 | 3 |
| T8 | T1 | T9 | 4 |
| T9 | T3, T8 | F1-F4 | 4 |
| F1-F4 | T7, T9 | — | FINAL |

### Agent Dispatch Summary

- **Wave 1**: **5 tasks** — T1 `quick`, T2 `quick`, T3 `unspecified-low`, T4 `unspecified-low`, T5 `quick`
- **Wave 2**: **1 task** — T6 `deep`
- **Wave 3**: **1 task** — T7 `deep`
- **Wave 4**: **2 tasks** — T8 `quick`, T9 `unspecified-high`
- **FINAL**: **4 tasks** — F1 `oracle`, F2 `unspecified-high`, F3 `unspecified-high`, F4 `deep`

---

## TODOs

- [x] 1. Expand PortInfo Struct + Add All New Key Bindings

  **What to do**:
  - In `proc.go`: Add 3 new fields to `PortInfo` struct:
    - `Connections int` — active TCP connection count (0 for UDP)
    - `Service string` — human-readable service name (e.g., "http", "ssh")
    - `Status string` — `"new"`, `"gone"`, or `""` for port change tracking
  - In `keys.go`: Add 5 new key bindings to the `keyMap` struct and `keys` var:
    - `ForceKill`: `key.WithKeys("X")`, help `"X"`, `"force kill (SIGKILL)"`
    - `Sort`: `key.WithKeys("s")`, help `"s"`, `"cycle sort"`
    - `Open`: `key.WithKeys("o")`, help `"o"`, `"open in browser"`
    - `ToggleUDP`: `key.WithKeys("t")`, help `"t"`, `"toggle TCP/UDP"`
    - `ToggleDedup`: `key.WithKeys("m")`, help `"m"`, `"toggle merge"`
  - **IMPORTANT**: `s` key might conflict with table keys. Verify by checking bubbles table v1 source — `s` is NOT used by default table keymap (only `up/down/k/j/g/G/d/u/ctrl+d/ctrl+u` are). Safe to use.

  **Must NOT do**:
  - Do NOT add any model struct fields here — model.go is handled in T6/T7
  - Do NOT rename or remove existing PortInfo fields
  - Do NOT change existing key bindings

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Two small files, additive changes only, ~20 lines total
  - **Skills**: `[]`
  - **Skills Evaluated but Omitted**:
    - `git-master`: Not needed — simple additive changes

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 2, 3, 4, 5)
  - **Blocks**: Tasks 6, 7, 8
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `keys.go:17-54` — Existing `keyMap` struct and `keys` var initialization pattern. Follow exact same `key.NewBinding(key.WithKeys(...), key.WithHelp(...))` pattern.
  - `proc.go:4-11` — Current PortInfo struct. Add new fields after `Type string`.

  **API/Type References**:
  - `github.com/charmbracelet/bubbles/key` — `key.NewBinding`, `key.WithKeys`, `key.WithHelp`

  **WHY Each Reference Matters**:
  - `keys.go` shows the exact binding pattern — executor must follow it identically
  - `proc.go` shows the struct that ALL parsers and model code depend on — field ordering matters for readability

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Build succeeds with new struct fields and key bindings
    Tool: Bash
    Preconditions: Working directory is /Users/msh/code/pp/passed
    Steps:
      1. Run: go build -o /dev/null ./...
      2. Run: go vet ./...
      3. Run: GOOS=linux go build -o /dev/null ./...
    Expected Result: All three commands exit 0 with no output
    Failure Indicators: Any compilation error mentioning undefined fields or unused imports
    Evidence: .sisyphus/evidence/task-1-build.txt

  Scenario: PortInfo struct has exactly 9 fields
    Tool: Bash
    Preconditions: proc.go has been modified
    Steps:
      1. Run: grep -c "^\s\+\w\+\s\+\(int\|string\)" proc.go
    Expected Result: Output is "9" (6 original + 3 new)
    Failure Indicators: Count is not 9
    Evidence: .sisyphus/evidence/task-1-struct-check.txt

  Scenario: keys.go has exactly 14 key bindings
    Tool: Bash
    Preconditions: keys.go has been modified
    Steps:
      1. Run: grep -c "key.NewBinding" keys.go
    Expected Result: Output is "14" (9 original + 5 new)
    Failure Indicators: Count is not 14
    Evidence: .sisyphus/evidence/task-1-keys-check.txt
  ```

  **Commit**: YES (group with Wave 1)
  - Message: `feat(ports): expand PortInfo struct and add new key bindings`
  - Files: `proc.go`, `keys.go`
  - Pre-commit: `go build -o /dev/null ./...`

---

- [x] 2. Create services.go — Service Name Map

  **What to do**:
  - Create new file `services.go` in project root with `package main`
  - Define `var serviceNames = map[int]string{...}` with ≤30 common port→service mappings
  - Define function `func serviceName(port int) string` that returns the service name or `""` if not in map
  - Include these ports at minimum: 21 (ftp), 22 (ssh), 25 (smtp), 53 (dns), 80 (http), 110 (pop3), 143 (imap), 443 (https), 993 (imaps), 995 (pop3s), 1433 (mssql), 1521 (oracle), 3000 (dev), 3306 (mysql), 4200 (ng-serve), 5000 (flask), 5432 (postgres), 5672 (amqp), 5900 (vnc), 6379 (redis), 8000 (dev), 8080 (http-alt), 8443 (https-alt), 8888 (jupyter), 9090 (prometheus), 9200 (elastic), 11211 (memcached), 27017 (mongodb)
  - **Cap at 30 entries maximum**

  **Must NOT do**:
  - Do NOT read from `/etc/services` or any external file
  - Do NOT add more than 30 entries
  - Do NOT create an interface or abstraction — just a map and a function

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Single new file, ~40 lines, hardcoded data
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 3, 4, 5)
  - **Blocks**: Task 6
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `proc.go:1-2` — File header pattern (`package main` only, no doc comment needed for internal file)

  **WHY Each Reference Matters**:
  - Ensures the new file follows the same minimal header convention as existing files

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: services.go compiles and has correct function signature
    Tool: Bash
    Preconditions: services.go created in /Users/msh/code/pp/passed
    Steps:
      1. Run: go build -o /dev/null ./...
      2. Run: go vet ./...
    Expected Result: Both exit 0
    Failure Indicators: Compilation error or vet warning
    Evidence: .sisyphus/evidence/task-2-build.txt

  Scenario: Service map has ≤30 entries
    Tool: Bash
    Preconditions: services.go exists
    Steps:
      1. Run: grep -c '^\s\+[0-9]' services.go
    Expected Result: Output is ≤30
    Failure Indicators: Count exceeds 30
    Evidence: .sisyphus/evidence/task-2-count.txt
  ```

  **Commit**: YES (group with Wave 1)
  - Message: `feat(ports): add service name lookup map`
  - Files: `services.go`
  - Pre-commit: `go build -o /dev/null ./...`

---

- [x] 3. Fix Linux Type Field + Add UDP Parser (proc_linux.go)

  **What to do**:
  - **Fix Type field detection**: In `parseSsOutput()`, after extracting the address, detect IPv4 vs IPv6:
    - If address contains `::` or was wrapped in `[]` → `Type: "IPv6"`
    - If address is `*` or contains only dots and digits → `Type: "IPv4"`
    - Fallback: `Type: "IPv4"` (most common)
  - **Add UDP parser function**: Create `GetUDPPorts() ([]PortInfo, error)` that runs `ss -ulnp` (note: `-u` for UDP, no `-t`) and parses output using same `parseSsOutput()` function but sets `Protocol: "UDP"`. UDP sockets show state `UNCONN` instead of `LISTEN` — both should be accepted.
  - **Add connection count function**: Create `GetConnectionCounts() (map[int]int, error)` that runs `ss -tnp state established` and counts established TCP connections per local port. Returns map[localPort]count.
  - **CRITICAL**: The existing `parseSsOutput` hardcodes `Protocol: "TCP"` on line 97. The UDP parser must override this to `"UDP"`. Options: (a) pass protocol as parameter, or (b) have GetUDPPorts post-process the results. Preferred: add `protocol string` parameter to `parseSsOutput`.

  **Must NOT do**:
  - Do NOT change the existing `GetListeningPorts()` function signature
  - Do NOT merge TCP and UDP parsing into one call — keep them separate
  - Do NOT add any model or TUI logic

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Parser changes in one file, following existing patterns, needs careful ss command knowledge
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2, 4, 5)
  - **Blocks**: Tasks 6, 9
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `proc_linux.go:20-35` — Existing `GetListeningPorts()` pattern: exec command, check error, parse output, sort
  - `proc_linux.go:37-102` — Existing `parseSsOutput()` parser: split lines, skip header, extract fields, regex PID/name
  - `proc_linux.go:14-17` — Regex patterns for PID and process name extraction

  **API/Type References**:
  - `proc.go:4-11` — PortInfo struct (will have new fields from T1: Connections, Service, Status)

  **External References**:
  - `ss` man page: `-u` flag for UDP, `-t` for TCP, `state established` filter for connection count
  - Linux `ss -ulnp` output format: same columns as `-tlnp` but state is `UNCONN` not `LISTEN`

  **WHY Each Reference Matters**:
  - `parseSsOutput` is the core parser — UDP reuse requires understanding its assumptions (hardcoded "TCP", state filtering)
  - The PortInfo struct shows what fields to populate — new `Connections` field will be filled by model.go, not here

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Linux build succeeds with new functions
    Tool: Bash
    Preconditions: proc_linux.go modified, proc.go has new fields from T1
    Steps:
      1. Run: GOOS=linux go build -o /dev/null ./...
    Expected Result: Exit 0, no errors
    Failure Indicators: Undefined function or type errors
    Evidence: .sisyphus/evidence/task-3-build.txt

  Scenario: parseSsOutput populates Type field correctly
    Tool: Bash
    Preconditions: Code is compiled
    Steps:
      1. Run: grep -n 'Type:' proc_linux.go
      2. Verify Type field is set conditionally (not hardcoded empty string)
    Expected Result: Type field assignment includes IPv4/IPv6 detection logic
    Failure Indicators: Type is still hardcoded as ""
    Evidence: .sisyphus/evidence/task-3-type-field.txt

  Scenario: GetUDPPorts function exists with correct signature
    Tool: Bash
    Preconditions: proc_linux.go modified
    Steps:
      1. Run: grep 'func GetUDPPorts' proc_linux.go
    Expected Result: Output contains "func GetUDPPorts() ([]PortInfo, error)"
    Failure Indicators: Function not found or wrong signature
    Evidence: .sisyphus/evidence/task-3-udp-func.txt

  Scenario: GetConnectionCounts function exists with correct signature
    Tool: Bash
    Preconditions: proc_linux.go modified
    Steps:
      1. Run: grep 'func GetConnectionCounts' proc_linux.go
    Expected Result: Output contains "func GetConnectionCounts() (map[int]int, error)"
    Failure Indicators: Function not found or wrong signature
    Evidence: .sisyphus/evidence/task-3-conn-func.txt
  ```

  **Commit**: YES (group with Wave 1)
  - Message: `feat(ports): fix Linux Type detection, add UDP and connection count parsers`
  - Files: `proc_linux.go`
  - Pre-commit: `GOOS=linux go build -o /dev/null ./...`

- [x] 4. Add UDP Parser + Connection Count (proc_darwin.go)

  **What to do**:
  - **Add UDP parser function**: Create `GetUDPPorts() ([]PortInfo, error)` that runs `lsof -iUDP -P -n -F pcfnPt` (note: NO `-sTCP:LISTEN` filter — UDP has no LISTEN state) and parses output using existing `parseLsofOutput()`. Post-process results to set `Protocol: "UDP"` on each entry (since lsof -F `P` field will return "UDP" naturally — verify this).
  - **Add connection count function**: Create `GetConnectionCounts() (map[int]int, error)` that runs `lsof -iTCP -P -n -sTCP:ESTABLISHED -F n` and counts connections per local port. Parse the `n` field to extract local port, build map[localPort]count.
  - **CRITICAL**: The existing `parseLsofOutput` does NOT hardcode Protocol — it reads the `P` field from lsof output (line 60: `case 'P': currentProtocol = value`). So for UDP, lsof will naturally return `P:UDP`. No modification to the parser needed — just use it directly.
  - **CRITICAL for connection count**: lsof ESTABLISHED output includes BOTH local and remote addresses. The `n` field shows the connection pair. Need to parse carefully: established connections show `n192.168.1.1:3000->10.0.0.1:54321`. Extract the LOCAL port (before `->`) and count.

  **Must NOT do**:
  - Do NOT modify existing `GetListeningPorts()` function
  - Do NOT modify existing `parseLsofOutput()` function  
  - Do NOT add any model or TUI logic

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Parser additions in one file, following existing patterns, needs careful lsof knowledge
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2, 3, 5)
  - **Blocks**: Tasks 6, 7
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `proc_darwin.go:13-31` — Existing `GetListeningPorts()` pattern: exec lsof, check error, parse, sort
  - `proc_darwin.go:33-91` — Existing `parseLsofOutput()`: field-prefix parsing (p=PID, c=command, t=type, P=protocol, n=name)
  - `proc_darwin.go:14` — lsof command flags: `-iTCP -P -n -sTCP:LISTEN -F pcfnPt`

  **API/Type References**:
  - `proc.go:4-11` — PortInfo struct (will have new fields from T1)

  **External References**:
  - lsof `-iUDP` flag selects UDP sockets
  - lsof `-sTCP:ESTABLISHED` filters to established TCP connections
  - lsof `-F n` with ESTABLISHED shows: `n<local>-><remote>` format

  **WHY Each Reference Matters**:
  - Existing parser already handles the `-F` output format — GetUDPPorts just calls the same parser with different lsof flags
  - Connection count needs different parsing because ESTABLISHED output has `->` separators

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: macOS build succeeds with new functions
    Tool: Bash
    Preconditions: proc_darwin.go modified, proc.go has new fields from T1
    Steps:
      1. Run: go build -o /dev/null ./...
      2. Run: go vet ./...
    Expected Result: Both exit 0
    Failure Indicators: Undefined function or type errors
    Evidence: .sisyphus/evidence/task-4-build.txt

  Scenario: GetUDPPorts returns valid data on macOS
    Tool: Bash
    Preconditions: Built binary exists
    Steps:
      1. Start a UDP listener: nc -u -l 19876 &
      2. Build and run a quick Go test: create a temp main that calls GetUDPPorts() and prints results
      3. Verify port 19876 appears with Protocol "UDP"
      4. Kill the nc process
    Expected Result: Port 19876 appears in output with Protocol UDP
    Failure Indicators: Port not found or Protocol is "TCP"
    Evidence: .sisyphus/evidence/task-4-udp-test.txt

  Scenario: GetConnectionCounts returns valid data on macOS
    Tool: Bash
    Preconditions: Built binary exists
    Steps:
      1. Start a TCP listener: nc -l 19877 &
      2. Connect to it: nc localhost 19877 &
      3. Create a temp main that calls GetConnectionCounts() and prints results
      4. Verify port 19877 has count ≥ 1
      5. Kill nc processes
    Expected Result: map includes key 19877 with value ≥ 1
    Failure Indicators: Port not in map or count is 0
    Evidence: .sisyphus/evidence/task-4-conn-test.txt
  ```

  **Commit**: YES (group with Wave 1)
  - Message: `feat(ports): add macOS UDP and connection count parsers`
  - Files: `proc_darwin.go`
  - Pre-commit: `go build -o /dev/null ./...`

---

- [x] 5. Add New/Disappeared Marker Styles (styles.go)

  **What to do**:
  - Add two new lipgloss styles to `styles.go`:
    - `newPortMarker`: green foreground for `●` marker — `lipgloss.NewStyle().Foreground(lipgloss.Color("42"))` (green)
    - `gonePortMarker`: red foreground for `○` marker — `lipgloss.NewStyle().Foreground(lipgloss.Color("196"))` (red)
  - These will be used in `portsToRows()` to render the STATUS column content

  **Must NOT do**:
  - Do NOT add per-row styling — only marker character styles
  - Do NOT modify existing styles
  - Do NOT add more than 2 new styles

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 2 lines added to an existing file
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2, 3, 4)
  - **Blocks**: Task 6
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `styles.go:5-32` — Existing style definitions. Follow exact same `var ( ... )` block pattern with `lipgloss.NewStyle()` chained methods.

  **WHY Each Reference Matters**:
  - Must match the existing style definition pattern exactly — chained builder, var block

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Build succeeds with new styles
    Tool: Bash
    Preconditions: styles.go modified
    Steps:
      1. Run: go build -o /dev/null ./...
    Expected Result: Exit 0
    Failure Indicators: Compilation error
    Evidence: .sisyphus/evidence/task-5-build.txt

  Scenario: Two new marker styles exist
    Tool: Bash
    Preconditions: styles.go modified
    Steps:
      1. Run: grep -c 'PortMarker' styles.go
    Expected Result: Output is "2" (newPortMarker + gonePortMarker)
    Failure Indicators: Count is not 2
    Evidence: .sisyphus/evidence/task-5-style-check.txt
  ```

  **Commit**: YES (group with Wave 1)
  - Message: `feat(ports): add new/disappeared port marker styles`
  - Files: `styles.go`
  - Pre-commit: `go build -o /dev/null ./...`

- [x] 6. Model Integration Batch 1 — New Columns, Service Names, Dedup Toggle, Sort Cycling, Status Markers, UDP Toggle

  **What to do**:
  This is the largest task — integrating 6 features into model.go. Follow this precise order:

  **Step 1: Expand model struct** (add new fields):
  ```go
  type model struct {
      // ... existing fields ...
      prevPorts     []PortInfo     // previous refresh for diff tracking
      sortMode      int            // 0=Port↑, 1=Port↓, 2=PID, 3=Process A-Z
      dedupEnabled  bool           // merge IPv4+IPv6 rows
      protoFilter   int            // 0=TCP, 1=UDP, 2=Both
      statusCounter int            // increment on each status msg to prevent stale clears
  }
  ```

  **Step 2: Add 2 new columns to table** in `NewModel()`:
  - Insert `{Title: "STATUS", Width: 3}` as FIRST column (before PORT)
  - Insert `{Title: "SERVICE", Width: 10}` after TYPE column
  - Full column order: STATUS, PORT, PID, PROCESS, PROTO, ADDRESS, TYPE, SERVICE
  - **CRITICAL**: This changes all column indices! Update `portsToRows()`, Kill handler (selectedRow[1]→[2] for PID), Copy handler (selectedRow[0]→[1] for port), etc.

  **Step 3: Update `portsToRows()`**:
  - Add STATUS column: render `●` (using `newPortMarker.Render("●")`) for `p.Status == "new"`, `○` (using `gonePortMarker.Render("○")`) for `p.Status == "gone"`, empty string otherwise
  - **WAIT — ANSI in cells breaks table**: Actually, DON'T use lipgloss Render in cell values. The bubbles/table truncation uses `runewidth.Truncate()` which counts ANSI escape codes as visible characters. Instead, use plain Unicode characters: just `"●"` and `"○"` without any styling. The characters themselves are visually distinct enough.
  - Add SERVICE column: call `serviceName(p.Port)` from services.go

  **Step 4: Port change tracking** in `portsMsg` handler:
  - Before updating `m.allPorts`, compare new ports with `m.prevPorts` to detect new/disappeared:
    - Build set of `(Port, PID)` from previous and current
    - New: in current but not in previous → set `Status = "new"`
    - Gone: in previous but not in current → create phantom PortInfo entries with `Status = "gone"`, append to current list
    - Neither: clear `Status = ""`
  - Store current (pre-diff) ports as `m.prevPorts` for next cycle
  - **Edge case**: On first refresh (prevPorts is nil), don't mark anything as new — everything is baseline
  - **Edge case**: Disappeared ports should only persist for ONE refresh cycle, then be removed

  **Step 5: Sorting** in `portsMsg` handler (after change tracking, before table update):
  - Apply sort based on `m.sortMode`:
    - 0: Port ascending (current default)
    - 1: Port descending
    - 2: PID ascending
    - 3: Process name A-Z (case-insensitive)
  - Add `s` key handler in `Update()`:
    - `m.sortMode = (m.sortMode + 1) % 4`
    - Re-sort and re-render table
    - Show status: `"Sort: Port ↑"` / `"Sort: Port ↓"` / `"Sort: PID"` / `"Sort: Process"`
  - Update column headers to show sort indicator: modify the column title in the sorted column to append ` ▲` or ` ▼`
  - **CRITICAL**: Sort must persist through 2-second refresh. Store `sortMode` in model, apply in portsMsg handler.

  **Step 6: Dedup toggle**:
  - Add `m` key handler in `Update()`:
    - `m.dedupEnabled = !m.dedupEnabled`
    - Re-render table
    - Show status: `"Dedup: ON"` / `"Dedup: OFF"`
  - Dedup logic (apply after sorting, before table update):
    - When `dedupEnabled`, group ports by `(Port, PID)` key
    - If a group has both IPv4 and IPv6 entries, merge into one with `Type: "4+6"`, keep first entry's Address
    - If a group has only one entry, keep as-is
  - **Edge case**: Don't dedup "gone" phantom entries with live entries
  - **CRITICAL**: Dedup must persist through refresh. Store `dedupEnabled` in model.

  **Step 7: UDP toggle**:
  - Add `t` key handler in `Update()`:
    - `m.protoFilter = (m.protoFilter + 1) % 3` — cycles: TCP → UDP → Both
    - Trigger immediate refresh: `cmds = append(cmds, fetchPortsCmd())`
    - Show status: `"Protocol: TCP"` / `"Protocol: UDP"` / `"Protocol: TCP+UDP"`
  - Modify `fetchPortsCmd()` to accept the proto filter (or read from a package-level var, or pass through a closure):
    - TCP (0): call `GetListeningPorts()` only (existing behavior)
    - UDP (1): call `GetUDPPorts()` only
    - Both (2): call both and merge results
  - **CRITICAL**: `fetchPortsCmd()` currently takes no arguments. Change it to accept the proto filter value. Options:
    - (a) Make it a method on model: NOT possible — Cmd functions can't access model
    - (b) Pass protoFilter as closure capture: `func fetchPortsCmd(proto int) tea.Cmd { return func() tea.Msg { ... } }`
    - Preferred: (b) — pass proto as parameter

  **Step 8: Fix key dispatch pattern**:
  - Currently, ALL key messages fall through to `m.table.Update(msg)` on line 211. This means custom keys like `s`, `m`, `t` also get passed to the table. While none of these currently conflict with table keys, it's wasteful and fragile.
  - Add `return m, tea.Batch(cmds...)` after each new key handler block (Sort, ToggleDedup, ToggleUDP, Open) to prevent fallthrough to table.Update.
  - **ALSO**: The existing Kill, Copy, and Refresh handlers DON'T return early either. Fix those too for consistency.
  - Exception: Keep the fallthrough for navigation keys (the `default` case at the end).

  **Step 9: Fix statusClearMsg race**:
  - Add `statusCounter int` to model. Increment on every status message set.
  - Change `clearStatusCmd()` to capture the current counter value. When `statusClearMsg` fires, only clear if the counter matches (i.e., no newer status was set).
  - This fixes the pre-existing bug where rapid status messages get cleared prematurely.
  - Change `statusClearMsg` to `type statusClearMsg struct{ gen int }` and check `msg.gen == m.statusCounter` before clearing.

  **Step 10: Update status bar**:
  - Show active toggle states in the default (no-message) status bar:
    - Format: `Last refresh: 15:04:05 | 12 ports | [UDP] [dedup] | Sort: Port ↑ | ? help  q quit`
    - Only show `[UDP]` when protoFilter != 0, `[dedup]` when dedupEnabled
  - Update filter status to include toggle indicators too

  **Must NOT do**:
  - Do NOT touch proc_darwin.go, proc_linux.go, keys.go, styles.go — those are done in T1-T5
  - Do NOT implement SIGKILL, open-in-browser, or connection count — those are in T7
  - Do NOT implement --port or --diff — those are in T8/T9
  - Do NOT add CONNECTIONS column yet — that's in T7
  - Do NOT use Bubbletea v2 APIs (tea.KeyPressMsg, charm.land/bubbletea/v2)
  - Do NOT use lipgloss.Render() inside table cell values (breaks runewidth truncation)
  - Do NOT add confirmation dialogs for any action

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Touches 10 interconnected areas of model.go, requires careful index management after column reorder, multiple state machines (sort, dedup, proto), and cross-cutting concerns (status bar, key dispatch). High risk of subtle bugs.
  - **Skills**: `[]`
  - **Skills Evaluated but Omitted**:
    - `playwright`: Not needed — no browser UI
    - `frontend-ui-ux`: Not a web UI task

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (sequential, after Wave 1)
  - **Blocks**: Task 7
  - **Blocked By**: Tasks 1, 2, 3, 4, 5

  **References**:

  **Pattern References**:
  - `model.go:25-38` — Current model struct definition. Add new fields after `showHelp bool`.
  - `model.go:40-73` — `NewModel()` — column definitions and table initialization. Must update columns here.
  - `model.go:95-216` — `Update()` — key handlers and message handlers. Add new key cases after existing ones but BEFORE table.Update fallthrough (line 211).
  - `model.go:105-114` — `portsMsg` handler — where to add change tracking, sorting, dedup logic.
  - `model.go:122-146` — Filter mode key handling — pattern for how to handle modal key states.
  - `model.go:165-189` — Kill handler — uses `selectedRow[1]` for PID. Must update index to `[2]` after adding STATUS column.
  - `model.go:191-209` — Copy handler — uses `selectedRow[0]` for port. Must update index to `[1]`.
  - `model.go:218-272` — `View()` — status bar rendering. Must add toggle indicators.
  - `model.go:253-268` — Status bar text construction — add toggle state display.
  - `model.go:274-287` — `portsToRows()` — must add STATUS and SERVICE columns.
  - `model.go:289-293` — `clearStatusCmd()` — must be updated for generation counter.
  - `model.go:309-323` — `filterPorts()` — extend to search Service field.

  **API/Type References**:
  - `proc.go` — PortInfo struct with new fields (Connections, Service, Status) from T1
  - `services.go` — `serviceName(port int) string` function from T2
  - `styles.go` — `newPortMarker`, `gonePortMarker` styles from T5

  **WHY Each Reference Matters**:
  - Column index changes are the #1 source of bugs — every selectedRow access must be updated
  - The portsMsg handler is the central integration point for sort/dedup/tracking
  - Status bar rendering needs to show toggle states without cluttering
  - filterPorts must be extended or new columns become unsearchable

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: TUI launches with new columns (STATUS, SERVICE)
    Tool: interactive_bash (tmux)
    Preconditions: go build -o ports ./... succeeds
    Steps:
      1. tmux new-session -d -s ports-test './ports'
      2. sleep 2
      3. tmux capture-pane -t ports-test -p > /tmp/ports-output.txt
      4. grep -c "STATUS.*PORT.*PID.*PROCESS.*PROTO.*ADDRESS.*TYPE.*SERVICE" /tmp/ports-output.txt
    Expected Result: Header row contains all 8 column names
    Failure Indicators: Missing columns or wrong order
    Evidence: .sisyphus/evidence/task-6-columns.txt

  Scenario: Sort cycling with 's' key
    Tool: interactive_bash (tmux)
    Preconditions: TUI is running in tmux session ports-test
    Steps:
      1. tmux send-keys -t ports-test 's'
      2. sleep 1
      3. tmux capture-pane -t ports-test -p > /tmp/ports-sort1.txt
      4. grep "Sort:" /tmp/ports-sort1.txt
      5. tmux send-keys -t ports-test 's'
      6. sleep 1
      7. tmux capture-pane -t ports-test -p > /tmp/ports-sort2.txt
      8. grep "Sort:" /tmp/ports-sort2.txt
    Expected Result: First press shows "Sort: Port ↓", second shows "Sort: PID"
    Failure Indicators: No sort indicator in status bar, or sort mode doesn't cycle
    Evidence: .sisyphus/evidence/task-6-sort.txt

  Scenario: Dedup toggle with 'm' key
    Tool: interactive_bash (tmux)
    Preconditions: TUI is running, ports with both IPv4+IPv6 entries visible
    Steps:
      1. tmux capture-pane -t ports-test -p > /tmp/ports-pre-dedup.txt
      2. Count rows with duplicate port numbers
      3. tmux send-keys -t ports-test 'm'
      4. sleep 1
      5. tmux capture-pane -t ports-test -p > /tmp/ports-post-dedup.txt
      6. Verify row count decreased and "Dedup: ON" appears
      7. tmux send-keys -t ports-test 'm'
      8. sleep 1
      9. tmux capture-pane -t ports-test -p > /tmp/ports-dedup-off.txt
      10. Verify "Dedup: OFF" appears and row count restored
    Expected Result: Toggle reduces duplicate rows, shows "4+6" in TYPE column for merged rows
    Failure Indicators: Row count unchanged, or no "4+6" type, or no status message
    Evidence: .sisyphus/evidence/task-6-dedup.txt

  Scenario: UDP toggle with 't' key
    Tool: interactive_bash (tmux)
    Preconditions: TUI is running
    Steps:
      1. tmux send-keys -t ports-test 't'
      2. sleep 2
      3. tmux capture-pane -t ports-test -p > /tmp/ports-udp.txt
      4. grep "Protocol: UDP" /tmp/ports-udp.txt
      5. Verify PROTO column shows "UDP" for all rows
      6. tmux send-keys -t ports-test 't'
      7. sleep 2
      8. tmux capture-pane -t ports-test -p > /tmp/ports-both.txt
      9. grep "Protocol: TCP+UDP" /tmp/ports-both.txt
    Expected Result: Toggle cycles through TCP → UDP → Both with status messages
    Failure Indicators: No protocol change, or PROTO column doesn't update
    Evidence: .sisyphus/evidence/task-6-udp-toggle.txt

  Scenario: New port detection with STATUS column marker
    Tool: interactive_bash (tmux)
    Preconditions: TUI is running
    Steps:
      1. Start a new listener: nc -l 19876 &
      2. Wait for refresh: sleep 3
      3. tmux capture-pane -t ports-test -p > /tmp/ports-new.txt
      4. grep "19876" /tmp/ports-new.txt
      5. Check for ● marker in the STATUS column of that row
    Expected Result: Port 19876 row has ● marker in STATUS column
    Failure Indicators: No marker, or marker in wrong column
    Evidence: .sisyphus/evidence/task-6-new-port.txt

  Scenario: Service name appears for known ports
    Tool: interactive_bash (tmux)
    Preconditions: TUI is running
    Steps:
      1. tmux capture-pane -t ports-test -p > /tmp/ports-services.txt
      2. Check if any known ports (e.g., 53, 80, 443, 8080) have service names in SERVICE column
    Expected Result: Known ports show service names (e.g., "dns", "http", "https", "http-alt")
    Failure Indicators: SERVICE column is empty for known ports
    Evidence: .sisyphus/evidence/task-6-services.txt

  Scenario: Kill still works after column reorder (regression test)
    Tool: interactive_bash (tmux)
    Preconditions: TUI is running
    Steps:
      1. Start a disposable process: nc -l 19877 &
      2. sleep 3
      3. Navigate to the nc row using j/k
      4. tmux send-keys -t ports-test 'x'
      5. sleep 1
      6. tmux capture-pane -t ports-test -p > /tmp/ports-kill-reorder.txt
      7. grep "Killed PID" /tmp/ports-kill-reorder.txt
    Expected Result: "Killed PID ..." status message appears (correct PID extracted from new column index)
    Failure Indicators: Wrong PID killed, or error message, or crash
    Evidence: .sisyphus/evidence/task-6-kill-regression.txt

  Scenario: Filter searches SERVICE column
    Tool: interactive_bash (tmux)
    Preconditions: TUI is running with some known ports
    Steps:
      1. tmux send-keys -t ports-test '/'
      2. tmux send-keys -t ports-test 'http'
      3. sleep 1
      4. tmux capture-pane -t ports-test -p > /tmp/ports-filter-service.txt
      5. Check filtered results include ports with "http" in SERVICE column
    Expected Result: Filter matches service names, not just port/process/address
    Failure Indicators: Filter only matches old columns
    Evidence: .sisyphus/evidence/task-6-filter-service.txt
  ```

  **Commit**: YES
  - Message: `feat(ports): integrate dedup, sort, service names, status markers, UDP toggle`
  - Files: `model.go`
  - Pre-commit: `go build -o ports ./... && go vet ./...`

- [x] 7. Model Integration Batch 2 — SIGKILL, Open-in-Browser, Connection Count, Help Overlay Update
- [x] 8. main.go — --port and --diff Flag Parsing
- [x] 9. Create diff.go — --diff Mode Logic

  **What to do**:
  - Create new file `diff.go` in project root with `package main`
  - Implement `func runDiffMode(portFilter int) int` that:
    1. Gets current listening ports via `GetListeningPorts()` (TCP only — diff mode doesn't need UDP toggle)
    2. If `portFilter > 0`, filters to only that port
    3. Loads previous state from cache file
    4. Compares current vs previous
    5. Prints diff output
    6. Saves current state as new cache
    7. Returns exit code: 0 = no changes, 1 = changes detected

  **Cache file handling**:
  - Use `os.UserCacheDir()` to get platform cache directory
  - Cache path: `{cacheDir}/ports/last.json`
  - Create directory with `os.MkdirAll({cacheDir}/ports, 0755)` if it doesn't exist
  - Cache format: JSON array of `PortInfo` (use `encoding/json` stdlib)
  - **First run** (no cache file): Print all ports as new (`+ port PID process`), save cache, exit 0

  **Diff logic**:
  - Build set of `(Port, PID, Process)` from previous and current
  - New ports (in current, not in previous): print `+ {port}\t{pid}\t{process}\t{address}`
  - Disappeared ports (in previous, not in current): print `- {port}\t{pid}\t{process}\t{address}`
  - Unchanged ports: don't print
  - Sort output: new ports first (sorted by port), then disappeared (sorted by port)

  **Output format**:
  ```
  ports diff (compared to last run at 15:04:05):
  + 3000	12345	node	*:3000
  + 8080	67890	python3	*:8080
  - 5432	11111	postgres	127.0.0.1:5432

  3 changes (2 new, 1 gone)
  ```
  - If no changes: `No changes since last run at 15:04:05`
  - If first run: `First run — baseline saved with N ports`
  - Save timestamp in cache alongside port data (add a wrapper struct)

  **Cache struct**:
  ```go
  type diffCache struct {
      Timestamp time.Time  `json:"timestamp"`
      Ports     []PortInfo `json:"ports"`
  }
  ```

  **Must NOT do**:
  - Do NOT start the TUI — this is CLI-only
  - Do NOT use any third-party JSON library — stdlib `encoding/json`
  - Do NOT add `--format` flag or any output format options
  - Do NOT add color to diff output — plain text with +/- prefixes
  - Do NOT import Bubbletea or lipgloss in diff.go

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: New file with file I/O, JSON marshaling, set comparison, formatted output. Medium complexity but well-scoped.
  - **Skills**: `[]`

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T8)
  - **Parallel Group**: Wave 4 (with Task 8)
  - **Blocks**: F1-F4
  - **Blocked By**: Task 3 (Linux parser for cross-compile), Task 8 (main.go calls runDiffMode)

  **References**:

  **Pattern References**:
  - `proc.go:4-11` — PortInfo struct. Used for JSON serialization in cache. Note: `encoding/json` needs exported fields (which PortInfo already has).
  - `proc_darwin.go:13-31` — `GetListeningPorts()` pattern. diff.go calls this directly.

  **External References**:
  - Go stdlib `encoding/json`: `json.Marshal`, `json.Unmarshal`
  - Go stdlib `os`: `os.UserCacheDir()`, `os.MkdirAll`, `os.ReadFile`, `os.WriteFile`
  - Go stdlib `fmt`: `fmt.Printf` for output

  **WHY Each Reference Matters**:
  - PortInfo struct must be JSON-serializable — all fields are exported, so it works out of the box
  - `os.UserCacheDir()` returns `~/Library/Caches` on macOS, `$XDG_CACHE_HOME` or `~/.cache` on Linux

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: First diff run saves baseline
    Tool: Bash
    Preconditions: No cache file exists (rm -rf "$(go run -e -c 'import "os"; d,_:=os.UserCacheDir(); print(d)')/ports")
    Steps:
      1. Remove cache: rm -rf ~/.cache/ports ~/Library/Caches/ports 2>/dev/null
      2. Run: ./ports --diff 2>&1
      3. Capture output and exit code
    Expected Result: Output contains "First run" or "baseline", exit code 0, cache file created
    Failure Indicators: Error about missing file, or TUI launches
    Evidence: .sisyphus/evidence/task-9-first-run.txt

  Scenario: Second diff run detects changes
    Tool: Bash
    Preconditions: First run completed (cache exists)
    Steps:
      1. Start a new listener: nc -l 19876 &
      2. sleep 1
      3. Run: ./ports --diff 2>&1
      4. Capture output and exit code
      5. Kill nc
    Expected Result: Output shows "+ 19876" line, exit code 1 (changes detected)
    Failure Indicators: No diff output, or exit code 0
    Evidence: .sisyphus/evidence/task-9-with-changes.txt

  Scenario: Diff with no changes
    Tool: Bash
    Preconditions: Two consecutive runs with same state
    Steps:
      1. Run: ./ports --diff (establish baseline)
      2. Run: ./ports --diff (same state)
      3. Capture output and exit code
    Expected Result: Output shows "No changes", exit code 0
    Failure Indicators: False positives showing changes
    Evidence: .sisyphus/evidence/task-9-no-changes.txt

  Scenario: --diff --port filters diff output
    Tool: Bash
    Preconditions: Cache exists, new listener on port 19876
    Steps:
      1. Start listener: nc -l 19876 &
      2. Start another: nc -l 19877 &
      3. Run: ./ports --diff --port 19876
      4. Capture output
      5. Kill both listeners
    Expected Result: Only port 19876 appears in diff, not 19877
    Failure Indicators: Port 19877 appears in output
    Evidence: .sisyphus/evidence/task-9-port-filter.txt

  Scenario: Build succeeds with diff.go
    Tool: Bash
    Steps:
      1. Run: go build -o ports ./...
      2. Run: go vet ./...
      3. Run: GOOS=linux go build -o /dev/null ./...
    Expected Result: All exit 0
    Evidence: .sisyphus/evidence/task-9-build.txt
  ```

  **Commit**: YES (group with Wave 4)
  - Message: `feat(ports): add --diff mode with cache-based change detection`
  - Files: `diff.go`
  - Pre-commit: `go build -o ports ./... && go vet ./...`

---

## Final Verification Wave (MANDATORY — after ALL implementation tasks)

> 4 review agents run in PARALLEL. ALL must APPROVE. Rejection → fix → re-run.

- [x] F1. **Plan Compliance Audit** — `oracle`
- [x] F2. **Code Quality Review** — `unspecified-high`
- [x] F3. **Real Manual QA** — `unspecified-high`
- [x] F4. **Scope Fidelity Check** — `deep`
  For each task: read "What to do", read actual diff (`git log --oneline`). Verify 1:1 — everything in spec was built (no missing), nothing beyond spec was built (no creep). Check "Must NOT do" compliance. Detect cross-task contamination: Task N touching Task M's files. Flag unaccounted changes.
  Output: `Tasks [N/N compliant] | Contamination [CLEAN/N issues] | Unaccounted [CLEAN/N files] | VERDICT`

---

## Commit Strategy

- **Wave 1**: `feat(ports): expand struct, keys, services, parsers, styles` — proc.go, keys.go, services.go, proc_darwin.go, proc_linux.go, styles.go
- **Wave 2**: `feat(ports): integrate dedup, sort, service, status markers, UDP toggle` — model.go
- **Wave 3**: `feat(ports): add SIGKILL, open-browser, conn count, help update` — model.go
- **Wave 4**: `feat(ports): add --port and --diff CLI flags` — main.go, diff.go

---

## Success Criteria

### Verification Commands
```bash
go build -o ports ./...    # Expected: clean build, exit 0
go vet ./...               # Expected: no issues
GOOS=linux go build -o /dev/null ./...  # Expected: cross-compile success
./ports                    # Expected: TUI with 8 columns (PORT, PID, PROCESS, PROTO, ADDRESS, TYPE, SERVICE, CONNS)
./ports --port 3000        # Expected: TUI with pre-filtered to port 3000
./ports --diff             # Expected: one-shot CLI output with +/- markers, exit 0 or 1
```

### Final Checklist
- [ ] All 10 features present and functional
- [ ] All "Must NOT Have" patterns absent from codebase
- [ ] Help overlay shows all key bindings (x, X, /, s, o, t, m, c, r, ?, q)
- [ ] Dedup toggle works: `m` merges IPv4+IPv6 by (Port, PID), shows Type "4+6"
- [ ] Sort cycling works: `s` cycles Port↑ → Port↓ → PID → Process A-Z with indicator
- [ ] `--diff` saves to `os.UserCacheDir()/ports/last.json` and prints changes
- [ ] STATUS column shows `●` for new, `○` for disappeared ports
- [ ] SIGKILL (`X`) and SIGTERM (`x`) both work with correct status messages
- [ ] UDP toggle (`t`) cycles TCP → UDP → Both
- [ ] Connection count shows integer for TCP, `-` for UDP

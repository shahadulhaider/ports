# Restructure ports for Open Source Distribution

## TL;DR

> **Quick Summary**: Reorganize the flat `package main` Go TUI project into a proper `cmd/` + `internal/` layout, update the module path to `github.com/shahadulhaider/ports`, and add OSS distribution files (LICENSE, README, CONTRIBUTING, Makefile, .gitignore).
> 
> **Deliverables**:
> - Restructured codebase: `cmd/ports/main.go` + `internal/{tui,scanner,diff}/`
> - Module path: `github.com/shahadulhaider/ports`
> - OSS files: .gitignore, LICENSE (GPL v3), README.md, CONTRIBUTING.md, Makefile
> - GitHub repo created at `github.com/shahadulhaider/ports`
> 
> **Estimated Effort**: Medium
> **Parallel Execution**: NO — 2 sequential tasks (restructure must pass before OSS files)
> **Critical Path**: Task 1 (restructure) → Task 2 (OSS files) → Task 3 (GitHub repo + push)

---

## Context

### Original Request
User wants to reorganize the flat root-level Go files into a proper project structure and prepare the project for open source distribution. Also wants .gitignore and proper OSS files.

### Interview Summary
**Key Discussions**:
- **Module path**: `github.com/shahadulhaider/ports` (GitHub user: shahadulhaider)
- **License**: GPL v3
- **Layout**: `cmd/ports/main.go` pattern (standard Go convention)
- **Package split**: By concern — `internal/tui/`, `internal/scanner/`, `internal/diff/`
- **OSS files**: README.md, CONTRIBUTING.md, Makefile (NO GitHub Actions, NO goreleaser, NO issue templates)
- **Release automation**: Not now — later

### Metis Review
**Identified Gaps** (addressed):
- **`NewModel()` return type**: Must change from `model` (unexported) to `tea.Model` interface — otherwise `cmd/ports/main.go` can't use it. Resolved: return `tea.Model`.
- **Symbols needing export**: `serviceName()` → `ServiceName()`, `runDiffMode()` → `RunDiffMode()`
- **Dead code preservation**: `GetConnectionCounts()`, `headerStyle`, `newPortMarker`, `gonePortMarker` — all unused but MUST be preserved (no behavior changes)
- **No Windows support**: `syscall.Kill` is Unix-only. Makefile targets darwin + linux only.
- **No git remote**: `go install` won't work until repo is pushed. README includes both install methods.
- **Cache file compatibility**: `PortInfo` has no JSON struct tags, uses Go default names. Moving to `scanner.PortInfo` doesn't change field names → existing `~/.cache/ports/last.json` files stay compatible.

---

## Work Objectives

### Core Objective
Restructure the flat `package main` codebase into a proper Go project layout with `cmd/` + `internal/` packages, and add open source distribution files.

### Concrete Deliverables
- `cmd/ports/main.go` — entry point
- `internal/tui/` — model.go, keys.go, styles.go (package `tui`)
- `internal/scanner/` — proc.go, proc_darwin.go, proc_linux.go, services.go (package `scanner`)
- `internal/diff/` — diff.go (package `diff`)
- `.gitignore`, `LICENSE`, `README.md`, `CONTRIBUTING.md`, `Makefile`
- GitHub repository created and code pushed

### Definition of Done
- [ ] `go build ./cmd/ports/` exits 0
- [ ] `go vet ./...` exits 0
- [ ] Cross-compile all 4 targets (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64) exits 0
- [ ] No `.go` files remain at project root (only go.mod, go.sum)
- [ ] All OSS files present and correct
- [ ] `make build` succeeds
- [ ] GitHub repo exists at github.com/shahadulhaider/ports

### Must Have
- Module path: `github.com/shahadulhaider/ports`
- `NewModel()` returns `tea.Model` (not exported struct)
- `ServiceName()` exported (was `serviceName()`)
- `RunDiffMode()` exported (was `runDiffMode()`)
- GPL v3 license with correct copyright
- .gitignore includes `.sisyphus/`, binary, `.DS_Store`
- README with install instructions and keybindings

### Must NOT Have (Guardrails)
- MUST NOT clean up dead code (`GetConnectionCounts`, `headerStyle`, `newPortMarker`, `gonePortMarker`)
- MUST NOT change any logic or behavior — only package declarations, imports, symbol capitalization
- MUST NOT split or merge files beyond the defined moves (model.go stays as one file)
- MUST NOT add build tags to files that don't have them
- MUST NOT add JSON struct tags to `PortInfo` (breaks existing cache files)
- MUST NOT create Windows build targets
- MUST NOT add tests, goreleaser, GitHub Actions, or issue templates
- MUST NOT refactor internal logic while restructuring (no "while I'm here" improvements)
- MUST NOT modify go.sum manually — only via `go mod tidy`

---

## Verification Strategy (MANDATORY)

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed. No exceptions.

### Test Decision
- **Infrastructure exists**: NO
- **Automated tests**: None
- **Framework**: None
- **Verification**: Build, vet, cross-compile, file existence checks

### QA Policy
Every task includes agent-executed QA scenarios.
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **Build/Compile**: Use Bash — `go build`, `go vet`, cross-compile
- **File structure**: Use Bash — `test -f`, `test ! -f`, `grep`
- **Makefile**: Use Bash — `make build`, `make vet`, `make clean`

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start immediately — code restructure):
└── Task 1: Restructure Go code into cmd/internal layout [deep]

Wave 2 (After Wave 1 — OSS files):
└── Task 2: Create OSS distribution files [writing]

Wave 3 (After Wave 2 — GitHub):
└── Task 3: Create GitHub repo + push [quick]

Critical Path: Task 1 → Task 2 → Task 3
No parallelization opportunity (sequential dependency chain)
```

### Dependency Matrix

| Task | Depends On | Blocks | Wave |
|------|-----------|--------|------|
| 1 | — | 2, 3 | 1 |
| 2 | 1 | 3 | 2 |
| 3 | 2 | — | 3 |

### Agent Dispatch Summary

- **Wave 1**: 1 task — T1 → `deep`
- **Wave 2**: 1 task — T2 → `writing`
- **Wave 3**: 1 task — T3 → `quick`

---

## TODOs

- [x] 1. Restructure Go code into cmd/ + internal/ layout

  **What to do**:
  1. Create directories: `cmd/ports/`, `internal/tui/`, `internal/scanner/`, `internal/diff/`
  2. Move `proc.go`, `proc_darwin.go`, `proc_linux.go`, `services.go` → `internal/scanner/`
     - Change `package main` → `package scanner` in each file
     - In `services.go`: rename `serviceName()` → `ServiceName()` and update the sole call site within `serviceNames` — wait, `serviceNames` is the map, `serviceName()` is the function. Only the function name changes.
     - Preserve `//go:build darwin` and `//go:build linux` tags exactly as-is
  3. Move `model.go`, `keys.go`, `styles.go` → `internal/tui/`
     - Change `package main` → `package tui` in each file
     - In `model.go`:
       - Change `NewModel()` return type from `model` to `tea.Model`
       - Add import `"github.com/shahadulhaider/ports/internal/scanner"`
       - Replace all `PortInfo` → `scanner.PortInfo` (in type annotations, function signatures, variable declarations)
       - Replace `GetListeningPorts()` → `scanner.GetListeningPorts()`
       - Replace `GetUDPPorts()` → `scanner.GetUDPPorts()`
       - Replace `serviceName(` → `scanner.ServiceName(`
       - The `portsMsg` type becomes `type portsMsg []scanner.PortInfo`
       - The `fetchPortsCmd` function's internal calls change to `scanner.GetListeningPorts()` and `scanner.GetUDPPorts()`
     - `keys.go` and `styles.go`: only package declaration changes (no cross-package refs)
  4. Move `diff.go` → `internal/diff/`
     - Change `package main` → `package diff`
     - Rename `runDiffMode()` → `RunDiffMode()`
     - Add import `"github.com/shahadulhaider/ports/internal/scanner"`
     - Replace all `PortInfo` → `scanner.PortInfo`
     - Replace `GetListeningPorts()` → `scanner.GetListeningPorts()`
  5. Create new `cmd/ports/main.go`:
     ```go
     package main

     import (
         "flag"
         "fmt"
         "os"

         tea "github.com/charmbracelet/bubbletea"
         "github.com/shahadulhaider/ports/internal/diff"
         "github.com/shahadulhaider/ports/internal/tui"
     )

     func main() {
         var portFlag int
         var diffFlag bool
         flag.IntVar(&portFlag, "port", 0, "pre-filter to specific port number")
         flag.BoolVar(&diffFlag, "diff", false, "show changes since last run and exit")
         flag.Parse()

         if diffFlag {
             os.Exit(diff.RunDiffMode(portFlag))
         }

         p := tea.NewProgram(tui.NewModel(portFlag), tea.WithAltScreen())
         if _, err := p.Run(); err != nil {
             fmt.Fprintf(os.Stderr, "Error: %v\n", err)
             os.Exit(1)
         }
     }
     ```
  6. Delete original root `.go` files: `main.go`, `model.go`, `keys.go`, `styles.go`, `proc.go`, `proc_darwin.go`, `proc_linux.go`, `services.go`, `diff.go`
  7. Delete compiled `ports` binary from root
  8. Update `go.mod`: change `module ports` → `module github.com/shahadulhaider/ports`
  9. Run `go mod tidy`
  10. Verify: build, vet, cross-compile all 4 targets

  **Must NOT do**:
  - MUST NOT clean up dead code (`GetConnectionCounts`, `headerStyle`, `newPortMarker`, `gonePortMarker`)
  - MUST NOT change any logic — only package declarations, imports, symbol capitalization
  - MUST NOT split or merge files (model.go stays as one file in internal/tui/)
  - MUST NOT add build tags to files that don't have them
  - MUST NOT add JSON struct tags to PortInfo
  - MUST NOT refactor internal logic

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Careful multi-file refactor with interrelated Go package system changes. Must understand import paths, symbol visibility, build tags, and the `tea.Model` interface return type issue. A cheap model will miss subtle compile errors.
  - **Skills**: [`git-master`]
    - `git-master`: Need a clean atomic commit of the restructure
  - **Skills Evaluated but Omitted**:
    - `playwright`: No browser interaction
    - `frontend-ui-ux`: Not a frontend task

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 1 (solo)
  - **Blocks**: Task 2, Task 3
  - **Blocked By**: None

  **References**:

  **Pattern References** (existing code to follow):
  - `model.go:46-89` — `NewModel()` function — must change return type from `model` to `tea.Model`
  - `model.go:22` — `type portsMsg []PortInfo` — becomes `type portsMsg []scanner.PortInfo`
  - `model.go:95-120` — `fetchPortsCmd()` — internal calls change to `scanner.GetListeningPorts()` / `scanner.GetUDPPorts()`
  - `model.go:408-433` — `portsToRows()` — `serviceName(p.Port)` → `scanner.ServiceName(p.Port)`
  - `model.go:449-463` — `filterPorts()` — `serviceName(p.Port)` → `scanner.ServiceName(p.Port)`
  - `diff.go:1` — `package main` → `package diff`; `runDiffMode` → `RunDiffMode`
  - `services.go:3` — `serviceName()` → `ServiceName()`
  - `proc_darwin.go:1` — `//go:build darwin` tag must be preserved exactly
  - `proc_linux.go:1` — `//go:build linux` tag must be preserved exactly
  - `main.go:11-27` — Current main.go — reference for the new cmd/ports/main.go

  **API/Type References**:
  - `proc.go:3-13` — `PortInfo` struct — stays in scanner package, already exported
  - `keys.go:5-20` — `keyMap` struct — stays in tui package, unexported (fine, only used by `keys` var)
  - `styles.go:1-38` — All style vars — stay in tui package, unexported (only used by model.go)

  **External References**:
  - `tea.Model` interface: model must satisfy `Init()`, `Update()`, `View()` — current code already does

  **WHY Each Reference Matters**:
  - `model.go:46` — The `NewModel()` return type is the critical change: `model` → `tea.Model`. Without this, `cmd/ports/main.go` can't reference the unexported struct.
  - `model.go:22,95-120,408-433,449-463` — All places where scanner package symbols are referenced. Miss one → compile error.
  - `diff.go` — Two changes: package declaration and function export. Miss the export → `cmd/ports/main.go` can't call it.
  - `services.go` — One rename: `serviceName` → `ServiceName`. Used in model.go → must be exported for cross-package call.
  - Build tags — If accidentally removed, platform-specific files will cause duplicate symbol errors.

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: Build succeeds after restructure
    Tool: Bash
    Preconditions: All files moved, packages updated, go.mod updated, go mod tidy run
    Steps:
      1. Run: go build -o ports ./cmd/ports/
      2. Assert exit code 0
      3. Run: go vet ./...
      4. Assert exit code 0
    Expected Result: Both commands exit 0
    Failure Indicators: Any compile error or vet warning
    Evidence: .sisyphus/evidence/task-1-build.txt

  Scenario: Cross-compilation works for all 4 targets
    Tool: Bash
    Preconditions: Build succeeds locally
    Steps:
      1. GOOS=darwin GOARCH=amd64 go build -o /dev/null ./cmd/ports/
      2. GOOS=darwin GOARCH=arm64 go build -o /dev/null ./cmd/ports/
      3. GOOS=linux GOARCH=amd64 go build -o /dev/null ./cmd/ports/
      4. GOOS=linux GOARCH=arm64 go build -o /dev/null ./cmd/ports/
    Expected Result: All 4 exit 0
    Failure Indicators: Any cross-compile error (likely syscall issues)
    Evidence: .sisyphus/evidence/task-1-xcompile.txt

  Scenario: File structure is correct
    Tool: Bash
    Preconditions: Restructure complete
    Steps:
      1. test ! -f ./model.go && test ! -f ./keys.go && test ! -f ./styles.go
      2. test ! -f ./proc.go && test ! -f ./proc_darwin.go && test ! -f ./proc_linux.go
      3. test ! -f ./services.go && test ! -f ./diff.go && test ! -f ./main.go
      4. test -f ./cmd/ports/main.go
      5. test -f ./internal/tui/model.go && test -f ./internal/tui/keys.go && test -f ./internal/tui/styles.go
      6. test -f ./internal/scanner/proc.go && test -f ./internal/scanner/proc_darwin.go
      7. test -f ./internal/scanner/proc_linux.go && test -f ./internal/scanner/services.go
      8. test -f ./internal/diff/diff.go
    Expected Result: All tests exit 0
    Failure Indicators: Any file missing or old file remaining
    Evidence: .sisyphus/evidence/task-1-structure.txt

  Scenario: Package declarations and module path correct
    Tool: Bash
    Steps:
      1. grep -q 'module github.com/shahadulhaider/ports' go.mod
      2. head -1 ./internal/scanner/proc.go | grep -q "package scanner"
      3. head -1 ./internal/tui/model.go | grep -q "package tui"
      4. head -1 ./internal/diff/diff.go | grep -q "package diff"
      5. head -1 ./cmd/ports/main.go | grep -q "package main"
    Expected Result: All greps exit 0
    Failure Indicators: Wrong package name or module path
    Evidence: .sisyphus/evidence/task-1-packages.txt

  Scenario: Binary runs correctly (smoke test)
    Tool: Bash
    Steps:
      1. ./ports --help 2>&1
      2. Assert output contains "port" and "diff" flags
    Expected Result: Help output shows both flags, no panic
    Evidence: .sisyphus/evidence/task-1-smoke.txt
  ```

  **Evidence to Capture:**
  - [ ] task-1-build.txt — build + vet output
  - [ ] task-1-xcompile.txt — cross-compile results
  - [ ] task-1-structure.txt — file existence checks
  - [ ] task-1-packages.txt — package/module verification
  - [ ] task-1-smoke.txt — binary --help output

  **Commit**: YES
  - Message: `refactor: restructure project into cmd/ + internal/ layout`
  - Files: All moved/deleted/created Go files, go.mod, go.sum
  - Pre-commit: `go build ./cmd/ports/ && go vet ./...`

- [ ] 2. Create OSS distribution files

  **What to do**:
  1. Create `.gitignore`:
     ```
     # Binary
     ports
     /cmd/ports/ports

     # Build output
     dist/

     # Sisyphus (planning/orchestration)
     .sisyphus/

     # OS
     .DS_Store
     Thumbs.db

     # IDE
     .idea/
     .vscode/
     *.swp
     *.swo
     *~
     ```
  2. Create `LICENSE`: Full GPL v3 text. Copyright line: `Copyright (C) 2025 Shahadul Haider`
  3. Create `README.md`:
     - Project name + one-line description ("A fast, keyboard-driven TUI for exploring listening ports on macOS and Linux")
     - Feature list (all 10 enhancements)
     - Installation: `go install github.com/shahadulhaider/ports/cmd/ports@latest` AND build from source
     - Usage: `ports`, `ports --port 3000`, `ports --diff`
     - Keybindings table (all keys from help overlay)
     - Requirements: macOS or Linux, Go 1.21+
     - License: GPL v3
  4. Create `CONTRIBUTING.md`:
     - Prerequisites (Go 1.21+, macOS or Linux)
     - Build from source steps
     - `make` targets
     - Project structure overview
     - Coding guidelines (flat internal packages, no Windows, no unit tests currently, Bubbletea v1)
  5. Create `Makefile`:
     ```makefile
     BINARY := ports
     PKG := ./cmd/ports
     VERSION ?= dev

     .PHONY: build install clean vet lint cross

     build:
     	go build -o $(BINARY) $(PKG)

     install:
     	go install $(PKG)

     clean:
     	rm -f $(BINARY)
     	rm -rf dist/

     vet:
     	go vet ./...

     lint:
     	@which staticcheck > /dev/null 2>&1 || { echo "Install: go install honnef.co/go/tools/cmd/staticcheck@latest"; exit 1; }
     	staticcheck ./...

     cross:
     	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64 $(PKG)
     	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY)-darwin-arm64 $(PKG)
     	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64 $(PKG)
     	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY)-linux-arm64 $(PKG)
     ```

  **Must NOT do**:
  - MUST NOT add goreleaser config
  - MUST NOT add GitHub Actions / CI workflows
  - MUST NOT add issue templates
  - MUST NOT add CHANGELOG.md
  - MUST NOT add Windows targets in Makefile

  **Recommended Agent Profile**:
  - **Category**: `writing`
    - Reason: All files are documentation/prose/config. No complex logic.
  - **Skills**: [`git-master`]
    - `git-master`: Clean commit of all OSS files
  - **Skills Evaluated but Omitted**:
    - `playwright`: No browser interaction
    - `frontend-ui-ux`: Not a frontend task

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2 (solo)
  - **Blocks**: Task 3
  - **Blocked By**: Task 1 (README install path must match module path; Makefile must reference `./cmd/ports/`)

  **References**:

  **Pattern References**:
  - `keys.go:22-79` — All keybindings with help text — use for README keybindings table
  - `model.go:388-401` — Help overlay content — use for README keybindings (canonical list)
  - `main.go` (new, in cmd/ports/) — Flag definitions — use for README usage section
  - `go.mod` — Module path — use for README install command

  **External References**:
  - GPL v3 full text: https://www.gnu.org/licenses/gpl-3.0.txt
  - Go project layout: https://github.com/golang-standards/project-layout

  **WHY Each Reference Matters**:
  - `keys.go` + help overlay — Canonical source for the keybindings table in README. Don't invent; copy from source.
  - `go.mod` module path — The `go install` command in README must exactly match the module path.

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: All OSS files exist and are correct
    Tool: Bash
    Steps:
      1. test -f .gitignore && test -f LICENSE && test -f README.md && test -f CONTRIBUTING.md && test -f Makefile
      2. grep -q "GNU GENERAL PUBLIC LICENSE" LICENSE
      3. grep -q "Version 3" LICENSE
      4. grep -q "go install github.com/shahadulhaider/ports/cmd/ports" README.md
      5. grep -q "ports" .gitignore
      6. grep -q ".sisyphus" .gitignore
    Expected Result: All checks pass
    Evidence: .sisyphus/evidence/task-2-oss-files.txt

  Scenario: Makefile targets work
    Tool: Bash
    Steps:
      1. make build — assert exit 0 and binary exists
      2. make vet — assert exit 0
      3. make clean — assert exit 0 and binary removed
    Expected Result: All make targets succeed
    Evidence: .sisyphus/evidence/task-2-makefile.txt
  ```

  **Evidence to Capture:**
  - [ ] task-2-oss-files.txt — file existence and content checks
  - [ ] task-2-makefile.txt — make target results

  **Commit**: YES
  - Message: `chore: add OSS distribution files`
  - Files: .gitignore, LICENSE, README.md, CONTRIBUTING.md, Makefile
  - Pre-commit: `make build && make vet`

- [ ] 3. Create GitHub repo and push

  **What to do**:
  1. Create a public GitHub repository using `gh repo create shahadulhaider/ports --public --description "A fast, keyboard-driven TUI for exploring listening ports on macOS and Linux" --source .`
  2. Push all code to the repository
  3. Verify the repo is accessible

  **Must NOT do**:
  - MUST NOT enable GitHub Actions
  - MUST NOT create releases
  - MUST NOT add topics/tags (user can do this later)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Single command to create repo and push. Trivial.
  - **Skills**: [`git-master`]
    - `git-master`: Push to remote
  - **Skills Evaluated but Omitted**:
    - All others: Not relevant

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3 (solo)
  - **Blocks**: None
  - **Blocked By**: Task 2 (need all OSS files committed before pushing)

  **References**:
  - `go.mod` — Module path must match GitHub URL

  **Acceptance Criteria**:

  **QA Scenarios (MANDATORY):**

  ```
  Scenario: GitHub repo exists and is accessible
    Tool: Bash
    Steps:
      1. gh repo view shahadulhaider/ports --json url -q '.url'
      2. Assert output contains "github.com/shahadulhaider/ports"
    Expected Result: Repo URL returned
    Evidence: .sisyphus/evidence/task-3-github.txt
  ```

  **Commit**: NO (push only)

---

## Final Verification Wave

> After ALL tasks complete, run these verification commands:

```bash
# Full build
go build -o ports ./cmd/ports/ && echo "BUILD OK"

# Vet
go vet ./... && echo "VET OK"

# Cross-compile (4 targets)
GOOS=darwin GOARCH=amd64 go build -o /dev/null ./cmd/ports/ && echo "XCOMPILE darwin/amd64 OK"
GOOS=darwin GOARCH=arm64 go build -o /dev/null ./cmd/ports/ && echo "XCOMPILE darwin/arm64 OK"
GOOS=linux GOARCH=amd64 go build -o /dev/null ./cmd/ports/ && echo "XCOMPILE linux/amd64 OK"
GOOS=linux GOARCH=arm64 go build -o /dev/null ./cmd/ports/ && echo "XCOMPILE linux/arm64 OK"

# Module is clean
go mod tidy && git diff --exit-code go.mod go.sum && echo "MOD TIDY OK"

# No old files at root
test ! -f ./model.go && test ! -f ./keys.go && test ! -f ./styles.go && \
test ! -f ./proc.go && test ! -f ./proc_darwin.go && test ! -f ./proc_linux.go && \
test ! -f ./services.go && test ! -f ./diff.go && echo "OLD FILES REMOVED OK"

# New structure exists
test -f ./cmd/ports/main.go && \
test -f ./internal/tui/model.go && test -f ./internal/tui/keys.go && test -f ./internal/tui/styles.go && \
test -f ./internal/scanner/proc.go && test -f ./internal/scanner/proc_darwin.go && \
test -f ./internal/scanner/proc_linux.go && test -f ./internal/scanner/services.go && \
test -f ./internal/diff/diff.go && echo "NEW STRUCTURE OK"

# OSS files exist
test -f .gitignore && test -f LICENSE && test -f README.md && test -f CONTRIBUTING.md && test -f Makefile && echo "OSS FILES OK"

# Module path correct
grep -q 'module github.com/shahadulhaider/ports' go.mod && echo "MODULE PATH OK"

# Package declarations correct
head -1 ./internal/scanner/proc.go | grep -q "package scanner" && echo "SCANNER PKG OK"
head -1 ./internal/tui/model.go | grep -q "package tui" && echo "TUI PKG OK"
head -1 ./internal/diff/diff.go | grep -q "package diff" && echo "DIFF PKG OK"
head -1 ./cmd/ports/main.go | grep -q "package main" && echo "CMD PKG OK"

# Makefile works
make build && echo "MAKEFILE BUILD OK"
make vet && echo "MAKEFILE VET OK"
make clean && echo "MAKEFILE CLEAN OK"

# Smoke test
./ports --help 2>&1 || true  # verify it doesn't panic

# GitHub repo exists
gh repo view shahadulhaider/ports --json url -q '.url' && echo "GITHUB REPO OK"
```

---

## Commit Strategy

- **Commit 1**: `refactor: restructure project into cmd/ + internal/ layout` — all Go file moves, package changes, import updates, go.mod
- **Commit 2**: `chore: add OSS distribution files` — .gitignore, LICENSE, README.md, CONTRIBUTING.md, Makefile
- **Commit 3**: (implicit) — push to GitHub

---

## Success Criteria

### Verification Commands
```bash
go build ./cmd/ports/          # Expected: exit 0
go vet ./...                   # Expected: exit 0
make build                     # Expected: exit 0, produces binary
./ports --help 2>&1 || true    # Expected: shows flags, no panic
gh repo view shahadulhaider/ports  # Expected: repo exists
```

### Final Checklist
- [ ] All "Must Have" present
- [ ] All "Must NOT Have" absent
- [ ] Build + vet + cross-compile pass
- [ ] OSS files present and correct
- [ ] Makefile targets work
- [ ] GitHub repo created and code pushed

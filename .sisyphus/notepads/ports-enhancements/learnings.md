# Learnings — ports-enhancements

## [2026-03-08] Session Start

### Critical Constraints (from Phase 1)
- **Bubbletea v1 ONLY** — use `tea.KeyMsg`, NOT `tea.KeyPressMsg` (v2 beta)
- **lsof -F machine-readable format** — NOT human-readable table (truncates command names)
- **Exec-based clipboard** — `pbcopy`/`xclip`, NOT `golang.design/x/clipboard` (requires CGO)
- **Kill key**: `x` (lowercase) = SIGTERM, `X` (uppercase) = SIGKILL
- **Go build tags**: `//go:build darwin` / `//go:build linux` (new style, NOT `// +build`)
- **Flat package main** — no interfaces, no sub-packages
- **No unit tests** — QA via tmux only
- **Module path**: `ports` (simple, not full GitHub path)

### Architecture Decisions (from Metis review)
- **Row highlighting**: STATUS column with `●`/`○` text markers — bubbles/table v1 has NO per-row styling API
- **Dedup key**: `m` (merge) — `d` conflicts with bubbles/table HalfPageDown
- **Dedup key**: `(Port, PID)` — different processes on same port should NOT merge
- **Sort persists through refresh** — store sortMode in model struct, apply in portsMsg handler
- **Dedup persists through refresh** — store dedupEnabled in model struct
- **Status clear race** — fix with generation counter (statusCounter int in model)
- **fetchPortsCmd** — must accept proto filter as closure parameter (not method on model)
- **ANSI in table cells BREAKS runewidth.Truncate** — use plain Unicode chars `●`/`○` without lipgloss.Render()
- **Linux Type field** — was empty in proc_linux.go, must detect from address format
- **os.UserCacheDir()** — for cross-platform cache path (not hardcoded ~/.cache)

### Column Order (final, after T6+T7)
STATUS(3), PORT(8), PID(8), PROCESS(25), PROTO(6), ADDRESS(20), TYPE(6), SERVICE(10), CONNS(6)
Total: 9 columns

### Column Index Map (after STATUS column added in T6)
- selectedRow[0] = STATUS
- selectedRow[1] = PORT  ← Copy handler uses this
- selectedRow[2] = PID   ← Kill handler uses this
- selectedRow[3] = PROCESS
- selectedRow[4] = PROTO
- selectedRow[5] = ADDRESS
- selectedRow[6] = TYPE
- selectedRow[7] = SERVICE
- selectedRow[8] = CONNS  ← Added in T7

### lsof Commands
- TCP listeners: `lsof -iTCP -P -n -sTCP:LISTEN -F pcfnPt`
- UDP listeners: `lsof -iUDP -P -n -F pcfnPt` (NO state filter — UDP has no LISTEN state)
- TCP established: `lsof -iTCP -P -n -sTCP:ESTABLISHED -F n`

### ss Commands (Linux)
- TCP listeners: `ss -tlnp`
- UDP listeners: `ss -ulnp`
- TCP established: `ss -tnp state established`

### Known Gotchas
- lsof exit code 1 with empty output = no listeners (not an error)
- tmux `?` key needs quoting: `tmux send-keys -t session '?'`
- Filter + kill interaction: after clearing filter, selected row changes
- `textinput` import pulls `github.com/atotto/clipboard` — needs explicit `go get`
- statusBarStyle has hardcoded Width(80) in styles.go but model.go overrides at render time with `.Width(w)` — this is fine

# Decisions — ports-enhancements

## [2026-03-08] Session Start

### User Decisions
- **Dedup key**: `m` (merge) — user chose this over `d` (conflict) and `D` (uppercase)
- **Row highlighting**: STATUS column with text markers — user chose over custom renderer and status-bar-only
- **Dedup default**: OFF — both IPv4+IPv6 rows shown by default, `m` toggles merge
- **Diff storage**: `os.UserCacheDir()/ports/last.json`

### Auto-Resolved Decisions
- **Dedup merge key**: `(Port, PID)` — different processes on same port don't merge
- **"New" marker persistence**: 1 refresh cycle only
- **Disappeared ports**: phantom rows for 1 cycle, then removed
- **Sort indicator**: in column headers (▲/▼)
- **--diff exit code**: 0 = no changes, 1 = changes detected
- **--diff first run**: show all as new, save baseline, exit 0
- **UDP + connection count**: show `-` for UDP rows in CONNS column
- **Filter searches new columns**: yes, SERVICE and CONNECTIONS included
- **Connection count**: always-on column (no toggle needed)
- **Open browser**: always `http://` (no HTTPS detection)
- **Sort modes**: exactly 4 (Port↑, Port↓, PID, Process A-Z), circular
- **Service name map**: ≤30 entries, hardcoded, no external lookups

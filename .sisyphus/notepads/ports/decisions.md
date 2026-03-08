# Decisions — ports project

## [2026-03-08] Architecture Decisions
- Flat `package main` structure — no sub-packages, no interfaces
- Files: main.go, model.go, proc.go, proc_darwin.go, proc_linux.go, styles.go, keys.go
- Refresh interval: 2 seconds (hardcoded)
- TCP LISTEN only (no UDP)
- IPv4+IPv6 shown as separate rows (simplest for MVP)
- Sort: port number ascending only
- Kill: SIGTERM only, no confirmation dialog

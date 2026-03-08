# Issues — ports project

## [2026-03-08] Known Gotchas
- v2 import trap: AI agents often generate `charm.land/bubbletea/v2` code — verify after every task
- lsof human-readable trap: command names truncated to 9 chars — must use -F format
- `k` key conflict: vim table nav uses j/k — kill must be `x`
- CGO trap: `golang.design/x/clipboard` requires CGO — use exec pbcopy/xclip
- tea.KeyPressMsg is v2 API — use tea.KeyMsg (v1)
- IPv4+IPv6 dedup: same port appears twice in lsof output (once per protocol) — expected behavior
- statusBarStyle has hardcoded Width(80) — Task 4 must override with dynamic terminal width via lipgloss .Width(m.width)

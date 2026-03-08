# ports

A fast, keyboard-driven TUI for exploring listening ports on macOS and Linux.

## Features

- Live auto-refresh every 2 seconds
- Filter ports by number, process name, or service
- Kill processes with SIGTERM (`x`) or SIGKILL (`X`)
- Open ports in browser (`o` → `http://localhost:<port>`)
- Copy port info to clipboard (`c`)
- Sort by port, PID, or process name (`s`)
- Toggle TCP / UDP / both (`t`)
- Merge duplicate IPv4+IPv6 entries (`m`)
- Service name hints (http, postgres, redis, etc.)
- Connection count column
- Status markers: `●` new ports, `○` disappeared ports
- `--port` flag to pre-filter on startup
- `--diff` mode for one-shot CLI diff

## Requirements

- macOS or Linux
- Go 1.21 or newer (for building from source)

## Installation

### go install (recommended)

```bash
go install github.com/shahadulhaider/ports/cmd/ports@latest
```

### Build from source

```bash
git clone https://github.com/shahadulhaider/ports
cd ports
make build
./ports
```

## Usage

```bash
# Launch TUI
ports

# Pre-filter to a specific port
ports --port 3000

# Show changes since last run (non-interactive)
ports --diff
```

## Keybindings

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate rows |
| `/` | Filter ports |
| `Esc` | Clear filter |
| `x` | Kill process (SIGTERM) |
| `X` | Force kill (SIGKILL) |
| `o` | Open in browser |
| `c` | Copy to clipboard |
| `s` | Cycle sort mode |
| `t` | Toggle TCP / UDP / Both |
| `m` | Toggle merge IPv4+IPv6 |
| `r` | Refresh now |
| `?` | Toggle help overlay |
| `q` / `Ctrl+C` | Quit |

## License

GNU General Public License v3.0 — see [LICENSE](LICENSE) for details.

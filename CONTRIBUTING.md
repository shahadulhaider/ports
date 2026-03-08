# Contributing to ports

## Prerequisites

- Go 1.21 or newer
- macOS or Linux (Windows is not supported)
- `git`

## Build from Source

```bash
git clone https://github.com/shahadulhaider/ports
cd ports
make build
./ports
```

## Make Targets

| Target | Description |
|--------|-------------|
| `build` | Build the binary to `./ports` |
| `install` | Install via `go install` |
| `clean` | Remove build artifacts |
| `vet` | Run `go vet ./...` |
| `lint` | Run `staticcheck ./...` |
| `cross` | Cross-compile for darwin+linux |

## Project Structure

```
cmd/ports/        Entry point (main.go)
internal/tui/     Bubbletea TUI model, keys, styles
internal/scanner/ Port scanning (lsof/ss parsers, PortInfo struct)
internal/diff/    --diff mode logic
```

## Coding Guidelines

- **Platform**: macOS and Linux only. No Windows support.
- **Framework**: Bubbletea v1 only (not v2 beta). Use `tea.KeyMsg`, not `tea.KeyPressMsg`.
- **Package layout**: Flat `internal/` packages — no nested sub-packages.
- **Tests**: No unit test suite currently. QA is manual via the TUI.
- **Style**: Follow existing patterns in each package. Keep it simple.

## Submitting Changes

1. Fork the repository
2. Create a feature branch: `git checkout -b my-feature`
3. Make your changes and verify: `make build && make vet`
4. Commit with a clear message
5. Open a pull request

## License

By contributing, you agree your contributions will be licensed under GPL v3.

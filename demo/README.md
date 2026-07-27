# Demo recordings

The GIFs and screenshots in the project README are generated here, so they can be
re-recorded deterministically whenever the UI changes.

## Why a container?

Recording `ports` on a real machine would publish whatever happens to be
listening at the time — personal apps, real PIDs, and an arbitrary set of rows
that changes between takes.

Instead, [`Dockerfile`](Dockerfile) builds an image whose network namespace
contains nothing but a synthetic dev stack (nginx, node, postgres, redis,
mongod, ...) started by [`entrypoint.sh`](entrypoint.sh). Every process is a
copy of [`listener/`](listener/) renamed to the service it stands in for, so
`ss -tlnp` reports the expected name and the scanner is exercised for real.

Nothing is faked at the rendering layer: the table you see is the real TUI
reading real kernel state.

## Requirements

- [VHS](https://github.com/charmbracelet/vhs) (`brew install vhs`)
- Docker
- Optional: `gifsicle` and `pngquant` to compress the output

## Regenerating

Run from the repository root, not from this directory:

```bash
docker build -f demo/Dockerfile -t ports-demo .

for t in demo filter sort kill diff; do
  vhs "demo/tapes/$t.tape"
done
```

Then compress:

```bash
for f in demo/out/*.gif; do
  gifsicle -O3 --lossy=70 --colors 128 "$f" -o "$f.opt" && mv "$f.opt" "$f"
done
for f in demo/out/*.png; do
  pngquant --quality=60-85 --speed 1 --force --output "$f" "$f"
done
```

## Layout

| Path | Purpose |
|------|---------|
| `Dockerfile` | Builds `ports` plus the listener harness into a Debian image |
| `entrypoint.sh` | Starts the synthetic dev stack, then drops into a shell |
| `listener/` | Tiny Go program that binds ports and holds them open |
| `tapes/common.tape` | Shared VHS settings (theme, font, viewport) |
| `tapes/*.tape` | One tape per recording |
| `out/` | Generated GIFs and PNGs referenced by the root README |

## Notes

- `tapes/common.tape` sets the viewport to 1320x600, which fits the full table
  without dead space. If columns are added to the TUI, widen it there.
- Tapes hide the `docker run` step with `Hide`/`Show` so recordings open
  directly on a clean prompt.
- `kill.tape` presses `Down 13` to land on redis-server. If the service list in
  `entrypoint.sh` changes, that offset needs updating.

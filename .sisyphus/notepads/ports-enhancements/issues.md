# Issues — ports-enhancements

## [2026-03-08] Session Start

### Pre-existing Issues (from Phase 1)
- **statusBarStyle hardcoded Width(80)** in styles.go — but model.go overrides at render time. Not a bug.
- **Multiple clearStatusCmd timers** — pre-existing race condition. T6 must fix with generation counter.
- **Linux Type field empty** — proc_linux.go line 98 has `Type: ""`. T3 must fix this.

### Potential Issues to Watch
- **Column index drift**: T6 adds STATUS column (shifts all indices by 1). T7 adds CONNS column (shifts again). Every selectedRow access must be updated.
- **fetchPortsCmd signature change**: T6 changes it to accept proto int. T7 must use the new signature.
- **NewModel signature change**: T8 changes it to accept initialPort int. Must coordinate with T6's changes.
- **Wave 4 dependency**: T8 and T9 can run in parallel, but T9 needs T8's runDiffMode stub to compile.

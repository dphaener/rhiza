# tui-launch-usage-verified

evidence · `01M686HGS9GBRF6PCZ34V043P9` · actor [claude-code-22977427](../actor/claude-code-22977427.md)

## result

PASS

## statement

At 28ce530 ran go test ./cmd/rhiza -run TestScripts/tui-launch-usage -count=1 -v: PASS. The script runs ! exec rhiza tui launch, checks stderr has 'usage: rhiza tui app <name>' and empty stdout, then runs it under sh and checks exit=2 (also for bare rhiza tui). cmd/rhiza/tui.go tuiCommand only dispatches 'app <name>'; anything else prints the usage and returns 2, so no picker opens. make check: fmt, vet and tests pass; self-check fails only on stale renders of this run's own actor records (outside the leaf's diff).

## Links

### Out

- verifies [tui-launch-usage](../capability/rhiza/rex-app-actions/tui-launch-usage.md)

## History

- `01M6489HCR688XKH99YV9E2VGY` Remove the all-apps picker, as decision rex-picker-removed says: rhiza tui launch and the rhiza_launch action written by rhiza rex install go, and the capabilities that describe the picker are retired or revised, so each app opens only from its own key in the rhiza key mode or its own action in Rex's command palette.

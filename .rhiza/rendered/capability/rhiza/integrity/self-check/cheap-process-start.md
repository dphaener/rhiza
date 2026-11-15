# cheap-process-start

capability · `01M7EQF5XD1RT48N8TWXXJAXT3` · parent [self-check](../self-check.md)

## acceptance

go build -trimpath -o $TMPDIR/rhiza-x ./cmd/rhiza && GODEBUG=inittrace=1 $TMPDIR/rhiza-x version 2>&1 | grep runewidth shows its init under 2 ms (2026-11-15 on v0.0.27: 22 ms); go.mod requires github.com/mattn/go-runewidth v0.0.31 // indirect, has no replace directive, and git diff main -- go.mod changes only that line; go mod why -m github.com/mattn/go-runewidth names the same importers (cmd/rhiza, charmbracelet/colorprofile, charmbracelet/x/ansi); go test -trimpath ./internal/tui/... passes, so text widths are unchanged. Failure case: go mod verify passes, so go.sum holds the new version's hashes. Speedup measure: /usr/bin/time -f '%e s' sh -c 'for i in $(seq 200); do $TMPDIR/rhiza-x version >/dev/null; done' against the same loop with a build of 88516d956. Touches go.mod, go.sum.

## statement

A rhiza process starts without go-runewidth's eager table init, so a start that never measures text width costs about 17 ms instead of 39 ms, across the 6,600+ rhiza processes cmd/rhiza's scripts start.

## Links

### Out

- decided-by [runewidth-pin-not-fork](../../../../decision/runewidth-pin-not-fork.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)

### In

- verifies [cheap-process-start-verified](../../../../evidence/cheap-process-start-verified.md)
- verifies [cheap-process-start-qa-1](../../../../evidence/cheap-process-start-qa-1.md)

## History

- `01M7EQ42Y1M8ENZCQ202D761JM` Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).

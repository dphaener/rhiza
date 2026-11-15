# cheap-process-start-verified

evidence · `01M7FKYB5M2FGSWJ043QX25T4Z` · actor [claude-code-d855450e](../actor/claude-code-d855450e.md)

## result

PASS

## statement

Verified at a9a119bcc: go build -trimpath -o $TMPDIR/rhiza-x ./cmd/rhiza exit 0; GODEBUG=inittrace=1 rhiza-x version | grep runewidth shows 'init github.com/mattn/go-runewidth @@2.4 ms, 0.030 ms clock, 48 bytes, 3 allocs' (88516d956 build: 23 ms clock, 56472 bytes); git diff main -- go.mod changes only go-runewidth v0.0.27 -> v0.0.31 // indirect, no replace directive; go.sum swaps only the two runewidth lines; go mod why -m names cmd/rhiza, charmbracelet/colorprofile, charmbracelet/x/ansi; go mod verify: all modules verified; go test -trimpath -p 2 ./internal/tui/... ok (9 packages). 200-start loop /usr/bin/time, alternated under the machine lock: 88516d956 build (detached worktree under $TMPDIR) 7.28 s, 7.40 s; this build 2.20 s, 2.53 s. Diff touches only go.mod, go.sum and the builder's actor record; no stubs or TODOs.

## Links

### Out

- verifies [cheap-process-start](../capability/rhiza/integrity/self-check/cheap-process-start.md)

## History

- `01M7EQ42Y1M8ENZCQ202D761JM` Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).

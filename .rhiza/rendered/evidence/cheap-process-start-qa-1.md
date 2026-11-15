# cheap-process-start-qa-1

evidence · `01M7FN89B4RMA2B2E05RMFKWRS` · actor [claude-code-9160dfa9](../actor/claude-code-9160dfa9.md)

## result

PASS

## statement

QA at c9e38ec85: go build -trimpath -o $TMPDIR/rhiza-x ./cmd/rhiza; GODEBUG=inittrace=1 rhiza-x version shows go-runewidth init 0.023-0.038 ms clock, 48 bytes (88516d956 build: 22 ms clock, 56472 bytes). 200-start loop alternated under the machine lock: old 7.14 s, 7.09 s; new 1.97 s, 2.09 s. go mod why -m names cmd/rhiza, charmbracelet/colorprofile, charmbracelet/x/ansi; go mod verify all modules verified; no replace directive; go test -trimpath -p 2 ./internal/tui/... ok. Old and new binaries give byte-identical output for slice --change (default and --budget 800), status, a bad change list call (exit 2 both) and a missing change set (exit 4, clean message). StringWidth of ASCII, CJK, emoji, flags, box drawing is unchanged between v0.0.27 and v0.0.31 by default; under RUNEWIDTH_EASTASIAN=1 or a ja_JP locale, e-acute goes from 2 to 1 (recorded as a minor finding).

## Links

### Out

- verifies [cheap-process-start](../capability/rhiza/integrity/self-check/cheap-process-start.md)

## History

- `01M7EQ42Y1M8ENZCQ202D761JM` Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).

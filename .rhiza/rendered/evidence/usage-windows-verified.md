# usage-windows-verified

evidence · `01M6J1GQ704YXYBSMV9CXTGTH5` · actor [claude-code-affa7d7b](../actor/claude-code-affa7d7b.md)

## result

PASS

## statement

Verifier claude-code-affa7d7b at db286c40 read the diff and ran go test -count=1 -v ./internal/usage. internal/usage is a pure package: no clock, file or adapter reads. Windows(kind, []adapter.UsageRecord, []Limit) takes the same record type that usage-window-tokens' UsageRecords returns, and Current gives the live window. All 6 tests PASS, each looping over Kinds (Session 5h, Weekly 7d) with fixed times from t0=2026-11-02T08:00Z. TestFirstWindow: start at the first record, reset at start+length, Used=15. TestRecordAfterEndStartsNext: a record at exactly start+length starts window 2, while a record one second earlier, given out of order, stays in window 1 (Used 11/7). TestGapStartsAtNextRecord: after a 3x-length gap the next window starts at the next record, and Current reports no window after reset. TestLimitResetMovesStart: a limit hit moves window 1 to end at the limit's reset and start one length earlier, and the next record starts window 2 at that reset. Also tested: a limit whose reset is not after its hit is ignored, and an empty input or unknown kind gives no windows. In make check, fmt-check, vet and go test ./... pass. self-check, run with RHIZA_ROOT/CHANGE/RUN/ACTOR unset, fails only on stale .rhiza/rendered files (actor/claude-code-5548f8eb.md and index.md) left by the builder's committed actor record. That is outside this leaf's acceptance and is fixed by running rhiza render at close-out.

## Links

### Out

- verifies [usage-windows](../capability/rhiza/orchestration/usage-governor/usage-windows.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate

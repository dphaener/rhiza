# usage-hold-verified-2

evidence · `01M6JK20J5YTHBZ43Z11Z8VVK6` · actor [claude-code-ac9c5478](../actor/claude-code-ac9c5478.md)

## result

PASS

## statement

At c2cd1c1e (review fix cycle 2) ran `go test ./internal/api -run 'TestUsage|TestHealth' -count=1 -v`: TestUsageLine, TestUsageHoldProjection, TestUsageWindowsHitKinds, TestUsageWindowsAnchored, TestUsageHoldStopsLaunches and TestUsageHoldStopsEveryLaunch all PASS, with a fake clock and fake token records. Read the tests against the acceptance: under the line (limit 1000, line 800, use 600) starts; at limit 900 (line 720) use 600 alone is under and one projected run (750) holds, with Used 600; with no finished runs only use so far counts (801 holds); a held daemonPoll1 launches no pick, resume, queued start or continue, or plan request, and leaves a live run's journal untouched and not ended; the poll after the hold clears starts the held change set and launches each queued item; the window has no hold 6h after its reset; a learned limit with no setting holds, a set session limit of 5000 wins over a learned 700, and no setting and no learned limit never holds. Read the diff: windows are now anchored by a saved per-kind window in daemon state, only hits of a window's own kind move it (TestUsageWindowsHitKinds), and records are read once per poll. These fix the confirmed review findings. fmt-check and vet pass. go test ./... fails only on TestReadBudget timings and a plan-workflow testscript timeout on a machine at load 17 with 10 other test runs; plan-workflow passed when re-run alone, and neither touches this diff.

## Links

### Out

- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate

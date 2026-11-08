# health-schedules-part-verified

evidence · `01M6WW1F2VV82H2P1N6ZHFQQEV` · actor [claude-code-4888938c](../actor/claude-code-4888938c.md)

## result

PASS

## statement

At b6d10bd3 (attempt 2) ran go test -trimpath ./internal/api -run TestHealthSchedulesPart -v -count=1: all five subtests pass — no schedules has no schedules part; a not-yet-run schedule is ok with 'none yet' and 'next run 2026-11-08T11:00:00Z'; a failed last run warns with next 'rhiza run log <id>'; a run that moved HEAD warns naming the main checkout; no daemon warns. Each subtest now runs the health op through schedHealth, which requires ExitCode()==0 and the schedules part never bad, on a fixture rendered so rhiza check passes (fixing attempt 1's finding f1). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(describe-golden|health.*)$' -count=1: pass; describe-golden carries the new schedules part enum value. Read the diffs of a2ec10f4 and b6d10bd3: the fix commit changes only the test and stages the builder's own actor record.

## Links

### Out

- verifies [health-schedules-part](../capability/rhiza/workflows/scheduled-triggers/health-schedules-part.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).

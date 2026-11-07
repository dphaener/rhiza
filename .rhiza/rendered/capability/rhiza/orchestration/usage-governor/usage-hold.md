# usage-hold

capability · `01M6HJ10G74HCGZ37HH47P4ABP` · parent [usage-governor](../usage-governor.md)

## acceptance

Tests in internal/api's daemon tests with a fake clock and fake token use show: under the line it starts; one projected run over the line holds picks, resumes, queued work and plans and stops no live run; with no finished runs only use so far counts; a live run adds only the mean less its tokens so far, never below zero; the next poll after the reset starts the held change set; with no live snapshot a learned limit with no setting holds, a set rhiza.daemon.usage.session wins over a learned session limit, and a window with no live limit, no setting and no learned limit never holds; a live limit wins over both the setting and the learned limit (go test -trimpath ./internal/api -run 'TestUsageHold')

## statement

At each poll the daemon holds everything it would launch (picks, resumes of interrupted runs, queued starts and continues, feature and bug plan requests) while a window's projected use exceeds its limit minus the reserve. Projected use is the window's use so far, plus for each live run what it is still expected to burn (the mean tokens of the last 10 finished runs less that run's tokens so far, never below zero; usage-live-run-remaining), plus one full mean for the next. A window's limit is its live limit when a live snapshot gives one (usage-live-limit, which also fixes the window's start and reset), else its setting when set, else its learned limit, and a window with none never holds. Live runs, runs the owner starts directly and milestone requests' plan runs (milestone-request-past-usage-hold) are untouched, and the hold clears at the first poll after the window resets

## Links

### Out

- depends-on [usage-windows](usage-windows.md)
- depends-on [usage-settings](usage-settings.md)
- depends-on [usage-learned-limit](usage-learned-limit.md)
- decided-by [usage-governor-hold-rules](../../../../decision/usage-governor-hold-rules.md)
- decided-by [milestone-interview-starts-at-once](../../../../decision/milestone-interview-starts-at-once.md)
- decided-by [usage-governor-limits-fix-scope](../../../../decision/usage-governor-limits-fix-scope.md)
- decided-by [live-usage-source](../../../../decision/live-usage-source.md)

### In

- depends-on [usage-hold-reported](usage-hold-reported.md)
- verifies [usage-hold-verified](../../../../evidence/usage-hold-verified.md)
- verifies [usage-hold-qa-1](../../../../evidence/usage-hold-qa-1.md)
- verifies [usage-hold-verified-2](../../../../evidence/usage-hold-verified-2.md)
- verifies [usage-hold-qa-2](../../../../evidence/usage-hold-qa-2.md)
- verifies [usage-hold-qa-3](../../../../evidence/usage-hold-qa-3.md)
- verifies [pr-238-usage-governor](../../../../evidence/pr-238-usage-governor.md)
- verifies [pr-291-a-milestone-request](../../../../evidence/pr-291-a-milestone-request.md)
- verifies [usage-hold-01m6t1-qa-1](../../../../evidence/usage-hold-01m6t1-qa-1.md)
- verifies [pr-345-fix-the-usage-governor-s](../../../../evidence/pr-345-fix-the-usage-governor-s.md)
- verifies [pr-349-fix-the-usage-governor-s](../../../../evidence/pr-349-fix-the-usage-governor-s.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
- `01M6T1ZR3ZCDCYF96CSR2AYN7D` Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens
- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available

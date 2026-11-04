# usage-hold

capability · `01M6HJ10G74HCGZ37HH47P4ABP` · parent [usage-governor](../usage-governor.md)

## acceptance

Tests in internal/api's daemon tests with a fake clock and fake token use show: under the line it starts; one projected run over the line holds picks, resumes, queued work and plans and stops no live run; with no finished runs only use so far counts; the next poll after the reset starts the held change set; a learned limit with no setting holds, a set rhiza.daemon.usage.session wins over a learned session limit, and a window with no setting and no learned limit never holds

## statement

At each poll the daemon holds everything it would launch (picks, resumes of interrupted runs, queued starts and continues, plan requests) while a window's use so far plus the mean tokens of the last 10 finished runs, once for each live run and once for the next, exceeds its limit minus the reserve. A window's limit is its setting when set, else its learned limit, and a window with neither never holds. Live runs and runs the owner starts directly are untouched, and the hold clears at the first poll after the window resets

## Links

### Out

- depends-on [usage-windows](usage-windows.md)
- depends-on [usage-settings](usage-settings.md)
- depends-on [usage-learned-limit](usage-learned-limit.md)
- decided-by [usage-governor-hold-rules](../../../../decision/usage-governor-hold-rules.md)

### In

- depends-on [usage-hold-reported](usage-hold-reported.md)
- verifies [usage-hold-verified](../../../../evidence/usage-hold-verified.md)
- verifies [usage-hold-qa-1](../../../../evidence/usage-hold-qa-1.md)
- verifies [usage-hold-verified-2](../../../../evidence/usage-hold-verified-2.md)
- verifies [usage-hold-qa-2](../../../../evidence/usage-hold-qa-2.md)
- verifies [usage-hold-qa-3](../../../../evidence/usage-hold-qa-3.md)
- verifies [pr-238-usage-governor](../../../../evidence/pr-238-usage-governor.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate

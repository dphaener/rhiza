# usage-hold-verified

evidence · `01M6J9RCA147DAD744R86R2MPY` · actor [claude-code-fcb9edd5](../actor/claude-code-fcb9edd5.md)

## result

PASS

## statement

Verified at 1ce8e3c0 (attempt 2). Ran go test ./internal/api -run TestUsageHold -count=1 -v: TestUsageHoldProjection, TestUsageHoldStopsLaunches and TestUsageHoldStopsEveryLaunch PASS. Read the diff of 3cf094cf and 1ce8e3c0. With a fake clock and fake token records (hold-fake adapter), usageHoldOf: no limit never holds; limit 1000/reserve 20% with use 600 does not hold (under the line); with no finished runs, use 801 alone holds (projected 801); mean of the last 10 finished runs (150, the two older 9999 runs excluded, cache reads uncounted) projects 750 under and 900 over once a live run is added; limit 700 holds on one projected run (750); a set session limit 5000 wins over a learned 700; a learned 700 with no setting holds; weekly holds on its own; 6h later (after the reset) nothing holds. daemonPoll1 with a hold launches nothing on a fixture with an interrupted resumable run, a queued start, a queued continue, a plan request and a live run on another host: no Detach, no daemon-resumed entry, the queues are unchanged, the live run gets no entry and does not end, and the target is skipped as usage-reserve. The next poll with the hold cleared resumes both runs, starts the queued start and the plan request (4 detaches) and the held change set's pick starts in TestUsageHoldStopsLaunches. This answers the attempt 1 findings f1 and f2. make fmt-check vet self-check pass, and the 7 actor files rendered in 1ce8e3c0 each have a tracked record. go test ./internal/api -timeout 40m passes except TestReadBudget, and go test ./... also failed slice TestBudget, modelapp TestGeneratedModel and script plan-workflow (338s timeout): all are timing budgets or timeouts in code this leaf does not touch, run under load average 10-13 on 4 cores

## Links

### Out

- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate

# health-orphan-log-entries-verified

evidence · `01M76V5ZJ4JV7R0X32HP5SMJCE` · actor [claude-code-a59a0be0](../actor/claude-code-a59a0be0.md)

## result

PASS

## statement

At e09d453f9, attempt 2: go test -trimpath ./internal/api -run TestHealthMainOrphanLog -v -count=1 passes TestHealthMainOrphanLog (an untracked entry with no branch makes healthMain bad, Next and detail name the .rhiza/log/<cs>/ entry, its change set and rhiza model commit-untracked; a stale close-out keeps it bad instead of softening it to warn; after daemonCommitUntracked puts it on a model branch the part is ok with a pending entry listed) and TestHealthMainOrphanLogLiveRun (an entry a live run's worktree holds is pending, the part is ok, Next is empty, no 'on no branch' detail). go vet ./internal/api is clean; go test -run 'Health|Orphan|Untracked|ModelBranch|ModelPR|Check' ./internal/api passes in 27s. HealthResponse.ExitCode is 1 when any part is bad, so an orphan exits 1. Diff read: classifyUntrackedLog splits orphans from pending with the same branch/worktree logic as orphanLogEntries, and check.go still uses that logic. The whole internal/api package hit go test's default 600s timeout on a host with load ~12 on 4 cores; make test uses -timeout 20m and CI shards the package, so this is environmental and not caused by the change.

## Links

### Out

- verifies [health-orphan-log-entries](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-orphan-log-entries.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

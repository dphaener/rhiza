# health-orphan-log-entries-verified-3

evidence · `01M76WRRAFKC7Z85J17Q3DGMXS` · actor [claude-code-ab186019](../actor/claude-code-ab186019.md)

## result

PASS

## statement

Verifier claude-code-ab186019 at bc6f4805f: go test -trimpath -count=1 ./internal/api -run 'TestHealthMainOrphanLog' passes (TestHealthMainOrphanLog: orphan entry makes main bad with Next 'rhiza model commit-untracked' and detail naming .rhiza/log/<cs>/ and the change set; entry on a model branch leaves main ok with summary 'pending' and a 'pending entry' detail; TestHealthMainOrphanLogLiveRun: entry held by a live run leaves main ok with the pending entry listed, no Next). Review cycle 2 findings checked: f1 run-start-new-main-clean now asserts no 'orphan' plus the pending line, and passes (go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(run-start-new-main-clean|model-commit-untracked|health.*)$', 11 scripts PASS, GOTMPDIR outside the repo and free of the word 'orphan', since $WORK paths appear in health's output); f2 the leaf's statement no longer claims the Inbox, only rhiza health and daemon status. Read healthMain (internal/api/health.go:495-560) and classifyUntrackedLog (internal/api/check.go:618).

## Links

### Out

- verifies [health-orphan-log-entries](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-orphan-log-entries.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

# health-orphan-log-entries-verified-5

evidence · `01M76Z72EWQHGXQ1GXGHEP60EK` · actor [claude-code-40b4bec6](../actor/claude-code-40b4bec6.md)

## result

PASS

## statement

At 769366365, go test -trimpath -count=1 ./internal/api -run 'TestHealthMainOrphanLog' passes (TestHealthMainOrphanLog, TestHealthMainOrphanLogLiveRun). Read the tests: an orphan entry makes healthMain bad with Next 'rhiza model commit-untracked' and detail naming .rhiza/log/<cs>/ and the change set; a stale close-out never softens the orphan's bad; after daemonCommitUntracked the entry is held by a model branch and the part is ok with a 'pending entry' detail; an entry a live run's worktree holds leaves the part ok with the pending entry listed and no 'on no branch'; and the new case (Decision orphan-log-needs-origin-main, the confirmed cycle-2 finding) uses modelBranchFixture(t,false) and asserts an untracked entry with no origin/main leaves the part ok with no Next and no orphan or commit-untracked detail. Testscripts run-start-new-main-clean, health*, check-orphan-log and model-commit-untracked pass (an earlier run-start-new-main-clean failure was my own $WORK path containing 'orphan', gone with a neutral TMPDIR). internal/api -run 'Health|ModelBranch|CommitUntracked|Orphan|Untracked|Check' passes; the only FAIL line was the root-state guard on ~/code/.rhiza-root/usage.live.json, which the live daemon rewrites.

## Links

### Out

- verifies [health-orphan-log-entries](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-orphan-log-entries.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

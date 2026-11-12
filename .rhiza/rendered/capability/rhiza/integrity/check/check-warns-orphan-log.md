# check-warns-orphan-log

capability · `01M76FAABBDGMZA9F3BBBA3H0D` · parent [check](../check.md)

## acceptance

A testscript (cmd/rhiza/testdata/script/check-orphan-log.txtar) leaves an orphan entry in a main checkout and asserts rhiza check exits 0 with a warning naming it, and the same entry in a linked worktree gives no warning

## statement

rhiza check run in a main checkout warns, without failing, about each untracked .rhiza/log entry that no pending model branch or live run holds, naming rhiza model commit-untracked; in a run worktree or any other linked worktree it says nothing about untracked entries

## Links

### Out

- decided-by [decision-main-writes-get-a-branch](../../../../decision/decision-main-writes-get-a-branch.md)

### In

- blocks [review-01m76f247z4y9t4j86ryf70hnt-r3-f2](../../../../question/review-01m76f247z4y9t4j86ryf70hnt-r3-f2.md)
- verifies [check-warns-orphan-log-verified](../../../../evidence/check-warns-orphan-log-verified.md)
- verifies [check-warns-orphan-log-qa-1](../../../../evidence/check-warns-orphan-log-qa-1.md)
- verifies [check-warns-orphan-log-qa-2](../../../../evidence/check-warns-orphan-log-qa-2.md)
- verifies [pr-464-nothing-rhiza-writes-is-ever](../../../../evidence/pr-464-nothing-rhiza-writes-is-ever.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

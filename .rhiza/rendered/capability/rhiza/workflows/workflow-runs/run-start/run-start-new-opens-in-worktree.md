# run-start-new-opens-in-worktree

capability · `01M76FAA1QXNBQDGDWWN6BS1KX` · parent [run-start](../run-start.md)

## acceptance

A testscript (cmd/rhiza/testdata/script/run-start-new-main-clean.txtar) starts a --new plan run and a --milestone plan run and asserts the run's worktree holds and commits the open entry, rhiza change show of the new change set succeeds from the main checkout, and rhiza health and rhiza check report the main copy as held by the run, not an orphan; go test ./internal/api -run 'TestNewPlanChange' covers a dry run still predicting the same change set ID and worktree

## statement

rhiza run start --new (a request's plan run, a milestone record) creates the run's branch and worktree first and writes the new change set's open entry in that worktree, where the run commits it on its branch; it also writes the same bytes into the main checkout's .rhiza/log as a copy held by the live run, so the daemon, roadmap, Inbox and rhiza change show in the main checkout see the planning change set at once, health and check count it as held, close-out drops it, and an abandoned run leaves it to rhiza model commit-untracked

## Links

### Out

- decided-by [decision-main-writes-get-a-branch](../../../../../decision/decision-main-writes-get-a-branch.md)

### In

- verifies [run-start-new-opens-in-worktree-verified](../../../../../evidence/run-start-new-opens-in-worktree-verified.md)
- verifies [run-start-new-opens-in-worktree-qa-1](../../../../../evidence/run-start-new-opens-in-worktree-qa-1.md)
- verifies [run-start-new-opens-in-worktree-verified-2](../../../../../evidence/run-start-new-opens-in-worktree-verified-2.md)
- verifies [run-start-new-opens-in-worktree-qa-2](../../../../../evidence/run-start-new-opens-in-worktree-qa-2.md)
- verifies [pr-464-nothing-rhiza-writes-is-ever](../../../../../evidence/pr-464-nothing-rhiza-writes-is-ever.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

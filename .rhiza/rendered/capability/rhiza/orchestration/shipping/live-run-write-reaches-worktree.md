# live-run-write-reaches-worktree

capability · `01M76FAA3B1FNWA3CMWCB4KQDD` · parent [shipping](../shipping.md)

## acceptance

go test ./internal/api -run 'TestLiveRunWrite' accepts a change set from the main checkout while its run's worktree exists and asserts the entry appears in the worktree's .rhiza/log and is committed by the run's next commit; a write for a change set with no live run is not copied to any worktree

## statement

A model write made in the main checkout (change new, withdraw, accept, propose, reject, apply, link, a gate answer, from the CLI, Inbox or web app) for a change set with a live run is also written into that run's worktree, where the run commits it on its branch at its next commit, so the main checkout's copy is held by the run and dropped by close-out

## Links

### Out

- decided-by [decision-main-writes-get-a-branch](../../../../decision/decision-main-writes-get-a-branch.md)

### In

- verifies [live-run-write-reaches-worktree-verified](../../../../evidence/live-run-write-reaches-worktree-verified.md)
- verifies [live-run-write-reaches-worktree-qa-1](../../../../evidence/live-run-write-reaches-worktree-qa-1.md)
- verifies [pr-464-nothing-rhiza-writes-is-ever](../../../../evidence/pr-464-nothing-rhiza-writes-is-ever.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

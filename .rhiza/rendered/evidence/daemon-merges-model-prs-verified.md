# daemon-merges-model-prs-verified

evidence · `01M76TGCM4FZBXPEBJCTGTGFCS` · actor [claude-code-780c23d1](../actor/claude-code-780c23d1.md)

## result

PASS

## statement

At 0d719dfe4 I ran go test -trimpath -count=1 ./internal/api -run TestDaemonModelPR: PASS for Green (a fake gh's pr merge squash-merges into origin/main; branches deleted, HEAD == origin/main, no untracked .rhiza/log left, entries present, no inbox item), red and conflict (nothing merged, branches kept, an inbox model-pr item names the change set), and GhFailureNeverFails. go test ./internal/daemon -run TestSelfUpdate: all 7 PASS, including TestSelfUpdateDropsIncomingLogCopies (an untracked copy byte-identical to what origin/main adds is dropped, the fast-forward succeeds and the file becomes tracked, and a differing untracked file is kept). go vet and gofmt are clean, and TestScripts/describe-golden passes. In the diff, daemonMergeModelPRs is called from the daemon poll after daemonCommitUntracked; it runs gh pr merge --squash --match-head-commit, then git push --delete, then pullMainAfterModelMerge (dropIncomingLogCopies, then pull --ff-only, and main-stale if blocked); SelfUpdate.Dir is set in both the root and the project paths.

## Links

### Out

- verifies [daemon-merges-model-prs](../capability/rhiza/orchestration/daemon/daemon-merges-model-prs.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

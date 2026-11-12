# daemon-merges-model-prs-verified-5

evidence · `01M76YZVNFZD46AAMQ2ET5PFQW` · actor [claude-code-56a18c26](../actor/claude-code-56a18c26.md)

## result

PASS

## statement

Verifier at e8e289be4 (attempt 5): go test -trimpath -count=1 ./internal/api -run 'TestDaemonModelPR|TestModelNetThrottled|TestBoundedExec' passes: with a fake gh and no run, Green squash-merges the model PR, leaves no model branch, no untracked .rhiza/log copy, HEAD == origin/main with the change set's entries present and no inbox item; red and conflict merge nothing, keep the branches, and inboxModelPRs lists an InboxModelPR item naming the change set; a gh failure never fails the poll. go test -trimpath -count=1 ./internal/daemon -run TestSelfUpdate passes, including TestSelfUpdateDropsIncomingLogCopies. Diff read: the poll now runs daemonCommitUntracked and daemonMergeModelPRs at most once per modelNetEvery (1m) per project, through boundedExecCtx, which kills each git/gh program after 30s with GIT_TERMINAL_PROMPT=0 set over the environment (startProgram's programEnv merges it, so PATH/HOME are kept). Nothing in the merge or commit-untracked path matches on stderr text, so dropping stderr from errors only shortens failure details. This addresses the cycle's confirmed quality finding f1 (untimed, unthrottled network work every 5s). go vet clean on both packages.

## Links

### Out

- verifies [daemon-merges-model-prs](../capability/rhiza/orchestration/daemon/daemon-merges-model-prs.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

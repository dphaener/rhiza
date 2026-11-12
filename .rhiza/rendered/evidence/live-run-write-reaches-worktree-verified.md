# live-run-write-reaches-worktree-verified

evidence · `01M76SF6DF2ZTNFHDCFBNNKXNF` · actor [claude-code-9cc99d10](../actor/claude-code-9cc99d10.md)

## result

PASS

## statement

At a02ef9dd1 ran `go test -trimpath -count=1 ./internal/api -run TestLiveRunWrite -v`: TestLiveRunWriteReachesWorktree, TestLiveRunWriteCommittedByRunsNextCommit and TestLiveRunWriteNoRunNoCopy all PASS. An accept made in the main checkout while run R1's worktree exists is copied byte for byte into the worktree's .rhiza/log together with its uncommitted open and propose, and the worktree's log folds cleanly with the change set accepted (fixes attempt 1's f2). With real git and a linked worktree, the run's next commit (CommitCheckWrites through commitOnly, which now adds mirroredPaths: untracked entries of the run's change set identical to the main checkout's) holds all three entries on the run's branch, HEAD:<path> equal to the main copy, and leaves nothing untracked (fixes f1). A change set with no live run, and a dry run, copy nothing to any worktree. Read the diff (live_run_mirror.go copyMissingEntries/mirroredPaths, run_leaves.go commitOnly); go vet ./internal/api is clean. A wider ./internal/api -run 'Commit|Evidence|Check|Mirror|LiveRun' run passed every test, then failed only the known .rhiza-root/usage.live.json state-folder guard, which this change does not touch.

## Links

### Out

- verifies [live-run-write-reaches-worktree](../capability/rhiza/orchestration/shipping/live-run-write-reaches-worktree.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

# commit-untracked-log-verified-5

evidence · `01M76YHHY0NXA3AMWQF28B7Y9W` · actor [claude-code-49b8d5df](../actor/claude-code-49b8d5df.md)

## result

PASS

## statement

At 5c8ba463b (RHIZA_ROOT/RHIZA_RUN unset, GOTMPDIR outside the repo): go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/model-commit-untracked$' passes. The script uses a local bare origin and a fake gh, leaves two orphan change sets plus one held by a live run's worktree, runs rhiza model commit-untracked, and asserts one model branch and one gh pr create per orphan folder. That is 3, because rhiza init's actor folder is a third orphan. It also asserts the held entry stays untracked and on no branch, and that a second run is quiet. go test -trimpath -count=1 ./internal/api -run TestDaemonCommitsUntracked -v passes (daemonCommitUntracked makes a model/<cs>- branch with one PR per branch, and a second tick ships nothing). Run from a scratch worktree outside ~/code, because inside ~/code the rootguard tripped on the live daemon's usage.live.json. In the diff, daemonCommitUntracked is called from every daemon poll (internal/api/daemon_projects.go:424), behind modelBranchEnabled. Known minor review findings (branch count 3 vs 2, 8-char prefix check) stand, but they do not defeat the acceptance.

## Links

### Out

- verifies [commit-untracked-log](../capability/rhiza/orchestration/shipping/commit-untracked-log.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

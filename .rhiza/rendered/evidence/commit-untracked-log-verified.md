# commit-untracked-log-verified

evidence · `01M76SVFHDJAKSNXT8SBAWX4ZG` · actor [claude-code-c7ea6cab](../actor/claude-code-c7ea6cab.md)

## result

PASS

## statement

Verifier claude-code-c7ea6cab at 6c092c4b4: go test -trimpath ./cmd/rhiza -run 'TestScripts/(model-commit-untracked|describe-golden)$' passes. In the kept $WORK: the bare origin has one model/<cs>-<entry> branch per orphan folder (the two orphan change sets plus the init actor's log folder, so the script asserts 3, not the acceptance's 2) and the fake gh logged one pr create per branch; the run-held change set 01M5PWGK… is on no branch and still untracked; a second run prints 'no untracked log entries to commit'. go test -trimpath ./internal/api -run 'TestDaemonCommitsUntracked|TestModelBranch' passes: a daemon tick (daemonCommitUntracked, called from each project poll, so the first poll after start too) ships an orphan onto a model branch with a PR and a second tick ships nothing. gofmt and go vet are clean. Diff read: runModelCommitUntracked groups orphanLogEntries by change set and ships each through shipModelWrite, which shipOnModelBranch now shares; model branches are now named by full IDs. Weak spot: the held-entry check compares 8-char ID prefixes, which passes only because the seeded IDs differ there. One run of internal/api failed on the known .rhiza-root/usage.live.json write by the enclosing root and passed on rerun.

## Links

### Out

- verifies [commit-untracked-log](../capability/rhiza/orchestration/shipping/commit-untracked-log.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

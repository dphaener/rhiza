# commit-untracked-log-verified-9

evidence · `01M7706FR8ACSRQEJ41AFHXEGK` · actor [claude-code-89a3edf6](../actor/claude-code-89a3edf6.md)

## result

PASS

## statement

Verified at 2ebbea723 (attempt 9). go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/model-commit-untracked$' passes: the bare origin and fake gh get one model branch and one PR per orphan change set, and the held entry stays untracked on no branch. That count is 3, because rhiza init's actor folder is a third orphan; this is the known minor carried finding. go test -trimpath -count=1 ./internal/api -run 'TestDaemonCommitsUntracked' passes: the round ships the orphan entry on model/<cs>-* with one PR, and a second call ships nothing. TestModelNetThrottled passes: due on the first poll, not at +5s or +59s, due again at +1m. Read daemon_projects.go:429-470: poll runs daemonCommitUntracked and daemonMergeModelPRs only when modelNetDue, through boundedExecCtx (modelNetTimeout 30s, GIT_TERMINAL_PROMPT=0). This matches the leaf's revised statement (377e4c386) and Decision commit-untracked-throttled-cadence. make self-check in this checkout fails to fold 01M76ZNP56RR5HYFB099XN640P because its predecessor 01M76ZKWG11Y4FK6T8A8Y93JMV (the run's cadence ruling) is still untracked in the run worktree and is committed at close-out. In a scratch clone with that entry and 01M76ZJC3V8YNSQA5DEE0DQPSG added, this commit's build reports no fold or docs failures, only .rhiza/rendered/ files, which close-out renders.

## Links

### Out

- verifies [commit-untracked-log](../capability/rhiza/orchestration/shipping/commit-untracked-log.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

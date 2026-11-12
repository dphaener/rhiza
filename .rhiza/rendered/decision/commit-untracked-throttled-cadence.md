# commit-untracked-throttled-cadence

decision · `01M76ZKWCQ3ZMWYAG25X9H3C11`

## reasoning

Finding f1 is a wording mismatch, not a behaviour defect. internal/api/daemon_projects.go:429 runs daemonCommitUntracked and daemonMergeModelPRs only when modelNetDue: first poll, then once a minute, under a 30s per-program bound. Commit e8e289be4 added that on purpose: a stalled ls-remote/fetch/push/gh on every 5s poll held up every project's poll and spent the gh token each tick, the same reason the self-update fetch is throttled. Reverting it to every tick would bring that back. The change set's intent (nothing left uncommitted on main) is still met with a one-minute delay, since orphans are picked up on start and within a minute after. So the leaf's text, not the code, should change; three review cycles failed only because the leaf still promises every tick.

## statement

Keep the throttle from e8e289be4 and change leaf @commit-untracked-log instead: its statement should say the daemon commits orphan entries on its first poll and then at most once a minute (modelNetEvery), each git/gh call bounded by modelNetTimeout, not 'on each tick'; the acceptance's daemon test asserts a due round (first poll) does the commit, not every tick

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m76njm-1](../question/leaf-failed-01m76njm-1.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

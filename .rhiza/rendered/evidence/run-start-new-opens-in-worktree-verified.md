# run-start-new-opens-in-worktree-verified

evidence · `01M76QHP7ZAZ2HPQA5AM1C96YK` · actor [claude-code-d9c8e62c](../actor/claude-code-d9c8e62c.md)

## result

PASS

## statement

At 2c86968ef, go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/run-start-new-main-clean passes. The script starts a --new plan run and a --new --milestone run, each through a gate step (rhiza run gate). For each run, held.sh checks that every entry of the new change set in main's .rhiza/log is byte-identical in ../repo.<slug>, is committed at the run branch's HEAD with the same bytes, and that the worktree's status is clean (prints 'held and committed: 1'). rhiza change show <cs> from main prints the title. rhiza check prints no commit-untracked or untracked-log-entry warning. rhiza health prints 'ok main on main, clean' with no orphan or untracked. go test ./internal/api -run TestNewPlanChange passes: a dry run predicts the change set ID and worktree the real run then uses. The diff of internal/api/run.go draws the entry in a dry workspace, adds the branch and worktree, then writes the same entry to the worktree and main. Attempt 1's findings f1-f3 (no commit, no health, milestone run not checked) are all addressed.

## Links

### Out

- verifies [run-start-new-opens-in-worktree](../capability/rhiza/workflows/workflow-runs/run-start/run-start-new-opens-in-worktree.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

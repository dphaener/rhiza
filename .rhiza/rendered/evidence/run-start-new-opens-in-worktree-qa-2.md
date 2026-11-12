# run-start-new-opens-in-worktree-qa-2

evidence · `01M76WZNW4ETA8BHXCZHA1CXGT` · actor [claude-code-df2a667e](../actor/claude-code-df2a667e.md)

## result

PASS

## statement

QA cycle 2, hands-on with rhiza built from this branch at c57d3a309 (go build into the session scratchpad), in a scratch Rhiza root with a project repo, a local bare origin, real git, a fake rex/claude and a scripted gh logging calls (pr view fails, pr create prints a URL), with no RHIZA_TEST_NOW so the model-branch step is live: rhiza run start --foreground --workflow plan --new 'Plan reports' (and again with --milestone) succeeded; the open entry in the main checkout was byte-identical to the one in ../repo.<slug> and committed at the run branch's HEAD with a clean worktree; rhiza change show of the new change set worked from main; rhiza check printed ok with no untracked-entry warning and rhiza health's main part was 'ok ... 2 untracked log entries pending on a model branch or run'. No model/* branch reached origin and gh was never called (the cycle 1 spurious model PR is gone).

## Links

### Out

- verifies [run-start-new-opens-in-worktree](../capability/rhiza/workflows/workflow-runs/run-start/run-start-new-opens-in-worktree.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

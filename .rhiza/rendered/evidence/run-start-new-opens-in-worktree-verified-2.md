# run-start-new-opens-in-worktree-verified-2

evidence · `01M76WB1MHTY51VN163ZFF87K5` · actor [claude-code-053f410c](../actor/claude-code-053f410c.md)

## result

PASS

## statement

Verifier, attempt 2, at 2990cabf5. Ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/run-start-new-main-clean$': PASS. The script starts a --new plan run and a --milestone plan run. For each one it checks that the run's worktree holds the open entry byte for byte and commits it at HEAD with a clean tree, that rhiza change show of the new change set succeeds from main, that check reports no commit-untracked or untracked entry, and that health shows main ok with no orphan, the held copy on its 'untracked log entries pending on a model branch or run' line, and no 'on no branch'. With RHIZA_TEST_MODEL_BRANCH=1 it also checks that no model/* branch is pushed and that the fake gh is never called. This fixes attempt 1's failing '! stdout untracked'. Ran go test -trimpath -count=1 ./internal/api -run TestNewPlanChange: PASS (a dry run predicts the same change set ID and worktree and writes nothing to main). Neighbouring scripts also pass: run-start-new, model-commit-untracked, health*. Read the diff: newPlanChange writes the main copy through workspace.writeEntryUnshipped, which turns off the mirror and modelBranch hooks only for that write, so only the run's branch ships it.

## Links

### Out

- verifies [run-start-new-opens-in-worktree](../capability/rhiza/workflows/workflow-runs/run-start/run-start-new-opens-in-worktree.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).

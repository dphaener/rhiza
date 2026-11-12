# workflow-cleanup-steps-verified-2

evidence · `01M7768MPQ3YC50B455805DVP3` · actor [claude-code-0dfba7ab](../actor/claude-code-0dfba7ab.md)

## result

PASS

## statement

Verifier, attempt 2, at 4a122043c. go test -trimpath ./internal/workflow/... ./internal/run/... passed (TestParseCleanup covers the [[cleanup]] section). go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-cleanup-steps|run-abandon.*|abandon.*|run-retry.*|workflow-.*)$' passed. run-cleanup-steps.txtar shows a failing step's cleanup steps ran in order in the run's worktree with the run's env (RHIZA_STEP, RHIZA_CHANGE) and the run stayed failed. It shows a cleanup step exiting 1 journalled as fail while the run ended succeeded and the next cleanup step still ran. It shows an abandoned run's cleanup ran. Built-in workflows declare no [[cleanup]]. Fixes for attempt 1's findings, from the diff and the tests: abandoning an already-failed run no longer re-runs cleanup, because the done set now resets only on TypeRetried, and the script asserts 2 cleanup entries after abandon. abandonedCleanup skips cleanup and logs 'cleanup not run: worktree ... is gone' instead of falling back to the main checkout, and the script asserts no cleanup entry and no side effect.

## Links

### Out

- verifies [workflow-cleanup-steps](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-cleanup-steps.md)

## History

- `01M763YS08ZM13JZ065Q95YA2E` Workflow cleanup phase: a workflow file declares command steps that run once when a run of it ends, done, failed or abandoned; a failing cleanup step is logged and never changes the run's result; workflow check and describe cover it.

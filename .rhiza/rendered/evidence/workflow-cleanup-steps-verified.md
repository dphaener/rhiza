# workflow-cleanup-steps-verified

evidence · `01M775122EPDAMZQK45HP1PRX4` · actor [claude-code-59e5ff33](../actor/claude-code-59e5ff33.md)

## result

PASS

## statement

Verifier at 3f4897766: read the diff (workflow.Cleanup parse/validate, run engine cleanup after finish and on abandon, cleanup journal entry, describe/check). Ran go test -trimpath ./internal/workflow/... ./internal/run/... (ok, incl. TestParseCleanup covering names, default cleanup-<n>, missing/blank run, bad name, duplicate, unknown key, Validate) and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-cleanup-steps|describe-golden)$' -count=1 -v (both PASS). run-cleanup-steps shows a failing step's run end failed with its cleanup journalled pass and run in the worktree with the run env, a cleanup exiting 1 journalled fail while the run still ends succeeded and the next cleanup still runs, an abandoned run running its cleanup, and check refusing a cleanup without run. git grep finds no [[cleanup]] in built-in or repo workflows.

## Links

### Out

- verifies [workflow-cleanup-steps](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-cleanup-steps.md)

## History

- `01M763YS08ZM13JZ065Q95YA2E` Workflow cleanup phase: a workflow file declares command steps that run once when a run of it ends, done, failed or abandoned; a failing cleanup step is logged and never changes the run's result; workflow check and describe cover it.

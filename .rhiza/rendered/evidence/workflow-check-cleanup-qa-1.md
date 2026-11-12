# workflow-check-cleanup-qa-1

evidence · `01M7758GWE2EGAC1F0Q5A9K5XW` · actor [claude-code-3c7e77a4](../actor/claude-code-3c7e77a4.md)

## result

PASS

## statement

QA at 333b5bf65 with a scratch build: workflow check passes workflows with valid [[cleanup]] sections and refuses (exit 1) a cleanup with kind = "agent", prompt, on.pass (unknown key "cleanup.kind"/"cleanup.prompt"/"cleanup.on.pass") and name "Bad Name" (want lowercase words). workflow describe prints 'cleanup <name>: <run>' lines in file order and --json carries cleanup [{name, run}], with cleanup-<n> for an unnamed step. docs/reference/workflow-file.md and cli.md document the cleanup key and field. go test ./internal/workflow -run TestParseCleanup and TestScripts/(run-cleanup-steps|workflow-check|workflow-describe) pass.

## Links

### Out

- verifies [workflow-check-cleanup](../capability/rhiza/workflows/workflow-authoring/workflow-check-cleanup.md)

## History

- `01M763YS08ZM13JZ065Q95YA2E` Workflow cleanup phase: a workflow file declares command steps that run once when a run of it ends, done, failed or abandoned; a failing cleanup step is logged and never changes the run's result; workflow check and describe cover it.

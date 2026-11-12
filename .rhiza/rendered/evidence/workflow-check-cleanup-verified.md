# workflow-check-cleanup-verified

evidence · `01M774KRQZD99FKCX62NMS02AT` · actor [claude-code-70f30e59](../actor/claude-code-70f30e59.md)

## result

PASS

## statement

At 14e7e4aeb (with leaf commit 3f4897766): go test -trimpath ./internal/workflow passes, incl. TestParseCleanup covering a valid [[cleanup]] section and invalid ones (no run, blank run, bad name, duplicate, unknown key) via Parse and Validate. TestScripts/workflow-check, workflow-describe, describe-golden and run-cleanup-steps pass: workflow-check scripts a valid cleanup section and rejects no-run, agent-kind and on.pass cleanup steps; workflow-describe shows 'cleanup <name>: <run>' lines in order and a cleanup field in --json. describe-golden carries the cleanup field; make docs leaves docs/reference unchanged (cli.md and workflow-file.md already document cleanup). A scratch build's workflow check prints unknown key "cleanup.kind"/"cleanup.prompt"/"cleanup.on.pass" and 'run is missing' for the invalid files.

## Links

### Out

- verifies [workflow-check-cleanup](../capability/rhiza/workflows/workflow-authoring/workflow-check-cleanup.md)

## History

- `01M763YS08ZM13JZ065Q95YA2E` Workflow cleanup phase: a workflow file declares command steps that run once when a run of it ends, done, failed or abandoned; a failing cleanup step is logged and never changes the run's result; workflow check and describe cover it.

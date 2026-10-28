# workflow-step-verified

evidence · `01M608S29PVAE8CZA7TN1DSGG4` · actor [claude-code-ad85f0f2](../actor/claude-code-ad85f0f2.md)

## result

PASS

## statement

Verifier claude-code-ad85f0f2 at fec4749 (diff 3be07f6..fec4749 read). go test ./internal/workflow ./internal/run -count=1: TestWorkflowStepErrors rejects an unknown workflow, skip naming a step the inner lacks, a workflow including itself directly (inner -> inner) and through another (test -> inner -> test); TestParseErrors rejects a workflow step without workflow and workflow/skip on a command step; all PASS. TestWorkflowStepRunsInnerSteps (fake runner, outer step w over three-step inner skipping two) journals inner/one:pass, inner/three:pass, after:pass in order, inner/two never run, outer on.pass taken: PASS. TestWorkflowStepInnerFailureResumesThere fails the run at inner/three, then Engine.Retry (what rhiza run resume --retry calls via runEngineRetry, reparsing from the run's recorded Includes) re-runs only inner/three then after: PASS. Resolution goes through workflowResolver (repo .rhiza/workflows/<name>.toml, else built-in), same as run start. make fmt-check, vet and go test ./... pass; make self-check reports only stale render for actors registered by this run's later sessions (claude-code-8f2a5c63, claude-code-ad85f0f2), not by the leaf's commit.

## Links

### Out

- verifies [workflow-step](../capability/ship-through-rhiza/workflow-step.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.

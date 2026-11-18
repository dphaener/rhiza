# workflow-retry-from-key

capability · `01M7QQ8DPFS4RYM9NQ682371XB` · parent [workflow-step](../workflow-step.md)

## acceptance

Tests in internal/workflow parse a step with retry_from and an included workflow whose retry_from comes out prefixed (plan/interview); (failure case) retry_from naming no step of the workflow, or itself, is refused at parse with the step and key named; rhiza workflow describe shows the key (a cmd/rhiza testscript or the describe test); internal/guide/text/workflow-keys.md and docs/reference/workflow-file.md document it. Run: go test -trimpath -p 2 ./internal/workflow/... and the touched testscript

## files

internal/workflow/workflow.go, internal/workflow/include.go, internal/workflow/workflow_test.go, internal/workflow/include_test.go, internal/api/workflow*.go (describe only), internal/guide/text/workflow-keys.md, docs/reference/workflow-file.md, cmd/rhiza/testdata/script/workflow-describe*.txtar

## must_not

Change how a run retries (the engine leaf does); change any built-in workflow file; change other step keys

## statement

A workflow step may set retry_from, the step of the same workflow a failed run restarts at when it fails at this step; it parses, is prefixed like on.* targets when its workflow is included by a workflow step, is validated, shown by rhiza workflow describe and documented

## Links

### Out

- decided-by [retry-from-step-key](../../../decision/retry-from-step-key.md)

### In

- depends-on [retry-restarts-at-named-step](../../rhiza/workflows/workflow-runs/resume-failed-run/retry-restarts-at-named-step.md)
- verifies [workflow-retry-from-key-verified](../../../evidence/workflow-retry-from-key-verified.md)
- verifies [workflow-retry-from-key-qa-1](../../../evidence/workflow-retry-from-key-qa-1.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox

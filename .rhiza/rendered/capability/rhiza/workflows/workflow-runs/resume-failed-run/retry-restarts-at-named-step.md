# retry-restarts-at-named-step

capability · `01M7QQ8DR1NT9E6KAPQSF5614F` · parent [resume-failed-run](../resume-failed-run.md)

## acceptance

Tests in internal/run fail a run at a step with retry_from and see the retried entry name the restart step and the failed step and the next visit enter the restart step with its max_visits counted afresh; a test fails a plan-bug-shaped run (a workflow step including plan) at plan/check over its cap and sees the retry enter plan/interview; (failure case) a step without retry_from still restarts at itself and a failed check fed by a wait still re-waits (RetryStep's existing tests pass unchanged); a test in internal/workflow asserts every review-loop step of the built-in plan and plan-milestone sets retry_from = interview. Run: go test -trimpath -p 2 ./internal/run/... ./internal/workflow/... and cmd/rhiza TestScripts/run-resume-retry

## files

internal/run/engine.go, internal/run/journal.go, internal/run/engine_test.go, internal/run/journal_test.go, internal/workflow/builtin/plan.toml, internal/workflow/builtin/plan-milestone.toml, internal/workflow/workflow_test.go, cmd/rhiza/testdata/script/run-resume-retry.txtar

## must_not

Change the brief (the brief leaf does) or the watchdog; change any review-loop step's on.* routing or max_visits; change ship-change

## statement

Retrying a failed run (run.Engine.Retry, so rhiza run resume --retry, rhiza run start on a failed change set and the daemon) restarts it at the failed step's retry_from when set, else where it does today; the retried entry also journals the failed step; the built-in plan and plan-milestone workflows set retry_from = "interview" on their review-loop steps

## Links

### Out

- depends-on [workflow-retry-from-key](../../../../ship-through-rhiza/workflow-step/workflow-retry-from-key.md)
- decided-by [retry-from-step-key](../../../../../decision/retry-from-step-key.md)
- decided-by [plan-review-retries-at-interview](../../../../../decision/plan-review-retries-at-interview.md)

### In

- depends-on [retry-brief-carries-failure](retry-brief-carries-failure.md)
- verifies [retry-restarts-at-named-step-verified](../../../../../evidence/retry-restarts-at-named-step-verified.md)
- verifies [retry-restarts-at-named-step-qa-1](../../../../../evidence/retry-restarts-at-named-step-qa-1.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox

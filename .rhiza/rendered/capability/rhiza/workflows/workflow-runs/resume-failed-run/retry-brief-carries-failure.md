# retry-brief-carries-failure

capability · `01M7QQ8DVDM757QQGNGSG05RM7` · parent [resume-failed-run](../resume-failed-run.md)

## acceptance

Tests in internal/run render a brief with a previous attempt (step, outcome, output tail, findings with check, severity and text) and one without it unchanged byte for byte; a test in internal/api (scoped with -run) retries a run failed at a review-check-shaped command step after a findings round and sees the restarted agent step's brief.md hold the refusal's output and the round's findings; (failure case) a run never retried, and a session started before the retry, get no such section. Run: go test -trimpath -p 2 ./internal/run/... and go test -trimpath -p 2 ./internal/api -run <the new tests>

## files

internal/run/brief.go, internal/run/brief_test.go, internal/run/journal.go, internal/api/run_brief.go, internal/api/run_brief_test.go, internal/guide/text/contract.md

## must_not

Change handback.md or a spawned child's brief; write anything into .rhiza/log; change where a retry restarts

## statement

Every agent step session started after a run's last retried entry finds a Previous attempt section in its brief.md: the failed step, outcome, cause or output tail and the findings of the latest round recorded before the retry

## Links

### Out

- depends-on [retry-restarts-at-named-step](retry-restarts-at-named-step.md)
- decided-by [retry-brief-previous-attempt](../../../../../decision/retry-brief-previous-attempt.md)

### In

- depends-on [watchdog-retries-real-failure](../../../orchestration/daemon/daemon-watchdog/watchdog-retries-real-failure.md)
- verifies [retry-brief-carries-failure-verified](../../../../../evidence/retry-brief-carries-failure-verified.md)
- verifies [retry-brief-carries-failure-qa-1](../../../../../evidence/retry-brief-carries-failure-qa-1.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox

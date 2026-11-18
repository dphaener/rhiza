# watchdog-retries-real-failure-verified

evidence · `01M7QS90HFD0PGWTG7NP84T0VY` · actor [claude-code-240ba28d](../actor/claude-code-240ba28d.md)

## result

PASS

## statement

At 4667de181, under capped scopes: go test -trimpath -p 2 ./internal/api -run 'Watchdog|Inbox' passed (TestWatchdogFailedRetries: one --retry of a real failure charged as an attempt with Spent 0, no inbox item while due, second failure left unretried and listed naming the retry; TestWatchdogFailedRetriesNever: failedRetries 0, attended under autonomousOnly, live, abandoned and disabled never resumed; TestWatchdogFailedRetriesTransientKeepsBudget: transient charged to retryBudget only; TestWatchdogFailedRetryMeasure: plan run fails at review-check, retry restarts at interview with the refusal and Failed step in brief.md, second failure in rhiza inbox with the retry named). Also passed: ./internal/api -run 'Daemon|Notify|Health|Failed', ./internal/daemon/..., go vet, cmd/rhiza TestScripts/watchdog-policy (failedRetries listed at both scopes, set, -1 refused) and the daemon/watchdog/inbox/run-resume-retry scripts; rhiza docs, ship generate and render exit 0 with no tracked diff. Docs in config.md and runs-and-daemon.md. Classifier untouched; PR-idle and budget-hold guards only exclude failed-retry-only memory, so the other budgets are unchanged. No TODOs.

## Links

### Out

- verifies [watchdog-retries-real-failure](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-retries-real-failure.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox

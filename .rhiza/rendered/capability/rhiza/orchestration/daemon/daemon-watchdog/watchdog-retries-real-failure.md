# watchdog-retries-real-failure

capability · `01M7QQ8DY3RPXE969S13KP6WWJ` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

Tests in internal/api (scoped with -run) give the watchdog a failed autonomous run classified real and see one resume with --retry and an attempt charged, then a second failure of the same run left unretried and listed in rhiza inbox with the retry named; failedRetries 0 retries nothing; the inbox leaves out a failed run with a retry due; a transient failure still uses retryBudget; a test of the measure fails a plan run at review-check, sees the daemon's retry restart at interview with the refusal in the brief, and a second failure reach the inbox; (failure case) an attended run under autonomousOnly, a live or abandoned run and a disabled watchdog are never retried; the setting is listed by rhiza daemon settings (watchdog-policy testscript) and documented. Run: go test -trimpath -p 2 ./internal/api -run 'Watchdog|Inbox' and cmd/rhiza TestScripts/watchdog-policy

## files

internal/api/daemon_watchdog.go, internal/api/watchdog_policy.go, internal/api/inbox.go, internal/api/daemon_watchdog_test.go, internal/api/inbox_test.go, internal/daemon/state.go, cmd/rhiza/testdata/script/watchdog-policy.txtar, docs/reference/config.md, docs/runs-and-daemon.md

## must_not

Change the failure classifier or its prompt; change the stalled-step or CI passes' budgets; retry an abandoned, live or attended (under autonomousOnly) run

## statement

The daemon's watchdog retries a failed run whose failure is real or unjudged with rhiza run resume --retry up to failedRetries (default 1), journalling each as a watchdog attempt; once they are spent the run goes to the inbox with the retries tried, and no failed run is in the inbox while a retry is still due

## Links

### Out

- depends-on [retry-brief-carries-failure](../../../workflows/workflow-runs/resume-failed-run/retry-brief-carries-failure.md)
- decided-by [daemon-failed-retries-setting](../../../../../decision/daemon-failed-retries-setting.md)

### In

- verifies [watchdog-retries-real-failure-verified](../../../../../evidence/watchdog-retries-real-failure-verified.md)
- verifies [watchdog-retries-real-failure-qa-1](../../../../../evidence/watchdog-retries-real-failure-qa-1.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox

# watchdog-classifies-failures-verified

evidence · `01M77JH8J7H4WV8QW9C9WMJH1F` · actor [claude-code-2352b3eb](../actor/claude-code-2352b3eb.md)

## result

PASS

## statement

Verifier at 54110a0ef: go test -trimpath -count=1 ./internal/api -run 'Watchdog|Classif' passes. With an injected classifier, TestWatchdogTransient retries the transient failures (work, lock, gh) and leaves the real one to the inbox, logging 'classified real' with its reason; TestWatchdogTransientClassifierError leaves both a classifier error and a missing classifier to the inbox, never resuming, and logs the reason. grep finds no transientSignatures in code. Attempt-1 findings fixed: classifyOnce remembers a real or unjudged verdict in WatchdogRun.Judged at the journal's last entry, so a second poll neither classifies nor logs again, and a new failure is judged anew (tested); watchdogRequests now checks budget and queue before classifying and logs real verdicts (TestWatchdogRequests). Memory carries across polls via watchdogPasses for non-live runs. go vet and gofmt clean; internal/daemon tests pass.

## Links

### Out

- verifies [watchdog-classifies-failures](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-classifies-failures.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

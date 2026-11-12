# watchdog-classifies-failures-verified-2

evidence · `01M77KJ37SJMTBE9YYMV2DPC8X` · actor [claude-code-7b8d33b6](../actor/claude-code-7b8d33b6.md)

## result

PASS

## statement

Verifier, attempt 2, at dd206f2dc. Ran go test -trimpath ./internal/api -run 'Watchdog|Classif' -count=1: all pass. TestWatchdogTransient, with an injected classifier, retries the transient failures (work, lock, gh) and leaves the real one to the inbox with 'classified real' and its reason logged. TestWatchdogTransientClassifierError leaves a classifier error ('classifier failed') and a missing classifier to the inbox. grep finds no transientSignatures in any Go file. Attempt 1's finding (unbounded classifier stalling the daemon tick) is fixed by a0179336b: shellWithStdin runs under a 60s classifierTimeout in its own process group, killed on expiry, and the error makes ClassifyFailure return a non-transient verdict. go test -trimpath ./cmd/rhiza -run 'TestShellWithStdin' -count=1 passes: TestShellWithStdinTimesOut kills 'sleep 20' in about 0.2s with a timeout error. go vet on both packages is clean.

## Links

### Out

- verifies [watchdog-classifies-failures](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-classifies-failures.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

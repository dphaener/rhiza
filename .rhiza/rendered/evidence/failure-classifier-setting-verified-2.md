# failure-classifier-setting-verified-2

evidence · `01M77KGNCHQVCCKR70R876SAFR` · actor [claude-code-98a3f588](../actor/claude-code-98a3f588.md)

## result

PASS

## statement

Verifier (attempt 2) at a0179336b: read the diff (shellWithStdin now runs the classifier under a 60s context deadline in its own process group, killing the group on timeout and returning an error, which ClassifyFailure reads as real). Ran go test -trimpath ./internal/api -run Classif: TestClassifyFailure covers transient, real, garbage, empty output, a failing command and off, with real on anything but a clear transient; all PASS. Ran go test -trimpath ./cmd/rhiza -run TestShellWithStdin: a hanging 'sleep 20' classifier is killed after the shortened 200ms timeout with a timed-out error, and a 'cat' classifier returns its stdin; PASS. Ran TestScripts/daemon-settings (lists rhiza.daemon.failureClassifier, empty from Setup, with the claude -p --model haiku default in --json), TestScripts/run-failure-classifier (empty from Setup: failing step runs once and no classifier call; a fake transient classifier reruns once and sees the step output on stdin; a real one does not rerun) and TestScripts/build-leaf-fails: all PASS. Setup in cmd/rhiza/main_test.go sets rhiza.daemon.failureClassifier empty. A broader internal/api -run Watchdog pass was flagged only by the root-state guard because the live root daemon rewrote ~/code/.rhiza-root/usage.live.json during the run, which this change does not touch.

## Links

### Out

- verifies [failure-classifier-setting](../capability/rhiza/orchestration/daemon/daemon-settings/failure-classifier-setting.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

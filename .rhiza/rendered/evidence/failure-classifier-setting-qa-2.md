# failure-classifier-setting-qa-2

evidence · `01M77KWG18Q7SFFT3R84VEE7C8` · actor [claude-code-20068c74](../actor/claude-code-20068c74.md)

## result

PASS

## statement

QA cycle 2 at a3db0fa7f with the built binary in a scratch root and HOME: rhiza daemon set rhiza.daemon.failureClassifier to a command that backgrounds a 'sleep 300' child, then sleeps 300 and would print 'transient: slow'. rhiza run start --foreground on a workflow whose command step fails took 61s and printed 'failure classifier: real: classifier failed: timed out after 1m0s'; the step ran once (no rerun), and both the classifier shell and its backgrounded child were gone afterwards (kill -0 fails; no 'sleep 300' left in pgrep). go test -trimpath ./cmd/rhiza -run 'TestShellWithStdin|TestScripts/run-failure-classifier' passes.

## Links

### Out

- verifies [failure-classifier-setting](../capability/rhiza/orchestration/daemon/daemon-settings/failure-classifier-setting.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

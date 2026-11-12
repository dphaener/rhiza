# failure-classifier-setting-qa-1

evidence · `01M77K5FQ95Q7ZTYF4M0AQR2JE` · actor [claude-code-097fafbe](../actor/claude-code-097fafbe.md)

## result

PASS

## statement

QA at 3a745c744 with the built binary in a scratch Rhiza root, no RHIZA_* env. rhiza daemon settings lists rhiza.daemon.failureClassifier with the claude -p --model haiku default (default, machine); rhiza daemon set ... '' writes "rhiza.daemon.failureClassifier" = "" to rhiza.toml and shows it as "" (rhiza.toml, machine). Against a failing command step: empty ran it once with no classifier line; a fake classifier saw the step's output ('boom') on stdin; 'transient: flaky net', 'TRANSIENT - lock held' and a verdict after blank lines were read as transient; 'real: test fails', a free-text sentence ('I think it is probably transient') and a classifier exiting 3 were all real (no clear verdict / classifier failed: exit status 3). With the setting unset, the real default (Claude headless, Haiku) called a git index.lock 'File exists' error transient and a Go '--- FAIL: TestAdd ... want 2' real, each in about 7s with a one-line reason. A 'sleep 20' classifier held the step 21s: the call has no time limit.

## Links

### Out

- verifies [failure-classifier-setting](../capability/rhiza/orchestration/daemon/daemon-settings/failure-classifier-setting.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

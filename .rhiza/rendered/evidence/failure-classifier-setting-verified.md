# failure-classifier-setting-verified

evidence · `01M77GXG6PJY8XGRYNDFJQAKXN` · actor [claude-code-1b4e23b2](../actor/claude-code-1b4e23b2.md)

## result

PASS

## statement

At 237e9ff9b, read the diff and ran go test -trimpath ./internal/api -run 'TestClassify|Setting': TestClassifyFailure's fake classifier commands for transient, real, garbage output, empty output, a near-miss word and a failing command all PASS, and only a clear 'transient' line comes out transient (the rest are real with a reason). The scripts daemon-settings (lists rhiza.daemon.failureClassifier = "" from Setup's rhiza.toml, the claude -p --model haiku default in --json, and set/unset), run-failure-classifier (empty from Setup: the failing command step runs once and nothing is called; a fake transient classifier reruns it and its output reaches the classifier on stdin; a fake real verdict does not rerun it), build-leaf-fails and describe-golden all PASS. Ran 16 scripts that rewrite rhiza.toml or rerun failed steps with a logging fake claude first on PATH: all PASS and claude was never called. Outside this leaf: TestShipCheckStabilityRecord/red_then_green_on_rerun fails the same way at the parent 424e22486, and TestInboxGateQuestionsBudget passes but takes 135s.

## Links

### Out

- verifies [failure-classifier-setting](../capability/rhiza/orchestration/daemon/daemon-settings/failure-classifier-setting.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

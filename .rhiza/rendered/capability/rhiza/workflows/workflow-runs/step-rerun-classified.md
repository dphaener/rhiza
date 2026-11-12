# step-rerun-classified

capability · `01M7642N3Y65X0MX8J6J2WSJVH` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Run tests with an injected classifier show a command step rerun once on transient and not on real.

## statement

The step runner reruns a failed command step once only when the classifier says transient, instead of matching Go's vanished-$WORK signature (step-retry-vanished-go-work is retired).

## Links

### Out

- decided-by [d-llm-classifier](../../../../decision/d-llm-classifier.md)

### In

- verifies [step-rerun-classified-verified](../../../../evidence/step-rerun-classified-verified.md)
- verifies [step-rerun-classified-qa-1](../../../../evidence/step-rerun-classified-qa-1.md)
- verifies [pr-480-failures-classified-by-a](../../../../evidence/pr-480-failures-classified-by-a.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

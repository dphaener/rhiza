# ci-logs-to-fixer

capability · `01M7642N4ZK43G1W494MBFE7VE` · parent [shipping](../shipping.md)

## acceptance

internal/api ship tests with fake gh output for a failing Jest job record a finding naming the saved log, and the saved log holds the failure lines.

## statement

When CI still fails after its reruns, ship saves the failed jobs' logs, trimmed to the failing part, in the run's folder and the fix step's finding points the fixer at them; the Go, RSpec, Minitest and pytest log parsers are gone.

## Links

### Out

- decided-by [d-llm-classifier](../../../../decision/d-llm-classifier.md)

### In

- verifies [ci-logs-to-fixer-verified](../../../../evidence/ci-logs-to-fixer-verified.md)
- verifies [ci-logs-to-fixer-qa-1](../../../../evidence/ci-logs-to-fixer-qa-1.md)
- verifies [pr-480-failures-classified-by-a](../../../../evidence/pr-480-failures-classified-by-a.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

# watchdog-classifies-failures

capability · `01M7642N2Q7TPP6BCV43AGR07T` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

internal/api watchdog tests with an injected classifier show a transient failure retried, a real one and a classifier error left to the inbox; no transientSignatures remain.

## statement

The watchdog retries a failed step within its budget only when the classifier says transient, and leaves it to the inbox when the classifier says real or cannot run; the hard-coded transient signature list is gone and each verdict and reason is logged.

## Links

### Out

- decided-by [d-llm-classifier](../../../../../decision/d-llm-classifier.md)

### In

- verifies [watchdog-classifies-failures-verified](../../../../../evidence/watchdog-classifies-failures-verified.md)
- verifies [watchdog-classifies-failures-qa-1](../../../../../evidence/watchdog-classifies-failures-qa-1.md)
- verifies [watchdog-classifies-failures-verified-2](../../../../../evidence/watchdog-classifies-failures-verified-2.md)
- verifies [watchdog-classifies-failures-qa-2](../../../../../evidence/watchdog-classifies-failures-qa-2.md)
- verifies [pr-480-failures-classified-by-a](../../../../../evidence/pr-480-failures-classified-by-a.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

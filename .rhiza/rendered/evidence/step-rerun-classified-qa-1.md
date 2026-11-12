# step-rerun-classified-qa-1

evidence · `01M77K5FZ2JPN0T06ZXX00X4BC` · actor [claude-code-097fafbe](../actor/claude-code-097fafbe.md)

## result

PASS

## statement

QA at 3a745c744 with the built binary: rhiza run start --foreground on a workflow whose command step fails on its first run. A transient verdict printed 'step failure classified transient; running it once more' and the step ran twice and passed; a real verdict, garbage output, a failing classifier and an empty setting each left it at one run, failed. A step that always fails with a transient classifier ran exactly twice, then failed: one rerun only. With the real Haiku default, an index.lock error was rerun and a Go test assertion failure was not.

## Links

### Out

- verifies [step-rerun-classified](../capability/rhiza/workflows/workflow-runs/step-rerun-classified.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

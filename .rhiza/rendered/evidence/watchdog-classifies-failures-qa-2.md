# watchdog-classifies-failures-qa-2

evidence · `01M77KWG4S3DE85Z92JZFK8CMM` · actor [claude-code-20068c74](../actor/claude-code-20068c74.md)

## result

PASS

## statement

QA cycle 2 at a3db0fa7f with the built binary: an autonomous run whose command step fails, then the same hanging classifier (backgrounded sleep 300 child, sleep 300, 'transient: slow') and rhiza daemon --once. The daemon poll finished after 61s and logged 'watchdog: run <id> step a failure classified real, left to the inbox: classifier failed: timed out after 1m0s' (plus a poll-over-limit warning); the classifier and its child were killed. A second rhiza daemon --once finished in under a second and did not call the classifier again (no pid file written). internal/api Watchdog|Classify|Rerun tests pass from a worktree outside ~/code.

## Links

### Out

- verifies [watchdog-classifies-failures](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-classifies-failures.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

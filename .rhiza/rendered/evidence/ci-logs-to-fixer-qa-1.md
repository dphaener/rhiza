# ci-logs-to-fixer-qa-1

evidence · `01M77K5GBM1MRFJKP82NWCZV26` · actor [claude-code-097fafbe](../actor/claude-code-097fafbe.md)

## result

PASS

## statement

QA at 3a745c744 with the built binary and a fake gh whose failed job 'jest' prints a real-shaped gh --log-failed Jest log (job/step/timestamp prefixes, 3000 noise lines then 'FAIL src/garden.test.js', the expect/Expected/Received lines, 'Tests: 1 failed' and '##[error]Process completed with exit code 1.'). First rhiza ship check --settled --json reran run 42; the second listed tests ['CI run 42 (jest)'] and a finding 'failed on <sha> after its reruns; its failing log lines are saved in .git/rhiza/ship/<cs>-ci-logs/42.log: read them, reproduce the failure locally and fix it'. The saved file was 11 lines, prefixes stripped, holding the FAIL line, the assertion and the summary. A log failing at its first lines followed by 3000 PASS lines saved 20 lines keeping both 'FAIL src/early.test.js' and the closing summary. No Go/RSpec/Minitest/pytest parser remains in non-test code (grep).

## Links

### Out

- verifies [ci-logs-to-fixer](../capability/rhiza/orchestration/shipping/ci-logs-to-fixer.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

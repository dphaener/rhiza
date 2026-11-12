# pr-480-failures-classified-by-a

evidence · `01M77N145YN0Y869ASH1T8D3S0` · actor [dphaener](../actor/dphaener.md)

## head

52ace47f7cd7400cf202985136244126a33e607d

## pr

https://github.com/dphaener/rhiza/pull/480

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #480 (https://github.com/dphaener/rhiza/pull/480) passed CI at 52ace47

## status

verified

## Links

### Out

- verifies [failure-classifier-setting](../capability/rhiza/orchestration/daemon/daemon-settings/failure-classifier-setting.md)
- verifies [watchdog-classifies-failures](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-classifies-failures.md)
- verifies [step-rerun-classified](../capability/rhiza/workflows/workflow-runs/step-rerun-classified.md)
- verifies [ci-logs-to-fixer](../capability/rhiza/orchestration/shipping/ci-logs-to-fixer.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

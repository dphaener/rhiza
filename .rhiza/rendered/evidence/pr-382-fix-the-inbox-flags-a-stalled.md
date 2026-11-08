# pr-382-fix-the-inbox-flags-a-stalled

evidence · `01M6XGBESV07WH7N6Y9E9FY57J` · actor [dphaener](../actor/dphaener.md)

## head

05a94dc92429d858877ef1b4819a4edd668548ae

## pr

https://github.com/dphaener/rhiza/pull/382

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #382 (https://github.com/dphaener/rhiza/pull/382) passed CI at 05a94dc

## status

verified

## Links

### Out

- verifies [ship-record-names-run](../capability/rhiza/orchestration/daemon/daemon-watchdog/ship-record-names-run.md)
- verifies [watchdog-pr-owned-by-run](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-pr-owned-by-run.md)

## History

- `01M6XB5KQ5CH2DT5NW7BKETMS8` Fix The inbox flags a stalled run against the wrong run, and tells the owner to merge a PR that a live run owns: item 0568a71 (2026-11-08 ~08:00Z) says 'stalled run 01M6VY7E4JWY (the plan run of change set 01M6VY7H9PYX, which succeeded at 2026-11-07) ... the run ended succeeded but its PR is green at f3a3250 and unmerged: merge it by hand'. PR #376 at f3a3250 belongs to ship-change run 01M6WMX236BTT04TYQDACGKBZ4 of the same change set and branch, which was being resumed toward its merge step; a hand merge then breaks that run. The stalled-PR check should match a PR only to the run that opened or pushed it (its ship record), and skip any change set with a live or queued run.

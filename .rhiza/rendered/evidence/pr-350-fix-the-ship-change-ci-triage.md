# pr-350-fix-the-ship-change-ci-triage

evidence · `01M6THD8AXYMW5PBP2DXGEWDET` · actor [dphaener](../actor/dphaener.md)

## head

a19368db74b94b9c0a045f2469d20686b968eeba

## pr

https://github.com/dphaener/rhiza/pull/350

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #350 (https://github.com/dphaener/rhiza/pull/350) passed CI at a19368d

## status

verified

## Links

### Out

- verifies [ship-check-runner-test-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-runner-test-failures.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests

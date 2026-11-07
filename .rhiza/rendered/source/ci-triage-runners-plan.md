# ci-triage-runners-plan

source · `01M6TBD5MV606EMNJPGW6R277R`

## quotes

```json
[]
```

## raw_sha256

cd4bb96534955463cc8df28f45c99ef72e4c94209f8f8868c61621d144cb442c

## session

bfd0cb38-b7d7-452d-a819-01a1b4b37d39

## summary

Autonomous plan run for change set 01M6TB6N: ship check's parseFailedTests reads only Go '--- FAIL:' lines, so nebari#72's RSpec failure ended the run at ci-triage. The planner asked which intent to serve, which runners to parse and how to name a non-Go failing test; the arbiter ruled (decision-ci-triage-runners): serve m27-nebari-on-rhiza; also parse RSpec 'rspec ./path:line', Rails Minitest '(bin/)rails test path:line' and pytest 'FAILED path::test' lines, nothing broader; name each by the runner's rerun command, Go names unchanged. No owner words: the owner did not take part.

## title

Planning ci-triage runner parsing, autonomous run with the arbiter

## Links

### In

- derived-from [decision-ci-triage-runners](../decision/decision-ci-triage-runners.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests

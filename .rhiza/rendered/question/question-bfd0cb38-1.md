# question-bfd0cb38-1

question · `01M6TBA4SSA2NGV3ZMPT0SV8P3`

## session

bfd0cb38-b7d7-452d-a819-01a1b4b37d39

## statement

Change set 01M6TB6N (ci-triage can't read RSpec failures on nebari#72) serves no intent yet. Which intent should it serve? / ship check's parseFailedTests reads only Go '--- FAIL:' lines. Which runners' failure lines should it also read from gh run view --log-failed? / How should a non-Go failing test be named in ship check's tests and findings?

## Links

### In

- answers [decision-ci-triage-runners](../decision/decision-ci-triage-runners.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests

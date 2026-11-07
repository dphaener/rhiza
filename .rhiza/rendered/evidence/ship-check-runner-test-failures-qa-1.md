# ship-check-runner-test-failures-qa-1

evidence · `01M6TFK7G4428YTQ7N1M5TARWY` · actor [claude-code-d35fa30c](../actor/claude-code-d35fa30c.md)

## result

PASS

## statement

QA cycle 1: built this branch's rhiza (go build -trimpath) and ran it in a scratch repo with a bare origin and a gh shim whose 'run view --log-failed' prints real logs fetched with gh from dphaener/nebari's failed Actions runs. With run 42339573426's log (nebari#72's red head, 5295 lines, 'Failed examples:' then 'rspec ./spec/views/tasks/list_spec.rb:290 # tasks/_list ...'), the first ship check --settled --json reran the run and the second listed tests ['rspec ./spec/views/tasks/list_spec.rb:290'] and one finding with that test and no comment text; ship check --test-failures exited 0. An infra-only log ('Error: Process completed with exit code 1.') gave detail 'failed: test', no tests, and --test-failures exited 1 with 'no failing test was recorded for the head'. Defect seen: runs 42331018773 and 42329799773 print RSpec's id form 'rspec ./spec/system/garden_location_spec.rb[1:1:11] # ...', which ship check names 'rspec ./spec/system/garden_location_spec.rb[1:1:11' (closing bracket dropped), recorded as a qa finding.

## Links

### Out

- verifies [ship-check-runner-test-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-runner-test-failures.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests

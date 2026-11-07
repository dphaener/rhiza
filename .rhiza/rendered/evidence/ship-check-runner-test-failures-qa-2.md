# ship-check-runner-test-failures-qa-2

evidence · `01M6TGQ1E71JZHJMAJC5HZA87H` · actor [claude-code-5b8e4a42](../actor/claude-code-5b8e4a42.md)

## result

PASS

## statement

QA cycle 2 at cb434652: built this branch's rhiza (go build -trimpath) and ran it in a scratch repo with a bare origin and a gh shim whose 'run view --log-failed' prints real logs fetched with gh from dphaener/nebari's failed Actions runs. Each head: first ship check --settled --json reran run 42, the second listed tests and findings. Run 42339573426 (nebari#72) gave ['rspec ./spec/views/tasks/list_spec.rb:290']; run 42331018773 gave ['rspec ./spec/system/garden_location_spec.rb[1:1:11]'] and run 42329799773 gave ['rspec ./spec/system/climate_profile_spec.rb[1:1:6]', 'rspec ./spec/system/garden_location_spec.rb[1:1:11]'], closing bracket kept, so cycle 1's id-form defect is fixed; each with one finding per test and no '# ...' comment text, and ship check --test-failures exited 0. gh-prefixed synthetic logs: Minitest gave 'bin/rails test test/models/x_test.rb:12', pytest 'FAILED tests/test_x.py::test_y - AssertionError' gave 'pytest tests/test_x.py::test_y' (both --test-failures exit 0); a Jest 'FAIL src/x.test.js' plus '1 failure' log and an infra-only 'Error: Process completed with exit code 1.' log gave no tests, and --test-failures exited 1 with 'no failing test was recorded for the head: an infra or unparsable failure.'

## Links

### Out

- verifies [ship-check-runner-test-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-runner-test-failures.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests

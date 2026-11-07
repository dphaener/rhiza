# ship-check-runner-test-failures

capability · `01M6TBDK41SM1XY2928W79FEYZ` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

TestParseFailedTests in internal/api parses sample gh --log-failed output, each line carrying gh's job/step/timestamp prefix: an RSpec log with a 'Failed examples:' block (from nebari#72's shape, 'rspec ./spec/views/tasks/list_spec.rb:290 # Tasks list ...') into 'rspec ./spec/views/tasks/list_spec.rb:290', a Minitest log into 'bin/rails test test/models/x_test.rb:12', a pytest summary into 'pytest tests/test_x.py::test_y', a repeated line once, the existing Go cases unchanged, and logs with only '1 failure', a Jest 'FAIL src/x.test.js' line or an infra error into none; a testscript run-ship-rspec-failure with a fake gh whose failed run's log prints an RSpec 'Failed examples:' block shows rhiza ship check --json listing 'rspec ./spec/views/tasks/list_spec.rb:290' in tests and as a finding's test, and rhiza ship check --test-failures exiting 0; describe-golden regenerated for the tests field's description; run with go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship|describe-golden)' only

## statement

rhiza ship check reads failing tests of RSpec, Rails Minitest and pytest from a failed Actions run's gh run view --log-failed output as well as Go's: besides Go '--- FAIL: <name>' lines (leaves only, bare names as today), it matches, at the start of the log text after gh's job/step/timestamp prefix, RSpec's 'rspec ./<path>:<line>' lines, Minitest's 'rails test <path>:<line>' and 'bin/rails test <path>:<line>' lines and pytest's 'FAILED <path>::<test>' summary lines, and names each such failing test by the rerun command its runner printed, trimmed of any trailing comment ('rspec ./spec/views/tasks/list_spec.rb:290', 'bin/rails test test/models/x_test.rb:12', 'pytest tests/test_x.py::test_y'), deduplicated, unchanged in ship check's tests and each finding's Test, so ci-triage passes a Rails or Python project's red CI to ci-fix; a log with none of these lines (counts such as '1 failure', Jest's file-level 'FAIL <file>', infra errors) still ends as an infra or unparsable failure, and the tests field's description names Go, RSpec, Minitest and pytest

## Links

### Out

- depends-on [ship-check-ci-test-findings](ship-check-ci-test-findings.md)
- decided-by [decision-ci-triage-runners](../../../../../decision/decision-ci-triage-runners.md)

### In

- verifies [ship-check-runner-test-failures-verified](../../../../../evidence/ship-check-runner-test-failures-verified.md)
- verifies [ship-check-runner-test-failures-qa-1](../../../../../evidence/ship-check-runner-test-failures-qa-1.md)
- verifies [ship-check-runner-test-failures-verified-2](../../../../../evidence/ship-check-runner-test-failures-verified-2.md)
- verifies [ship-check-runner-test-failures-qa-2](../../../../../evidence/ship-check-runner-test-failures-qa-2.md)
- verifies [pr-350-fix-the-ship-change-ci-triage](../../../../../evidence/pr-350-fix-the-ship-change-ci-triage.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests

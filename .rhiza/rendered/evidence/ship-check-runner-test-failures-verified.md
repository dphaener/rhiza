# ship-check-runner-test-failures-verified

evidence · `01M6TFC3YBTV4EYEAWEJJZPMB0` · actor [claude-code-ad2d1657](../actor/claude-code-ad2d1657.md)

## result

PASS

## statement

Verifier claude-code-ad2d1657 at 502573ff: read the diff (internal/api/ship.go parseFailedTests plus logText, which strips gh's job/step/timestamp prefix and ANSI codes, then anchors rspec/rails test/bin/rails test/FAILED ...::... matches at the start of the text; Go '--- FAIL:' leaf handling untouched; Tests desc names Go, RSpec, Minitest and pytest). TestParseFailedTests passes: gh-prefixed RSpec 'Failed examples:' block -> 'rspec ./spec/views/tasks/list_spec.rb:290' (comment trimmed, repeat once), Minitest -> 'bin/rails test test/models/x_test.rb:12', pytest -> 'pytest tests/test_x.py::test_y'; '1 failure', Jest 'FAIL src/x.test.js' and an infra error -> none; existing Go cases unchanged (diff only adds lines). go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship|describe-golden)' ok (219s): the new run-ship-rspec-failure script shows ship check --json listing the rspec test and a finding test, --test-failures exits 0, and describe-golden matches the new description. go test -trimpath ./internal/api/... passes (649s) in a scratch worktree at this commit with RHIZA_ROOT/CHANGE/RUN unset. In the run checkout it hit the 10m timeout under load ~10 on 4 cores, and then TestSessionRunRecordsLimitHit failed there because of the run environment; that test passes at both this commit and its parent once outside it, and the commit touches no usage code

## Links

### Out

- verifies [ship-check-runner-test-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-runner-test-failures.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests

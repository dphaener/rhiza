# ship-check-runner-test-failures-verified-2

evidence · `01M6TGGYV4S9XAA733ETVE3YDN` · actor [claude-code-1844e741](../actor/claude-code-1844e741.md)

## result

PASS

## statement

Verifier, attempt 2, at b9f87b04: read the leaf's diff from the merge-base (internal/api/ship.go, ship_test.go). The fix commit's rspecLine `^rspec (\./\S+?(?::\d+|\[[\d:,]+\]))(?:\s|$)` keeps RSpec's example-id form whole, so the attempt-1 finding (closing ']' dropped) is fixed. go test -trimpath ./internal/api -run TestParseFailedTests -v PASS; it covers prefixed RSpec 'Failed examples:' lines -> 'rspec ./spec/views/tasks/list_spec.rb:290' (plus ANSI, deduped, and the id forms [1:1:11] and [1:2,1:3]), Minitest 'bin/rails test test/models/x_test.rb:12' and 'rails test ...', pytest FAILED -> 'pytest tests/test_x.py::test_y', the Go nested and dedup cases unchanged, and none for '1 failure', Jest 'FAIL src/x.test.js', infra errors and a mid-line rspec mention. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship|describe-golden)' ok (225s): run-ship-rspec-failure lists the rspec command in tests and as a finding's test, and --test-failures exits 0; describe-golden's tests description names Go, RSpec, Minitest and pytest. The full ./internal/api/... run had one failure, TestSessionRunRecordsLimitHit (usage_learn_test.go:119). It fails only in checkouts under ~/code: the same HEAD passes in a scratch worktree outside ~/code, and so does origin/main. The change set doesn't touch that code. The shipReadFailures doc comment still names only '--- FAIL:' lines, a minor stale comment outside acceptance.

## Links

### Out

- verifies [ship-check-runner-test-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-runner-test-failures.md)

## History

- `01M6TB6N9BYGB9S6ZF3HS4EDWM` Fix the ship-change ci-triage step can't read RSpec failures: on nebari#72 (run 01M6T4EH) one RSpec example failed (rspec ./spec/views/tasks/list_spec.rb:290, '1 failure' in the job log) but ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure' and failed the run, so a Rails project's red CI never reaches ci-fix; ci-triage must parse RSpec's 'Failed examples:' / 'rspec ./path:line' lines (and other common runners' output) as failing tests

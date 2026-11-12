# ci-logs-to-fixer-verified

evidence · `01M77JE5XEWKXR33GSNQGV9DWY` · actor [claude-code-9d8c013b](../actor/claude-code-9d8c013b.md)

## result

PASS

## statement

At 5b00592 (attempt 2) I re-ran the acceptance: go test -trimpath -count=1 ./internal/api -run 'TestShip|TestTrimFailureLog' -v. Every test passed, including TestShipCheckSavesFailedLog: fake gh returns a failing Jest job's log, and after the head's one rerun ship check records one finding naming ci-logs/…42.log, and that saved file holds 'FAIL src/garden.test.js' and 'expect(received).toBe(2)'. TestTrimFailureLog and TestShipCheckStabilityRecord also passed. Attempt 1's findings are fixed. f1: the flake record now names 'test / Run' and the test is updated. f2: failureMarker now matches GNU make's '*** [Makefile:163: fmt-check] Error 1' and gofmt's 'unformatted files:', and new cases cover them plus an RSpec 'Failed examples:' block. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-ci-fix|run-ship-ci-lint-fix|run-ship-rspec-failure|ship-check-ci-test-findings|describe-golden)$' printed ok. grep finds no Go/RSpec/Minitest/pytest log parser left in non-test code. The package run exited non-zero only from internal/rootguard: ~/code/.rhiza-root/usage.live.json changed, a file the live root daemon and the Claude status line rewrite every ~15s. It is missing from rootguard's owned list, a pre-existing gap unrelated to this leaf.

## Links

### Out

- verifies [ci-logs-to-fixer](../capability/rhiza/orchestration/shipping/ci-logs-to-fixer.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.

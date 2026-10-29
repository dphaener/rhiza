# read-apps-match-cli-verified-2

evidence · `01M6305WBGMZFQNMH2YG65J9TQ` · actor [claude-code-655422e2](../actor/claude-code-655422e2.md)

## result

PASS

## statement

At 4457f88, attempt 2: go test ./cmd/rhiza -run TestScripts/tui-read-apps -count=1 passes. On a fixture model the script cmps `rhiza tui app <name> --once` against `rhiza <name> --json` for roadmap, ready, show, why and slice, plus a new slice --change case with another change set current, and asserts no daemon.lock, no daemon.state.json and no socket under $WORK. The new case run against parent dd3d5b0 fails at line 50 (stdout differs), so it covers quality finding f1. By hand, with RHIZA_ROOT/CHANGE/RUN unset, `rhiza slice --change 01M62JMSW4FH8RAWSQXHC7ZPND --json` and `rhiza tui app slice --change 01M62JMSW4FH8RAWSQXHC7ZPND --once` both exit 0 with byte-identical output (23791 bytes). Diff read: run passes g to tuiCommand, which passes it to tuiApp, and newCtx now takes g instead of globals{}. go test ./..., fmt-check and vet pass; make self-check fails only on stale rendered actor files from this run's sessions, not on code.

## Links

### Out

- verifies [read-apps-match-cli](../capability/rhiza/read-apps-match-cli.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps

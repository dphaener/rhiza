# inbox-stalled-items-qa-2

evidence · `01M6AVR4428K6TA6A9RC1HRRBM` · actor [claude-code-2d4e2697](../actor/claude-code-2d4e2697.md)

## result

PASS

## statement

Built binary at 8276602a against this repo's own live ship-change run 01M6AREK66 (review/review, the cycle-1 gap): rhiza inbox --json at report+21m lists two items of kind stalled, one per working check session (quality c98484a2 and qa 2d4e2697), each with run, change, step review/review, session, at, and next 'rhiza session message <session> <text>'; at report+19m none. Human rhiza inbox prints '<id> stalled run 01M6AREK66..., change set 01M6A0K4..., step review/review / session <id> has not reported since <at> / next: rhiza session message <id> <text>' per session. go test -run TestScripts/inbox-stalled passes (blocked and idle give none) and go test ./internal/tui/inboxapp passes (Stalled group).

## Links

### Out

- verifies [inbox-stalled-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-stalled-items.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox

# inbox-stalled-items-qa-1

evidence · `01M6ATQ7CD67FEF55BA5R319F2` · actor [claude-code-7b8ceb12](../actor/claude-code-7b8ceb12.md)

## result

FAIL

## statement

Same scratch run with the built binary at 973ce724, real clock: hook back-dated to 04:00Z, rhiza inbox prints '03aba66 stalled run 01M6ATPA..., change set 01M6ATPA..., step build / session 107f052b-... has not reported since 2026-11-01T04:00:00Z / next: rhiza session message 107f052b-... <text>' and --json has one item kind stalled with run, session, step, at 04:00:00Z and that next; turned blocked: no stalled item. Fresh report: none. Gap: in this repo's own ship-change run 01M6AREK66 (current step review/review), the built binary's run status --json gives stalled null and inbox --json lists no stalled item at RHIZA_TEST_NOW=2026-11-01T09:00:00Z although check session 7b8ceb12 is a live session whose hook is working and last reported 05:52Z: runStalledNow reads only st.Session[st.Current], and ship-change records its builder, verifier and check sessions as leaf-session/check-session entries, so they never stall.

## Links

### Out

- verifies [inbox-stalled-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-stalled-items.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox

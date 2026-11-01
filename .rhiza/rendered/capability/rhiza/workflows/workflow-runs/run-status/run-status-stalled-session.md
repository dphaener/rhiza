# run-status-stalled-session

capability · `01M6A0NWP7AH4JGMDT37F6X7DG` · parent [run-status](../run-status.md)

## acceptance

A Go test of the stall cutoff: working at 19 minutes is not stalled, at 21 minutes is, blocked, idle and usage-paused sessions never are; a testscript with a hand-written session record and RHIZA_TEST_NOW shows 'stalled' in run status only past the cutoff

## statement

rhiza run status marks a live session stalled when its hook status is working and its last report is more than 20 minutes old, unless it is paused for a usage limit

## Links

### Out

- decided-by [stalled-after-20-minutes](../../../../../decision/stalled-after-20-minutes.md)

### In

- verifies [run-status-stalled-session-verified](../../../../../evidence/run-status-stalled-session-verified.md)
- verifies [run-status-stalled-session-qa-1](../../../../../evidence/run-status-stalled-session-qa-1.md)
- verifies [run-status-stalled-session-verified-2](../../../../../evidence/run-status-stalled-session-verified-2.md)
- verifies [run-status-stalled-session-qa-2](../../../../../evidence/run-status-stalled-session-qa-2.md)
- verifies [pr-158-stalled-sessions-surfaced](../../../../../evidence/pr-158-stalled-sessions-surfaced.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox

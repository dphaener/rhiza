# inbox-stalled-items

capability · `01M6A0NY4FZ07ZQMR23E700DTQ` · parent [inbox-list](../inbox-list.md)

## acceptance

inbox-stalled.txtar: with a working session record last reported at T and a live run journal, rhiza inbox --json has no stalled item at T+19m and one of kind stalled at T+21m; a blocked or idle record gives none; the inbox app test shows the stalled group

## statement

rhiza inbox lists each stalled session of a live run as a stalled item naming the session and its run, with rhiza session message as the next step, and the inbox app shows them

## Links

### Out

- decided-by [stalled-after-20-minutes](../../../../../decision/stalled-after-20-minutes.md)

### In

- verifies [inbox-stalled-items-verified](../../../../../evidence/inbox-stalled-items-verified.md)
- verifies [inbox-stalled-items-qa-1](../../../../../evidence/inbox-stalled-items-qa-1.md)
- verifies [inbox-stalled-items-verified-2](../../../../../evidence/inbox-stalled-items-verified-2.md)
- verifies [inbox-stalled-items-qa-2](../../../../../evidence/inbox-stalled-items-qa-2.md)
- verifies [pr-158-stalled-sessions-surfaced](../../../../../evidence/pr-158-stalled-sessions-surfaced.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox

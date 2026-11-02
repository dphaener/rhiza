# inbox-nudge-stalled

capability · `01M6CS2KAT7R41TVN9WFDEFBRF` · parent [inbox-stalled-items](../inbox-stalled-items.md)

## acceptance

An inbox app unit test on a fake runner: n on a stalled item sends the carry-on message to the session; for a session whose last pane text shows a usage-limit prompt it sends Escape before the message; n on other kinds does nothing

## statement

On a stalled item the inbox app's n nudges the session: when it is parked on a usage-limit prompt it sends Escape first, then a message asking it to carry on, through the session's adapter

## Links

### Out

- decided-by [autonomy-app-controls](../../../../../../decision/autonomy-app-controls.md)

### In

- verifies [inbox-nudge-stalled-verified](../../../../../../evidence/inbox-nudge-stalled-verified.md)
- verifies [inbox-nudge-stalled-qa-1](../../../../../../evidence/inbox-nudge-stalled-qa-1.md)
- verifies [inbox-nudge-stalled-verified-2](../../../../../../evidence/inbox-nudge-stalled-verified-2.md)
- verifies [inbox-nudge-stalled-qa-2](../../../../../../evidence/inbox-nudge-stalled-qa-2.md)
- verifies [pr-195-autonomy-app](../../../../../../evidence/pr-195-autonomy-app.md)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions

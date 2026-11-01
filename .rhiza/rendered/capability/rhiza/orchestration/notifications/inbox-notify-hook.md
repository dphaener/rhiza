# inbox-notify-hook

capability · `01M6BDAMMR2G2KYT2ZTF2HK9XX` · parent [notifications](../notifications.md)

## acceptance

daemon-notify-hook.txtar: a fake command appending stdin to a file, over the four items daemon-notify.txtar uses, writes four JSON lines with matching kinds and env; a second daemon --once writes none; an item that clears and returns runs it again; a command exiting 1 is logged with the item ID and the poll succeeds; a command past the timeout (shortened in tests through a RHIZA_TEST_ variable honoured only with RHIZA_TEST_NOW) is killed and logged; without Rex the hook still runs once per item; git config rhiza.notifyCommand is used when the flag is absent; describe-golden.txtar is regenerated

## statement

The daemon runs the owner's notify command (--notify-command, or git config rhiza.notifyCommand) once per new inbox item, with or without Rex, giving it the item as JSON on stdin and its kind, ID, change set and notice text in RHIZA_INBOX_* variables; a failing or timed-out command is logged and not retried

## Links

### Out

- decided-by [notify-command-flag-or-git-config](../../../../decision/notify-command-flag-or-git-config.md)

### In

- verifies [inbox-notify-hook-verified](../../../../evidence/inbox-notify-hook-verified.md)
- verifies [inbox-notify-hook-qa-1](../../../../evidence/inbox-notify-hook-qa-1.md)
- verifies [pr-172-inbox-notify-hook](../../../../evidence/pr-172-inbox-notify-hook.md)

## History

- `01M6BD9HT87Z8WCD31EKQQ47R1` The daemon runs an owner-configured command once per new inbox item, with or without Rex

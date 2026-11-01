# notify-command-flag-or-git-config

decision · `01M6BD9ZSCTSFJ66DZEYPD6P0P`

## statement

The daemon runs a notify command, from --notify-command or else git config rhiza.notifyCommand, through sh -c in the repo root once per new inbox item, with the item's JSON on stdin and RHIZA_INBOX_KIND, RHIZA_INBOX_ID, RHIZA_INBOX_CHANGE and RHIZA_INBOX_TITLE set; it records hooked items in its own state apart from Rex notifications so it runs with or without Rex; a failing or timed-out (10 s) command is logged with the item ID and not retried

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m14-orchestrator-picks](../source/m14-orchestrator-picks.md)

### In

- decided-by [inbox-notify-hook](../capability/rhiza/orchestration/notifications/inbox-notify-hook.md)

## History

- `01M6BD9FT5SGV1CMZXV25XZCMR` Plan milestone 14: find it, hear about it, picked by the orchestrator under the owner's standing instruction

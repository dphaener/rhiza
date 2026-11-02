# terminal-current-next

capability · `01M6CAAEEEEZ1XHKEA811180P7` · parent [current-change](../current-change.md)

## acceptance

status-terminal-current.txtar: land change set A with --change, write A into .rhiza/local/change by hand; rhiza status shows (landed) and next: rhiza change use --none; rhiza status --hook prints the hint; a model write exits with a conflict carrying the same next; after rhiza change use --none, status shows none current

## statement

When the current change set has landed, been withdrawn or been rejected, rhiza status and the session hook say so with next: rhiza change use --none, and a write refused for it carries the same next step

## Links

### Out

- decided-by [terminal-current-hint](../../../../decision/terminal-current-hint.md)

### In

- verifies [terminal-current-next-verified](../../../../evidence/terminal-current-next-verified.md)
- verifies [terminal-current-next-qa-1](../../../../evidence/terminal-current-next-qa-1.md)
- verifies [pr-186-status-tells-the-truth](../../../../evidence/pr-186-status-tells-the-truth.md)

## History

- `01M6CA9QMHN6N95AS3F935ZVN2` daemon status keeps the daemon's state through a re-exec, and a terminal current change set comes with the command that clears it

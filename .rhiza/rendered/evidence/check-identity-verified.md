# check-identity-verified

evidence · `01M5V9A637FQP6C6SYNVCPHS1K`

## result

PASS under make check on branch cs/rename-to-rhiza

## statement

The check unit tests in internal/check and internal/api, and the check-move.txtar testscript, show check judging log entries by ID and bytes, so an entry moved from .genesis/log to .rhiza/log passes and an edited or deleted one fails. They pass.

## tests

go test: internal/check (all), internal/api TestCheckInputPreRename; testscript: check-move.txtar

## Links

### Out

- verifies [check-entry-identity](../capability/rhiza/integrity/check/check-entry-identity.md)

## History

- `01M5TW3X2FX0DVZ9XZGXMQVC2Q` The rename to Rhiza: the module, repo, binary, .rhiza/ folder, RHIZA_* variables, skills, guide, cache and root handle all become rhiza, with check identifying entries by ID and bytes so the log can move.

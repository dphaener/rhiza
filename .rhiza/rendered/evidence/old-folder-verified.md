# old-folder-verified

evidence · `01M5V9A69GJR7P3AR7V4DHEHZ4`

## result

PASS under make check on branch cs/rename-to-rhiza

## statement

The old-folder.txtar testscript and the store tests show a repo with only a .genesis folder getting the old-folder notice (conflict, exit 3, Next git mv .genesis .rhiza) instead of acting as uninitialised. They pass.

## tests

go test: internal/store (ErrOldFolder); testscript: old-folder.txtar

## Links

### Out

- verifies [old-folder-notice](../capability/rhiza/model/repo-init/old-folder-notice.md)

## History

- `01M5TW3X2FX0DVZ9XZGXMQVC2Q` The rename to Rhiza: the module, repo, binary, .rhiza/ folder, RHIZA_* variables, skills, guide, cache and root handle all become rhiza, with check identifying entries by ID and bytes so the log can move.

# note-add

capability · `01M6BZWDXPERND94DYYXHTQ8WK` · parent [change-notes](../change-notes.md)

## acceptance

change-notes.txtar: note add on a node, on a change set, and on a change set with --op each pass fold-ok and rhiza check; a note on an unknown target exits 2; --op naming an operation not in the target change set exits 2

## statement

rhiza note add <target> <text> [--op id] writes a note discussing a node or change set (and naming an operation with --op) into the current change set in one entry

## Links

### Out

- decided-by [notes-are-note-nodes](../../../../decision/notes-are-note-nodes.md)

### In

- verifies [note-add-verified](../../../../evidence/note-add-verified.md)
- verifies [note-add-qa-1](../../../../evidence/note-add-qa-1.md)
- verifies [note-add-verified-2](../../../../evidence/note-add-verified-2.md)
- verifies [pr-197-notes-on-a-change](../../../../evidence/pr-197-notes-on-a-change.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

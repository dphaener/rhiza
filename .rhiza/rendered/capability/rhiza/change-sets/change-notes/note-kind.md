# note-kind

capability · `01M6BZWB6CR3ERK0Z0DSVVVPEP` · parent [change-notes](../change-notes.md)

## acceptance

Go tests in internal/op and internal/check: a note with a discusses link to a node and one to a change set fold and check clean; a discusses link from a non-note is refused; the field schema knows note's fields; describe-golden.txtar and the guide's count of kinds and link types are updated

## statement

The model has a note node kind, with its body in statement and an optional op field, and a discusses link type from a note to any node or change set

## Links

### Out

- decided-by [notes-are-note-nodes](../../../../decision/notes-are-note-nodes.md)

### In

- verifies [note-kind-verified](../../../../evidence/note-kind-verified.md)
- verifies [note-kind-qa-1](../../../../evidence/note-kind-qa-1.md)
- verifies [pr-197-notes-on-a-change](../../../../evidence/pr-197-notes-on-a-change.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

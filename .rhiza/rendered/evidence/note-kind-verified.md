# note-kind-verified

evidence · `01M6DZBJ9DB091QYS2MR699RFW` · actor [claude-code-5ad3f5d4](../actor/claude-code-5ad3f5d4.md)

## result

PASS

## statement

At commit 363a99f8 I ran go test ./internal/op ./internal/check -run 'TestRunNote|TestNoteOps|TestDiscusses|TestNoteFields' -v -count=1 and all passed. TestRunNote (internal/check) folds a log where a note discusses a capability and its own change set, and Run reports OK with no findings. A capability discussing a note gives one RuleLink fail ('"discusses" does not go from capability to note'). TestNoteOps (internal/op) decodes, validates and round-trips a note node.create and a discusses link.create. TestNoteFields covers knownFields[note] = op, statement. Reading the diff of 3727472c: describe-golden.txtar is updated, guide apply.md and rhiza.md list note among the kinds and discusses among the link types, and the op.go counts say ten kinds and fifteen link types. make check exited 0. This fixes attempt-1 finding f1.

## Links

### Out

- verifies [note-kind](../capability/rhiza/change-sets/change-notes/note-kind.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

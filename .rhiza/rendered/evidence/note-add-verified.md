# note-add-verified

evidence · `01M6DZVHJ8G7WZ965HC45TRFD0` · actor [claude-code-49d98852](../actor/claude-code-49d98852.md)

## result

PASS

## statement

At 6addeb28 (detached checkout, RHIZA_ROOT/CHANGE/RUN unset): go test ./cmd/rhiza -run 'TestScripts/(change-notes|describe-golden)$' -count=1 passes; change-notes.txtar runs note add on a node, on a change set, and on a change set with --op, each followed by fold-ok, render and rhiza check, plus the unknown-target and foreign --op failures. make check exits 0 (self-check clean, so the attempt-1 render finding is fixed). By hand in a scratch repo with every RHIZA_* var unset: note add on a capability wrote one entry (a note node and a link) and exited 0; note add nonesuch exited 2 ('no node or change set'); --op 01ZZZZZZZZZZZZZZZZZZZZZZZZ on the change set exited 2 ('is not an operation of change set'); render then rhiza check printed ok.

## Links

### Out

- verifies [note-add](../capability/rhiza/change-sets/change-notes/note-add.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

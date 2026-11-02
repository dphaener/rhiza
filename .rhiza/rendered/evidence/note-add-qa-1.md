# note-add-qa-1

evidence · `01M6E0H0P7Q0C1WJXS7BQ1CSEA` · actor [claude-code-e798503e](../actor/claude-code-e798503e.md)

## result

PASS

## statement

QA at 071ed6b8, binary built with go build into a scratch dir, run in a fresh git repo with RHIZA_ROOT/CHANGE/RUN/ACTOR unset. rhiza note add store '...' wrote one entry (node note-store plus a link) and exited 0; note add <change set id> wrote a note and exited 0; note add <change set id> --op <its node.create op> exited 0; note add nonesuch exited 2 with 'no node or change set "nonesuch"'; --op 01ZZZZZZZZZZZZZZZZZZZZZZZZ exited 2 with 'is not an operation of change set'; render then rhiza check printed ok (exit 0). make check exit 0.

## Links

### Out

- verifies [note-add](../capability/rhiza/change-sets/change-notes/note-add.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

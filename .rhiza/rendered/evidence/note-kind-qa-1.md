# note-kind-qa-1

evidence · `01M6E0H2CEWBPJY4Z70MVB75GW` · actor [claude-code-e798503e](../actor/claude-code-e798503e.md)

## result

PASS

## statement

QA at 071ed6b8, binary built with go build into a scratch dir, run in a fresh git repo with RHIZA_ROOT/CHANGE/RUN/ACTOR unset. change show --json listed a node.create of kind note with its body in fields.statement and a link of type discusses from the note to the capability; a note on a change set with --op carried op in change show --json; apply of a discusses link from capability store to a note exited 2 with 'link type "discusses" does not go from capability to note'; rhiza guide lists note among the kinds; describe --json documents discusses; rhiza check clean. make check exit 0.

## Links

### Out

- verifies [note-kind](../capability/rhiza/change-sets/change-notes/note-kind.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

# notes-shown-qa-1

evidence · `01M6E0H44P8AV7EGNM13RK5PNT` · actor [claude-code-e798503e](../actor/claude-code-e798503e.md)

## result

PASS

## statement

QA at 071ed6b8, binary built with go build into a scratch dir, run in a fresh git repo with RHIZA_ROOT/CHANGE/RUN/ACTOR unset. After two notes on capability store and two on the change set (one with --op), rhiza show store printed a Notes section with 'Is the log append-only?' then 'Second thought.', each as '- <id> by darin-haener at 2026-11-02T11:22:..Z: <text>'; show store --json had node.notes with actor, at, id, text in that order. rhiza change show printed both change-set notes in order, the second marked '(on op 01M6DZXGM8GHA7)'; change show --json notes carried actor, at, id, op, text.

## Links

### Out

- verifies [notes-shown](../capability/rhiza/change-sets/change-notes/notes-shown.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

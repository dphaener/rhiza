# notes-shown-verified-2

evidence · `01M6E299349SDCBKRN2CP1991H` · actor [claude-code-f79773a9](../actor/claude-code-f79773a9.md)

## result

PASS

## statement

At c3b2557f ran go test ./cmd/rhiza -run TestScripts/change-notes -count=1: PASS. The script adds two notes on node store, then show store lists 'Is the log append-only?' before 'Second thought.' with '- <id> by <actor> at <RFC3339>Z:' lines and a notes array in --json; change show lists the change set's 'Please split this.' and the --op note with actor and time, in text and --json. The new case writes a note about that change set from a second current change set; it shows in change show of the discussed set (text and --json) and not in the writer's own Notes. Read the diff: changeNotes now reads cs's view plus the current one's, de-duplicating by note ID. make fmt-check vet test all pass. make self-check fails only on stale .rhiza/rendered actor/index files left by the run's own model writes, not on this code.

## Links

### Out

- verifies [notes-shown](../capability/rhiza/change-sets/change-notes/notes-shown.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

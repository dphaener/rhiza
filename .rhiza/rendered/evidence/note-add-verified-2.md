# note-add-verified-2

evidence · `01M6E1ANN7TV18TKJZ6C4PR7C9` · actor [claude-code-563f881d](../actor/claude-code-563f881d.md)

## result

PASS

## statement

At 50167ed1 (attempt 2): go test ./cmd/rhiza -run TestScripts/change-notes -count=1 passes (note add on a node, on a change set and with --op each followed by fold-ok, render and check; unknown target and foreign --op fail). By hand with a scratch build in a fresh repo: note add on a node, a change set and a change set with --op of one of its operations each exit 0 and write node+link in one entry with handle note-<short id> (e.g. note-01m6e18wbm1qs4); rhiza render then rhiza check is ok; note add nonesuch exits 2 (no node or change set); --op 01ZZZZZZZZZZZZZZZZZZZZZZZZ exits 2 (not an operation of change set). Diff read: noteHandle now derives from the drawn node ID via an unexported AddRequest.handleFor hook, matching decision notes-are-note-nodes. make check: go tests pass; self-check fails only on stale .rhiza/rendered actor/index files from the run's own chore writes, unrelated to this leaf.

## Links

### Out

- verifies [note-add](../capability/rhiza/change-sets/change-notes/note-add.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

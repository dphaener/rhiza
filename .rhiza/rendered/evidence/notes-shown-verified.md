# notes-shown-verified

evidence · `01M6DZ3CA0RP63JHWV8KRTXAYN` · actor [claude-code-5e334fa7](../actor/claude-code-5e334fa7.md)

## result

PASS

## statement

At 6de27e61, go test ./cmd/rhiza -run 'TestScripts/(change-notes|describe-golden)$' -count=1 passes, and so does go test ./... -count=1. change-notes.txtar adds a second note on store, and checks that rhiza show store lists 'Is the log append-only?' before 'Second thought.' under ## Notes, with '- <id> by <actor> at <RFC3339>Z:' lines. It checks the notes array in show --json, and that change show lists 'Please split this.' then 'About this op.' with by/at. By hand, a scratch build in a fresh repo had two notes on a node one second apart: show printed '- … by darin-haener at 2026-11-02T11:07:40Z: First.' then '…11:07:41Z: Second.'. change show listed the change set's own note 'On the set.' with actor and time. show --json gave notes with id/text/actor/at in order. Diff read: internal/api/notes_show.go (showNotes, changeNotes, writeNotes, sorted by time then ID), wired into show.go and change.go. make check's fmt, vet and tests pass. Its self-check fails only on rendered actor files for run sessions registered after this commit.

## Links

### Out

- verifies [notes-shown](../capability/rhiza/change-sets/change-notes/notes-shown.md)

## History

- `01M6BVTCJCSG1JKEERK7Y3ACYJ` Notes attached to a node, a change set or one of its operations hold the discussion, shown by show and change show

# spawned-child-question-arbitrated-qa-1

evidence · `01M7M99RAPBPGZKXFX9RH0YD8G` · actor [claude-code-2166ab54](../actor/claude-code-2166ab54.md)

## result

PASS

## statement

Built ./cmd/rhiza (HEAD 289c5fa77) to a scratch path; ran it in a scratch Rhiza root with a fake claude (real hook path) and fake rex, runs not started --autonomous. Conduct step whose conductor spawned a builder child that asked: asked entry keyed on the child session, arbiter pass, the answer sent to the child (messages/<child>/0001.md), answered and unblocked entries, then spawn-ended, and the step passed on the conductor's own outcome; no advisor and no blocked notice. A child asking two questions in turn got both arbitrated and answered (2 asked/arbiter/answered entries). Failure case: arbiter wrote no Decision on the child's question, so the step ended in error naming the question and the run failed. A child spawned in a plain step was not watched: no asked or arbiter entries. Scoped tests: go test ./internal/run -run TestSpawnedChildQuestion ok.

## Links

### Out

- verifies [spawned-child-question-arbitrated](../capability/rhiza/workflows/arbiter-sessions/spawned-child-question-arbitrated.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand

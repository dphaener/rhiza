# roadmap-app-list-qa-1

evidence · `01M6765P9XT45579T5D2410BY2` · actor [claude-code-a83eb6cd](../actor/claude-code-a83eb6cd.md)

## result

PASS

## statement

Built this branch's rhiza (go build -o $TMPDIR/rhiza-qa ./cmd/rhiza) and, in a scratch git repo with a proposed and an accepted-but-question-blocked change set, ran rhiza tui app roadmap in a real Rex tab driven with rex send-key. Saw the Backlog heading and one row per change set with short ID, [state], title and '(not ready: ...)' reasons, no JSON; j/k and down/up moved the cursor between change set rows; r reloaded keeping the cursor (a new change set appeared); q quit with exit 0. rhiza tui app roadmap --once output was byte-equal (cmp) to rhiza roadmap --json.

## Links

### Out

- verifies [roadmap-app-list](../capability/rhiza/roadmap-app-list.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab

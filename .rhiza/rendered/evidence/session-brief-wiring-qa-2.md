# session-brief-wiring-qa-2

evidence · `01M7B24D7FX8HXS9FXEG2TB745` · actor [claude-code-a1cb0102](../actor/claude-code-a1cb0102.md)

## result

PASS

## statement

QA cycle 2: agent, builder, check, critic and both fixer sessions each got RHIZA_BRIEF and RHIZA_SLICE (no RHIZA_HANDBACK, no handback.md); slice.json byte-equal to rhiza slice --json run inside the session; an unwritable runs/<id> (a file at the path) fails the run naming the session directory path and starts no agent. go test ./internal/api -run TestSessionBrief and TestScripts/brief-files ok.

## Links

### Out

- verifies [session-brief-wiring](../capability/rhiza/workflows/workflow-runs/session-brief-wiring.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps

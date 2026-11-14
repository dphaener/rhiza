# session-brief-wiring-qa-1

evidence · `01M7AWMP88STC84PMD24CEZW7V` · actor [claude-code-1707fde5](../actor/claude-code-1707fde5.md)

## result

FAIL

## statement

QA: agent, builder, verifier, check, critic, fixers all got RHIZA_BRIEF and RHIZA_SLICE, no RHIZA_HANDBACK; slice.json byte-equal to rhiza slice --change --json run by the agent at start; session start outside a run got none; an unwritable session dir failed the run naming the path with no agent started. But brief.md lists Decisions none and Constraints none when the change set has its own decision and constraint (finding f1, blocker).

## Links

### Out

- verifies [session-brief-wiring](../capability/rhiza/workflows/workflow-runs/session-brief-wiring.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps

# step-session-waits-holds-qa-2

evidence · `01M7QKYRF83M1DQQ3BQ3DC84J0` · actor [claude-code-608090ea](../actor/claude-code-608090ea.md)

## result

PASS

## statement

QA cycle 2, scratch build: with a memory-low hold saved, rhiza run start --foreground of a one-agent-step workflow journaled one spawn-held entry for the step with no parent and never called the agent (no claude invocation); run resume while still held added no second entry; after the hold was released, resume started the step's session, whose session entry released the hold, and the run succeeded

## Links

### Out

- verifies [step-session-waits-holds](../capability/rhiza/workflows/workflow-runs/step-session-waits-holds.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it

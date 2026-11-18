# conduct-no-wait-guide-qa-1

evidence · `01M7QFVMCKP412GQQHEYS3R2ZA` · actor [claude-code-c1a3776c](../actor/claude-code-c1a3776c.md)

## result

PASS

## statement

QA with a fresh build: rhiza guide conduct tells the conductor to spawn with --no-wait, start independent children together, end its turn without sleep or poll and act on each wake-up via rhiza run status --json, keeping the never sleep/poll rule; rhiza guide conduct-roles documents --no-wait, the wake message and the waiting default for other callers. internal/guide tests, TestDocsPrimitives and TestScripts/guide-install pass.

## Links

### Out

- verifies [conduct-no-wait-guide](../capability/rhiza/workflows/workflow-runs/conduct-no-wait-guide.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)

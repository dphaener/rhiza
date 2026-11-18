# spawn-wake-conductor-qa-1

evidence · `01M7QFVMAC94WYEP2R3FYQRZAJ` · actor [claude-code-c1a3776c](../actor/claude-code-c1a3776c.md)

## result

PASS

## statement

QA with a fresh build: two no-wait children (exit 1, said done) each got one spawn-ended entry and one wake; the parent's message file names child, role, leaf, outcome, detail and rhiza run status --json; the said-done child's pane closed, the exited child's left open. A working parent got both pointers queued at once; an exited or prompt-less parent's wake failure was logged, not retried, and the run succeeded. A waiting child's end woke no one. TestSpawnWake, TestSpawnResume, TestRunSpawnWake pass.

## Links

### Out

- verifies [spawn-wake-conductor](../capability/rhiza/workflows/workflow-runs/spawn-wake-conductor.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)

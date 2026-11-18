# spawn-no-wait-qa-1

evidence · `01M7QFVM8AF1SXGVN7AVCMBF0H` · actor [claude-code-c1a3776c](../actor/claude-code-c1a3776c.md)

## result

PASS

## statement

QA with a fresh build in a scratch root (fake rex running the child in the background): rhiza run spawn --no-wait returned at once with rc 0, session, brief and no_wait:true, no outcome; journal spawn entry no_wait:true, spawn-ended only when the child ended. Missing brief rc 2, failed rex start rc 1 with no spawn entry; the default waiting spawn unchanged (pass (exit 0), no no_wait, no wake). Scoped tests TestRunSpawn, TestSpawn|TestFold, TestScripts/(run-spawn|describe-golden) pass.

## Links

### Out

- verifies [spawn-no-wait](../capability/rhiza/workflows/workflow-runs/spawn-no-wait.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)

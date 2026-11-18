# spawn-no-wait-wake

decision · `01M7QDZZMB67DX13M3TJ8Y67BD`

## reasoning

Opt-in keeps run.spawn's current contract (a blocking call that returns the outcome, exit 1 on fail) that the conductor skill, the describe golden and existing scripts rely on, so nothing breaks and no breaking: note is needed; the conductor skill can switch to --no-wait. The spawning CLI process exits under --no-wait, so it cannot send the message; the engine already polls every live spawn through ChildWatcher in reconcileSpawns and journals each end exactly once, which is the natural single point to send the parent a rhiza session message, and it works daemon or daemonless as long as the run's driver runs. For questions, internal/run/engine.go askChildren already routes a conductor child's question to the arbiter and delivers the answer to the child; waking the conductor for it would only interrupt it with nothing to do, so the intent's 'blocked question' clause is satisfied by the arbiter handling it, and the child's eventual end is what wakes the conductor.

## statement

Returning at start is opt-in: rhiza run spawn gains --no-wait, which returns once the child session is started and journaled (its response has the session and no outcome); without it spawn waits and prints the outcome as today. The run engine sends the wake-up: when reconcileSpawns journals a --no-wait child's spawn-ended entry (AppendSpawnEnded added it, so once per child), it runs rhiza session message to the child's parent session naming the child, its role, leaf, outcome and detail. A child's question does not wake the conductor: in a conductor step askChildren already has the arbiter answer it and delivers the answer to the child, so the conductor has nothing to act on; only the child's end (said done or exit) wakes it, and an arbiter failure on a child's question already ends the step. A child's other block (such as a permission prompt) stays the owner's, as now.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-dc492322-1](../question/question-dc492322-1.md)

### In

- decided-by [spawn-no-wait](../capability/rhiza/workflows/workflow-runs/spawn-no-wait.md)
- decided-by [spawn-wake-conductor](../capability/rhiza/workflows/workflow-runs/spawn-wake-conductor.md)
- decided-by [conduct-no-wait-guide](../capability/rhiza/workflows/workflow-runs/conduct-no-wait-guide.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)

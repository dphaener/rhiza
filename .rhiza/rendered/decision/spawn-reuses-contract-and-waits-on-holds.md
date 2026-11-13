# spawn-reuses-contract-and-waits-on-holds

decision · `01M7AHR2JAKP00Q8TYA0222B74`

## reasoning

The sibling change set's leaves already define brief.md with 'the contract section', handback.md 'on a hand-back only', and the step's handback text (run-brief-files, agent-invoke-template); building them twice here would duplicate code and let the two drift, so this change set should depend on that one and reuse its writers. The direction decision (conductor-milestone-direction) asks for a dumb daemon of few primitives, which also favours one brief writer. On waiting: maxRuns counts runs (internal/api/daemon.go picks runs under maxRuns via daemon.RoundRobin); counting a conductor's children against it would deadlock a run at maxRuns=1, since the conductor holds the only slot while waiting on its child. Usage and load holds protect the shared machine and the plan's limits regardless of which run asks, so a child must respect them, and run-spawn's acceptance already expects 'a hold while a usage hold is on'. Reading 'counts against the slot' as the parent's slot keeps that phrase true without a deadlock.

## statement

The contract change set ('Every agent session a run starts gets one gate-checked contract…') owns the pieces both name: run-brief-files writes brief.md with the contract section and handback.md, and agent-invoke-template owns the step's handback text and its placeholders. This change set depends on it and calls those writers: run-spawn writes the spawned child's brief through run-brief-files' writer (so it carries the contract section), and run-handback calls the handback.md writer and sends the text agent-invoke-template defines, adding only the routing, resume and handback_cap. rhiza run spawn waits on the usage holds and the machine (load) hold, like a new run start, but not on rhiza.daemon.maxRuns or rhiza.project.maxRuns: the parent run already holds its slot, and a spawned child counts as part of that run, not as a second run. While a hold is on, spawn waits (journaled as held) rather than failing.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-90768e1a-1](../question/question-90768e1a-1.md)
- derived-from [plan-daemon-offers-spawn](../source/plan-daemon-offers-spawn.md)

### In

- decided-by [run-spawn](../capability/rhiza/workflows/workflow-runs/run-spawn.md)
- decided-by [run-handback](../capability/rhiza/workflows/workflow-runs/run-handback.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox

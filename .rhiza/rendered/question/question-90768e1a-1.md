# question-90768e1a-1

question · `01M7AHK97SND7NV5B8ANW9WJSS`

## session

90768e1a-8457-4f12-95ad-01a1d5147827

## statement

Change set 'The daemon offers spawn' overlaps sibling change set 'Every agent session a run starts gets one gate-checked contract…' (leaves run-brief-files and agent-invoke-template), and no depends-on link orders them. Here, run-spawn adds 'the contract section' to the spawned brief, and run-handback 'writes handback.md' and sends the step's 'handback text'. Who owns those pieces? / run-spawn says a spawned session 'counts against the slot and usage holds'. Today rhiza.daemon.maxRuns counts runs, and the daemon's machine (load) hold and usage holds gate new run starts. What should rhiza run spawn wait on?

## Links

### In

- answers [spawn-reuses-contract-and-waits-on-holds](../decision/spawn-reuses-contract-and-waits-on-holds.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox

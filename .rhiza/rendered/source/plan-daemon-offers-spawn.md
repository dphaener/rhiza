# plan-daemon-offers-spawn

source · `01M7AHV4XYBX80SJC5VM6N9KZH`

## quotes

```json
[
  {
    "by": "arbiter (assumed, would have escalated)",
    "text": "This change set depends on it and calls those writers: run-spawn writes the spawned child's brief through run-brief-files' writer (so it carries the contract section), and run-handback calls the handback.md writer and sends the text agent-invoke-template defines, adding only the routing, resume and handback_cap.",
    "where": "answer to question 01M7AHK97S, 2026-11-13T13:32Z"
  },
  {
    "by": "arbiter (assumed, would have escalated)",
    "text": "rhiza run spawn waits on the usage holds and the machine (load) hold, like a new run start, but not on rhiza.daemon.maxRuns or rhiza.project.maxRuns: the parent run already holds its slot, and a spawned child counts as part of that run, not as a second run. While a hold is on, spawn waits (journaled as held) rather than failing.",
    "where": "answer to question 01M7AHK97S, 2026-11-13T13:32Z"
  }
]
```

## raw_sha256

61c9ddd180f763aca185696d0fa5bff0fc023fdc88886f84ba762922ff81451b

## session

90768e1a-8457-4f12-95ad-01a1d5147827

## summary

Autonomous plan run on change set 'The daemon offers spawn', whose five leaves came from milestone 36's planning. Research: workflow step kinds live in internal/workflow/workflow.go (agent, wait, command, arbiter, leaves, checks, signoff, workflow; no gate); the engine (internal/run/engine.go) runs steps through the Runner, implemented by stepRunner in internal/api/run.go, whose StartLeaf (internal/api/run_leaves.go) starts a run's sessions through session.start; runs are driven by a detached 'rhiza run start --foreground' process, not the daemon process; commands like 'rhiza run review-check' append to a run's journal from their own process without the run lock; usage and machine holds (internal/api/usage_hold.go, daemonMachineHold) gate only new run starts today; ship.gate (internal/api/ship_gate.go) already keeps a gate's full output in a log file in the run directory; adapters can resume a session in the same conversation (adapter.Resume); sessions are messaged by session.message; run.status's RunInfo (internal/api/run.go) and the web run page (internal/web/rundetail.go) list sessions, leaves and checks but no spawned children. Overlaps found: the sibling change set 'Every agent session a run starts gets one gate-checked contract' owns brief.md with its contract section, handback.md and the handback text; draft change set 01M79KHBKT owns rhiza workflow check's warning about uncapped hand-back loops, so that clause left run-handback's acceptance. run-handback held two concerns (routing and capping, delivery and resume), so it was split into handback-route and run-handback. Two questions went to the arbiter (owner absent): who owns the contract pieces, and what spawn waits on. The arbiter ruled (Decision spawn-reuses-contract-and-waits-on-holds) that this change set depends on the contract change set and calls its writers, and that spawn waits on the usage and machine holds but not on maxRuns. Derived here: rhiza run spawn starts the child from its own process through the run's session-start path and journals into the run's journal without the lock, as run review-check does; the run's driver records children that ended while nothing watched them.

## title

Plan of The daemon offers spawn (autonomous plan run)

## Links

### In

- derived-from [spawn-reuses-contract-and-waits-on-holds](../decision/spawn-reuses-contract-and-waits-on-holds.md)
- derived-from [spawn-starts-from-caller](../decision/spawn-starts-from-caller.md)
- derived-from [handback-split-route-deliver](../decision/handback-split-route-deliver.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox

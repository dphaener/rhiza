# conductor-questions-plan

source · `01M7M5CYJJ84ERCPH5KHTA8BZT`

## quotes

```json
[]
```

## raw_sha256

8483c8fb6e6503ecd0fa4df962732029ee445711a36c1d99c5a2cfecb74bf378

## session

4060ff02-8832-4dfd-8bd3-01a1e8513fe0

## summary

Autonomous plan run, no owner present. Research: internal/run/engine.go arbitrates a question only when State.Autonomous, set only by run start --autonomous; conductor runs 01M7DVV7 and 01M7E56M (ship-rhiza, every rhiza.autonomy stage auto) started without it, so the conductor's AskUserQuestion launched an advisor and blocked on the owner. Hand-queued run start and daemonStart dropping --autonomous for an owner sign-off hit the same path. A child session started with rhiza run spawn is not watched for questions by anything: its denied question in an autonomous run is never arbitrated. The planner asked which fix to build; the arbiter decided (Decision conductor-questions-to-arbiter, wording (assumed, would have escalated)): questions from the conductor step and from its spawned sessions always go to the arbiter, with no advisor or owner notice, whatever --autonomous says; State.Autonomous keeps governing only the owner's gates and sign-off; do not force --autonomous on conductor runs.

## title

Conductor questions stall: autonomous plan of bug 01M7M51

## Links

### In

- derived-from [conductor-questions-to-arbiter](../decision/conductor-questions-to-arbiter.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand

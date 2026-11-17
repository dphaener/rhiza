# question-4060ff02-1

question · `01M7M5904F9KDG6PYQJ266WA1M`

## session

4060ff02-8832-4dfd-8bd3-01a1e8513fe0

## statement

Bug 01M7M51 (conductor questions stall): the engine sends a question to the arbiter only when State.Autonomous is set, and that is set only by `run start --autonomous`. Conductor runs 01M7DVV7 and 01M7E56M (ship-rhiza, every rhiza.autonomy stage `auto`) started without the flag, so the conductor's AskUserQuestion launched an advisor and the run waited for the owner. The same happens for a hand-queued `rhiza run start <cs>`, and when daemonStart drops --autonomous because the sign-off stage is `owner`. Yet `rhiza guide conduct` says the conductor's questions always go to the arbiter. Which fix should the plan build?

## Links

### In

- answers [conductor-questions-to-arbiter](../decision/conductor-questions-to-arbiter.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand

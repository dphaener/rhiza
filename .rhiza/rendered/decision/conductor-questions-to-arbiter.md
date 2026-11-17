# conductor-questions-to-arbiter

decision · `01M7M5A6P8PEBTC0P87WMJ8J8K`

## reasoning

rhiza guide conduct already promises that a conductor's question goes to the arbiter and its Decision stands, with no owner to ask; the code breaks that promise because recordHook and blockedNotice key arbiter routing on rec.Autonomous, which only run start --autonomous sets (session.go via run_spawn.go and run.go). The three stalling paths named in the question (daemon runs started without the flag, hand-queued run start, and daemonStart dropping --autonomous when rhiza.autonomy.signoff is owner) all share that one cause. Forcing --autonomous instead would also skip the owner's sign-off and gates, overriding an explicit owner setting, and would not cover a hand-started run. Routing by the conductor role fixes all three paths at one point, keeps the sign-off setting meaningful, and matches the guide, so no doc change is needed.

## statement

Decouple question routing from State.Autonomous in conductor runs: a question asked by the conductor session (the conduct step) or by a session it spawns with rhiza run spawn always goes to the arbiter (AskReply with the arbiter reason, an arbiter Decision, no advisor, no owner notice), whether or not the run started with --autonomous. State.Autonomous keeps governing only the owner's gates and sign-off, so daemonStart may still drop --autonomous for an owner sign-off and a hand-queued rhiza run start still stops at its sign-off, but neither stalls on a planning question. Do not fix it by forcing --autonomous on every conductor run.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-4060ff02-1](../question/question-4060ff02-1.md)
- derived-from [conductor-questions-plan](../source/conductor-questions-plan.md)

### In

- decided-by [conductor-question-arbitrated](../capability/rhiza/workflows/arbiter-sessions/conductor-question-arbitrated.md)
- decided-by [conductor-sessions-launch-autonomous](../capability/rhiza/workflows/arbiter-sessions/conductor-sessions-launch-autonomous.md)
- decided-by [spawned-child-question-arbitrated](../capability/rhiza/workflows/arbiter-sessions/spawned-child-question-arbitrated.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand

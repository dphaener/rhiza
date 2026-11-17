# sessionless-ruling-not-undelivered

decision · `01M7MPK6TPAPT4DB3KKJ1CT0GB`

## reasoning

internal/api/run.go sets Answered only when the engine delivers to an asking session (an answered entry), which never happens for a session-less ask (engine.go arbitrateAsked: 'No session asked, so nothing is delivered'), so every review ruling reads 'not yet delivered', the symptom the intent quotes.

## statement

Run status says 'not yet delivered' only for a question a session asked and is waiting on. A question no session asked (a review question journalled by rhiza run review-check) has nothing to deliver: its passed ruling prints as '— decision <handle>' alone. RunQuestion gains an omitempty session field (the asking session) so the markdown can tell them apart; the --json change is additive.

## Links

### Out

- derived-from [plan-run-rulings-research](../source/plan-run-rulings-research.md)

### In

- decided-by [run-status-sessionless-ruling](../capability/rhiza/workflows/built-in-workflows/plan-workflow/run-status-sessionless-ruling.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.

# main-ci-one-wait-step

decision · `01M7BZMBERXQM2M2PPEYW8KR5Y`

## reasoning

Read the question, the leaf builtin-conductor-default (its acceptance fixes exactly 10 steps, main's CI one step running rhiza ship check --main, no on.fail route), decision builtin-ship-change-open-points (10 steps within the cap; a red main stops the run as classic), classic's wait-main/main-green (internal/workflow/builtin/ship-change.toml), shipCheckMainDispatch (internal/api/ship_main.go) and the engine's wait loop (internal/run/engine.go). The finding is right that a command step checks once and reads pending CI as red. Adding classic's wait-main would make 11 steps and break the decided 10. But the engine's wait already tells a stop exit from a condition that does not hold (Conditioner, StopExit returns result{stopped}), and ship check --main without --settled returns Passed=false (exit 1) while unsettled, a stop (exit 75) once CI failed after its rerun, and pass (writing st.Main) when green. So one wait step on that command does what classic's wait-main plus main-green do together: it waits for pending CI, passes on green, and stops on red with the stop: line, keeping the count, the command and the no-hand-back rule the leaf already states.

## statement

Main's CI stays one step, so the builtin keeps exactly 10 steps, but that step is of kind wait, not command: until = "rhiza ship check --main" (without --settled), with a timeout (1h, as classic's wait-main) and polling interval/max_interval as classic, on.pass = close-out, and no on.fail, on.timeout or hand-back route. While CI dispatched on the merge commit is pending the check exits 1 and the wait keeps polling; when the merge commit's tree equals the tested head it passes at its first check; when main's CI is red after the rerun the check exits with the stop exit and the run stops with its stop: line saying CI failed on main after this change set's merge; a timeout ends the run too. TestBuiltinShipChange asserts that step's kind is wait and its until is rhiza ship check --main.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7ah6rayrzax6z6r63q0m2qk-r6-f5](../question/review-01m7ah6rayrzax6z6r63q0m2qk-r6-f5.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

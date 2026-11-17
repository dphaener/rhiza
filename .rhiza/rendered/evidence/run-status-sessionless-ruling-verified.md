# run-status-sessionless-ruling-verified

evidence · `01M7MQCFT93A2NNENHJAWXPVAD` · actor [claude-code-7021e610](../actor/claude-code-7021e610.md)

## result

PASS

## statement

At 0d2f9b2fe re-ran the Checks: block: go test -p 2 -trimpath -count=1 ./internal/api -run TestRunInfoAcceptsAndQuestions exit 0 (PASS); go test -p 2 -trimpath -count=1 ./cmd/rhiza -run TestScripts/describe-golden exit 0 (describe-golden PASS, golden carries the new session field and answered description in every RunQuestion schema); rhiza ship generate exit 0 (0 changed); rhiza render exit 0 (only this verifier's own actor rendered, reverted). Diff: RunQuestion gains Session string json session,omitempty, runInfo copies q.Session, Human prints ', not yet delivered' only when !Answered && Session != ''. Test covers the session-less passed Q5 printing '- review-q (Q5): Split? — decision use-z' with no 'not yet delivered' and no session key in its JSON, and the failure case Q4 (session S4, passed, undelivered) still printing ', not yet delivered' with session in its JSON. No engine delivery, journal entries or other status sections changed; no stubs or TODOs.

## Links

### Out

- verifies [run-status-sessionless-ruling](../capability/rhiza/workflows/built-in-workflows/plan-workflow/run-status-sessionless-ruling.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.

# run-status-sessionless-ruling

capability · `01M7MPK782XNNC37XVT5G84EX9` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/api -run 'TestRunInfoAcceptsAndQuestions' passes with a session-less passed question printing '— decision <handle>' and no 'not yet delivered', and its JSON having no session field; go test -trimpath ./cmd/rhiza -run 'TestScripts/describe-golden' passes with the golden regenerated. Failure case: a session's passed question not yet delivered still prints ', not yet delivered' and carries its session.

## files

internal/api/run.go, internal/api/run_test.go, cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Change how or when the engine delivers answers, the arbitrated or answered journal entries, or any other run status section.

## statement

rhiza run status prints a passed ruling on a question no session asked as '— decision <handle>' without ', not yet delivered', which stays only for a session's question whose answer the engine has not delivered; RunQuestion carries the asking session in an omitempty session field, and in an omitempty change field the change set its asked entry journals (where the question and its Decision live), which the plan reviser reads to apply a ruling.

## Links

### Out

- decided-by [sessionless-ruling-not-undelivered](../../../../../decision/sessionless-ruling-not-undelivered.md)

### In

- verifies [run-status-sessionless-ruling-verified](../../../../../evidence/run-status-sessionless-ruling-verified.md)
- verifies [run-status-sessionless-ruling-qa-1](../../../../../evidence/run-status-sessionless-ruling-qa-1.md)
- verifies [run-status-sessionless-ruling-qa-2](../../../../../evidence/run-status-sessionless-ruling-qa-2.md)
- verifies [run-status-sessionless-ruling-qa-3](../../../../../evidence/run-status-sessionless-ruling-qa-3.md)
- verifies [run-status-sessionless-ruling-qa-4](../../../../../evidence/run-status-sessionless-ruling-qa-4.md)
- verifies [run-status-sessionless-ruling-qa-5](../../../../../evidence/run-status-sessionless-ruling-qa-5.md)
- verifies [run-status-sessionless-ruling-qa-6](../../../../../evidence/run-status-sessionless-ruling-qa-6.md)
- verifies [run-status-sessionless-ruling-qa-7](../../../../../evidence/run-status-sessionless-ruling-qa-7.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.

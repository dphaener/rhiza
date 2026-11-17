# review-question-in-node-change-set-verified-2

evidence · `01M7MS3D427T29E5W5734AMSJR` · actor [claude-code-4c53132c](../actor/claude-code-4c53132c.md)

## result

PASS

## statement

Re-verified at HEAD 2c96751c5 (attempt 1 retry). Re-ran every Checks: command of f2230efd0, each exit 0: go test -trimpath ./internal/api -run TestReviewCheckQuestionInNodeChangeSet; ./internal/api -run 'Review|Rulings|Verdict|Findings'; ./internal/run -run 'TestCommandAskedArbitratedInChange|TestCommandAsked|TestArbitrate' (both InChange subtests pass: named->RHIZA_CHANGE=CS2, none->run's CS1); ./internal/run whole; ./cmd/rhiza TestScripts/(run-review-check|run-review-check-landed|run-review-check-handles|run-arbitrate-asked)$; TestScripts/(describe-golden|inbox-question-advisor|run-gate|plan-workflow-ci-infra-wait|plan-workflow|plan-workflow-pr)$; rhiza ship generate (0 changed) and rhiza render with a scratch build, no code diff. The earlier failure is resolved: go test -trimpath ./cmd/rhiza -run 'TestScripts/plan-workflow-autonomous-left$' passes at HEAD (2c96751c5 moved its fake reviser's decided-by link into the draft). Diff read: reviewQuestionTarget checks in-flight plan change sets creating the node before the open run change set, then record; asked entry carries Change only when not the run's; engine arbitrateIn starts the arbiter with the entry's change, falling back to st.Change. Test covers open (draft leaf -> draft with blocks, asked.Change=draft; own node -> run cs, Change empty) and landed (record write, Change empty). Only the leaf's files plus .rhiza log/rendered changed; no TODOs or stubs.

## Links

### Out

- verifies [review-question-in-node-change-set](../capability/review-question-in-node-change-set.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.

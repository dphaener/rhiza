# plan-gate-open-items-verified

evidence · `01M69CG13SNAXXPVV531THJE4J` · actor [claude-code-744428af](../actor/claude-code-744428af.md)

## result

PASS

## statement

At 916e15f4: go test ./cmd/rhiza -run 'TestScripts/(run-gate|describe-golden)$' -count=1 passes and make check exits 0. run-gate.txtar adds a draft change set 'Serve E one' serving want-e with planner-written question pq-e, plus a run whose last review round leaves the unanswered question review-<run>-r3-f1. A temporary copy of the script that dumped the output showed run gate --json returning "open":[{handle review-01m5pwj1-r3-f1},{handle pq-e}] and the human text '2 question(s) the plan left are open' listing both. After Decision ans-e1 answers both (one --record, one in the plan change set), the script checks run gate --json exits 0 with "open": []. describe-golden.txtar is regenerated with the open array (id required, handle optional) and open added to required. In the diff, gateOpenQuestions in internal/api/run_review.go collects reviewQuestions plus every question node-create op in the plan's change sets, skips retired nodes and those reviewAnswered covers, and lists each once.

## Links

### Out

- verifies [plan-gate-open-items](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-gate-open-items.md)

## History

- `01M625W0AJZ6AMCCWGFXQW19FN` A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.

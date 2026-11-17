# plan-workflow-applies-rulings-verified

evidence · `01M7MS09N0Y7C07YSQ7JF1816P` · actor [claude-code-c0697627](../actor/claude-code-c0697627.md)

## result

PASS

## statement

Verified commit 2c96751c5 (Leaf: plan-workflow-applies-rulings). Re-ran its Checks: go test -trimpath -p 2 ./internal/workflow -run 'TestPlanBuiltin|TestBuiltinPlanMilestone|TestBuiltinPlanReview' exit 0 (plan-milestone covered by DeepEqual of steps with plan); go test ./cmd/rhiza -run TestScripts/plan-workflow-autonomous-left$ exit 0; TestScripts/plan-workflow-autonomous-unapplied$ exit 0; the combined (plan-workflow|plan-workflow-pr|plan-workflow-autonomous|plan-review-visit-round|plan-milestone-backlog|...-left|...-unapplied)$ run exit 0 with -count=1 (a first run under load failed plan-workflow-autonomous.txtar:68 on journal order, a tokens event after ended in the interview step, which is outside this diff; it passed 3/3 alone and in the full re-run); rhiza ship generate exit 0, 0 files changed; rhiza render exit 0. Diff: plan.toml and plan-milestone.toml route check on.pass to rulings (run rhiza run review-check --rulings, max_visits 2, on.pass gate, on.fail apply-rulings) and add apply-rulings (the exact reviser prompt, until true, until_when done, max_visits 1, on.pass rulings); review/critic/check/revise steps and caps and the post-gate ship steps are unchanged. The five scripts assert 'check (visit N): pass, then rulings' then 'rulings (visit 1): pass, then gate'. plan-workflow-autonomous-left asserts rulings visit 1 fail, apply-rulings visit 1, rulings visit 2 pass, gate, and its applied.log asserts both export-csv (in the serving draft) and keep-flat (in the run's change set) were applied, each shown decided-by a settle- Decision, so the leaf ruling is not narrowed away. Failure case plan-workflow-autonomous-unapplied.txtar: the run fails, no gate or open-pr step, rulings visit 2 fails with no second apply-rulings, both change sets stay draft, and run status names both unapplied Decisions (export-csv, keep-flat). No TODOs or stubs. Extra file cmd/rhiza/main_test.go adds only the processScripts entry the new script needs.

## Links

### Out

- verifies [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.

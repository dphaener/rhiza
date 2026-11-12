# plan-prepush-fix-verified

evidence · `01M773WVXP7JAW95RMGGS2EDAQ` · actor [claude-code-65b5a3d6](../actor/claude-code-65b5a3d6.md)

## result

PASS

## statement

Verifier claude-code-65b5a3d6 at 2554db589: go test -trimpath ./internal/workflow/... ok; new TestPlanPrePushFix checks, for plan and plan-milestone, that open-pr and rebase on.fail = prepush-triage, prepush-triage runs `rhiza ship gate --failed` with only on.pass = prepush-fix (no on.fail), and prepush-fix is an agent step with max_visits 2, only on.pass = open-pr, and a prompt telling it to commit; TestPlanBuiltin's order and route table updated likewise. A build of this commit's `rhiza workflow describe plan` and `describe plan-milestone` both list prepush-triage (-> prepush-fix) and prepush-fix (-> open-pr) with open-pr/rebase on.fail -> prepush-triage. go test -trimpath ./internal/api -run TestShipGate ok: TestShipGateFailed exits 1 with no kept head (a rebase conflict records none), TestShipGateFailedOnPrePushFailure exits 1 after a gh or push error (nothing recorded), 0 on the recorded pre-push-failed head, 1 once HEAD moves. make self-check fails here only on stale rendered actor pages and index.md for this run's committed actor records (the known in-run render staleness, rendered at close-out); in a scratch worktree at this commit, `rhiza render` then `rhiza check` reports no fail. Diff read: plan.toml and plan-milestone.toml changes are identical; cli.md --failed text updated.

## Links

### Out

- verifies [plan-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/plan-prepush-fix.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

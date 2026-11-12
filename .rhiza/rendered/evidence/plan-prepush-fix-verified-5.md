# plan-prepush-fix-verified-5

evidence · `01M776NMCCKJSQ4ETVWR0HQ2YD` · actor [claude-code-1302d46f](../actor/claude-code-1302d46f.md)

## result

PASS

## statement

At 66cd47b7 (fixer cycle 3): go test -trimpath -count=1 ./internal/workflow/... passes (TestPlanPrePushFix checks plan and plan-milestone open-pr/rebase on.fail -> prepush-triage running ship gate --failed with only on.pass prepush-fix; prepush-fix agent, max_visits 2, on.pass open-pr, prompt says commit). rhiza workflow describe plan (and plan-milestone) lists prepush-triage and prepush-fix. go test -trimpath -count=1 ./internal/api -run TestShipGate passes (TestShipGateFailed: no kept head after a conflict exits 1; TestShipGateFailedOnPrePushFailure: a gh/push error records nothing and exits 1, the recorded pre-push head passes). make self-check now runs go test ./internal/workflow -run ^TestRepoWorkflowsPrePushFix$ first (seen in its output); in a scratch worktree at HEAD after rhiza render, make self-check exits 0 (in the run checkout it fails only on stale .rhiza/rendered actor files and index.md from this run's records, which render-pr fixes). Mutation check: a .rhiza/workflows/plan.toml copy of the built-in with open-pr's on.fail removed fails TestRepoWorkflowsPrePushFix (repo plan open-pr on.fail = ""; want prepush-triage), and an unchanged copy passes; main has no repo plan/plan-milestone copies today. The leaf's acceptance in the model was extended to name this guard, as decision prepush-fix-reaches-repo-workflows asks.

## Links

### Out

- verifies [plan-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/plan-prepush-fix.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

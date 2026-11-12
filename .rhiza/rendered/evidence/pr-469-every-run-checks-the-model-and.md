# pr-469-every-run-checks-the-model-and

evidence · `01M777H3FSJEHEDWBJX3K529DY` · actor [dphaener](../actor/dphaener.md)

## head

59e74343ae33d010eaddb8617e9d4fb7a48f7857

## pr

https://github.com/dphaener/rhiza/pull/469

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #469 (https://github.com/dphaener/rhiza/pull/469) passed CI at 59e7434

## status

verified

## Links

### Out

- verifies [prepush-check-setting](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-setting.md)
- verifies [prepush-check-before-push](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-before-push.md)
- verifies [prepush-check-in-run-status](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-in-run-status.md)
- verifies [ship-change-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-prepush-fix.md)
- verifies [plan-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/plan-prepush-fix.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

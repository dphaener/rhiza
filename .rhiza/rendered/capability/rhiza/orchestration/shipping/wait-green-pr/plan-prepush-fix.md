# plan-prepush-fix

capability · `01M76RT03BTXG6QDTFNE23XMA2` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/workflow/... passes, checking that plan.toml and plan-milestone.toml route open-pr's and rebase's on.fail to a ship gate --failed triage whose on.pass is prepush-fix and which has no on.fail, that prepush-fix has max_visits 2 and on.pass open-pr and its prompt tells the agent to commit its fix, and that rhiza workflow describe plan lists both steps; go test -trimpath ./internal/api -run TestShipGate passes, showing ship gate --failed failing on a ship rebase conflict and on a gh or push error; make self-check passes; make self-check also runs go test -trimpath ./internal/workflow -run TestRepoWorkflowsPrePushFix, which fails if a repo .rhiza/workflows/plan.toml or plan-milestone.toml lacks the prepush-triage and prepush-fix routing, and any such copy on main carries it

## statement

The built-in plan and plan-milestone workflows each send open-pr's and rebase's failure (both have no on.fail today) to a triage step running rhiza ship gate --failed, which passes only on a head the pre-push check failed on and then goes to a new agent step, prepush-fix (max_visits 2), that fixes the plan's model or rendered files from the failing lines with one rhiza apply or rhiza render, commits that fix to the run's branch, and goes back to open-pr; any other failure (a conflicting merge of main, a gh or push error) fails the triage with no on.fail and ends the run as today, and once prepush-fix's cap is spent the run fails as a failing step does today

## Links

### Out

- depends-on [prepush-check-before-push](prepush-check-before-push.md)
- depends-on [ship-change-prepush-fix](ship-change-prepush-fix.md)
- decided-by [pre-push-check-shape](../../../../../decision/pre-push-check-shape.md)

### In

- verifies [plan-prepush-fix-verified](../../../../../evidence/plan-prepush-fix-verified.md)
- verifies [plan-prepush-fix-qa-1](../../../../../evidence/plan-prepush-fix-qa-1.md)
- verifies [plan-prepush-fix-verified-2](../../../../../evidence/plan-prepush-fix-verified-2.md)
- verifies [plan-prepush-fix-qa-2](../../../../../evidence/plan-prepush-fix-qa-2.md)
- verifies [plan-prepush-fix-verified-5](../../../../../evidence/plan-prepush-fix-verified-5.md)
- verifies [pr-469-every-run-checks-the-model-and](../../../../../evidence/pr-469-every-run-checks-the-model-and.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

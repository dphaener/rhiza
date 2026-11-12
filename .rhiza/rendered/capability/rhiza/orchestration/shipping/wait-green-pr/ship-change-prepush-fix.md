# ship-change-prepush-fix

capability · `01M76RT010NG76M92D8AF73776` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/workflow/... passes, checking that ship-change.toml's open-pr, push-merge and land go on.fail to a ship gate --failed triage whose on.pass is gate-fix (gate-fix-merge for land) and which has no on.fail, that rebase still reaches gate-fix through rebase-triage and rebase-gate-failed, and that gate-fix's and gate-fix-merge's max_visits are unchanged; go test -trimpath ./internal/api -run TestShipGate passes, showing ship gate --failed passing on a head the pre-push check failed on and failing after an open-pr or land failure that was not the pre-push check (a gh or push error); make self-check runs go test -trimpath ./internal/workflow -run TestRepoWorkflowsPrePushFix, which fails if a repo .rhiza/workflows/ship-change.toml lacks that routing (the ship gate --failed triages passing to gate-fix and gate-fix-merge), and any such copy on main carries it

## statement

In the built-in ship-change workflow (and its wrappers such as bugfix), open-pr, push-merge and land, which have no on.fail today, each go on.fail to a triage step running rhiza ship gate --failed, as rebase does through rebase-triage and rebase-gate-failed; it passes only on a head the pre-push check (or local gate) failed on and then goes to the existing gate-fix fixer (gate-fix-merge after land) within its cap, and any other failure of those steps (gh, network, git) fails the triage with no on.fail and ends the run as today; ship gate --failed reports a head the pre-push check failed on, and the fixers' prompt reads the failing pre-push lines as well as a failing local gate's

## Links

### Out

- depends-on [prepush-check-before-push](prepush-check-before-push.md)
- decided-by [pre-push-check-shape](../../../../../decision/pre-push-check-shape.md)

### In

- depends-on [plan-prepush-fix](plan-prepush-fix.md)
- blocks [review-01m76rgz9k7zmdyec31e5kq7ta-r3-f1](../../../../../question/review-01m76rgz9k7zmdyec31e5kq7ta-r3-f1.md)
- verifies [ship-change-prepush-fix-verified](../../../../../evidence/ship-change-prepush-fix-verified.md)
- verifies [ship-change-prepush-fix-qa-1](../../../../../evidence/ship-change-prepush-fix-qa-1.md)
- verifies [ship-change-prepush-fix-qa-2](../../../../../evidence/ship-change-prepush-fix-qa-2.md)
- verifies [ship-change-prepush-fix-verified-5](../../../../../evidence/ship-change-prepush-fix-verified-5.md)
- verifies [pr-469-every-run-checks-the-model-and](../../../../../evidence/pr-469-every-run-checks-the-model-and.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

# ship-change-prepush-fix-qa-2

evidence · `01M775J1AA47VJ45FCR87WZNFJ` · actor [claude-code-9f99e82d](../actor/claude-code-9f99e82d.md)

## result

PASS

## statement

QA at f4f59c1c0 with a scratch build in a hermetic scratch repo with a bare origin: after ship rebase failed its pre-push check (exit 1, merge kept unpushed, ship state pre_push_failed set to the merged head), rhiza ship gate --failed exited 0, so ship-change's rebase-triage/rebase-gate-failed reaches gate-fix; once HEAD moved past it gate --failed exited 1. TestShipGate* and ./internal/workflow/... passed.

## Links

### Out

- verifies [ship-change-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-prepush-fix.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

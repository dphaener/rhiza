# ship-change-prepush-fix-qa-1

evidence · `01M774B03YPEGEXET5SJVGGJXR` · actor [claude-code-1b9be091](../actor/claude-code-1b9be091.md)

## result

PASS

## statement

QA at a191c9bf6 with a scratch build (go build -trimpath) in a hermetic scratch repo with a bare origin and a fake gh: rhiza workflow describe ship-change shows open-pr on.fail -> open-pr-triage (on.pass gate-fix), land/push-merge to their triage, gate-fix prompt naming the pre-push lines in run status --json. After the failed open-pr, .git/rhiza/ship state held pre_push_failed=651ba2c and rhiza ship gate --failed exited 0; after HEAD moved it exited 1. go test ./internal/workflow/... and TestShipGate* passed.

## Links

### Out

- verifies [ship-change-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-prepush-fix.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

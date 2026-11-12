# prepush-check-before-push

capability · `01M76RSZWZW5K5DG8AN1SP994H` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShipPrePushCheck|TestShipLocalGate|TestWatchdogPRNoCI' passes, showing in a run worktree with a committed non-canonical log entry, a committed stale render, or a failing prePushCheck command that ship open-pr, ship land, ship rebase and the watchdog's merge of origin/main each push nothing to origin and fail with the failing lines in their error and response (the watchdog leaving the branch at its pre-merge head); that a fix present only in the working tree (a render or log entry not committed, HEAD still stale) still fails the check and pushes nothing; that an untracked non-canonical .rhiza/log entry in the worktree with a clean committed head does not fail it; that the check's temporary checkout is removed afterwards; that with both passing they push as before (CI stable or unstable); and that localGate still runs only while CI is unstable

## statement

Before every git push a run makes (shipPush, used by ship open-pr and the land push; ship rebase's push of a merged head; and the watchdog's merge of origin/main into a run's branch, watchdogMergeMain in internal/api/watchdog_noci.go), all through one shared helper, whatever CI's state, Rhiza checks exactly the committed head being pushed: it makes a temporary clean checkout of that commit (a detached git worktree, removed afterwards) and runs rhiza check there (model, log canonicity, stale render) and then the project's rhiza.ship.prePushCheck command there, so uncommitted or untracked files in the run worktree (a stray .rhiza/log entry, a render or rhiza apply left uncommitted) neither pass nor fail it; if either fails it pushes nothing and fails the step with the failing lines and the head checked (the watchdog resets its merge as on a failed push and reports the failing lines), while rhiza.ship.localGate keeps running only while CI is unstable

## Links

### Out

- depends-on [prepush-check-setting](prepush-check-setting.md)
- decided-by [pre-push-check-shape](../../../../../decision/pre-push-check-shape.md)

### In

- depends-on [plan-prepush-fix](plan-prepush-fix.md)
- depends-on [ship-change-prepush-fix](ship-change-prepush-fix.md)
- depends-on [prepush-check-in-run-status](prepush-check-in-run-status.md)
- verifies [prepush-check-before-push-verified](../../../../../evidence/prepush-check-before-push-verified.md)
- verifies [prepush-check-before-push-qa-1](../../../../../evidence/prepush-check-before-push-qa-1.md)
- verifies [pr-469-every-run-checks-the-model-and](../../../../../evidence/pr-469-every-run-checks-the-model-and.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

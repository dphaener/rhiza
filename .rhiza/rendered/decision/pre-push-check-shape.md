# pre-push-check-shape

decision · `01M76RR6DE5Y21VWVD28RHDN40`

## reasoning

Intent: milestone-25 (the daemon builds and ships the backlog on its own) is still open and is what CI failures on model and generated-file errors stall; the predecessor ship-local-gate served it, and milestone-32 is about cleanup and language-neutral defaults, which this change only has to respect. Enforcement: the owner wants every push covered (fixer, ci-fix, merge-of-main pushes, plan runs). Workflow steps alone would need a gate before each push step in five builtin workflows and any custom one, and the three failing PRs came from plan runs that have no gate at all; pushes in code all go through shipPush (internal/api/ship_close.go, called by ship_openpr.go and the land) and ship_rebase.go's push, so a check there cannot be skipped by a workflow that forgets a step, and the push step's failure is already routed by on.fail. Name: existing per-project command settings are rhiza.ship.localGate and rhiza.ship.generate (ship_gate.go, ship_generate.go); rhiza.ship.prePushCheck follows that prefix, says when it runs, and names no language or tool, as milestone-32 asks. Plan fixer: plan.toml's revise step is tied to review findings and loops back to review, and its cap of 2 is shared with them; pre-push failures in plans are model errors (duplicate handle, non-canonical log entry, stale render) a dedicated small fixer step can repair with rhiza apply or rhiza render, matching the intent's point (4) that a failure goes back to a fixer within the step's cap instead of pushing. Head checked: the pushes send the committed HEAD while rhiza check reads the working tree, and run worktrees carry untracked stray .rhiza/log entries and fixes not yet committed, so the check runs on a clean checkout of the head being pushed and a fixer must commit its fix.

## statement

Change set 01M76RH0 serves @milestone-25, as its predecessor @ship-local-gate did, and only coordinates with @milestone-32 by keeping the new setting language neutral. The pre-push check is enforced in the product, in one shared helper run before every git push a run makes (shipPush, used by ship open-pr and the land push; ship rebase's push; and the watchdog's merge-of-main push in watchdogMergeMain), so plan, plan-milestone, ship-change and its wrappers such as bugfix all get it, including fixer, ci-fix and merge-of-main pushes; it runs rhiza check, then the project's command, in a temporary clean checkout of exactly the committed head being pushed (not the run worktree, whose uncommitted or untracked files would make the result differ from what CI sees), and a failure refuses the push and fails that step with the failing lines, recorded in the run journal and rhiza run status. The setting is rhiza.ship.prePushCheck, a project setting with an empty default (rhiza check alone), set to make self-check for this repo, beside rhiza.ship.localGate and rhiza.ship.generate. rhiza.ship.localGate keeps its meaning (only while CI is unstable). In ship-change a push step failed by the pre-push check goes on.fail through a rhiza ship gate --failed triage, as rebase does today, to the existing gate-fix fixer, within its cap; other failures of those steps end the run as today. The plan and plan-milestone workflows each get the same triage and one new agent fixer step for it (for example prepush-fix, max_visits 2) that fixes the plan's model or rendered files from the failing lines with one rhiza apply or rhiza render, commits that fix, and goes back to open-pr; conflicts and gh or push errors still end the run, and once its cap is spent the run fails as a failing step does today.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-aa618fe3-1](../question/question-aa618fe3-1.md)
- derived-from [pre-push-check-plan](../source/pre-push-check-plan.md)

### In

- decided-by [prepush-check-setting](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-setting.md)
- decided-by [prepush-check-before-push](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-before-push.md)
- decided-by [prepush-check-in-run-status](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-in-run-status.md)
- decided-by [ship-change-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-prepush-fix.md)
- decided-by [plan-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/plan-prepush-fix.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

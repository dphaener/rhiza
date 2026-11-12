# ship-change-prepush-fix-verified-5

evidence · `01M7767H08RZHHYQV8KWCC2BH8` · actor [claude-code-6d692889](../actor/claude-code-6d692889.md)

## result

PASS

## statement

Verifier, attempt 5, at 63e02a2a4. go test -trimpath ./internal/workflow/... passes (ship-change.toml routes open-pr, push-merge and land on.fail to ship gate --failed triages that go on.pass to gate-fix / gate-fix-merge with no on.fail; rebase still reaches gate-fix via rebase-triage and rebase-gate-failed; gate-fix caps unchanged). go test -trimpath ./internal/api -run 'TestShipGate|TestShipCheckFailsOn|TestShipPrePushCheck|TestShipCheck' passes, including TestShipGateFailedOnPrePushFailure (gate --failed passes on the pre-push-failed head and fails after a gh/push error) and the fix commit's new tests: TestShipPrePushCheckSettleConflict (ship check --settled on a CONFLICTING PR with a failing real pre-push check ends the wait settled, exit 0, no push, PrePushFailed recorded, then ship gate --failed passes on the kept head) and TestShipCheckFailsOnPrePushFailedHead (plain ship check exits 1 on the kept head, passes after a fix). Read the diff: shipSettleConflict now handles rb.PrePushCheck like a failed gate, and runShipCheck has a PrePushFailed twin of the GateFailed branch, closing the review finding. Run from a detached scratch worktree outside ~/code, since inside it the api test guard trips on ~/code/.rhiza-root/usage.live.json, which the live daemon rewrites every few seconds regardless of the tests.

## Links

### Out

- verifies [ship-change-prepush-fix](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-prepush-fix.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).

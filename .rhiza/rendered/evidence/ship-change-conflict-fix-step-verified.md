# ship-change-conflict-fix-step-verified

evidence · `01M6XDCSWZF1RPETGJSWS9AP0K` · actor [claude-code-979e10f9](../actor/claude-code-979e10f9.md)

## result

PASS

## statement

Verifier at 88d5cefd read the diff and re-ran the acceptance. go test -trimpath ./internal/workflow/... passes: TestShipChangeBuiltin asserts rebase on.fail rebase-triage; rebase-triage on.pass conflict-fix with no on.fail; pr-green and rebased-green on.fail conflict-triage; conflict-triage on.pass conflict-fix and on.fail ci-triage; conflict-fix on.pass render-pr; rebase-merge, landed-green and rebased-merge-green on.fail conflict-triage-landed; conflict-triage-landed on.pass conflict-fix-landed with no on.fail; conflict-fix-landed on.pass rebase-merge; and max_visits 3 on both agent fixers. internal/guide passes. In internal/api, TestShipCheckConflicts (no recorded conflict fails; a recorded conflict with 0 or 2 fixes passes; 3 merge commits with the run's Review: trailer fail with '3 conflict fixes'; a resolved merge fails) passes along with every TestShip*/TestSession* test (go test -count=1 -run 'TestShip|TestSession'). The full package run reported no failed test; it failed only on wall-clock and CPU budgets on a VM at load ~11 and on the root-state guard for ~/code/.rhiza-root/usage.live.json, which the live root daemon writes, neither of which this diff touches. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship|ship-check)' passes (266s). That includes the new run-ship-conflict-fix.txtar: a bare origin, fake gh and fake fixer cover two conflicts before land (wait-pr and rebase) each fixed once and pushed plainly, with pr-green passing with no second triage or fixer; a post-land conflict routed through conflict-triage-landed and conflict-fix-landed back to rebase-merge with one land; a refused push ending at rebase-triage with no fixer; and a fourth conflict failing the run at conflict-triage-landed with no merge. ship code has no force push.

## Links

### Out

- verifies [ship-change-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-conflict-fix-step.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.

# conflict-fix-design-clash-verified

evidence · `01M6XF7306VP9N02R52ED209NR` · actor [claude-code-d9421646](../actor/claude-code-d9421646.md)

## result

PASS

## statement

At a445502a ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/run-ship-conflict-fix' (ok, 242s), go test -trimpath ./internal/run/... ./internal/workflow/... (ok) and TestScripts/(run-ship-change-rebase|describe-golden) (ok). Read the diff: ship-change.toml's conflict-fix and conflict-fix-landed prompts now tell the fixer to leave the merge in progress and ask one blocking AskUserQuestion naming both change sets, clashing files and behaviours, and their until is now '! git rev-parse --verify --quiet MERGE_HEAD' (was true); no engine change. Script case 5 (autonomous): stderr shows conflict-fix blocked on question, the session asked a question the arbiter answers, the fake arbiter answered it with clash-ruling; the journal has one asked entry for conflict-fix with 'Design clash' and one answered; one send-key to the same block; clash.log and fixes.log each 1; conflict-fix visit 1 passes to render-pr, no visit 2; PR merged into origin main. Case 6 (non-autonomous): run status --json blocked with kind question and 'Design clash', current step conflict-fix, no fix commit and nothing pushed while waiting; after the owner's answer the same session finishes and conflict-fix passes. Case 7: fixer says done without committing, conflict-fix fails with until did not hold, run failed, no push, no pr merge.

## Links

### Out

- verifies [conflict-fix-design-clash](../capability/rhiza/orchestration/shipping/wait-green-pr/conflict-fix-design-clash.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.

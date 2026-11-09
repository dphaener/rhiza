# ship-merge-main-plan

source · `01M6WMQ02SNH8RVKY8WXYWN8DY`

## quotes

```json
[]
```

## raw_sha256

db1bfba98d301569da7baf812bae22fc52e724d6fb605cf75ffb825050faf201

## session

97efbc82-de6a-48f6-a9fa-01a1b9467030

## summary

Autonomous plan run, no owner. The planner read the change set's intent (run 01M6W0Z45JXE stopped at wait-pr on a rebase conflict in seven files), internal/workflow/builtin/ship-change.toml, internal/api/ship_rebase.go (rebase, force-push with lease, rebase-conflict question), shipSettleConflict in internal/api/ship.go, the run engine's StopExit and the wait step's pass/timeout outcomes, and the model (milestone-25, wait-settles-conflicting-pr, ship-change-ci-fix-cycle). It picked milestone-25 as the served intent, as the CI-fix change set did. It raised one question with four topics: merge or rebase, who pushes the fixer's commit, what a design clash does, and the cap on fix cycles. The arbiter answered with Decision ship-merge-main-conflict-fixer: merge only, plain push, a conflicting merge is left in the worktree and the step fails to a conflict-fix agent step. The fixer resolves, regenerates goldens and rendered files, runs the touched packages' tests and commits with a Review: trailer, never pushing; the run's render and push step pushes and records. On a design clash the fixer aborts and raises a question that the run waits on. At most 3 cycles. The planner then worked out the wiring: wait steps have no fail outcome, so a conflicting wait settles and the green step routes there through a conflict triage. After land, a separate fixer step returns to rebase-merge, not render-pr, so nothing lands twice.

## title

Plan: ship-change merges main and fixes conflicts

## Links

### In

- derived-from [ship-merge-main-conflict-fixer](../decision/ship-merge-main-conflict-fixer.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.

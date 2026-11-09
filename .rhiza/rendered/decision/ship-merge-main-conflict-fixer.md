# ship-merge-main-conflict-fixer

decision · `01M6WMNNQ3C861GJX73FRKRXRY`

## reasoning

Read the question, the change set intent and internal/workflow/builtin/ship-change.toml plus internal/api/ship.go. The intent says outright: merge origin/main, never rebase or force-push, give conflicts to a fixer, ask only on a real design clash, so the merge-only answer follows it directly. On who pushes: the intent says the fixer 'pushes' and 'the run records the push'. I chose the run pushing because that is the pattern already in place. ci-fix only commits and the run goes back through render-pr and the push step. ship state tracks Pushed and wait-pr compares the PR head with it, so a push from outside that step (a hand or agent push) can stall wait-pr. Rendering again just before every push also keeps `rhiza check` passing on the pushed branch. Letting the run push still meets the intent, since the run records it. For a design clash, the run already routes failures to questions and fresh rounds once a Decision answers. A half-resolved merge must not be pushed, so the fixer aborts and asks. That keeps the owner or arbiter in the loop only for real clashes, as the intent asks. Three cycles matches max_visits = 3 on the rebase and wait steps and ci-triage. Main moves under long runs, so one cycle is too few, but a cap stops a run that keeps chasing main from looping forever.

## statement

Ship-change brings origin/main into a change set's branch only by merging, never by rebasing or force-pushing: `rhiza ship rebase` (steps rebase, rebase-merge) and `ship check --settled` on a CONFLICTING PR run `git merge origin/main` on the run worktree. A clean merge is committed and pushed with a plain push (no --force, no lease) and recorded as pushed, as today. A conflicting merge is left in the worktree (not aborted), the command fails with a distinct result, and on.fail routes the run to a conflict-fix agent step. The fixer only resolves and commits: it resolves the conflicts with both sides' commits and change sets as context, regenerates golden files and rendered output, runs the touched packages' tests, and makes the merge commit with Change:/Actor: and a `Review: {run}` trailer, then `rhiza session done`. It does not push. The run pushes it: the fixer step goes on to the same render and push step the run already uses (as ci-fix does), which pushes plainly and records the pushed head in the ship state, then waits for CI again. On a real design clash (both sides change the same behaviour incompatibly) the fixer commits nothing and asks one blocking question in its session naming both change sets, the clashing files and the incompatible behaviours; the session waits while the arbiter (in an autonomous run) or the owner answers it, and the same session then resolves the merge following the answer. A half-resolved merge is never pushed. One ship-change run may make at most 3 conflict-fix cycles (max_visits = 3 on the conflict-fix step, matching the rebase and wait steps' 3); a fourth conflict fails the run.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-97efbc82-1](../question/question-97efbc82-1.md)
- derived-from [ship-merge-main-plan](../source/ship-merge-main-plan.md)

### In

- decided-by [ship-change-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-conflict-fix-step.md)
- decided-by [conflict-fix-design-clash](../capability/rhiza/orchestration/shipping/wait-green-pr/conflict-fix-design-clash.md)
- decided-by rebase-keeps-branch-merges (retired)
- decided-by wait-settles-conflicting-pr (retired)
- decided-by [ship-rebase-merges-only](../capability/rhiza/workflows/workflow-runs/rebase-when-behind/ship-rebase-merges-only.md)
- decided-by [wait-merges-conflicting-pr](../capability/rhiza/orchestration/shipping/wait-green-pr/wait-merges-conflicting-pr.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.

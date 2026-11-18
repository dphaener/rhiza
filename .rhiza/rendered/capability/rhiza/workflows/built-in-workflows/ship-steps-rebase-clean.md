# ship-steps-rebase-clean

capability · `01M7MNGT7KV1QW8BA3VVXX1DVW` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/workflow -run 'TestBuiltinShipChange|TestPlainRedLoopsShipChangeHasNone|TestMergeWithoutGate' passes, TestBuiltinShipChange asserting exactly the 12 steps conduct, gate, signoff, push, wait-ci, pr-green, merge, rebase, main-ci, close-out, ship, deploy-check in that order; wait-ci's until exactly rhiza ship check --settled, its timeout handing back to conduct; pr-green a command step running rhiza ship check, pass to merge, fail to handback:conduct; merge running exactly rhiza ship merge, with only pass to main-ci and fail to rebase; rebase running exactly rhiza ship rebase with stall = "70m", at least the gate step's stall, pass to wait-ci, fail to handback:conduct; max_visits 6 on wait-ci, pr-green, merge and rebase; handback_cap 3; main-ci, close-out, ship and deploy-check with no on.fail route and no hand-back; the conduct prompt naming the conflict fixer and not containing --no-push; every other non-pass route of conduct, gate, signoff, push, wait-ci and pr-green handing back to conduct; (failure case) TestBuiltinShipChange fails when merge's fail route is handback:conduct, when wait-ci's until carries --no-rebase, when rebase's fail goes anywhere but handback:conduct, or when rebase's stall is missing or below the gate step's stall. go test -trimpath ./cmd/rhiza -run TestScripts/workflow-new passes.

## files

internal/workflow/builtin/ship-change.toml, internal/workflow/workflow_test.go

## must_not

Change internal/api (ship rebase, ship check, ship merge, the --no-push and --no-rebase flags stay as they are), ship-change-classic.toml, plan.toml or the repo's .rhiza/workflows; give main-ci, close-out, ship or deploy-check an on.fail route or a hand-back; change handback_cap; touch the guide topics or docs (another leaf's)

## statement

The builtin ship-change rebases a clean merge of main in its own steps: wait-ci waits on rhiza ship check --settled (no --no-rebase), pr-green runs rhiza ship check and hands back on fail, merge runs rhiza ship merge alone and goes to a new rebase step on fail, and rebase runs rhiza ship rebase with stall = "70m" (at least the gate step's stall, since ship rebase runs the local gate on the merged head), passing back to wait-ci and handing back to conduct on fail; the conduct prompt no longer tells the conductor to run rhiza ship rebase --no-push.

## Links

### Out

- decided-by [ship-steps-rebase-clean-merge](../../../../decision/ship-steps-rebase-clean-merge.md)
- decided-by [ship-steps-rebase-stall-70m](../../../../decision/ship-steps-rebase-stall-70m.md)

### In

- depends-on [conduct-handback-ship-rebase](../../guide/agent-guide/conduct-handback-ship-rebase.md)
- depends-on [ship-one-ci-rebase-in-steps](../../orchestration/shipping/ship-one-ci-rebase-in-steps.md)
- blocks [review-01m7nwyk3pyjjy43bsg63bzyd1-r3-f1](../../../../question/review-01m7nwyk3pyjjy43bsg63bzyd1-r3-f1.md)
- verifies [ship-steps-rebase-clean-verified](../../../../evidence/ship-steps-rebase-clean-verified.md)
- verifies [ship-steps-rebase-clean-qa-1](../../../../evidence/ship-steps-rebase-clean-qa-1.md)
- verifies [ship-steps-rebase-clean-qa-2](../../../../evidence/ship-steps-rebase-clean-qa-2.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict

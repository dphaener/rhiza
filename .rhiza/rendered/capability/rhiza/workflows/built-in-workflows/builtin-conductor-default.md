# builtin-conductor-default

capability · `01M79KPE0KMZHN6W97EDAZFNT3` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/workflow -run TestBuiltinShipChange shows the steps and routes ship-steps-rebase-clean's acceptance gives, the conduct step declaring exactly the checks plan, code and qa, the gate step of kind gate with its four commands, the hand-back cap at 3 and no leaves or checks step; rhiza workflow describe ship-change lists those steps. (failure case) go test -trimpath ./internal/api -run TestBuildCheck shows build-check failing on a leaf with no Evidence, and TestBuiltinShipChange fails when merge's fail route is handback:conduct or when main-ci, close-out, ship or deploy-check has an on.fail route.

## files

internal/workflow/builtin/ship-change.toml, internal/workflow/workflow_test.go

## must_not

Change internal/api (rhiza run build-check, rhiza ship rebase, ship check and ship merge stay as they are), ship-change-classic.toml, plan.toml or the repo's .rhiza/workflows; give main-ci, close-out, ship or deploy-check an on.fail route or a hand-back; raise the hand-back cap above 3

## statement

The builtin ship-change is the conductor default: conduct (the conduct skill, its step declaring the review checks plan, code and qa), gate (a step of the gate step kind running rhiza run build-check, rhiza run review-check, rhiza check and rhiza ship gate; the pre-push check runs inside the push step's rhiza ship open-pr and inside rhiza ship rebase), signoff, push and open the PR, wait for CI at that head (rhiza ship check --settled, which merges main on a conflicting PR), pr-green (rhiza ship check), merge the tested head (rhiza ship merge), rebase (rhiza ship rebase, back to the CI wait), main's CI (one step), close-out, ship (rhiza ship command) and deploy-check (rhiza ship deploy), as ship-steps-rebase-clean gives them; every red result before the merge (conduct, gate, signoff, push, wait-ci and pr-green) hands back to conduct, at most three times; every merge fail (main moved, or a gh, fetch or git error) goes to the rebase step, which hands back only on a code conflict or a local gate or pre-push check failing on the merged head. The steps after the merge (main's CI, close-out, ship and deploy-check) have no on.fail route and no hand-back: nothing on the branch can fix a merged change set. A draft or a merge GitHub refused stops the run (rhiza ship merge's stop exit).

## Links

### Out

- depends-on [classic-tests-pinned](classic-tests-pinned.md)
- depends-on [review-three-checks](review-three-checks.md)
- depends-on [critic-on-major](critic-on-major.md)
- depends-on [run-handback](../workflow-runs/run-handback.md)
- depends-on [conduct-skill](../../guide/agent-guide/conduct-skill.md)
- depends-on [agent-step-checks](built-in-review/agent-step-checks.md)
- depends-on [review-check-declared-checks](built-in-review/review-check-declared-checks.md)
- depends-on [gate-step](../workflow-definitions/gate-step.md)
- decided-by [conductor-promote-builtin-now](../../../../decision/conductor-promote-builtin-now.md)
- decided-by [builtin-ship-change-open-points](../../../../decision/builtin-ship-change-open-points.md)
- decided-by [conductor-checks-gate-pick](../../../../decision/conductor-checks-gate-pick.md)
- decided-by [ship-steps-rebase-clean-merge](../../../../decision/ship-steps-rebase-clean-merge.md)

### In

- depends-on [ship-one-ci-round](../../orchestration/shipping/ship-one-ci-round.md)
- depends-on [rhiza-repo-conductor-default](../workflow-overrides/rhiza-repo-conductor-default.md)
- depends-on [docs-conductor-default](docs-conductor-default.md)
- depends-on [default-tests-follow-conductor](default-tests-follow-conductor.md)
- blocks [review-01m7ah6rayrzax6z6r63q0m2qk-r6-f4](../../../../question/review-01m7ah6rayrzax6z6r63q0m2qk-r6-f4.md)
- blocks [review-01m7ah6rayrzax6z6r63q0m2qk-r6-f5](../../../../question/review-01m7ah6rayrzax6z6r63q0m2qk-r6-f5.md)
- verifies [builtin-conductor-default-verified](../../../../evidence/builtin-conductor-default-verified.md)
- verifies [builtin-conductor-default-qa-2](../../../../evidence/builtin-conductor-default-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict

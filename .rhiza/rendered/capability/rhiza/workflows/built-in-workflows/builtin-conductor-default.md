# builtin-conductor-default

capability · `01M79KPE0KMZHN6W97EDAZFNT3` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/workflow -run TestBuiltinShipChange shows the builtin has exactly the 10 steps conduct, gate, signoff, push, CI wait, merge, main's CI (running rhiza ship check --main), close-out, ship (rhiza ship command) and deploy-check (rhiza ship deploy), in that order, the conduct step declaring exactly the checks plan, code and qa, the gate step being of kind gate with the commands rhiza run build-check, rhiza run review-check, rhiza check and rhiza ship gate, every failure outcome of the steps up to and including the merge routing to handback:conduct (the merge step's fail included, with no merge route to the push step), main's CI, close-out, ship and deploy-check having no on.fail route and no hand-back, the hand-back cap at 3, the conduct step's prompt telling the conductor to run rhiza ship rebase when the hand-back is rhiza ship merge's main-moved output, and no leaves or checks step remaining; go test -trimpath ./internal/api -run TestBuildCheck shows build-check failing on a leaf with no Evidence; rhiza workflow describe ship-change lists those steps.

## statement

The builtin ship-change is the conductor default: conduct (the conduct skill, its step declaring the review checks plan, code and qa), gate (a step of the gate step kind, gate-step, running rhiza run build-check, which fails unless every leaf has live Evidence by a verifier who wrote none of its commits; rhiza run review-check, failing on any open blocker or major finding of a declared check; rhiza check; and rhiza ship gate, the local gate while CI is unstable; the pre-push check runs where it does today, inside the push step's rhiza ship open-pr), signoff, push and open the PR, wait for CI at that head, merge the tested head, main's CI (one step), close-out, ship (rhiza ship command) and deploy-check (rhiza ship deploy), exactly 10 steps; every red result up to and including the merge hands back to conduct, at most three times. The steps after the merge (main's CI, close-out, ship and deploy-check) have no on.fail route and no hand-back, as in classic: nothing on the branch can fix a merged change set, so a red main CI stops the run with output saying main is red after this change set's merge, and a failed close-out, ship or deploy-check stops it too. The merge step has only pass and fail: a draft or a merge GitHub refused stops the run (rhiza ship merge's stop exit), and every other failure, main having moved past the tested head or a gh, fetch or git error, hands back to conduct with rhiza ship merge's output. On a main-moved hand-back the conductor merges origin/main into the branch with rhiza ship rebase, resolves and commits any code conflict it leaves in progress, and the run goes through gate, signoff, push and a second CI wait on the new head; the push step never re-pushes an unchanged head to retry a main-moved merge.

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

# builtin-ship-change-open-points

decision · `01M7AHS4650X7Q6THEFQR2D2CB`

## reasoning

Read the question, the change set's slice (leaves builtin-conductor-default, builtin-ship-change-classic, rhiza-repo-conductor-default, ship-one-ci-round; decision conductor-promote-builtin-now), the shadow .rhiza/workflows/ship-change.toml, ship.sh, bugfix.toml, dispatch.toml, rhiza guide workflows and internal/api/ship_build.go. (1) The eight listed steps plus ship and deploy make 10, so the cap leaves room; this repo's ship.sh runs make update after every merge and deployUrl is a user setting, so dropping them would silently break existing projects (a breaking change the milestone never asked for). (2) After the merge nothing on the branch can fix main, so a hand-back to conduct would loop on a merged change set; the classic behaviour (stop with a stop: line) is the safe, known one, and the hand-back rule in the intent is read as covering reds a conductor can act on. (3) Only the outer workflow's [env] applies to included steps (bugfix.toml's comment), and dispatch [ready] feature is the documented hook for a repo's default ready workflow, so a wrapper keeps ship-change resolving to the builtin (the leaf's acceptance) with no new setting; dropping the head-build steps fits the 10-step cap, and ship_build.go already has the build-on-demand pattern. (4) builtin-ship-change-classic promises the old workflow stays unchanged and pinnable, so its behaviour must stay tested; retargeting keeps that coverage at the lowest cost, while the leaves' acceptance names the tests that must follow the new default.

## statement

(1) The new builtin ship-change keeps ship (rhiza ship command) and deploy-check (rhiza ship deploy) after close-out, unchanged: conduct, gate, signoff, push/PR, CI, pr-green, merge, rebase, main CI, close-out, ship, deploy is exactly 12 steps (10 when first decided; ship-steps-rebase-clean-merge adds pr-green and rebase). (2) A red main CI after the merge stops the run, as classic does: main-green has no on.fail route and no hand-back, and its output says main is red after the merge of this change set; the hand-back rule covers every red before the merge. (3) Rhiza's shadow .rhiza/workflows/ship-change.toml is removed; its [env] (GOTMPDIR, SHIP_HEAD_RHIZA, SHIP_HEAD_BUILD) moves to a thin repo workflow .rhiza/workflows/ship-rhiza.toml that includes the builtin ship-change (like bugfix.toml does), and .rhiza/dispatch.toml sets [ready] feature = "ship-rhiza"; bugfix.toml keeps its own [env]. The builtin has no head-build steps: the ship commands that land or render run SHIP_HEAD_BUILD themselves before doing so when it is set, as ship rebase already does. (4) The classic tests stay and are pointed at ship-change-classic (by name, or with the workflow pinned in the script), since classic stays available and pinnable; only the tests of the default itself (bugfix-workflow-autonomous, daemon-root-ship-workflow, the daemon pick tests that assert the default's steps) move to the new builtin, and the new builtin gets its own tests (TestBuiltinShipChange, ship-change-one-ci). No classic test is deleted.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-9a4a7e78-1](../question/question-9a4a7e78-1.md)
- derived-from [builtin-ship-change-plan-interview](../source/builtin-ship-change-plan-interview.md)

### In

- decided-by [classic-tests-pinned](../capability/rhiza/workflows/built-in-workflows/classic-tests-pinned.md)
- decided-by [builtin-conductor-default](../capability/rhiza/workflows/built-in-workflows/builtin-conductor-default.md)
- decided-by [rhiza-repo-conductor-default](../capability/rhiza/workflows/workflow-overrides/rhiza-repo-conductor-default.md)
- decided-by [ship-head-build-on-demand](../capability/rhiza/orchestration/shipping/ship-head-build-on-demand.md)
- decided-by [default-tests-follow-conductor](../capability/rhiza/workflows/built-in-workflows/default-tests-follow-conductor.md)
- decided-by [conductor-land-render-head-build](../capability/rhiza/orchestration/shipping/conductor-land-render-head-build.md)
- decided-by [handback-target-survives-include](../capability/rhiza/workflows/built-in-workflows/handback-target-survives-include.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict

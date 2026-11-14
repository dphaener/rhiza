# ship-head-build-on-demand

capability · `01M7AKMZ9QVN9EJJDVT64WXM13` · parent [shipping](../shipping.md)

## acceptance

go test -trimpath ./internal/api -run TestShipHeadBuildOnDemand shows, with SHIP_HEAD_RHIZA naming a missing binary and SHIP_HEAD_BUILD set, ship's render and ship close-out running SHIP_HEAD_BUILD once in the worktree and then rendering with the built binary; with the binary built at an earlier commit and a new commit on the branch since, ship close-out running SHIP_HEAD_BUILD again before rendering, and with it built at the current HEAD, not running it; with SHIP_HEAD_BUILD empty, the missing binary failing with its 'not a built rhiza' error; and a failing SHIP_HEAD_BUILD failing with nothing rendered or pushed. go test -trimpath ./cmd/rhiza -run TestScripts/ship-land-branch-build passes, with a new case where SHIP_HEAD_RHIZA is missing and SHIP_HEAD_BUILD is set and ship land builds it, logs to head.log and lands. Touches internal/api/ship_build.go, internal/api/ship_build_test.go and cmd/rhiza/testdata/script/ship-land-branch-build.txtar.

## statement

When the run's workflow sets SHIP_HEAD_BUILD, the ship commands that land or render with the SHIP_HEAD_RHIZA binary (ship land, ship close-out and every ship render) run SHIP_HEAD_BUILD themselves (sh -c, in the worktree) and then use the built binary whenever that binary does not exist yet or was not built at the worktree's current HEAD commit (ship records the commit it built at beside the binary), as ship rebase already does, so a workflow with no head-build steps (Rhiza's ship-rhiza over the conductor default) renders and closes out with a build of the branch's current head (the conductor's own rhiza land and rhiza render reach that build by conductor-land-render-head-build) even after a hand-back committed on top of an earlier build. With SHIP_HEAD_BUILD empty, a missing binary still fails with nothing landed, rendered or pushed, and a failed build fails the same way, carrying the build's output.

## Links

### Out

- decided-by [builtin-ship-change-open-points](../../../../decision/builtin-ship-change-open-points.md)

### In

- depends-on [rhiza-repo-conductor-default](../../workflows/workflow-overrides/rhiza-repo-conductor-default.md)
- depends-on [conductor-land-render-head-build](conductor-land-render-head-build.md)
- verifies [ship-head-build-on-demand-verified](../../../../evidence/ship-head-build-on-demand-verified.md)
- verifies [ship-head-build-on-demand-qa-2](../../../../evidence/ship-head-build-on-demand-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

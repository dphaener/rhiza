# ship-head-build-on-demand-qa-2

evidence · `01M7DTKRBVSZC2QCSZ34FYM8Q3` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: In a scratch repo, render with SHIP_HEAD_RHIZA missing and SHIP_HEAD_BUILD set built the head once and wrote <bin>.built-at at HEAD. A second render did not rebuild; after a new commit it rebuilt (builds 1, 1, 2). A failing SHIP_HEAD_BUILD failed with its output and 'nothing was rendered or pushed', with nothing rendered. TestShipHeadBuildOnDemand and ship-land-branch-build pass.

## Links

### Out

- verifies [ship-head-build-on-demand](../capability/rhiza/orchestration/shipping/ship-head-build-on-demand.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

# conductor-land-render-head-build-qa-2

evidence · `01M7DTKQ4PG15QSSN2QSJJZRBF` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: In a scratch repo: SHIP_HEAD_RHIZA naming a fake head made render exec it with 'render --json' in the repo dir; with a real head build, render and land went through it and landed the change set (accepted to landed); with the head as the running binary, render ran in-process; a missing head with no SHIP_HEAD_BUILD made render and land fail with 'not a built rhiza', with nothing rendered or landed. TestLandRenderHeadBuild passes.

## Links

### Out

- verifies [conductor-land-render-head-build](../capability/rhiza/orchestration/shipping/conductor-land-render-head-build.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

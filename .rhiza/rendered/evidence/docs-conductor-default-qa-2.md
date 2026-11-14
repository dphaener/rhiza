# docs-conductor-default-qa-2

evidence · `01M7DTKQC0MQYAGER6E3BM46ZC` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: docs/runs-and-daemon.md no longer says 'build and then review' for ship-change. It, docs/custom-workflows.md and guide workflows name the conductor default's steps and ship-change-classic pinned via [ready]. TestWorkflowsText and internal/doclinks pass.

## Links

### Out

- verifies [docs-conductor-default](../capability/rhiza/workflows/built-in-workflows/docs-conductor-default.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

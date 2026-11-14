# ship-one-ci-round-qa-2

evidence · `01M7DTKRFWDKYGKPRT63AY23FN` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: The ship-change-one-ci script passes, as does ship-check-settles-conflict. ship rebase --no-push against a bare origin merged origin/main into the branch, printed 'Pushed nothing (--no-push)' and left origin's branch at the same SHA. On a code conflict it left the merge in progress (UU, exit 1) with origin unchanged, and on an up-to-date branch it said nothing to merge (exit 0).

## Links

### Out

- verifies [ship-one-ci-round](../capability/rhiza/orchestration/shipping/ship-one-ci-round.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

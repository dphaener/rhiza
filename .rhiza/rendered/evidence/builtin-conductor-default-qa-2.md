# builtin-conductor-default-qa-2

evidence · `01M7DTKPPSJJX217VQK6RSXYYS` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: workflow describe ship-change in a fresh repo lists conduct, gate, signoff, push, wait-ci, merge, main-ci, close-out, ship, deploy-check, with every non-pass route up to merge going to handback:conduct and none after it; workflow check ship-change is valid; go test -run 'TestBuiltinShipChange|TestBuildCheck' pass; run build-check and review-check outside a run are refused cleanly (exit 2).

## Links

### Out

- verifies [builtin-conductor-default](../capability/rhiza/workflows/built-in-workflows/builtin-conductor-default.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

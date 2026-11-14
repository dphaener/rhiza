# leafless-refusal-by-leaves-step-qa-2

evidence · `01M7DTKQQWN1EEZZMH3J8XZGE2` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: Built binary, leafless accepted change set: run start --dry-run with ship-change and with a repo workflow including ship-change starts (exit 0). With ship-change-classic, a repo workflow including build or ship-change-classic, and a dispatch.toml pinning [ready] feature = ship-change-classic, it is refused with exit 3 and 'has no leaves: <workflow> builds and reviews its leaves'. TestCheckShipLeaves and run-ship-needs-leaves pass.

## Links

### Out

- verifies [leafless-refusal-by-leaves-step](../capability/rhiza/workflows/built-in-workflows/leafless-refusal-by-leaves-step.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

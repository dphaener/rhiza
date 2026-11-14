# handback-target-survives-include-qa-2

evidence · `01M7DTKQFTJ8WXV53WXDE1GF6S` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: In this repo, workflow describe ship-rhiza and bugfix show every hand-back route as handback:ship-change/conduct. A repo ship-change.toml including ship-change is refused as including itself. TestIncludeHandbackTarget and TestHandbackIncluded pass.

## Links

### Out

- verifies [handback-target-survives-include](../capability/rhiza/workflows/built-in-workflows/handback-target-survives-include.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic

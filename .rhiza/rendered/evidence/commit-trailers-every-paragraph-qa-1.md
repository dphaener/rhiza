# commit-trailers-every-paragraph-qa-1

evidence · `01M62X6KT0Q1SA948N0HR27EBY` · actor [claude-code-6a92fa45](../actor/claude-code-6a92fa45.md)

## result

PASS

## statement

At 4c2751d, built ./cmd/rhiza and the merge-base's binary, and in a scratch git repo (clean env, own HOME) added capability leaf, committed leaf.go with 'Change: <cs> / Actor: builder / Leaf: leaf', a blank line, then 'Co-Authored-By: Agent <a@x>', and wrote Evidence verifying leaf as builder. rhiza check: the base binary printed only 'ok' (trailers missed); the head binary warned 'evidence-actor ... its actor builder wrote that capability in commit 45ebb6f70c7f', so Change:/Actor:/Leaf: were read from the split paragraph. With all trailers in one final paragraph the head binary warned the same way; with the split shape naming Leaf: other it did not warn.

## Links

### Out

- verifies [commit-trailers-every-paragraph](../capability/rhiza/integrity/commit-trailers-every-paragraph.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.

# build-leaf-split-trailers

capability · `01M62RE25K7P1CD4GD7SP27CNM` · parent [build-workflow](../build-workflow.md)

## acceptance

An engine test against a real git repo counts a builder commit with Leaf: in its first paragraph and Co-Authored-By after it as the leaf's attempt, and a leaf fixer's commit the same way; a commit naming another leaf still fails.

## statement

A builder's or a leaf fixer's new commit counts for its leaf when its Leaf: trailer names the leaf, even with Co-Authored-By in a separate final paragraph.

## Links

### Out

- depends-on [commit-trailers-every-paragraph](../../../integrity/commit-trailers-every-paragraph.md)
- decided-by [leaf-trailer-on-build-commits](../../../../../decision/leaf-trailer-on-build-commits.md)

### In

- verifies [build-leaf-split-trailers-verified](../../../../../evidence/build-leaf-split-trailers-verified.md)
- verifies [pr-61-leaf-trailer-split-paragraphs](../../../../../evidence/pr-61-leaf-trailer-split-paragraphs.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.

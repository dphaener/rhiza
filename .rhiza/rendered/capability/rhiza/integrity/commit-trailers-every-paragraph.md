# commit-trailers-every-paragraph

capability · `01M62RE1Q8DH5JWM6YHEGZ0G15` · parent [integrity](../integrity.md)

## acceptance

A unit test of the trailer reader finds Change:, Actor:, Leaf: and Review: in a commit body whose Co-Authored-By sits in a separate final paragraph, and still reads them when all trailers share the final paragraph.

## statement

rhiza check and the workflows read Change:, Actor:, Leaf: and Review: trailers from every paragraph of a commit body, so a Leaf: trailer counts when Co-Authored-By sits in a later paragraph.

## Links

### Out

- decided-by [change-trailers](../../../decision/change-trailers.md)

### In

- depends-on [build-leaf-split-trailers](../workflows/built-in-workflows/build-workflow/build-leaf-split-trailers.md)
- verifies [commit-trailers-every-paragraph-verified](../../../evidence/commit-trailers-every-paragraph-verified.md)
- verifies [commit-trailers-every-paragraph-qa-1](../../../evidence/commit-trailers-every-paragraph-qa-1.md)
- verifies [commit-trailers-every-paragraph-verified-2](../../../evidence/commit-trailers-every-paragraph-verified-2.md)
- verifies [pr-61-leaf-trailer-split-paragraphs](../../../evidence/pr-61-leaf-trailer-split-paragraphs.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.

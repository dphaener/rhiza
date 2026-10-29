# guides-one-final-paragraph-verified

evidence · `01M62WSFDSKBMSCD3Y610TPBGM` · actor [claude-code-07d2df3e](../actor/claude-code-07d2df3e.md)

## result

PASS

## statement

At 8fbb4af, go test ./internal/guide passes and make check exits 0. Diff read: build.md's commit example now ends with Co-Authored-By in the same paragraph as Change:/Actor:/Leaf:, and review.md's fixer instruction says to keep every trailer, Co-Authored-By included, in one final paragraph with no blank line. Mutations in a scratch copy: a blank line before Co-Authored-By in the build example makes TestBuildTrailersInOneFinalParagraph fail (final paragraph lacks Change:, Actor:, Leaf: export-csv). Rewriting the fixer instruction to put Co-Authored-By in its own paragraph after a blank line makes TestReviewFixerTrailersInOneFinalParagraph fail (no 'one final paragraph', no 'no blank line').

## Links

### Out

- verifies [guides-one-final-paragraph](../capability/rhiza/guide/guides-one-final-paragraph.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.

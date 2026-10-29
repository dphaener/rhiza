# commit-trailers-every-paragraph-verified-2

evidence · `01M62XV3SNSXCJEH0S6A7KKB2C` · actor [claude-code-5b9ed8ab](../actor/claude-code-5b9ed8ab.md)

## result

PASS

## statement

At 22bf687 (the fix commit after review), go test ./internal/api -run 'TestParseTrailerLog|LeafCommit|ReviewCommit' -count=1 -v passes. In TestParseTrailerLog, case ddd reads Change: CS4, Actor: a, Leaf: split and Review: R1 from a paragraph followed by a separate Co-Authored-By paragraph. Case eee reads Change: CS5 and Leaf: one when they share the final paragraph with Co-Authored-By, and case aaa reads Change:, Actor: and Leaf: from a single final paragraph. Actor: and Review: are not exercised in the shared-paragraph shape, but the same line loop reads all four keys. Reading the diff: parseTrailerLog now reads only the trailing run of whole trailer paragraphs (trailerRun and trailerParagraph). Prose lines (ggg, hhh) and a squash merge's earlier commits' trailers (iii) are not read, and folded values are unfolded (jjj). That answers the review's major finding: rhiza check built from 22bf687 (RHIZA_ROOT, RHIZA_CHANGE and RHIZA_RUN unset) prints 0 evidence-actor warnings, against the 202 the review reported. go vet ./... and gofmt -l are clean, and go test ./... passes in every package.

## Links

### Out

- verifies [commit-trailers-every-paragraph](../capability/rhiza/integrity/commit-trailers-every-paragraph.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.

# guides-one-final-paragraph-qa-1

evidence · `01M62X6MAEDFKH7753P2BKTPSS` · actor [claude-code-6a92fa45](../actor/claude-code-6a92fa45.md)

## result

PASS

## statement

At 4c2751d, ran the built binary's rhiza guide build: its builder example ends Change:/Actor:/Leaf:/Co-Authored-By: in one paragraph, followed by 'Keep every trailer, Co-Authored-By: included, in the one final paragraph, with no blank line between them.' rhiza guide review's Fixer section says the same for the fixer's trailers. go test ./internal/guide -count=1: ok.

## Links

### Out

- verifies [guides-one-final-paragraph](../capability/rhiza/guide/guides-one-final-paragraph.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.

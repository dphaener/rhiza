# why-cited-quotes-verified

evidence · `01M6AC3669E1DGPZ2FYS2ZQ64R` · actor [claude-code-18b0777f](../actor/claude-code-18b0777f.md)

## result

PASS

## statement

At 77f9a776 ran go test ./internal/api -run TestWhy -count=1: TestWhyCitedQuotes passes (source with two quotes, decision quote 'must keep it' cited across a double space and newline, why returns only 'we must  keep\nit small'). Ran go test ./cmd/rhiza -run TestScripts/why: the new why.txtar step shows 'we must  keep' and not 'we shall grow', PASS. Read the diff: whyCitedQuotes in internal/api/lineage.go filters a source's quotes by whitespace-normalised substring of the decision's quote, only for sources reached from a decision. make check exits 0.

## Links

### Out

- verifies [why-cited-quotes](../capability/rhiza/lineage/why-walk/why-cited-quotes.md)

## History

- `01M5SWTW02ZB4590CKA09AFSQJ` `genesis why` on a decision shows only the quotes that decision cites, not every quote of its Source, so a decision's lineage stays short.

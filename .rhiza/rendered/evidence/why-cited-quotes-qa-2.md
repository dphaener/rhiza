# why-cited-quotes-qa-2

evidence · `01M6AE9AGR8Z3DKJR1DQPWW14G` · actor [claude-code-08fe33ac](../actor/claude-code-08fe33ac.md)

## result

PASS

## statement

Built 22c94aa1 into a scratch binary. On this repo's model, rhiza why name-is-genesis and six-layers now show 1 of founding-conversation's 40 object quotes (quotes_left_out 39, 41 lines instead of 235), go-one-binary 3 (its quote '3. Yes' is in 3 of them, left out 37). In a scratch repo, source talk with object quotes [{text:'alpha  one\nhere',by:owner},{text:'beta two'},{text:'gamma three'}]: why d1 (quote 'alpha one here') kept only the first quote, by field intact, left out 2, text output 'Quotes left out: 2'; why d2 kept 'beta two'; why c2, decided by d1 and d2, kept both cited quotes, left out 1; why d3 (quote matching nothing) and why c4 (a capability derived from talk) kept all 3, left out 0. TestWhyCitedQuotes*, TestScripts/why and make fmt-check vet test pass.

## Links

### Out

- verifies [why-cited-quotes](../capability/rhiza/lineage/why-walk/why-cited-quotes.md)

## History

- `01M5SWTW02ZB4590CKA09AFSQJ` `genesis why` on a decision shows only the quotes that decision cites, not every quote of its Source, so a decision's lineage stays short.

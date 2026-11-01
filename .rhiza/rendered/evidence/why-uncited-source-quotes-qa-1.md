# why-uncited-source-quotes-qa-1

evidence · `01M6ACTYR23FAX85JAWSTY1RHS` · actor [claude-code-1a59192b](../actor/claude-code-1a59192b.md)

## result

PASS

## statement

Built 4df6a76b into a scratch binary; in a scratch repo with a three-string-quote source: rhiza why on a decision whose quote matches nothing, a decision with no quote, and a capability derived from the source each printed all three quotes with no 'Quotes left out' line and quotes_left_out 0 in --json; a trimmed decision printed 'Quotes left out: 1' in text and quotes_left_out 1 in --json. go test ./internal/api -run TestWhy and TestScripts/(why|describe-golden) pass.

## Links

### Out

- verifies [why-uncited-source-quotes](../capability/rhiza/lineage/why-walk/why-uncited-source-quotes.md)

## History

- `01M5SWTW02ZB4590CKA09AFSQJ` `genesis why` on a decision shows only the quotes that decision cites, not every quote of its Source, so a decision's lineage stays short.

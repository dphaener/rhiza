# why-cited-quotes-qa-1

evidence · `01M6ACTXBRC0P6Y53PSR26JBNF` · actor [claude-code-1a59192b](../actor/claude-code-1a59192b.md)

## result

FAIL

## statement

Built 4df6a76b into a scratch binary. In a scratch repo, a source added with --json-field quotes as plain strings ['We must  keep\nit small, always.','We shall grow later.','Cheap is fine: must keep it lean too? no'] and decision small (quote 'must keep it') derived from it: rhiza why small printed only the two matching quotes and 'Quotes left out: 1', --json quotes_left_out 1. But on this repo's own model, rhiza why name-is-genesis / go-one-binary / six-layers still print all 40 quotes of founding-conversation (235 lines, quotes_left_out 0), though name-is-genesis's quote 'yes, Genesis it is, move on to the log format' is the text of exactly one of them: sources written by rhiza source add store quotes as objects {text,by,where}, and whyCitedQuotes unmarshals quotes only as []string, so it never trims a real source. Also rhiza why c2 for a capability decided by two decisions citing different quotes of one source showed the source once, trimmed to the first decision's quotes, hiding the second decision's cited quote.

## Links

### Out

- verifies [why-cited-quotes](../capability/rhiza/lineage/why-walk/why-cited-quotes.md)

## History

- `01M5SWTW02ZB4590CKA09AFSQJ` `genesis why` on a decision shows only the quotes that decision cites, not every quote of its Source, so a decision's lineage stays short.

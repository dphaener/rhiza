# why-uncited-source-quotes-verified

evidence · `01M6ACPQTDKRVFD7ECVSBXFVEF` · actor [claude-code-9c5d12ea](../actor/claude-code-9c5d12ea.md)

## result

PASS

## statement

At 88414375, go test ./internal/api -run TestWhy passes: TestWhyCitedQuotesKeepsAll covers no match, no quote and empty quote (all quotes kept, left out 0), and TestWhyNonDecisionKeepsAllQuotes covers a capability walk plus addSource via capability, question and source nodes carrying a matching quote (both quotes kept, left out 0). Mutating the guard in addSource to drop `via.Kind == op.KindDecision` makes TestWhyNonDecisionKeepsAllQuotes fail for all three kinds, so the attempt-1 finding is fixed. TestWhyCitedQuotes checks QuotesLeftOut 1 on a trimmed case; why.txtar asserts 'Quotes left out: 1' in text and '"quotes_left_out": 1' in --json; describe-golden.txtar was regenerated with quotes_left_out and the describe-golden and why scripts pass. go test ./... and go vet pass. The make self-check fails only on stale run actor renders (.rhiza/rendered/actor/*, index.md), and the same failure shows at the previously verified commit 3078bc96, so this leaf does not cause it.

## Links

### Out

- verifies [why-uncited-source-quotes](../capability/rhiza/lineage/why-walk/why-uncited-source-quotes.md)

## History

- `01M5SWTW02ZB4590CKA09AFSQJ` `genesis why` on a decision shows only the quotes that decision cites, not every quote of its Source, so a decision's lineage stays short.

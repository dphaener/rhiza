# why-cited-quotes-verified-2

evidence · `01M6ADVGCGX1TFPF2NXF1060HV` · actor [claude-code-e9ffa7cc](../actor/claude-code-e9ffa7cc.md)

## result

PASS

## statement

Verifier (attempt 2, cycle 2) at 4dff386b: go test ./internal/api -run TestWhy -count=1 passes, incl. TestWhyCitedQuotes (two-quote source, decision citing one returns only that quote, quotes_left_out 1) and new TestWhyCitedQuotesShared/Objects; go test ./cmd/rhiza -run TestScripts/why passes (why.txtar shows the cited quote and 'Quotes left out: 1'). Review finding f1 reproduced with a built binary in a scratch repo: source talk [alpha one, beta two, gamma three], d1/d2 derived-from talk, d2 supersedes d1 — 'rhiza why d2' and 'why' on a capability decided-by both now show alpha one and beta two with 'Quotes left out: 1'; a capability also derived-from talk shows all three. make fmt-check vet test exit 0; make self-check fails only on stale .rhiza/rendered actor/index files from the run's own model writes, not on code. Diff read: addSource now unions cited quotes per source via whyCite, object quotes match by text.

## Links

### Out

- verifies [why-cited-quotes](../capability/rhiza/lineage/why-walk/why-cited-quotes.md)

## History

- `01M5SWTW02ZB4590CKA09AFSQJ` `genesis why` on a decision shows only the quotes that decision cites, not every quote of its Source, so a decision's lineage stays short.

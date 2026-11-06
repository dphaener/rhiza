# import-autonomous-decisions-verified

evidence · `01M6RMQ7Q4AQQH3AD0NNRJCV2H` · actor [claude-code-2eeb6b9b](../actor/claude-code-2eeb6b9b.md)

## result

PASS

## statement

At 5a20f3f8 I read the diff: internal/guide/text/import.md step 4 now states the load-bearing rule (bring in only a decision that shapes a capability or constraint) and the exact wording form '<as the document states it>, made autonomously by the agent and not confirmed by the owner', with no quote field and a derived-from link to the document's Source. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(import|guide)' -count=1: import-autonomous-decisions, import-records-built, guide and guide-install all PASS. The new script checks rhiza guide import prints the rule and form, writes a decision by that form, and asserts rhiza why shows the wording and '- keep-in-ram derived from decisions-md' with no quote, and rhiza show --json has no quote field. go test ./internal/guide/... passes.

## Links

### Out

- verifies [import-autonomous-decisions](../capability/rhiza/guide/import-autonomous-decisions.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows

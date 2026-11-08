# conflict-question-blocks-change-verified-2

evidence · `01M6D3CVB63BQX9NQ7VKMDE0XW` · actor [claude-code-0997b575](../actor/claude-code-0997b575.md)

## result

PASS

## statement

At 642735de I read the leaf diffs (9394ffd1: shipRaiseQuestion adds a blocks link from the raised question to the change set, and an existing identical question without the link gains it; 642735de: the test drops the link and re-raises). go test ./cmd/rhiza -run TestScripts/run-ship-change-rebase -count=1 passed. Case 2 of that script has a code conflict on feature.txt and asserts: the question shows a blocks link; rhiza inbox --json lists it as kind blocker with blocks 'this change set'; raising again on the same conflict leaves the .rhiza/log bytes unchanged; after unlinking, raising again restores the link with only one question; rhiza question answer rebase-conflict-ship-two 'Keep main' --change <cs> exits 0 and folds clean. make fmt-check vet test exit 0. make self-check fails only on stale rendered actor files from the review run's chore commits, not on this leaf.

## Links

### Out

- verifies conflict-question-blocks-change (retired)

## History

- `01M6D12NNWQZHEWFBX60JS04AG` A run's rebase-conflict question blocks its change set, so the inbox and question answer see it

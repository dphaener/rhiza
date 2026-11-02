# whatif-reads-in-budget-verified-2

evidence · `01M6CRTA7VMS9PC0PJH9NJEZGR` · actor [claude-code-356ffe40](../actor/claude-code-356ffe40.md)

## result

PASS

## statement

At c974e89e (attempt 2): go test -count=1 ./internal/api -run TestRead -v passed. TestReadBudget timed every budgetReads read in what-if:<fixture draft> on the 50,000-node model with medians show 5.8ms, ls 6.2ms, tree 7.7ms, why 6.7ms, log 7.4ms, roadmap 10.3ms, ready 10.6ms, all under 50 ms. TestReadBackendsAgree (stored vs what-if over Node, Children, Ancestors, Links, Subtree(Light), List, Resolve, Link, ResolveLink, NodesLight, LinksAt, LinksTo, EntriesOf, ChangeSets, Entries) passed. go test ./internal/index -run TestPropertyOverlayAgrees passed 100 rapid cases (183 nodes, 97 overlays) comparing the overlay with model.WhatIf. BenchmarkShowWhatIf's 'not budgeted (D14)' note is removed in budget_test.go. read.go answers what-if views through index.WhatIfOverlay. make check passed fmt, vet and test; self-check flagged only the render of this verifier's own session actor.

## Links

### Out

- verifies [whatif-reads-in-budget](../capability/rhiza/views/fast-reads/whatif-reads-in-budget.md)

## History

- `01M5SWTFB3ZAX4YKX4KA2ZVSZR` An incremental what-if folds one change set over the indexed as-is, so what-if reads meet the 50 ms target too.

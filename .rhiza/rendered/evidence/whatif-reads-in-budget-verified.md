# whatif-reads-in-budget-verified

evidence · `01M6CAHMW858KAVR2VC135BG4A` · actor [claude-code-0ffe3c72](../actor/claude-code-0ffe3c72.md)

## result

PASS

## statement

At 11ef4621: go test ./internal/api -run TestReadBudget -count=1 -v passed; it now times all 7 budgetReads (show, ls, tree, why, log, roadmap, ready) in what-if:<the fixture's draft> on the 50,000-node model, medians 7.0-11.4 ms on the first attempt, well under 50 ms. The 'not budgeted (D14)' note on BenchmarkShowWhatIf is removed (diff of budget_test.go). TestReadBackendsAgree (read_test.go stored-versus-what-if comparison over Node, Children, Ancestors, Links, Subtree, SubtreeLight, List, Resolve, Link, ResolveLink, NodesLight, LinksAt, LinksTo, EntriesOf) passed, now asserting the what-if reader uses the overlay (ov) backend. internal/index TestPropertyOverlayAgrees (rapid, 100 generated logs, 99 overlays, 209 nodes) passed, comparing overlay answers with model.WhatIf. make check exited 0.

## Links

### Out

- verifies [whatif-reads-in-budget](../capability/rhiza/views/fast-reads/whatif-reads-in-budget.md)

## History

- `01M5SWTFB3ZAX4YKX4KA2ZVSZR` An incremental what-if folds one change set over the indexed as-is, so what-if reads meet the 50 ms target too.

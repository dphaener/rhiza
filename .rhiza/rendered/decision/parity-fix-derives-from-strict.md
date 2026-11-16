# parity-fix-derives-from-strict

decision · `01M7GGERTNQMVS8C3Y5BVRJS1M`

## question

How is a disagreement the parity test finds fixed: patch each rule in model.go and write_indexed.go by hand, or derive both from the strict engine?

## reasoning

The change set's evidence is that the model rules are implemented three times; patching each case by hand keeps three copies and the next divergence. The strict views and the index's strict overlay already exist (model.WhatIfFrom, model.WhatIfDelta, index.WhatIfOverlay), so deriving from them is less code. The interim model is deleted by the milestone's projection change set, so effort spent re-implementing rules in it is thrown away. Building as-is once and seeding each what-if keeps the daemon's per-candidate cost from growing, which a plain swap to model.WhatIf would not.

## statement

A disagreement is fixed by having the non-strict side take its result from the strict engine, not by re-implementing the strict rule a third time: newModelSets builds its nodes, links, versions and conflicts from internal/model's as-is and what-if views (AsIsState once, WhatIfFrom per change set where a caller builds several), and writeIndex.buildModel lays the change set over the index's strict what-if (WhatIfOverlay's path, whatIfSets then model.WhatIfDelta over the stored seed, with its problems, through an index entry point that takes the change set as over, as WhatIfSets does, so a dry run's pending entries are laid over too) instead of writeLayer's interim apply. The daemon's per-candidate models (daemon_fold.go, proposals.go, inbox.go) must not rebuild as-is once per candidate.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [parity-harness-plan](../source/parity-harness-plan.md)

### In

- decided-by [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)
- decided-by [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.

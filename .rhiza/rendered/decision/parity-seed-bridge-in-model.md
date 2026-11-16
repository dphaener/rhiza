# parity-seed-bridge-in-model

decision · `01M7GHVXD5JEMBCN02TEQ6NJ4D`

## reasoning

The finding is right in the code: AsIsState (internal/model/seed.go:47) returns (Model, []Problem, State) and WhatIfFrom (seed.go:115) takes a Seed; the only Seed builder is index.loadSeed, which reads back rows that writeAsIsState (internal/index/whatif.go:38) split out of AsIsState's result. Rebuilding that in internal/api would be a fourth copy of seed rules, which parity-fix-derives-from-strict set out to avoid, and calling model.WhatIf per candidate would rebuild as-is once per candidate, which that Decision forbids for daemon_fold.go, proposals.go and inbox.go. The builder inside internal/model already holds everything a Seed needs (aspects per record, placed marks, problems and held conflicts), so a constructor there is a few lines next to AsIsState and keeps the seed rules in the one package that is the spec. The leaf's ban on changing internal/model was meant to protect the rules (the strict engine is the spec), not to forbid an additive constructor, so the ban is narrowed to rules. No test checks WhatIfFrom against WhatIf today, so the constructor gets its own equivalence test. Moving the index onto the same constructor is out of this change set's scope.

## statement

The bridge is a constructor in internal/model: model.AsIsSeed(sets []changeset.ChangeSet) (Model, []Problem, Seed) builds the as-is view once and returns it with the Seed WhatIfFrom takes (Known from every change set's ID, Nodes and Links with their History aspects and Placed, Problems without the duplicate-handle ones, State without Aspects), from the same builder AsIsState uses. newModelSets' callers call AsIsSeed once and WhatIfFrom per change set; internal/api assembles no Seed by hand. internal/model joins parity-interim-from-strict's touch list for this constructor only, a non-rule change: it adds no rule and changes no existing function's result. Its own test in internal/model checks, on random logs, that WhatIfFrom over AsIsSeed's Seed equals WhatIf over the whole log for every in-flight change set, and that AsIsSeed's model and problems equal AsIs's. Change set 01M7GFT8AD (check without the what-if rebuilds) plans the same model.AsIsSeed(sets) (Model, []Problem, Seed) in internal/model/seed.go (its model-seed-views leaf) and a copy in Seed.builder so views built from one seed never change it; parity-interim-from-strict uses that name and signature so the two share one constructor: whichever lands first adds it, and the other builds on it instead of adding a second.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7gg7gxsa97hesrnfebfcz4t-r3-f2](../question/review-01m7gg7gxsa97hesrnfebfcz4t-r3-f2.md)

### In

- decided-by [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.

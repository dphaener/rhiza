# review-01m7gg7gxsa97hesrnfebfcz4t-r3-f2

question · `01M7GHSMDA01JXB1K046DQSK0B`

## statement

Plan review round 3 left finding f2 (contradiction on parity-fix-derives-from-strict): The Decision (and parity-interim-from-strict, which follows it) says newModelSets builds from 'AsIsState once, WhatIfFrom per change set', assuming AsIsState's result feeds WhatIfFrom. It does not: model.AsIsState (internal/model/seed.go:47) returns (Model, []Problem, State), while model.WhatIfFrom (seed.go:115) takes a model.Seed (Known IDs, SeedNodes with Aspects and Placed, as-is Problems without the duplicate-handle ones, State), and the only code that builds a Seed is index.loadSeed from the database (internal/index/whatif.go:137). The interim leaf touches only internal/api and must not change internal/model, so it would have to reassemble a Seed by hand in internal/api (splitting State.Aspects/Placed back onto nodes and links and filtering the duplicate-handle problems), a fourth copy of seed rules the Decision set out to avoid. The plan should name the bridge: a model constructor (e.g. a Seed returned by or built from AsIsState, adding internal/model to the leaf's touch list as a non-rule change) or an explicit api-side assembly with its own test. The critic confirmed it: Checked internal/model/seed.go: AsIsState (47) returns (Model, []Problem, State), while WhatIfFrom (115) takes a Seed (Known, Nodes/Links as SeedNode/SeedLink carrying Aspects and Placed, Problems without duplicate-handle ones, State without Aspects). No code builds a Seed from AsIsState's result; the only constructor in the tree is index.loadSeed from the database (internal/index/whatif.go:137). @parity-fix-derives-from-strict and parity-interim-from-strict prescribe 'AsIsState once, WhatIfFrom per change set' as if one fed the other, while the leaf may not change internal/model and touches only internal/api, so it would have to reassemble the seed rules by hand in api, the copy the Decision set out to avoid. The plan contradicts the code and names no bridge; major stands. How should it be settled?

## Links

### In

- answers [parity-seed-bridge-in-model](../decision/parity-seed-bridge-in-model.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.

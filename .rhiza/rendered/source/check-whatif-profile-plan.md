# check-whatif-profile-plan

source · `01M7GGD0KJ19CHX2Y4WSGQK1ER`

## quotes

```json
[]
```

## raw_sha256

3014df041ff17734dfb31f0c5c4006ba5a04b0961c962834dfc8e9b45afe3e17

## session

ea75812d-c986-498a-88e5-01a1e1043949

## summary

Autonomous plan run (no owner). The arbiter asked that the builder take the profile and allowed one profiled rhiza check now. That profile was taken on 2026-11-15 in the main checkout (24,560 entries, 127 MB). rhiza check took 13.8 s wall, with 18.5 s of CPU samples. model.WhatIf was 46% of the samples (check's in-flight problems pass 23%, its in-flight links pass 25%), model.asIs 26%, Model.clone in finishModel 14%. store.ReadLog was 7% (ReadEntries 6.4%, entry.Decode 4.3%), changeset.FoldLog 3%, textRefs 6%, the render comparison 2.4%. The repeated what-if builds are therefore the main cost, and the parallel decode can save about a second at most. No in-memory constructor of model.Seed exists today: only internal/index builds one, from its stored rows (loadSeed), so check needs a model function that builds the Seed from the change sets. Design choices were made by the planner and are recorded as assumed.

## title

Planning check without what-if rebuilds, autonomous, with a profile

## Links

### In

- derived-from [check-seed-in-memory](../decision/check-seed-in-memory.md)
- derived-from [check-time-target](../decision/check-time-target.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.

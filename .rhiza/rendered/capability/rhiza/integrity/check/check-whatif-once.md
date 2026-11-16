# check-whatif-once

capability · `01M7GFXTMA9VHC2FAQB33P1KA1` · parent [check](../check.md)

## acceptance

Tests: go test -trimpath -p 2 -count=1 ./internal/check passes. A new TestWhatIfOnceMatches in internal/check, over a log with a landed base, two in-flight change sets where one depends on the other, one in-flight change set whose ops no longer apply, a duplicate handle, an in-flight link whose ends its type does not allow and an in-flight node whose text refers to a node that does not exist, shows check's in-flight findings listed exactly. Failure case: an in-flight change set whose what-if has problems, a bad in-flight link and a broken in-flight text ref are each reported once, not twice. Profile, at this leaf's base before any code change and again after: go build -trimpath -o $TMPDIR/rhiza-x ./cmd/rhiza, then in this repo's main checkout RHIZA_CPUPROFILE=$TMPDIR/check.prof $TMPDIR/rhiza-x check and go tool pprof -top -cum $TMPDIR/rhiza-x $TMPDIR/check.prof. The commit's Checks: block gives /usr/bin/time -f '%e s' for both builds, and the cumulative shares of store.ReadLog, changeset.FoldLog, model.WhatIf/WhatIfFrom, model.asIs and checker.render before and after. After the change, model.WhatIf and model.ToBe no longer appear under check.Run, and model.asIs is called once per check. Touches internal/check/check.go, internal/check/textrefs.go and their tests. Must not change what check reports, its order or its --json shape.

## statement

rhiza check builds the as-is view once with model.AsIsSeed, and from that one seed the to-be view with model.ToBeFrom (replacing checker.model's model.ToBe call, internal/check/check.go:378) and each in-flight change set's what-if once with model.WhatIfFrom; the in-flight problems pass (checker.model, :384), the in-flight links pass (checker.inFlightLinks, :456) and the text-ref pass (checker.textRefs, internal/check/textrefs.go:49, run from check.Run at :282) share those what-if views instead of each calling model.WhatIf

## Links

### Out

- depends-on [model-seed-views](model-seed-views.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [check-seed-in-memory](../../../../decision/check-seed-in-memory.md)
- decided-by [split-check-whatif-once](../../../../decision/split-check-whatif-once.md)

### In

- depends-on [read-entries-parallel](read-entries-parallel.md)
- blocks [review-01m7gg7hsbsb356bvggqgh8yhb-r3-f1](../../../../question/review-01m7gg7hsbsb356bvggqgh8yhb-r3-f1.md)
- verifies [check-whatif-once-verified](../../../../evidence/check-whatif-once-verified.md)
- verifies [check-whatif-once-qa-1](../../../../evidence/check-whatif-once-qa-1.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.

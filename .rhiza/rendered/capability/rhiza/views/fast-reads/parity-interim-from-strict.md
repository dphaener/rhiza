# parity-interim-from-strict

capability · `01M7GGYX1AJ24Y40XJK9GRY8P8` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestPropertyProjectionParity|TestParity)' passes with no entry of the known-cases table naming the interim model and no skipped draws, at -rapid.checks=1000 as well as the default; go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestWrite|TestDaemonPicks|TestDaemonPlansProposals|TestDaemonCandidates|TestInbox)' and go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^(daemon-picks|daemon-plans-proposals)$' pass. Fix by building newModelSets from internal/model's AsIsSeed and WhatIfFrom as @parity-fix-derives-from-strict says, with newModelSets' callers building as-is once. Each case found on the interim side has its own TestParity<Case> test with the smallest log that shows it, and a note on this leaf (rhiza note add parity-interim-from-strict ...) naming the log shape, both results and which side changed; when the harness finds no interim case beyond the named one and the version rule, the note says so. Failure case: the regression test for the named case on the interim model (a node whose create no longer applies in the strict what-if) fails on the code before the fix (show it in the commit's Checks: block by running it at the parent commit). The bridge is model.AsIsSeed(sets) (Model, []Problem, Seed) in internal/model/seed.go, as @parity-seed-bridge-in-model says: newModelSets' callers call it once and WhatIfFrom per change set, and internal/api assembles no Seed by hand. go test -trimpath -count=1 -p 2 ./internal/model passes with a new TestAsIsSeedWhatIfFrom which, on random logs (a landed base, in-flight change sets one depending on another, one whose ops no longer apply, a duplicate handle), checks that AsIsSeed's model and problems equal AsIs's and that WhatIfFrom over its Seed equals WhatIf over the whole log (model and problems) for every in-flight change set, all WhatIfFrom calls over one shared seed and compared only after every call has run, and the seed compared with a fresh AsIsSeed after them; if a view changes the shared seed, Seed.builder copies each node's and link's History, Fields and aspects. Touches internal/api/model.go, its tests, internal/api/daemon_fold.go, proposals.go and inbox.go only where newModelSets' callers must build as-is once, and internal/model/seed.go and its tests for AsIsSeed (and that copy) only. Must not change internal/api/write_indexed.go (@parity-write-index-from-strict fixes it), internal/model or internal/index rules (the strict engine is the spec; AsIsSeed adds no rule and changes no existing function's result) or a --json shape. Writes no longer refusing a node whose two landed heads touch disjoint fields is the rule @parity-version-rule settles, not a case to raise; any other case that would change user-visible output on this repo's own log is raised as a question blocking this leaf instead.

## statement

Every case where the parity test finds the interim model (newModel and newModelSets in internal/api/model.go) disagreeing with the strict engine is fixed on the interim side, starting with the one write_indexed.go's model comment names (the interim model keeps a node whose create no longer applies, the strict what-if drops it) and the version and conflict rule of @parity-version-rule, and each case found on this side has a named regression test and a note on this leaf

## Links

### Out

- depends-on [projection-parity-property](projection-parity-property.md)
- decided-by [strict-engine-is-spec](../../../../decision/strict-engine-is-spec.md)
- decided-by [parity-version-rule](../../../../decision/parity-version-rule.md)
- decided-by [parity-fix-derives-from-strict](../../../../decision/parity-fix-derives-from-strict.md)
- decided-by [parity-seed-bridge-in-model](../../../../decision/parity-seed-bridge-in-model.md)

### In

- depends-on [api-tests-on-index](../../integrity/self-check/api-tests-on-index.md)
- depends-on [parity-write-index-from-strict](parity-write-index-from-strict.md)
- verifies [parity-interim-from-strict-verified](../../../../evidence/parity-interim-from-strict-verified.md)
- verifies [parity-interim-from-strict-qa-1](../../../../evidence/parity-interim-from-strict-qa-1.md)
- verifies [parity-interim-from-strict-qa-2](../../../../evidence/parity-interim-from-strict-qa-2.md)
- discusses [note-01m7gn4h8wddvb](../../../../note/note-01m7gn4h8wddvb.md)
- discusses [note-01m7gn4kne8y9y](../../../../note/note-01m7gn4kne8y9y.md)
- discusses [note-01m7gn4nj603gs](../../../../note/note-01m7gn4nj603gs.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.

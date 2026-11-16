# parity-write-index-from-strict

capability · `01M7GGYX9DWE3H3TZ3G0TA4T8P` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestPropertyProjectionParity|TestParity)' passes with the known-cases table empty and no skipped draws, at -rapid.checks=1000 as well as the default; go test -trimpath -count=1 -p 2 ./internal/api -run '^TestWrite' and go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^set$' pass. Fix by having writeIndex.buildModel lay the change set over the index's strict what-if and its problems, in place of writeLayer's interim apply (write_indexed.go:164-209) and versionOf (write_indexed.go:332), as @parity-fix-derives-from-strict says: index.WhatIfOverlay as it stands (overlay.go:38) folds cs from its stored entries only (whatIfSets(tx, cs, nil)) and would drop a dry run's pending entries, so the fix adds to internal/index an overlay entry point that takes the folded change set as over, as WhatIfSets(cs, over) does (whatIfSets(tx, cs, over), then model.WhatIfDelta over the stored seed and OverlayOf), and buildModel passes it the change set folded from ix.entries(w, current), pending entries included; a parity draw with dry-run pending entries matches the strict fold. Each case found on the index side has its own TestParity<Case> test with the smallest log that shows it, and a note on this leaf (rhiza note add parity-write-index-from-strict ...) naming the log shape, both results and which side changed; when the harness finds no index case beyond the named one and the version rule, the note says so. Failure case: the regression test for the named case on the write model fails on the code before the fix (show it in the commit's Checks: block by running it at the parent commit), and write_test.go's LandedConflict (two heads writing one field) still refuses set, mv and retire. Touches internal/api/write_indexed.go and its tests, and internal/index/overlay.go for that entry point only. Must not change internal/api/model.go (@parity-interim-from-strict fixes it), internal/model or internal/index rules (the strict engine is the spec) or a --json shape. Writes no longer refusing a node whose two landed heads touch disjoint fields is the rule @parity-version-rule settles, not a case to raise; any other case that would change user-visible output on this repo's own log is raised as a question blocking this leaf instead.

## statement

Every case where the parity test finds write_indexed.go's write model (writeIndex.model) disagreeing with the strict engine is fixed on the index side, starting with the one its model comment names (it keeps a node whose create no longer applies, the strict what-if drops it) and the version and conflict rule of @parity-version-rule, and each case found on this side has a named regression test and a note on this leaf

## Links

### Out

- depends-on [projection-parity-property](projection-parity-property.md)
- depends-on [parity-interim-from-strict](parity-interim-from-strict.md)
- decided-by [strict-engine-is-spec](../../../../decision/strict-engine-is-spec.md)
- decided-by [parity-version-rule](../../../../decision/parity-version-rule.md)
- decided-by [parity-fix-derives-from-strict](../../../../decision/parity-fix-derives-from-strict.md)

### In

- depends-on [api-tests-on-index](../../integrity/self-check/api-tests-on-index.md)
- verifies [parity-write-index-from-strict-verified](../../../../evidence/parity-write-index-from-strict-verified.md)
- verifies [parity-write-index-from-strict-qa-1](../../../../evidence/parity-write-index-from-strict-qa-1.md)
- verifies [parity-write-index-from-strict-qa-2](../../../../evidence/parity-write-index-from-strict-qa-2.md)
- discusses [note-01m7gpv9w8mr88](../../../../note/note-01m7gpv9w8mr88.md)
- discusses [note-01m7gpvcdnzke7](../../../../note/note-01m7gpvcdnzke7.md)
- discusses [note-01m7gpve6s3vax](../../../../note/note-01m7gpve6s3vax.md)
- discusses [note-01m7gpvg63bf3s](../../../../note/note-01m7gpvg63bf3s.md)
- discusses [note-01m7gpvj2yy7w8](../../../../note/note-01m7gpvj2yy7w8.md)
- discusses [note-01m7gpvkt9m0kt](../../../../note/note-01m7gpvkt9m0kt.md)
- discusses [note-01m7gpvnek94q7](../../../../note/note-01m7gpvnek94q7.md)
- discusses [note-01m7gpvq6601ma](../../../../note/note-01m7gpvq6601ma.md)
- discusses [note-01m7gpvrt7kgct](../../../../note/note-01m7gpvrt7kgct.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.

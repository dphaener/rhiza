# parity-write-index-from-strict-verified

evidence · `01M7GQ6266JGBNFTRZ49H1N9Y0` · actor [claude-code-d7628b58](../actor/claude-code-d7628b58.md)

## result

PASS

## statement

Verifier, detached checkouts at f66f6c0a3 and its parent 5e991186a, every go test under a capped systemd scope with -p 2. At the parent with parity_index_test.go copied in, go test -trimpath -count=1 -p 2 ./internal/api -run '^TestParityIndex' exits 1 with all nine TestParityIndex cases failing. At f66f6c0a3: -run '^TestParityIndex' exit 0 (9 pass); -run '^(TestPropertyProjectionParity|TestParity)' exit 0 (100 draws, 42 with two landed heads, 76 with a rebuild, 37 with pending entries, 3619 views, no known-case lines, no skips); the same at -rapid.checks=1000 exit 0 (1000 draws, 410 with pending entries, 36466 views); -run '^TestWrite' exit 0; -run '^TestSet/LandedConflict' exit 0 (write_test.go unchanged, it still checks set, mv and retire refuse with exit 3 'conflicting changes'); ./internal/index/... exit 0; ./cmd/rhiza TestScripts/^set$ and TestScripts/^describe-golden$ exit 0. Diff read against 5e991186a: knownCases is []knownCase{} and the pending draw compares the index's write view with dry-run pending entries to the strict fold (vacuity check requires pending draws); internal/index/overlay.go adds only WhatIfOverlayOver(cs, over) passing over to whatIfSets, WhatIfOverlay delegating with nil; no change under internal/model, nor to internal/api/model.go or any other internal/index file; the response types are unchanged (describe-golden passes). writeLayer and the interim versionOf are gone; buildModel passes the change set folded from ix.entries(w, current), pending included, and Version/Conflict follow Decision parity-version-rule. The leaf has nine notes, one per case, each naming log shape, both results and the index side as changed, the last saying no further index case; the only own-log change claimed (question-answer-op no longer Conflict) is the disjoint-field heads rule of parity-version-rule. No stubs or TODOs.

## Links

### Out

- verifies [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.

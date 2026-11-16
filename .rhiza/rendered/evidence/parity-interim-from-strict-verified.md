# parity-interim-from-strict-verified

evidence · `01M7GNKVPJZ3CCYBNRERFZ5ZC0` · actor [claude-code-45ce1497](../actor/claude-code-45ce1497.md)

## result

PASS

## statement

Verified at 7e9489243 in a detached checkout (diff against a735000dc). TestParityCreateNoLongerApplies copied into the parent a735000dc fails (exit 1: interim present, strict absent for the child). At 7e9489243 these pass: go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestPropertyProjectionParity|TestParity|TestModel)' (rapid 100 OK); the same parity run at -rapid.checks=1000 (rapid 1000 OK, no skipped draws); ./internal/model including TestAsIsSeedWhatIfFrom (116 what-ifs, 15 with an in-flight dependency, 228 skipped ops, 42 duplicate handles; one shared seed, compared after every call, seed DeepEqual a fresh AsIsSeed); ./internal/api -run '^(TestWrite|TestDaemonPicks|TestDaemonPlansProposals|TestDaemonCandidates|TestInbox)'; cmd/rhiza TestScripts daemon-picks, daemon-plans-proposals and set; rhiza ship generate and rhiza render with a scratch build leave the tree clean. Diff read: knownCases has no interim-side entry (every entry requires SideA index); internal/model changes only add AsIsSeed, split AsIsState's state into builder.state with the same result, and copy History, Fields and aspects in Seed.builder, adding no rule; internal/api builds no Seed by hand; daemon_fold.go and proposals.go build one modelBase for every candidate, and inbox.go's one as-is newModelSets call needs no change; write_indexed.go, internal/index and cmd are untouched, so no --json shape changed; the three leaf notes name the log shape, both results and the side that changed, and say the harness found no other interim case.

## Links

### Out

- verifies [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.

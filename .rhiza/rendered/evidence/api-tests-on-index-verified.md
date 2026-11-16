# api-tests-on-index-verified

evidence · `01M7GTFR1974549EGP8B0NDXEB` · actor [claude-code-9dd24b50](../actor/claude-code-9dd24b50.md)

## result

PASS

## statement

Verified commit 24ef387fc against d343922ed. Diff: only internal/api/*_test.go plus .rhiza log/rendered changed, no non-test code. testCtx(t, root) gives every Ctx IndexPath=filepath.Join(t.TempDir(),"index.db") (fresh per call, so no index shared between parallel tests); testCtxNoIndex is used only by snapshot build/copy (copy then sets its test's TempDir index), drawsOf and runRead (which sets the generated repo's index), each with a comment. Every IndexPath="" (change_indexed_test x4, write_indexed_test, inbox_fold_test, read_test) has a comment saying it exercises the fallback/fold side on purpose. No Test func deleted; go test -list '.*' ./internal/api counts 782 Test lines at both d343922ed and HEAD. Ran (all in a capped systemd scope, -p 2): go test -race -count=5 -run 'TestDaemonWebListenerFollowsPort$' exit 0; go test -count=5 -run 'TestArbiter|TestStartArbiter|TestRunCycle|TestRunLeaves' exit 0; go test -count=1 -run 'TestChangeShowIndexed|TestWritesThroughIndexMatchFold|TestInboxReadsChangeSetsOnce|TestReadOpenErrors|TestDaemonProfile' exit 0 (TestDaemonProfile tests pass without -cpuprofile); whole package once: go test -trimpath -count=1 -p 2 -shuffle=7331 ./internal/api exit 0 (227s). Profile not re-run (whole-package budget); Checks: block gives store.ReadLog cum share before 24.87% and after 23.94%, which falls. rhiza check fails only on this verifier's own unrendered actor record, as expected in a run.

## Links

### Out

- verifies [api-tests-on-index](../capability/rhiza/integrity/self-check/api-tests-on-index.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.

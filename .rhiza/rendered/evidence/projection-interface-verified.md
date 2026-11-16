# projection-interface-verified

evidence · `01M7H3D0X0RKZFY9QF9TEHKNXT` · actor [claude-code-79c1e985](../actor/claude-code-79c1e985.md)

## result

PASS

## statement

Verified c8d19b528 in its worktree. All Checks: commands re-run under the capped scope: go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestProjection|TestPropertyProjectionParity|TestWrite)' exit 0; TestProjection* -v: DryRunSeesPending, OpensOnce, NilWithoutIndex, IndexNotADatabase, IndexCannotOpen all PASS; -count=3 -run '^TestDaemonIdlePollSkipsFold$' exit 0; go vet ./internal/api exit 0; rhiza ship generate and rhiza render exit 0, with no diff beyond the verifier's own actor. Diff read: projection.go is the new interface (changeSetIDs, changeSets, entries, problems, model, failure, wrote, close), and writeIndex is its only implementation, plus failedProjection for an open error. workspace.go opens it lazily through projectionSlot, shared by dry-run copies. The writes (write.go, apply.go), run_model_cache.go, inboxGateRead and gateAnswerCtx (question.go) go through it. The parity harness's indexView takes a projection and calls p.model(pending, cs). The tests cover every acceptance point: stamp unchanged in a dry run, one open across three rounds (openedProjection hook in api.go), nil projection with empty IndexPath, a garbage index rebuilt with the same answer, and an unopenable index returning the open's error with no log read. Must-nots hold: read.go, model.go, status, render, change ls and daemon fold are untouched; no IndexPath branch or interim model is deleted; no TODOs. Small scope note: question.go, apply.go, write.go, change_indexed.go and api.go are outside the listed files, but each change is the mechanical w.ix-to-projection swap or the acceptance's test hook. The builder's full-run TestDaemonIdlePollSkipsFold failure predates this commit: on 03b972942 the test is already t.Parallel() and diffs the package-global planReadCount/planFoldCount (idle_poll_test.go's TestRootIdlePoll stays serial for exactly that reason), and this commit touches neither plan.go, daemon_fold*.go nor that test.

## Links

### Out

- verifies [projection-interface](../capability/rhiza/views/fast-reads/projection-interface.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.

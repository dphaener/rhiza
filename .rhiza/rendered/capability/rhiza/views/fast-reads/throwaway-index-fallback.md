# throwaway-index-fallback

capability · `01M7GFZTA52SPV1ZK9DP14MXXK` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run '^TestProjection' passes, with TestProjectionNilWithoutIndex replaced by a test that with ctx.IndexPath empty a status, a show and an add give the same answers as with an index and leave no temporary folder behind, and TestProjectionIndexCannotOpen changed to show the same with IndexPath under a read-only folder. Failure case: when the temporary folder cannot be created either (TMPDIR unwritable) the call fails with an error naming both the index path and the temporary folder. Follows run_leaves.go:685's os.MkdirTemp pattern. Touches internal/api/projection.go, write_indexed.go and projection_test.go. Must not change the default index location or cmd/rhiza.

## statement

When a call's index path is empty or the index there cannot be opened or created (no user cache folder, a read-only cache), internal/api's projection open (openWriteIndex) builds a throwaway index in a temporary folder for that call, removed when the projection closes, so every read of every root goes through the projection and none falls back to folding the log

## Links

### Out

- decided-by [projection-no-index-fallback](../../../../decision/projection-no-index-fallback.md)
- decided-by [throwaway-index-in-api](../../../../decision/throwaway-index-in-api.md)

### In

- depends-on [interim-model-deleted](interim-model-deleted.md)
- depends-on [folds-from-projection](folds-from-projection.md)
- depends-on [writes-from-projection](writes-from-projection.md)
- depends-on [other-roots-own-projection](other-roots-own-projection.md)
- verifies [throwaway-index-fallback-verified](../../../../evidence/throwaway-index-fallback-verified.md)
- verifies [throwaway-index-fallback-qa-1](../../../../evidence/throwaway-index-fallback-qa-1.md)
- verifies [throwaway-index-fallback-qa-2](../../../../evidence/throwaway-index-fallback-qa-2.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).

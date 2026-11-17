# throwaway-index-fallback-verified

evidence · `01M7K3D2D5AD60FEW2AAD6MF0Y` · actor [claude-code-ae1aa029](../actor/claude-code-ae1aa029.md)

## result

PASS

## statement

At b1ae2e86f, capped (MemoryMax=6G, TasksMax=300), RHIZA_/SHIP_/REX_ unset: go test -trimpath -count=1 -p 2 ./internal/api -run '^TestProjection' exit 0 (TestProjectionThrowawayWithoutIndex replaces TestProjectionNilWithoutIndex: empty IndexPath, status/show/dry-run add JSON equal the indexed answers, temp folder empty after; TestProjectionIndexCannotOpen uses IndexPath under a 0555 folder, same comparison; TestProjectionNoTempDir: read-only index folder and TMPDIR, error contains both paths). TestPropertyProjectionParity exit 0; the Checks: block's -run '^(TestInbox|TestProjection|TestStatus|TestChangeLs|TestReadOpenErrors)' exit 0; the whole ./internal/api package exit 0 (221s); go vet ./internal/api exit 0. Diff read: openWriteIndex falls back to os.MkdirTemp(tempDir(ctx), "rhiza-projection-") as commitOnHead does, close removes it, useProjection no longer returns nil for an empty IndexPath, and the reader's index goes through openWriteIndex. Neither cmd/rhiza nor the default index location changed. No TODOs or stubs. The extra test-file edits (change, inbox_fold, read, status) only update earlier expectations of log reads or an error on an empty IndexPath.

## Links

### Out

- verifies [throwaway-index-fallback](../capability/rhiza/views/fast-reads/throwaway-index-fallback.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).

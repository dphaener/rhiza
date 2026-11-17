# interim-model-deleted-qa-2

evidence · `01M7KGW1HGBBSGXMMVYHBGWZYN` · actor [claude-code-69572f11](../actor/claude-code-69572f11.md)

## result

PASS

## statement

Scratch binary built from 96a0eabf3 (go build -trimpath -o $TMPDIR/rhiza-qa2), run in a scratch clone with TMPDIR set to an empty folder, once with no index path (XDG_CACHE_HOME=relative) and once with a read-only cache folder. Every read below runs through the projection with no fold fallback. go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestPropertyProjectionParity|TestProjection.*)$' passes (8.6 s). Error paths all answer cleanly with 0 leftovers: change show of an unknown ID gives not_found (next rhiza change ls), session show and run status of a missing ID give not_found, and slice --change unknown gives not_found.

## Links

### Out

- verifies [interim-model-deleted](../capability/rhiza/views/fast-reads/interim-model-deleted.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).

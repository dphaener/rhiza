# throwaway-index-fallback-qa-1

evidence · `01M7KB215YE4MSJ7KNFCPZJW76` · actor [claude-code-49a7d3f2](../actor/claude-code-49a7d3f2.md)

## result

PASS

## statement

QA, branch HEAD a2ef70341 built with go build -trimpath to $TMPDIR/rhiza-qa-orp, run in a hermetic scratch repo. With HOME and XDG_CACHE_HOME unset (empty IndexPath), status --json and show --json were byte-identical to the indexed run, add wrote a node that show then read, and TMPDIR was left empty. XDG_CACHE_HOME on a read-only folder: status --json identical, add succeeded, nothing left in the cache or TMPDIR. Failure: TMPDIR read-only or missing with no index path exits 1 'no index path; throwaway index in <tmp>: ...'; read-only cache plus read-only TMPDIR exits 1 naming both the index path and the temporary folder. An index file overwritten with garbage was rebuilt and show answered. The base build (9064447e1) refused the no-index case with 'no index path', so this is new behaviour. go test -run '^(TestProjection|TestPropertyProjectionParity)' ./internal/api ok.

## Links

### Out

- verifies [throwaway-index-fallback](../capability/rhiza/views/fast-reads/throwaway-index-fallback.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).

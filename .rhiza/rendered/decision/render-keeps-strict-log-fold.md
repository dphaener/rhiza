# render-keeps-strict-log-fold

decision · `01M7HCGGR3MYZ8WAHVE7RVG89G`

## reasoning

render must write what the log says, and the index never re-reads a path it has seen, so HoldsLog is the guard against a stale projection. That fold applies the strict engine's rules, so one set of model rules still holds; rebuilding the index there instead is index work outside this change set's intent.

## statement

render keeps its strict fold of the log (changeset.FoldLog then model.AsIs) when the index does not hold exactly the log's entries (an entry file edited in place); only its empty-IndexPath branch goes

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [one-read-path-plan-2026-11-16](../source/one-read-path-plan-2026-11-16.md)

### In

- decided-by [folds-from-projection](../capability/rhiza/views/fast-reads/folds-from-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).

# throwaway-index-in-api

decision · `01M7HCGGNXMWPVRRE138CBKYEJ`

## reasoning

Reads of other roots (the daemon's projects, a run worktree at a gate, the web's answer) pick their own index path inside internal/api and today fall back to an empty IndexPath and the log; a fallback in cmd/rhiza would cover only the caller's root. AGENTS.md keeps cmd/rhiza a thin frame, and change set 01M7GFTS0N rebuilds cmd/rhiza's ctx.go globals next, so the fallback stays out of its way. It refines projection-no-index-fallback, which says what happens and not where.

## statement

The throwaway index is made where the projection opens (internal/api's openWriteIndex, under every root's reads) when the call's index path is empty or the index there cannot be opened or created, not in cmd/rhiza/ctx.go

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [one-read-path-plan-2026-11-16](../source/one-read-path-plan-2026-11-16.md)

### In

- decided-by [throwaway-index-fallback](../capability/rhiza/views/fast-reads/throwaway-index-fallback.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).

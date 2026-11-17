# projection-no-index-fallback

decision · `01M7GFT0E7FBH3Q4JN308G3VYX`

## reasoning

Today cmd/rhiza/ctx.go leaves IndexPath empty when index.DefaultPath fails, and every op then folds the log through the interim model; deleting that model needs another answer. Failing the command would break a repo that works today, which is a breaking change; a throwaway index keeps the output identical and the code single.

## statement

When a worktree's index cannot be opened or created (no user cache folder, a read-only cache), a command builds a throwaway index in a temporary folder for that one call and reads through it, so there is still one read path and no fold fallback; the cost is a full build on that call

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

### In

- decided-by [throwaway-index-fallback](../capability/rhiza/views/fast-reads/throwaway-index-fallback.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
